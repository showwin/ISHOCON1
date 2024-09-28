import datetime
import os
import pathlib

import MySQLdb
from flask import Flask, abort, redirect, render_template, request, session
from MySQLdb.constants import FIELD_TYPE
from MySQLdb.cursors import DictCursor

static_folder = pathlib.Path(__file__).resolve().parent / "public"
print(static_folder)
app = Flask(__name__, static_folder=str(static_folder), static_url_path="")

app.secret_key = os.environ.get("ISHOCON1_SESSION_SECRET", "showwin_happy")

_config = {
    "db_host": os.environ.get("ISHOCON1_DB_HOST", "127.0.0.1"),
    "db_port": int(os.environ.get("ISHOCON1_DB_PORT", "3306")),
    "db_username": os.environ.get("ISHOCON1_DB_USER", "ishocon"),
    "db_password": os.environ.get("ISHOCON1_DB_PASSWORD", "ishocon"),
    "db_database": os.environ.get("ISHOCON1_DB_NAME", "ishocon1"),
}


def config(key):
    if key in _config:
        return _config[key]
    else:
        raise "config value of %s undefined" % key


def db():
    if hasattr(request, "db"):
        return request.db
    else:
        db = MySQLdb.connect(
            host=config("db_host"),
            port=config("db_port"),
            user=config("db_username"),
            password=config("db_password"),
            database=config("db_database"),
            charset="utf8mb4",
            conv={
                FIELD_TYPE.TINY: int,
                FIELD_TYPE.SHORT: int,
                FIELD_TYPE.LONG: int,
                FIELD_TYPE.INT24: int,
            },
            cursorclass=DictCursor,
        )
        cur = db.cursor()
        cur.execute(
            "SET SESSION sql_mode='TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY'"
        )
        cur.execute("SET NAMES utf8mb4")
        request.db = db
        return request.db


@app.teardown_request
def close_db(exception=None):
    if hasattr(request, "db"):
        request.db.close()


def to_jst(datetime_utc):
    return datetime.datetime.strptime(
        datetime_utc, "%Y-%m-%d %H:%M:%S"
    ) + datetime.timedelta(hours=9)


def to_utc(datetime_jst):
    return datetime_jst - datetime.timedelta(hours=9)


def authenticate(email, password):
    cur = db().cursor()
    cur.execute("SELECT * FROM users WHERE email = %s", (email,))
    user = cur.fetchone()
    if user is None or user.get("password", None) != password:
        abort(401)
    else:
        session["user_id"] = user["id"]


def authenticated():
    if not current_user():
        abort(401)


def current_user():
    if "user_id" in session:
        cur = db().cursor()
        cur.execute("SELECT * FROM users WHERE id = %s", (str(session["user_id"]),))
        return cur.fetchone()
    else:
        return None


def update_last_login(user_id):
    cur = db().cursor()
    cur.execute(
        "UPDATE users SET last_login = %s WHERE id = %s",
        (
            datetime.datetime.now(),
            user_id,
        ),
    )
    db().commit()


def get_comments(product_id):
    cur = db().cursor()
    cur.execute(
        """
SELECT *
FROM comments as c
INNER JOIN users as u
ON c.user_id = u.id
WHERE c.product_id = {}
ORDER BY c.created_at DESC
LIMIT 5
""".format(product_id)
    )

    return cur.fetchall()


def get_comments_count(product_id):
    cur = db().cursor()
    cur.execute(
        "SELECT count(*) as count FROM comments WHERE product_id = {}".format(
            product_id
        )
    )
    return int(cur.fetchone()["count"])


def buy_product(product_id, user_id):
    cur = db().cursor()
    cur.execute(
        "INSERT INTO histories (product_id, user_id, created_at) VALUES ({}, {}, '{}')".format(
            product_id,
            user_id,
            to_utc(datetime.datetime.now()).strftime("%Y-%m-%d %H:%M:%S"),
        )
    )
    db().commit()


