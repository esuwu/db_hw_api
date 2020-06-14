const pgp = require('pg-promise');

let config;


config = {
	host: 'localhost',
	port: 5432,
	database: 'forum',
	user: 'me',
	password: 'postgres'
};


const db = pgp({})(config);
module.exports = db;