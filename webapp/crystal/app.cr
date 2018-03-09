require "kemal"
require "kemal-session"
require "ecr"

require "./comment"
require "./user"
require "./product"
require "./db"

Kemal::Session.config do |config|
	config.cookie_name = "showwin_happy"
	config.secret = "mysession"
end

public_folder "public"

get "/login" do |env|
	env.session.destroy
	message = "ECサイトで爆買いしよう！！！！"

	render "views/login.ecr"
end

post "/login" do |env|
	user = authenticate(env.params.body["email"], env.params.body["password"])
	if user
		env.session.int("uid", user.id)
		user.update_last_login

		env.redirect "/"
	else
		message = "ログインに失敗しました"
		render "views/login.ecr"
	end
end

get "/logout" do |env|
	env.session.destroy
	env.redirect "/login"
end

get "/" do |env|
	uid = env.session.int?("uid")
	c_user = User.new(uid) if uid

	page = (env.params.query["page"]? || "0").to_i

	products = get_products_with_comments_at(page)
	s_products = [] of ProductWithComments

	products.each do |p|
		s_description : String = p.description
		if p.description.size > 70
			s_description = p.description[0, 70] + "…"
		end

		new_c_w : Array(CommentWriter) = [] of CommentWriter
		p.comments.each do |c|
			if c.content.size > 25
				new_c_w.push(CommentWriter.new(c.content[0, 25] + "…", c.writer))
			else
				new_c_w.push(CommentWriter.new(c.content, c.writer))
			end
		end
		s_products.push(ProductWithComments.new(p.id, p.name, s_description, p.image_path, p.price, p.created_at, p.comment_count, new_c_w))
	end

	render "views/index.ecr", "views/layout.ecr"
end

get "/users/:user_id" do |env|
	cuid = env.session.int?("uid")
	c_user = User.new(cuid) if cuid

	uid = (env.params.url["user_id"]? || "0").to_i
	user = User.new(uid)

	products = user.buying_history

	total_pay = 0
	products.each do |product|
		total_pay += product.price
	end
	sd_products = [] of Product
	products.each do |p|
		sd_description : String = p.description
		if p.description.size > 70
			sd_description = p.description[0,70] + "…"
		end

		sd_products.push(Product.new(p.id, p.name, sd_description, p.image_path, p.price, p.created_at))
	end

	render "views/mypage.ecr", "views/layout.ecr"
end

get "/products/:product_id" do |env|
	pid = (env.params.url["product_id"]? || "0").to_i
	product = Product.new(pid)
	comments = get_comments(pid)

	cuid = env.session.int?("uid")
	c_user = User.new(cuid) if cuid
	bought : Bool = false
	if c_user
		bought = product.is_bought(c_user.id)
	end

	render "views/product.ecr", "views/layout.ecr"
end

post "/products/buy/:product_id" do |env|
	cuid = env.session.int?("uid")
	if !cuid || not_authenticated(cuid)
		message = "先にログインをしてください"
		render "views/login.ecr"
	else
		c_user = User.new(cuid)
		pid = env.params.url["product_id"]? || "0"
		c_user.buy_product(pid)

		env.redirect "/users/"+c_user.id.to_s
	end
end

post "/comments/:product_id" do |env|
	cuid = env.session.int?("uid")
	if !cuid || not_authenticated(cuid)
		message = "先にログインをしてください"
		render "views/login.ecr"
	else
		c_user = User.new(cuid)
		pid = env.params.url["product_id"]? || "0"
		comment = env.params.body["content"]? || ""
		c_user.create_comment(pid, comment)

		env.redirect "/users/"+c_user.id.to_s
	end
end

get "/initialize" do
	Database.exec "DELETE FROM users WHERE id > 5000"
	Database.exec "DELETE FROM products WHERE id > 10000"
	Database.exec "DELETE FROM comments WHERE id > 200000"
	Database.exec "DELETE FROM histories WHERE id > 500000"

	"Finish"
end


Kemal.run 8080