def already_bought(product_id):
    if not current_user():
        return False
    cur = db().cursor()
    cur.execute(
        "SELECT count(*) as count FROM histories WHERE product_id = %s AND user_id = %s",
        (product_id, current_user()["id"]),
    )
    return int(cur.fetchone()["count"]) > 0


def create_comment(product_id, user_id, content):
    cur = db().cursor()
    cur.execute(
        """
INSERT INTO comments (product_id, user_id, content, created_at)
VALUES ({}, {}, '{}', '{}')
""".format(
            product_id,
            user_id,
            content,
            to_utc(datetime.datetime.now()).strftime("%Y-%m-%d %H:%M:%S"),
        )
    )
    db().commit()


@app.errorhandler(401)
def authentication_error(error):
    return render_template("login.html", message="ログインに失敗しました"), 401


@app.errorhandler(403)
def authentication_error(error):
    return render_template("login.html", message="先にログインをしてください"), 403


@app.route("/login")
def get_login():
    session.pop("user_id", None)
    return render_template("login.html", message="ECサイトで爆買いしよう！！！！")


@app.route("/login", methods=["POST"])
def post_login():
    authenticate(request.form["email"], request.form["password"])
    update_last_login(current_user()["id"])
    return redirect("/")


@app.route("/logout")
def get_logout():
    session.pop("user_id", None)
    return redirect("/login")


@app.route("/")
def get_index():
    page = int(request.args.get("page", 0))
    cur = db().cursor()
    cur.execute(
        "SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET {}".format(page * 50)
    )
    products = cur.fetchall()

    for product in products:
        product["description"] = product["description"][:70]
        product["created_at"] = to_jst(product["created_at"].decode())
        product["comments"] = get_comments(product["id"])
        product["comments_count"] = get_comments_count(product["id"])

    return render_template("index.html", products=products, current_user=current_user())


@app.route("/users/<int:user_id>")
def get_mypage(user_id):
    cur = db().cursor()
    cur.execute(
        """
SELECT p.id, p.name, p.description, p.image_path, p.price, h.created_at
FROM histories as h
LEFT OUTER JOIN products as p
ON h.product_id = p.id
WHERE h.user_id = {}
ORDER BY h.id DESC
""".format(str(user_id))
    )

    products = cur.fetchall()

    total_pay = 0
    for product in products:
        total_pay += product["price"]
        product["description"] = product["description"][:70]
        product["created_at"] = to_jst(product["created_at"].decode())

    cur = db().cursor()
    cur.execute("SELECT * FROM users WHERE id = {}".format(str(user_id)))
    user = cur.fetchone()

    return render_template(
        "mypage.html",
        products=products,
        user=user,
        total_pay=total_pay,
        current_user=current_user(),
    )


@app.route("/products/<int:product_id>")
def get_product(product_id):
    cur = db().cursor()
    cur.execute("SELECT * FROM products WHERE id = {}".format(product_id))
    product = cur.fetchone()

    cur = db().cursor()
    cur.execute("SELECT * FROM comments WHERE product_id = {}".format(product_id))
    comments = cur.fetchall()

    return render_template(
        "product.html",
        product=product,
        comments=comments,
        current_user=current_user(),
        already_bought=already_bought(product_id),
    )


@app.route("/products/buy/<int:product_id>", methods=["POST"])
def post_products_buy(product_id):
    authenticated()
    buy_product(product_id, current_user()["id"])

    return redirect("/users/{}".format(current_user()["id"]))


@app.route("/comments/<int:product_id>", methods=["POST"])
def post_comments(product_id):
    authenticated()
    create_comment(product_id, current_user()["id"], request.form["content"])
    return redirect("/users/{}".format(current_user()["id"]))


@app.route("/initialize")
def get_initialize():
    cur = db().cursor()
    cur.execute("DELETE FROM users WHERE id > 5000")
    cur.execute("DELETE FROM products WHERE id > 10000")
    cur.execute("DELETE FROM comments WHERE id > 200000")
    cur.execute("DELETE FROM histories WHERE id > 500000")
    db().commit()
    return "Finish"


if __name__ == "__main__":
    app.run()
