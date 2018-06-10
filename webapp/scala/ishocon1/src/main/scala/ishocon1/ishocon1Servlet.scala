package ishocon1

import org.scalatra._
import org.scalatra.i18n.I18nSupport
import org.scalatra.forms.{views => _, _}
import slick.driver.MySQLDriver.api._
import javax.servlet.ServletContext
import scala.concurrent._
import scala.concurrent.duration.Duration
import java.time._
import java.sql._
import slick.jdbc.GetResult

case class User(id: Int, name: String, email: String, password: String, lastLogin: java.sql.Timestamp)
case class Product(id: Int, name: String, description: String, imagePath: String, price: Int, createdAt: java.sql.Timestamp)
case class Comment(id: Int, productId: Int, userId: Int, content: String, createdAt: java.sql.Timestamp)
case class CommentWithWriter(comment: Comment, writer: String)
case class ProductView(product: Product, commentCount: Int, comments: Vector[CommentWithWriter])

class ishocon1Servlet extends ScalatraServlet with FormSupport with I18nSupport{
  implicit val getUserResult = GetResult(r => User(r.nextInt, r.nextString, r.nextString, r.nextString, r.nextTimestamp))
  implicit val getProductResult = GetResult(r => Product(r.nextInt, r.nextString, r.nextString, r.nextString, r.nextInt, r.nextTimestamp))
  implicit val getCommentResult = GetResult(r => Comment(r.nextInt, r.nextInt, r.nextInt, r.nextString, r.nextTimestamp))
  implicit val getCommentWithWriterResult = GetResult(r => CommentWithWriter(Comment(r.nextInt, r.nextInt, r.nextInt, r.nextString, r.nextTimestamp), r.nextString))

  val dbType = "mysql"
  val dbHost = "localhost"
  val dbPort = 3306
  val dbUser = sys.env.getOrElse("ISHOCON1_DB_NAME", "ishocon")
  val dbPassword = sys.env.getOrElse("ISHOCON1_DB_NAME", "ishocon")
  val dbName = sys.env.getOrElse("ISHOCON1_DB_NAME", "ishocon1")
  val db = Database.forURL(
    "jdbc:%s://%s:%d/%s?characterEncoding=UTF-8&useSSL=false&useUnicode=true&useJDBCCompliantTimezoneShift=true&useLegacyDatetimeCode=false&serverTimezone=UTC".format(dbType, dbHost, dbPort, dbName),
    driver="com.mysql.cj.jdbc.Driver",
    user=dbUser,
    password=dbPassword)

  def timeNowDb(): Timestamp = {
    Timestamp.valueOf(LocalDateTime.now(ZoneId.of("Asia/Tokyo")))
  }

  def authenticate(email: String, password: String): Unit = {
    val query = db.run(sql"SELECT * FROM users WHERE email = $email".as[User])
    val users = Await.result(query, Duration.Inf)
    if((users.headOption match{
      case Some(u) => (u.password == password)
      case None => false
    }) == false)(throw new AuthenticationErrorException())
    
    session("uid") = users.head.id
  }

  def authenticated() = {
    if(currentUser.isEmpty)(throw new PermissionDeniedException())
  }

  def currentUser(): Option[User] = {
    if(!session.contains("uid"))(None)else{
      val userId = session("uid").asInstanceOf[Int]
      val query = db.run(sql"SELECT * FROM users WHERE id = $userId".as[User])
      val users = Await.result(query, Duration.Inf)
      users.headOption
    }
  }

  def updateLastLogin(userId: Int): Unit = {
    val now = timeNowDb
    val query = db.run(sqlu"UPDATE users SET last_login = $now WHERE id = $userId")
    Await.result(query, Duration.Inf)
  }

  def buyProduct(productId: Int, userId: Int): Unit = {
    val now = timeNowDb
    val query = db.run(sqlu"INSERT INTO histories (product_id, user_id, created_at) VALUES ($productId, $userId, $now)")
    Await.result(query, Duration.Inf)
  }

  def alreadyBought(productId: Int): Boolean = {
    currentUser match{
      case None => false
      case Some(user) =>{
        val userId = user.id
        val query = db.run(sql"SELECT count(*) as count FROM histories WHERE product_id = $productId AND user_id = $userId".as[Int])
        val count = Await.result(query, Duration.Inf).head
        count > 0
      }
    }
  }

  def createComment(productId: Int, userId: Int, content: String): Unit = {
    val now = timeNowDb
    val query = db.run(sqlu"INSERT INTO comments (product_id, user_id, content, created_at) VALUES ($productId, $userId, $content, $now)")
    Await.result(query, Duration.Inf)
  }

