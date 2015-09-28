require 'sinatra/base'
require 'mysql2'
require 'mysql2-cs-bind'
require 'erubis'

module Ishocon1
  class AuthenticationError < StandardError; end
  class PermissionDenied < StandardError; end
  # class ContentNotFound < StandardError; end
  module TimeWithoutZone
    def to_s
      strftime("%F %H:%M:%S")
    end
  end
  ::Time.prepend TimeWithoutZone
end

class Ishocon1::WebApp < Sinatra::Base
  use Rack::Session::Cookie
  set :erb, escape_html: true
  set :public_folder, File.expand_path('../public', __FILE__)
  set :session_secret, ENV['ISHOCON1_SESSION_SECRET'] || 'showwin_happy'
  set :protection, true

  helpers do
    def config
      @config ||= {
        db: {
          host: ENV['ISHOCON1_DB_HOST'] || 'localhost',
          port: ENV['ISHOCON1_DB_PORT'] && ENV['ISHOCON1_DB_PORT'].to_i,
          username: ENV['ISHOCON1_DB_USER'] || 'root',
          password: ENV['ISHOCON1_DB_PASSWORD'],
          database: ENV['ISHOCON1_DB_NAME'] || 'ishocon1'
        }
      }
    end

    def db
      return Thread.current[:ishocon1_db] if Thread.current[:ishocon1_db]
      client = Mysql2::Client.new(
        host: config[:db][:host],
        port: config[:db][:port],
        username: config[:db][:username],
        password: config[:db][:password],
        database: config[:db][:database],
        reconnect: true
      )
      client.query_options.merge!(symbolize_keys: true)
      Thread.current[:ishocon1_db] = client
      client
    end

    def authenticate(email, password)
      user = db.xquery('SELECT * FROM users WHERE email = ?', email).first
      raise Ishocon1::AuthenticationError unless user[:password] == password
      session[:user_id] = user[:id]
    end

    def authenticated!
      raise Ishocon1::PermissionDenied unless current_user
    end

    def current_user
      db.xquery('SELECT * FROM users WHERE id = ?', session[:user_id]).first
    end

    def buy_product(product_id, user_id)
      db.xquery('INSERT INTO histories (product_id, user_id, created_at) VALUES (?, ?, ?)', \
        product_id, user_id, Time.now)
    end

    def already_bought(product_id)
      return false unless current_user
      count = db.xquery('SELECT count(*) as count FROM histories WHERE product_id = ? AND user_id = ?', \
                        product_id, current_user[:id]).first[:count]
      count > 0
    end

    def create_comment(product_id, user_id, content)
      db.xquery('INSERT INTO comments (product_id, user_id, content, created_at) VALUES (?, ?, ?, ?)', \
        product_id, user_id, content, Time.now)
    end
  end

  error Ishocon1::AuthenticationError do
    session[:user_id] = nil
    halt 401, erubis(:login, layout: false, locals: { message: 'ログインに失敗しました' })
  end

  error Ishocon1::PermissionDenied do
    halt 403, erubis(:login, layout: false, locals: { message: '先にログインをしてください' })
  end

  get '/login' do
    session.clear
    erb :login, layout: false, locals: { message: 'ECサイトで爆買いしよう！！！！' }
  end

  post '/login' do
    authenticate params['email'], params['password']
    redirect '/'
  end

  get '/logout' do
    session[:user_id] = nil
    session.clear
    redirect '/login'
  end

  get '/' do
    page = params[:page].to_i || 0
    products = db.xquery("SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET #{page * 50}")
    erb :index, locals: { products: products }
  end

  get '/users/:user_id' do
    products_query = <<SQL
SELECT *
FROM histories as h
LEFT OUTER JOIN products as p
ON h.product_id = p.id
WHERE h.user_id = ?
ORDER BY h.id DESC
SQL
    products = db.xquery(products_query, params[:user_id])

    total_pay = 0
    products.each do |product|
      total_pay += product[:price]
    end

    user = db.xquery('SELECT * FROM users WHERE id = ?', params[:user_id]).first
    erb :mypage, locals: { products: products, user: user, total_pay: total_pay }
  end

  get '/products/:product_id' do
    product = db.xquery('SELECT * FROM products WHERE id = ?', params[:product_id]).first
    comments = db.xquery('SELECT * FROM comments WHERE product_id = ?', params[:product_id])
    erb :product, locals: { product: product, comments: comments }
  end

  post '/products/buy/:product_id' do
    authenticated!
    buy_product(params[:product_id], current_user[:id])
    redirect "/users/#{current_user[:id]}"
  end

  post '/comments/:product_id' do
    authenticated!
    create_comment(params[:product_id], current_user[:id], params[:content])
    redirect "/users/#{current_user[:id]}"
  end
end
