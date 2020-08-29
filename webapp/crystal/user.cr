class User
  getter id
  getter name

  def initialize(@id : Int32, @name : String, @email : String, @password : String, @last_login : Time)
  end

  def initialize(@id : Int32)
    @id, @name, @email, @password, @last_login = Database.query_one "SELECT * FROM users WHERE id = ? LIMIT 1", @id, as: {Int32, String, String, String, Time}
  end

  def buying_history
    products : Array(Product) = [] of Product
    Database.query "SELECT p.id, p.name, p.description, p.image_path, p.price, h.created_at " +
                   "FROM histories as h " +
                   "LEFT OUTER JOIN products as p " +
                   "ON h.product_id = p.id " +
                   "WHERE h.user_id = ? " +
                   "ORDER BY h.id DESC", @id do |rs|
      rs.each do
        id, name, description, image_path, price, created_at = rs.read(Int32, String, String, String, Int32, Time)
        products.push(Product.new(id, name, description, image_path, price, created_at))
      end
    end
    return products
  end

  def buy_product(pid : String)
    Database.exec "INSERT INTO histories (product_id, user_id, created_at) VALUES (?, ?, ?)", pid, @id, Time.local
  end

  def create_comment(pid : String, content : String)
    Database.exec "INSERT INTO comments (product_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", pid, @id, content, Time.local
  end

  def update_last_login
    Database.exec "UPDATE users SET last_login = ? WHERE id = ?", Time.local, @id
  end
end

def authenticate(email_arg : String, password_arg : String) : User?
  begin
    id, name, email, password, last_login = Database.query_one "SELECT * FROM users WHERE email = ? LIMIT 1", email_arg, as: {Int32, String, String, String, Time}
    if password != password_arg
      return nil
    end
    return User.new(id, name, email, password, last_login)
  rescue
    return nil
  end
end

def not_authenticated(uid : Int32) : Bool
  return !(uid > 0)
end
