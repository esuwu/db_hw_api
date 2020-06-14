const db = require('../init/config');
const errors = require('../tools/errors');

class threadsM {
	static async createPost(slugOrId, posts) {
		const numbersPattern = '^[0-9]+$';
		let threadQuery = 'SELECT * FROM threads ';
		if (slugOrId.match(numbersPattern)) {
			threadQuery += 'WHERE id = ${slugOrId}';
		} else {
			threadQuery += 'WHERE slug = ${slugOrId}';
		}

		let thread = {};
		try {
			thread = await db.one(threadQuery, {slugOrId: slugOrId});
		} catch (error) {
			throw new errors.NotFoundError();
		}

		if (posts.length === 0) {
			return [];
		}

		let query = 'INSERT INTO posts(username, message, thread_id, parent_id) VALUES ';
		for (let i = 0; i < posts.length; i++) {
			const post = posts[i];

			let user = {};
			try {
				user = await db.one('SELECT nickname FROM users WHERE nickname = ${post.author}', {post: post});
			} catch (error) {
				throw new errors.NotFoundError();
			}

			query += ` ('${user.nickname}', '${post.message}', ${thread.id}, `;
			if (!post.parent) {
				query += `NULL )`;
			} else {
				try {
					const parentPost = await db.one('SELECT * FROM posts WHERE id = ${post.parent} AND thread_id = ${thread.id}',
						{post: post, thread: thread});
					query += `${parentPost.id})`;
				} catch (error) {
					throw new errors.AlreadyExistsError();
				}
			}
			if (i < posts.length - 1) {
				query += ', ';
			}
		}
		query += 'RETURNING username as author, created, id, isEdited, message, parent_id as parent, thread_id as thread ';

		let savedPosts = await db.manyOrNone(query);
		savedPosts.forEach(savedPost => {
			savedPost.forum = thread.forum_slug;
		});

		try {
			await db.one('UPDATE forums SET posts = posts + ${posts.length} WHERE slug = ${thread.forum_slug} RETURNING *', {
				posts: posts, thread: thread
			});
		} catch (error) {
			console.log(error);
		}

		try {
			await this.addForumUser(posts, thread);
		} catch (error) {
			console.log(error);
		}

		return savedPosts;
	}

	static async addForumUser(posts, thread) {
		let query = 'INSERT INTO forum_user(forum_slug, user_id) VALUES ';

		posts.forEach((post, index) => {
			if (index !== 0) {
				query += ' , ';
			}
			query += ` ('${thread.forum_slug}', (SELECT id FROM users WHERE nickname = '${post.author}')) `;
		});

		query += ' ON CONFLICT DO NOTHING ';
		console.log(query);
		await db.manyOrNone(query);
	}

	static async getDetails(slugOrId) {
		const numbersPattern = '^[0-9]+$';
		let threadQuery = 'SELECT username as author, created, forum_slug as forum, id, message, title, slug, votes FROM threads ';
		if (slugOrId.match(numbersPattern)) {
			threadQuery += 'WHERE id = ${slugOrId}';
		} else {
			threadQuery += 'WHERE slug = ${slugOrId}';
		}

		try {
			return await db.one(threadQuery, {slugOrId: slugOrId});
		} catch (error) {
			throw new errors.NotFoundError();
		}
	}

	static async updateDetails(slugOrId, thread) {
		let query = 'UPDATE threads SET message = COALESCE(${thread.message}, message), title = COALESCE(${thread.title}, title)';
		const numbersPattern = '^[0-9]+$';
		if (slugOrId.match(numbersPattern)) {
			query += 'WHERE id = ${slugOrId} ';
		} else {
			query += 'WHERE slug = ${slugOrId} ';
		}
		query += 'RETURNING username as author, created, forum_slug as forum, id, message, title, slug';

		try {
			return await db.one(query, {thread: thread, slugOrId: slugOrId});
		} catch (error) {
			throw new errors.NotFoundError();
		}
	}

	static async voteBySlug(threadSlug, vote) {
		try {
			 await db.any('INSERT INTO votes(username, thread_id, voice) VALUES (' +
				'${vote.nickname}, ' +
				'(SELECT id FROM threads WHERE slug = ${threadSlug}), ' +
				'${vote.voice}) ' +
				'ON CONFLICT ON CONSTRAINT unique_vote ' +
				'DO UPDATE SET voice = ${vote.voice} ' +
				'WHERE votes.thread_id = (SELECT id FROM threads WHERE slug = ${threadSlug}) AND votes.username = ${vote.nickname} ' +
				'RETURNING votes.id',
				{threadSlug: threadSlug, vote: vote});
			 return await this.getDetails(threadSlug);
		} catch (error) {
			throw new errors.NotFoundError();
		}
	}

