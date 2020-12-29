import mysql from "mysql";
import util from "util";
import express from "express";
import session from "express-session";
import bodyParser from "body-parser";
import morgan from "morgan";
import { AddressInfo } from "net";
import { stringify } from "querystring";

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

declare module "express-session" {
  interface SessionData {
    uid: string;
  }
}

type User = {
  id: string;
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

async function getUser(userId: string) {
  const rows = (await query({
    sql: "SELECT * FROM users WHERE id = ? LIMIT 1",
    values: [userId],
  })) as User[];
  return rows[0];
}

async function currentUser(req: express.Request) {
  const userId = req.session.uid;
  if (!userId) {
    return undefined;
  }
  return await getUser(userId);
}

type ProductHistory = {
  // same as ProductRow
  id: string;
  name: string;
  description: string;
  image_path: string;
  price: number;
  // history
  created_at: Date;
};
async function buyingHistory(user: User) {
  return (await query({
    sql:
      "SELECT p.id, p.name, p.description, p.image_path, p.price, h.created_at " +
      "FROM histories as h " +
      "LEFT OUTER JOIN products as p " +
      "ON h.product_id = p.id " +
      "WHERE h.user_id = ? " +
      "ORDER BY h.id DESC",
    values: [user.id],
  })) as ProductHistory[];
}

type ProductRow = {
  id: string;
  // there're other fields used in views. see views/*.ejs
};
type Product = ProductRow & {
  commentsCount: number;
  comments: {
    name: string;
    content: string;
  }[];
};
async function getProducts(page: number) {
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
    const comments: Product["comments"] = [];
    if (commentsCount > 0) {
      const subrows = (await query({
        sql:
          "SELECT * FROM comments as c INNER JOIN users as u ON c.user_id = u.id WHERE c.product_id = ? ORDER BY c.created_at DESC LIMIT 5",
        values: [row.id],
      })) as Product["comments"];
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

async function getProduct(productId: string) {
  const rows = (await query({
    sql: "SELECT * FROM products WHERE id = ? LIMIT 1",
    values: [productId],
  })) as ProductRow[];
  return rows[0];
}

type Comment = {
  // has more fields used in view
};
async function getComments(productId: string) {
  const rows = (await query({
    sql: "SELECT * FROM comments WHERE product_id = ?",
    values: [productId],
  })) as Comment[];
  return rows[0];
}

async function isBought(productId: string, userId: string) {
  const cc = (await query({
    sql:
      "SELECT count(*) as count FROM histories WHERE product_id = ? AND user_id = ?",
    values: [productId, userId],
  })) as { count: number }[];
  return cc[0].count > 0;
}

async function buyProduct(productId: string, userId: string) {
  await query({
    sql:
      "INSERT INTO histories (product_id, user_id, created_at) VALUES (?, ?, ?)",
    values: [productId, userId, new Date()],
  });
}

async function createComment(
  productId: string,
  userId: string,
  content: string
) {
  await query({
    sql:
      "INSERT INTO comments (product_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
    values: [productId, userId, content, new Date()],
  });
}

app.post("/login", async (req, res) => {
  req.session.regenerate(() => {});

  let user: User | undefined;
  try {
    user = await authenticate(req.body.email, req.body.password);
  } catch (e) {
    res.render("./login.ejs", { message: "ログインに失敗しました" });
    return;
  }
  req.session.uid = user.id;
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

app.get("/users/:userId", async (req, res) => {
  type ProductHistoryWithStringCreatedAt = {
    // same as ProductRow
    id: string;
    name: string;
    description: string;
    image_path: string;
    price: number;
    // history
    created_at: string;
  };

  const cUser = await currentUser(req);
  const user = await getUser(req.params.userId);
  const products = (await buyingHistory(user)).map((p) => {
    // to format as "2006-01-02 15:04:05" in JST
    const created_at = new Date(p.created_at.getTime() + 9 * 3600 * 1000)
      .toISOString()
      .replace("T", " ")
      .replace(/\.\d+Z/, "");
    return { ...p, created_at: created_at };
  });

  let totalPay = 0;
  products.forEach((p) => {
    totalPay += p.price;
  });
  res.render("./mypage.ejs", {
    user: user,
    products: products,
    current_user: cUser,
    totalPay: totalPay,
  });
});

app.get("/products/:productId", async (req, res) => {
  const productId = req.params.productId;
  const product = await getProduct(productId);
  const comments = await getComments(productId);
  const user = await currentUser(req);
  const bought = user ? await isBought(productId, user.id) : false;
  res.render("./product.ejs", {
    current_user: user,
    product: product,
    comments: comments,
    already_bought: bought,
  });
});

app.post("/products/buy/:productId", async (req, res) => {
  const user = await currentUser(req);
  if (!user) {
    return res.render("./login.ejs", { message: "先にログインをしてください" });
  }
  await buyProduct(req.params.productId, user.id);
  res.redirect(303, `/users/${user.id}`);
});

app.post("/comments/:productId", async (req, res) => {
  const user = await currentUser(req);
  if (!user) {
    return res.render("./login.ejs", { message: "先にログインをしてください" });
  }
  await createComment(req.params.productId, user.id, req.body.content);
  res.redirect(303, `/users/${user.id}`);
});

const server = app.listen(8080, function () {
  const address = server.address() as AddressInfo;
  const host = address.address;
  const port = address.port;

  console.log("Example app listening at http://%s:%s", host, port);
});
