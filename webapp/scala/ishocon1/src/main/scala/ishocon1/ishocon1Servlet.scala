package ishocon1

import org.scalatra._

class ishocon1Servlet extends ScalatraServlet {

  get("/") {
    views.html.hello()
  }

}