	static async voteById(threadId, vote) {
		try {
			await db.any('INSERT INTO votes(username, thread_id, voice) VALUES (' +
				'${vote.nickname}, ' +
				'${threadId}, ' +
				'${vote.voice}) ' +
				'ON CONFLICT ON CONSTRAINT unique_vote ' +
				'DO UPDATE SET voice = ${vote.voice} ' +
				'WHERE votes.thread_id = ${threadId} AND votes.username = ${vote.nickname} ' +
				'RETURNING votes.id',
				{threadId: threadId, vote: vote});
			return await this.getDetails(threadId);
		} catch (error) {
			throw new errors.NotFoundError();
		}
	}

	static async getPostsBySlug(slug, params) {
		let thread;
		try {
			thread = await db.one('SELECT * FROM threads WHERE slug = ${slug}', {slug: slug});
		} catch (error) {
			throw new errors.NotFoundError();
		}

		return await this.getPosts(thread, params);
	}

	static async getPostsById(id, params) {
		let thread;
		try {
			thread = await db.one('SELECT * FROM threads WHERE id = ${id}', {id: id});
		} catch (error) {
			throw new errors.NotFoundError();
		}

		return await this.getPosts(thread, params);
	}

	static async getPosts(thread, params) {
		let posts = [];
		try {
			switch (params.sort) {
				case 'tree':
					posts = await this.getPostsTree(thread, params);
					break;
				case 'parent_tree':
					posts = await this.getPostsParentTree(thread, params);
					break;
				default:
					posts = await this.getPostsFlat(thread, params);
			}
		}
		catch (error) {
			throw new Error(error);
		}
		posts.forEach(post => {
			post.isEdited = post.isedited;
			post.forum = thread.forum_slug;
			delete post.isedited;
		});
		return posts;
	}

	static async getPostsFlat(thread, params) {
		let query = 'SELECT username as author, created, id, isedited, message, parent_id as parent, thread_id as thread ' +
			'FROM posts ' +
			'WHERE thread_id = ${thread.id} ';

		if (params.desc && params.desc === 'true') {
			if (params.since) {
				query += 'AND id < ${params.since} ';
			}

			query += 'ORDER BY created DESC, id DESC ';
		} else {
			if (params.since) {
				query += 'AND id > ${params.since} ';
			}

			query += 'ORDER BY created, id ';
		}

		if (params.limit) {
			query += 'LIMIT ${params.limit}';
		}

		return await db.manyOrNone(query, {thread: thread, params: params});
	}

	static async getPostsTree(thread, params) {
		let query = 'SELECT posts.username as author, posts.created, posts.id, posts.isedited, posts.message, ' +
			'posts.parent_id as parent, posts.path, posts.thread_id as thread ' +
			'FROM posts ';

		if (params.since) {
			if (params.desc && params.desc === 'true') {
				query += 'JOIN posts p ON p.id = ${since} ' +
					'WHERE posts.path <  p.path ';

			} else {
				query += 'JOIN posts p ON p.id = ${since} ' +
					'WHERE posts.path  > p.path ';
			}
			query += 'AND posts.thread_id = ${thread.id} ORDER BY posts.path ';
		} else {
			query += 'WHERE thread_id = ${thread.id} ORDER BY path ';
		}

		if (params.desc && params.desc === 'true') {
			query += ' DESC ';
		}

		if (params.limit) {
			query += 'LIMIT ${params.limit}';
		}

		console.log(query);

		return await db.manyOrNone(query, {thread: thread, params: params, since: +params.since});
	}

	static async getPostsParentTree(thread, params) {
		const limit = +params.limit || 10000;

		let query = 'WITH ranks AS ( ' +
			'    SELECT posts.username as author, posts.created, posts.id, posts.isedited, posts.message, posts.parent_id as parent, ' +
			'           posts.thread_id as thread, posts.path, ' +
			'           dense_rank() over (ORDER BY path[1] ' + (params.desc === 'true' ? ' DESC ' : '') + ' ) as rank ' +
			'    FROM posts' +
			'    WHERE posts.thread_id = ${thread.id} ' +
			') ' +
			'SELECT ranks.id, ranks.message, ranks.created, ranks.parent, ranks.thread, ranks.author, ranks.rank, ranks.path ' +
			'FROM ranks ';

		if (params.since) {
			query += 'JOIN ranks r ON r.id = ${since} ' +
				'WHERE ranks.rank <= ${limit} + r.rank ' +
				'AND (ranks.rank > r.rank OR ranks.rank = r.rank AND ranks.path > r.path) ' +
				'ORDER BY ranks.path[1]' + (params.desc === 'true' ? ' DESC ' : '') + ', array_remove(ranks.path, ranks.path[1])';
		} else {
			query += 'WHERE ranks.rank <= ${limit} ' +
				'ORDER BY ranks.path[1] '  + (params.desc === 'true' ? ' DESC ' : '') + ', array_remove(ranks.path, ranks.path[1])';
		}

		console.log(query);

		return await db.manyOrNone(query, {thread: thread, desc: params.desc, limit: limit, since: +params.since});
	}
}

module.exports = threadsM;