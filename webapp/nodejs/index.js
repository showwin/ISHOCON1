const util = require("util");

const express = require("express");
const session = require("express-session");
const bodyParser = require("body-parser");

const app = express();

app.set("views", __dirname + "/views");
app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.use(
  session({
    secret: "showwin_happy",
    name: "ishocon1_node_session",
    resave: false,
    saveUninitialized: false,
  })
);

const mysql = require("mysql");
const pool = mysql.createPool({
  connectionLimit: 20,
  host: process.env.ISHOCON1_DB_HOST || "localhost", // not supported in other implementations
  port: process.env.ISHOCON1_DB_PORT || 3306, // not supported in other implementations
  user: process.env.ISHOCON1_DB_USER || "ishocon",
  password: process.env.ISHOCON1_DB_PASSWORD || "ishocon",
  database: process.env.ISHOCON1_DB_NAME || "ishocon1",
});
const query = util.promisify(pool.query).bind(pool);

app.get("/login", (req, res) => {
  req.session.destroy();
  res.render("./login.ejs", { message: "ECサイトで爆買いしよう！！！！" });
});

async function authenticate(email, password) {
  const rows = await query("SELECT * FROM users WHERE email = ?", email);
  if (!rows[0] || rows[0].password !== password) {
    throw new Error("user not found");
  }
  return rows[0];
}

async function getUser(userId) {
  const rows = await query("SELECT * FROM users WHERE id = ? LIMIT 1", userId);
  return rows[0];
}

async function currentUser(req) {
  const userId = req.session.uid;
  if (!userId) {
    return undefined;
  }
  return await getUser(userId);
}

async function getProducts(page) {
  const rows = await query(
    "SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET ?",
    page * 50
  );
  const products = await Promise.all(
    rows.map(async (row) => {
      // console.log(row);
      const cc = await query(
        "SELECT count(*) as count FROM comments WHERE product_id = ?",
        row.id
      );
      commentsCount = cc[0].count;
      let comments = []; // fill default
      if (commentsCount > 0) {
        const subrows = await query(
          "SELECT * FROM comments as c INNER JOIN users as u ON c.user_id = u.id WHERE c.product_id = ? ORDER BY c.created_at DESC LIMIT 5",
          row.id
        );
        comments = subrows.map((subrow) => ({
          content: subrow.content,
          name: subrow.name,
        }));
      }
      return {
        ...row,
        comments_count: commentsCount,
        comments,
      };
    })
  );
  return products;
}

app.post("/login", async (req, res) => {
  req.session.regenerate(() => {});

  let user;
  try {
    user = await authenticate(req.body.email, req.body.password);
  } catch (e) {
    res.render("./login.ejs", { message: "ログインに失敗しました" });
    return;
  }
  req.session.uid = user.id;
  await query("UPDATE users SET last_login = ? WHERE id = ?", [
    new Date(),
    user.id,
  ]);

  res.redirect(303, "/");
});

app.get("/logout", (req, res) => {
  req.session.destroy();

  res.redirect(303, "/login");
});

app.get("/initialize", async (_, res) => {
  await query("DELETE FROM users WHERE id > 5000");
  await query("DELETE FROM products WHERE id > 10000");
  await query("DELETE FROM comments WHERE id > 200000");
  await query("DELETE FROM histories WHERE id > 500000");
  res.send("Finish");
});

app.get("/", async (req, res) => {
  const user = await currentUser(req);

  let page = parseInt(req.params.page);
  if (Number.isNaN(page)) {
    page = 0;
  }
  const products = await getProducts(page);
  res.render("./index.ejs", { products: products, current_user: user });
});

app.get("/users/:userId", (req, res) => {
  userId = parseInt(req.params.userId);
  // TODO implement
  res.send(`hello user ${userId}`);
});

app.get("/products/:productId", (req, res) => {
  productId = parseInt(req.params.productId);
  // TODO implement
  res.send(`product ${productId}`);
});

app.post("/products/buy/:productId", (req, res) => {
  productId = parseInt(req.params.productId);
  // TODO implement
  res.send(`you bought product ${productId}`);
});

app.post("/comments/:productId", (req, res) => {
  productId = parseInt(req.params.productId);
  // TODO implement
  res.send(`you commented product ${productId}`);
});

const server = app.listen(8080, function () {
  const host = server.address().address;
  const port = server.address().port;

  console.log("Example app listening at http://%s:%s", host, port);
});
