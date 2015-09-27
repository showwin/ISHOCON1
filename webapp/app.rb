require 'sinatra/base'
require 'mysql2'
require 'mysql2-cs-bind'

module Ishocon1
  # class AuthenticationError < StandardError; end
  # class PermissionDenied < StandardError; end
  # class ContentNotFound < StandardError; end
  # module TimeWithoutZone
  #  def to_s
  #    strftime("%F %H:%M:%S")
  #  end
  # end
  ::Time.prepend TimeWithoutZone
end

class Ishocon::WebApp < Sinatra::Base
  use Rack::Session::Cookie
  set :erb, escape_html: true
  set :public_folder, File.expand_path('./public', __FILE__)
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
      client = Mysql2::Client.new(
        host: config[:db][:host],
        port: config[:db][:port],
        username: config[:db][:username],
        password: config[:db][:password],
        database: config[:db][:database],
        reconnect: true
      )
      client.query_options.merge!(symbolize_keys: true)
      client
    end
  end

  get '/login' do
    session.clear
    erb :login, locals: { message: 'ECサイトで爆買いしよう！！！！' }
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
    erb :index
  end
end
