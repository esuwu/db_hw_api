const pgp = require('pg-promise');
const db = require('../init/config');
const errors = require('../tools/errors');

class usersM {
	static async createUser(nickname, profile) {
		try {
			return await db.oneOrNone('INSERT INTO users(nickname, email, fullname, about) ' +
				'VALUES(${nickname}, ${profile.email}, ${profile.fullname}, ${profile.about}) RETURNING *', {
				profile: profile, nickname: nickname
			});
		}
		catch (error) {
			const users = await db.manyOrNone(`select * from users where nickname = '${nickname}' or email = '${profile.email}'`);
			throw new errors.AlreadyExistsError(users);
		}
	}

	static async getProfile(nickname) {
		try {
			return await db.one(`SELECT * FROM users WHERE nickname = '${nickname}'`);
		}
		catch(error) {
			throw new errors.NotFoundError();
		}
	}

	static async updateProfile(nickname, profile) {
		try {
			return await db.one('UPDATE users SET about = COALESCE(${profile.about}, about), email = COALESCE(${profile.email}, email), fullname = COALESCE(${profile.fullname}, fullname) ' +
				'WHERE nickname = ${nickname} RETURNING *',{
				profile: profile, nickname: nickname
			});
		}
		catch(error) {
			if (error instanceof pgp.errors.QueryResultError) {
				throw new errors.NotFoundError();
			}
			throw new errors.AlreadyExistsError();
		}
	}
}

module.exports = usersM;