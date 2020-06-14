const pgp = require('pg-promise');
const db = require('../init/config');
const errors = require('../tools/errors');

class forumsModel {
	static async createForum(forum) {
		try {
			const user = await db.one(`SELECT * FROM users WHERE nickname = '${forum.user}'`);
			return await db.one('INSERT INTO forums(slug, title, username) ' +
				'VALUES(${forum.slug}, ${forum.title}, ${user.nickname}) RETURNING slug, posts, threads, title, username as user', {forum: forum, user: user});
		} catch (error) {
			if (error instanceof pgp.errors.QueryResultError) {
				throw new errors.NotFoundError();
			}

			const forums = await db.one('SELECT slug, posts, threads, title, username as user FROM forums WHERE slug = ${forum.slug}',
				{forum: forum});
			throw new errors.AlreadyExistsError(forums);
		}
	}

	static async createThread(forumSlug, thread) {
		try {
			const forum = await db.one('SELECT * FROM forums WHERE slug = ${forumSlug}', {forumSlug: forumSlug});
			const user = await db.one('SELECT * FROM users WHERE nickname = ${thread.author}', {thread: thread});

			return await db.one('INSERT INTO threads(username, forum_slug, slug, created, message, title) ' +
				'VALUES(${user.nickname}, ${forum.slug}, ${thread.slug}, ${thread.created}, ${thread.message}, ${thread.title}) ' +
				'RETURNING username as author, created, forum_slug as forum, id, message, title, slug, votes',
				{forum: forum, thread: thread, user: user});

		} catch (error) {
			if (error instanceof pgp.errors.QueryResultError) {
				throw new errors.NotFoundError();
			}

			const threads = await db.one('SELECT username as author, created, forum_slug as forum, id, message, title, slug, votes ' +
				'FROM threads WHERE title = ${thread.title} or slug = ${thread.slug}', {thread: thread});
			throw new errors.AlreadyExistsError(threads);
		}
	}

	static async getDetails(forumSlug) {
		try {
			return await db.one('SELECT slug, posts, threads, title, username as user FROM forums WHERE slug = ${forumSlug}', {forumSlug: forumSlug});
		} catch (error) {
			throw new errors.NotFoundError();
		}
	}

	static async getThreads(forumSlug, params) {
		try {
			const forum = await db.one('SELECT * FROM forums WHERE slug = ${forumSlug}', {forumSlug: forumSlug});
			let query = 'SELECT username as author, created, forum_slug as forum, id, message, title, slug, votes FROM threads WHERE forum_slug = ${forum.slug} ';

			if (params.desc && params.desc === 'true') {
				if (params.since) {
					query += 'AND created <= ${params.since} ';
				}
				query += 'ORDER BY created DESC ';
			} else {
				if (params.since) {
					query += 'AND created >= ${params.since} ';
				}
				query += 'ORDER BY created ';
			}

			if (params.limit) {
				query += 'LIMIT ${params.limit} ';
			}

			return await db.manyOrNone(query, {forum: forum, params: params});

		} catch(error) {
			throw new errors.NotFoundError();
		}
	}

	static async getUsers(forumSlug, params) {
		try {
			const forum = await db.one('SELECT * FROM forums WHERE slug = ${forumSlug}', {forumSlug: forumSlug});
			let query = 'SELECT u.email, u.nickname, u.fullname, u.about, u.id ' +
				'FROM forum_user fu ' +
				'JOIN users u ON u.id = fu.user_id ' +
				'WHERE fu.forum_slug = ${forum.slug} ';

			if (params.desc && params.desc === 'true') {
				if (params.since) {
					query += 'AND u.nickname < ${params.since} ';
				}
				query += 'ORDER BY u.nickname DESC ';
			} else {
				if (params.since) {
					query += 'AND u.nickname > ${params.since} ';
				}
				query += 'ORDER BY u.nickname ';
			}

			if (params.limit) {
				query += 'LIMIT ${params.limit} ';
			}
			console.log(query);
			return await db.manyOrNone(query, {forum: forum, params: params});
		} catch (error) {
			console.log(error);
			throw new errors.NotFoundError();
		}
	}
}

module.exports = forumsModel;