  case class AuthenticationErrorException() extends Exception()
  case class PermissionDeniedException() extends Exception()
  val errorHandling: PartialFunction[Throwable, Unit] = {
    case e:AuthenticationErrorException => {
      session.remove("uid")
      halt(401, views.html.login("ログインに失敗しました"))
    }
    case e:PermissionDeniedException => {
      halt(403, views.html.login("先にログインをしてください"))
    }
  }

  get("/login") {
    if(session.contains("uid"))session.remove("uid")
    views.html.login("ECサイトで爆買いしよう！！！！")
  }

  case class ValidationForm(
    email: String,
    password: String,
  )
  val loginform = mapping(
    "email" -> text(),
    "password" -> text()
  )(ValidationForm.apply)
  post("/login") {
    try{
      validate(loginform)(
        errors => (throw new  AuthenticationErrorException()),
        form => {
          authenticate(form.email, form.password)
          
          val User = currentUser.get
          updateLastLogin(User.id)
          redirect("/")
        }
      )
    } catch errorHandling
  }

  get("/logout") {
    if(session.contains("uid"))(session.remove("uid"))
    redirect("/login")
  }

  get("/") {
    val cUser = currentUser
    val page = params.getOrElse("page", "0").asInstanceOf[String].toInt * 50
    val query1 = db.run(sql"SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET $page".as[Product])
    val productsRaw = Await.result(query1, Duration.Inf)
    val productsView: Vector[ProductView] = productsRaw.map(p => {
      val pId = p.id
      val query2 = db.run(sql"SELECT count(*) as count FROM comments WHERE product_id = $pId".as[Int])
      val commentCount =  Await.result(query2, Duration.Inf).head

      val query3 = db.run(sql"SELECT * FROM comments as c INNER JOIN users as u ON c.user_id = u.id WHERE c.product_id = $pId ORDER BY c.created_at DESC LIMIT 5".as[CommentWithWriter])
      val commentWithWriters = Await.result(query3, Duration.Inf)

      ProductView(
        p.copy(description = p.description.take(70)),
        commentCount,
        commentWithWriters.map(
          commentWithWriter => commentWithWriter.copy(comment = commentWithWriter.comment.copy(content = commentWithWriter.comment.content.take(25)))
        )
      )
    })
    views.html.layout(cUser)(views.html.index(cUser, productsView))
  }

  get("/users/:userId") {
    val cUser = currentUser
    val userId = params("userId")
    val query1 = db.run(sql"""
SELECT p.id, p.name, p.description, p.image_path, p.price, h.created_at
FROM histories as h
LEFT OUTER JOIN products as p
ON h.product_id = p.id
WHERE h.user_id = $userId
ORDER BY h.id DESC
    """.as[Product])
    val products = Await.result(query1, Duration.Inf)

    val totalPay = products.map(product => product.price).sum

    val query2 = db.run(sql"SELECT * FROM users WHERE id = $userId".as[User])
    val user = Await.result(query2, Duration.Inf).head

    views.html.layout(cUser)(views.html.mypage(cUser, user, products))
  }

  get("/products/:productId") {
    val cUser = currentUser
    val productId = params("productId")
    val query1 = db.run(sql"SELECT * FROM products WHERE id = $productId".as[Product])
    val product = Await.result(query1, Duration.Inf).head

    val query2 = db.run(sql"SELECT * FROM comments WHERE product_id = $productId".as[Comment])
    val comments = Await.result(query2, Duration.Inf)
    views.html.layout(cUser)(views.html.product(product, alreadyBought(product.id)))
  }

  post("/products/buy/:productId") {
    try {
      authenticated
      val cUser = currentUser
      val productId = params("productId")
      buyProduct(productId.toInt, cUser.get.id)
      redirect("/users/"+(cUser.get.id.toString))
    } catch errorHandling
  }

  post("/comments/:productId") {
    try {
      authenticated
      val cUser = currentUser
      val productId = params("productId")
      val content = params("content")
      createComment(productId.toInt, cUser.get.id, content)
      redirect("/users/"+(cUser.get.id.toString))
    } catch errorHandling
  }

  get("/initialize") {
    val query1 = db.run(sqlu"DELETE FROM users WHERE id > 5000")
    Await.result(query1, Duration.Inf)
    val query2 = db.run(sqlu"DELETE FROM products WHERE id > 10000")
    Await.result(query2, Duration.Inf)
    val query3 = db.run(sqlu"DELETE FROM comments WHERE id > 200000")
    Await.result(query3, Duration.Inf)
    val query4 = db.run(sqlu"DELETE FROM histories WHERE id > 500000")
    Await.result(query4, Duration.Inf)
  }
}
