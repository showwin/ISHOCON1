from flask import Flask, request, abort, session, render_template, redirect
import MySQLdb.cursors
import os
import pathlib
import html
import urllib
import datetime

static_folder = pathlib.Path(__file__).resolve().parent / 'public'
print(static_folder)
app = Flask(__name__, static_folder=str(static_folder), static_url_path='')

app.secret_key = os.environ.get('ISHOCON1_SESSION_SECRET', 'showwin_happy')

_config = {
    'db_host': os.environ.get('ISHOCON1_DB_HOST', 'localhost'),
    'db_port': int(os.environ.get('ISHOCON1_DB_PORT', '3306')),
    'db_username': os.environ.get('ISHOCON1_DB_USER', 'ishocon'),
    'db_password': os.environ.get('ISHOCON1_DB_PASSWORD', 'ishocon'),
    'db_database': os.environ.get('ISHOCON1_DB_NAME', 'ishocon1'),
}


def config(key):
    if key in _config:
        return _config[key]
    else:
        raise "config value of %s undefined" % key


def db():
    if hasattr(request, 'db'):
        return request.db
    else:
        request.db = MySQLdb.connect(**{
            'host': config('db_host'),
            'port': config('db_port'),
            'user': config('db_username'),
            'passwd': config('db_password'),
            'db': config('db_database'),
            'charset': 'utf8mb4',
            'cursorclass': MySQLdb.cursors.DictCursor,
            'autocommit': True,
        })
        cur = request.db.cursor()
        cur.execute("SET SESSION sql_mode='TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY'")
        cur.execute('SET NAMES utf8mb4')
        return request.db


def authenticate(email, password):
    cur = db().cursor()
    cur.execute("SELECT * FROM users WHERE email = %s", (email, ))
    user = cur.fetchone()
    if user is None or user.get('password', None) != password:
        abort(401)
    else:
        session['user_id'] = user['id']


def update_last_login(user_id):
    cur = db().cursor()
    cur.execute('UPDATE users SET last_login = %s WHERE id = %s', (datetime.datetime.now(), user_id,))


def current_user():
    cur = db().cursor()
    cur.execute('SELECT * FROM users WHERE id = %s', str(session['user_id']))
    return cur.fetchone()


@app.errorhandler(401)
def authentication_error(error):
    return render_template('login.html', message='ログインに失敗しました'), 401


@app.teardown_request
def close_db(exception=None):
    if hasattr(request, 'db'):
        request.db.close()


@app.route('/login')
def get_login():
    session.pop('user_id', None)
    return render_template('login.html', messae='ECサイトで爆買いしよう！！！！')


@app.route('/login', methods=['POST'])
def post_login():
    authenticate(request.form['email'], request.form['password'])
    update_last_login(current_user()['id'])
    return redirect('/')


if __name__ == "__main__":
    app.run()
