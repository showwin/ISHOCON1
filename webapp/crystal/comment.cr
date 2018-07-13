class Comment
  def initialize(@id : Int32, @product_id : Int32, @userid : Int32, @content : String, @createdat : Time)
  end
end

def get_comments(pid : Int32)
  comments : Array(Comment) = [] of Comment
  Database.query "SELECT * FROM comments WHERE product_id = ?", pid do |rs|
    rs.each do
      id, product_id, user_id, content, createdat = rs.read(Int32, Int32, Int32, String, Time)
      comments.push(Comment.new(id, product_id, user_id, content, createdat))
    end
  end
  return comments
end
