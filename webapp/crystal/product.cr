require "mysql"

class Product
	getter id : Int32
	getter name : String
	getter description : String
	getter image_path : String
	getter price : Int32
	getter created_at : Time
	def initialize(@id : Int32, @name : String, @description : String, @image_path : String, @price : Int32, @created_at : Time)
	end
	def initialize(id : Int32)
		id, name, description, image_path, price, created_at = Database.query_one "SELECT * FROM products WHERE id = ? LIMIT 1", id, as: {Int32, String, String, String, Int32, Time}
		initialize(id, name, description, image_path, price, created_at)
	end

	def is_bought(uid : Int32)
		count = Database.query_one "SELECT count(*) as count FROM histories WHERE product_id = ? AND user_id = ?", @id, uid, as: {Int64}
		return count > 0
	end
end

class ProductWithComments
	getter id : Int32
	getter name : String
	getter description : String
	getter image_path : String
	getter price : Int32
	getter created_at : Time
	getter comment_count : Int64
	getter comments : Array(CommentWriter)
	def initialize(@id : Int32, @name : String, @description : String, @image_path : String, @price : Int32, @created_at : Time, @comment_count : Int64, @comments : Array(CommentWriter))
	end
	def initialize(product : Product)
		@id = product.id
		@name = product.name
		@description = product.description
		@image_path = product.image_path
		@price = product.price
		@created_at = product.created_at
		@comment_count = Database.query_one "SELECT count(*) as count FROM comments WHERE product_id = ?", product.id, as:{Int64}
		@comments = [] of CommentWriter
		if @comment_count > 0
			Database.query "SELECT * FROM comments as c INNER JOIN users as u ON c.user_id = u.id WHERE c.product_id = ? ORDER BY c.created_at DESC LIMIT 5", product.id do |rs|
				rs.each do
					i, i, i, content, t, i, writer, s, s, t = rs.read(Int32, Int32, Int32, String, Time, Int32, String, String, String, Time)
					@comments.push(CommentWriter.new(content, writer))
				end
			end
		end
	end
end

def get_products_with_comments_at(page : Int32) : Array(ProductWithComments)
	ret : Array(ProductWithComments) = [] of ProductWithComments
	Database.query "SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET ?", page*50 do |rs|
		rs.each do
			id, name, description, image_path, price, created_at = rs.read(Int32, String, String, String, Int32, Time)
			ret.push(ProductWithComments.new(Product.new(id, name, description, image_path, price, created_at)))
		end
	end
	return ret
end

class CommentWriter
	getter content : String
	getter writer : String
	def initialize(@content : String, @writer : String)
	end
end
