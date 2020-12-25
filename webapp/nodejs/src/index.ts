import mysql from "mysql";
import util from "util";
import express from "express";
import session from "express-session";
import bodyParser from "body-parser";
import morgan from "morgan";
import { AddressInfo } from "net";

const app = express();

app.set("views", "views");
app.set("view engine", "ejs");
app.use(express.static("public"));
app.use(morgan("tiny"));
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

const pool = mysql.createPool({
  connectionLimit: 20,
  host: process.env.ISHOCON1_DB_HOST || "localhost", // not supported in other implementations
  port: parseInt(process.env.ISHOCON1_DB_PORT ?? "3306"), // not supported in other implementations
  user: process.env.ISHOCON1_DB_USER || "ishocon",
  password: process.env.ISHOCON1_DB_PASSWORD || "ishocon",
  database: process.env.ISHOCON1_DB_NAME || "ishocon1",
});
const query = util.promisify(pool.query.bind(pool));

app.get("/login", (req, res) => {
  req.session.destroy(() => {});
  res.render("./login.ejs", { message: "ECサイトで爆買いしよう！！！！" });
});

type User = {
  password: string;
  // there're other fields used in views. see views/*.ejs
};
async function authenticate(email: string, password: string) {
  const rows = (await query({
    sql: "SELECT * FROM users WHERE email = ?",
    values: [email],
  })) as User[];
  if (!rows[0] || rows[0].password !== password) {
    throw new Error("user not found");
  }
  return rows[0];
}

async function getUser(userId) {
  const rows = (await query({
    sql: "SELECT * FROM users WHERE id = ? LIMIT 1",
    values: [userId],
  })) as User[];
  return rows[0];
}

async function currentUser(req) {
  const userId = req.session.uid;
  if (!userId) {
    return undefined;
  }
  return await getUser(userId);
}

type ProductRow = {
  id: string;
  // there're other fields used in views. see views/*.ejs
};
type Product = ProductRow & {
  commentsCount: number;
  comments: Comment[];
};
type Comment = {
  name: string;
  content: string;
};
async function getProducts(page) {
  const rows = (await query({
    sql: "SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET ?",
    values: [page * 50],
  })) as ProductRow[];
  const products: Product[] = [];
  for (const row of rows) {
    const cc = (await query({
      sql: "SELECT count(*) as count FROM comments WHERE product_id = ?",
      values: [row.id],
    })) as { count: number }[];
    const commentsCount = cc[0].count;
    const comments: Comment[] = [];
    if (commentsCount > 0) {
      const subrows = (await query({
        sql:
          "SELECT * FROM comments as c INNER JOIN users as u ON c.user_id = u.id WHERE c.product_id = ? ORDER BY c.created_at DESC LIMIT 5",
        values: [row.id],
      })) as Comment[];
      for (const subrow of subrows) {
        comments.push({
          content: subrow.content,
          name: subrow.name,
        });
      }
    }
    products.push({
      ...row,
      commentsCount,
      comments,
    });
  }
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
  req.session["uid"] = user.id;
  await query({
    sql: "UPDATE users SET last_login = ? WHERE id = ?",
    values: [new Date(), user.id],
  });

  res.redirect(303, "/");
});

app.get("/logout", (req, res) => {
  req.session.destroy(() => {});

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
  const userId = parseInt(req.params.userId);
  // TODO implement
  res.send(`hello user ${userId}`);
});

app.get("/products/:productId", (req, res) => {
  const productId = parseInt(req.params.productId);
  // TODO implement
  res.send(`product ${productId}`);
});

app.post("/products/buy/:productId", (req, res) => {
  const productId = parseInt(req.params.productId);
  // TODO implement
  res.send(`you bought product ${productId}`);
});

app.post("/comments/:productId", (req, res) => {
  const productId = parseInt(req.params.productId);
  // TODO implement
  res.send(`you commented product ${productId}`);
});

const server = app.listen(8080, function () {
  const address = server.address() as AddressInfo;
  const host = address.address;
  const port = address.port;

  console.log("Example app listening at http://%s:%s", host, port);
});
