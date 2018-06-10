package ishocon1

import org.scalatra.test.scalatest._

class ishocon1ServletTests extends ScalatraFunSuite {

  addServlet(classOf[ishocon1Servlet], "/*")

  test("GET / on ishocon1Servlet should return status 200") {
    get("/") {
      status should equal (200)
    }
  }

}
