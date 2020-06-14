const db = require('../init/config');

class serviceM {
	static async status() {
		const status = await db.one('SELECT ' +
			'(SELECT COUNT(*) FROM forums) as forum, ' +
			'(SELECT COUNT(*) FROM threads) as thread, ' +
			'(SELECT COUNT(*) FROM posts) as post, ' +
			'(SELECT COUNT(*) FROM users) as user');

		return {
			forum: +status.forum,
			thread: +status.thread,
			post: +status.post,
			user: +status.user
		};
	}

	static async clear() {
		return await db.any('TRUNCATE votes, posts, threads, forums, users, forum_user');
	}
}

module.exports = serviceM;