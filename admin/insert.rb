require 'mysql2'
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
  Thread.current[:ishocon1] = client
  client
end

def time_now
  Time.now.strftime("%Y-%m-%d %H:%M:%S")
end

def insert_user
  o = [('a'..'z'), ('A'..'Z'), ('0'..'9')].map { |i| i.to_a }.flatten
  db.query("INSERT INTO users (name, email, password, last_login) VALUES ('ishocon', 'ishocon@isho.con', 'ishoconpass', '#{time_now}')")
  query = ''
  4999.times do |i|
    if i%100 == 0
      db.query(query[0..-3]) unless query == ''
      query = 'INSERT INTO users (name, email, password, last_login) VALUES '
    end
    name = %w(さとう すずき たかはし たなか わたなべ いとう やまもと).sample + "#{i}号"
    email = "ishocon#{i}@isho.con"
    pass = (0...16).map { o[rand(o.length)] }.join
    last_login = time_now
    query << "('#{name}', '#{email}', '#{pass}', '#{last_login}'), "
  end
  db.query(query[0..-3])
end

def insert_products
  query = ''
  10000.times do |i|
    if i%100 == 0
      db.query(query[0..-3]) unless query == ''
      query = 'INSERT INTO products (name, description, image_path, price, created_at) VALUES '
    end
    name = %w(すごい やばい まずい おいしい きもちわるい ぴかぴかな まろやかな).sample + "商品-#{i}号"
    description = "これはかなり#{name}です。取り扱いには十分ご注意ください\n" * 25
    image_path = '/images/image' + (i % 5).to_s + '.jpg'
    price = rand(500) * 10
    created_at = time_now
    query << "('#{name}', '#{description}', '#{image_path}', #{price}, '#{created_at}'), "
  end
  db.query(query[0..-3])
end

def insert_comments
  query = ''
  20.times do |i|
    10000.times do |j|
      if j%1000 == 0
        db.query(query[0..-3]) unless query == ''
        query = 'INSERT INTO comments (product_id, user_id, content, created_at) VALUES '
      end
      product_id = j + 1
      user_id = rand(5000) + 1
      content = %w(
        これはすごいですね、感動しました。
        ほんとうに買ってよかったと思います。
        友人の友人にも勧めます。
        1週間ですぐに壊れました。
        買ってからずっと引き出しの中です。
        日本製の製品はさすがに質が高いですね。
        嬉しくて夜も眠れません。
      ).sample(3).join
      created_at = time_now
      query << "(#{product_id}, #{user_id}, '#{content}', '#{created_at}'), "
    end
    db.query(query[0..-3])
    query = ''
  end
end

def insert_histories
  query = ''
  5000.times do |i|
    100.times do |j|
      query = 'INSERT INTO histories (product_id, user_id, created_at) VALUES ' if j == 0
      product_id = rand(10000)
      user_id = i + 1
      created_at = time_now
      query << "(#{product_id}, #{user_id}, '#{created_at}'), "
    end
    db.query(query[0..-3])
  end
end

def insert
  insert_user
  insert_products
  insert_comments
  insert_histories
end

insert
