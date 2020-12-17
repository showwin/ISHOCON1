const util = require('util');

const express = require('express');
const session = require('express-session');
const bodyParser = require('body-parser')

const app = express();

app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');
app.use(express.static('public'));
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.use(session({ secret: 'showwin_happy', name: 'ishocon1_node_session'}))

const mysql = require('mysql');
const pool = mysql.createPool({
    connectionLimit: 20,
    host: process.env.ISHOCON1_DB_HOST || 'localhost', // not supported in other implementations
    port: process.env.ISHOCON1_DB_PORT || 3306, // not supported in other implementations
    user: process.env.ISHOCON1_DB_USER || 'ishocon',
    password: process.env.ISHOCON1_DB_PASSWORD || 'ishocon',
    database: process.env.ISHOCON1_DB_NAME || 'ishocon1'
});
const query = util.promisify(pool.query).bind(pool);

app.get('/initialize', async (_, res) => {
    await query('DELETE FROM users WHERE id > 5000');
    await query('DELETE FROM products WHERE id > 10000');
    await query('DELETE FROM comments WHERE id > 200000');
    await query('DELETE FROM histories WHERE id > 500000');
    res.send('Finish')
})

app.get('/', (_, res) => {
    res.send('hello')
})

const server = app.listen(8080, function () {
    const host = server.address().address;
    const port = server.address().port;

    console.log('Example app listening at http://%s:%s', host, port);
});