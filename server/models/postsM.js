const db = require('../../init/config');
const errors = require('../../tools/errors');

class postsM {
	static async getDetails(id, params) {
		try {
			let query = 'SELECT p.username as post_author, p.created as post_created, p.id as post_id, p.isedited as post_isedited, ' +
				'p.message as post_message, p.parent_id as post_parent, p.thread_id as post_thread, ' +
				't.forum_slug as post_forum ' +
				(params.user ? ', u.about as user_about, u.email as user_email, u.fullname as user_fullname, u.nickname as user_nickname ' : '') +
				(params.forum ? ', f.slug as forum_slug, f.posts as forum_posts, f.threads as forum_threads, f.title as forum_title, f.username as forum_user ' : '') +
				(params.thread ? ', t.username as thread_user, t.created as thread_created, t.forum_slug as thread_forum, t.id as thread_id, t.message as thread_message, ' +
				't.slug as thread_slug, t.title as thread_title, t.votes as thread_votes ' : '') +
				'FROM posts p JOIN threads t ON p.thread_id = t.id ' +
				(params.forum ? 'JOIN forums f ON t.forum_slug = f.slug ' : '') +
				(params.user ? 'JOIN users u ON p.username = u.nickname ' : '') +
				'WHERE p.id = ${id}';

			const relatedPost = await db.one(query, {id: id});
			let result = {
				post: {
					author: relatedPost.post_author,
					created: relatedPost.post_created,
					id: relatedPost.post_id,
					message: relatedPost.post_message,
					parent: relatedPost.post_parent,
					thread: relatedPost.post_thread,
					forum: relatedPost.post_forum,
					isEdited: relatedPost.post_isedited
				}
			};

			if (params.user) {
				result = {
					...result,
					author: {
						nickname: relatedPost.user_nickname,
						fullname: relatedPost.user_fullname,
						email: relatedPost.user_email,
						about: relatedPost.user_about
					}
				};
			}

			if (params.forum) {
				result = {
					...result,
					forum: {
						slug: relatedPost.forum_slug,
						posts: relatedPost.forum_posts,
						threads: relatedPost.forum_threads,
						title: relatedPost.forum_title,
						user: relatedPost.forum_user
					}
				}
			}

			if (params.thread) {
				result = {
					...result,
					thread: {
						author: relatedPost.thread_user,
						created: relatedPost.thread_created,
						forum: relatedPost.thread_forum,
						id: relatedPost.thread_id,
						message: relatedPost.thread_message,
						slug: relatedPost.thread_slug,
						title: relatedPost.thread_title,
						votes: relatedPost.thread_votes
					}
				};
			}

			return result;
		} catch (error) {
			throw new errors.NotFoundError();
		}
	}

	static async updateDetails(id, post) {
		try {
			let updatedPost = await db.one('UPDATE posts ' +
				'SET message = COALESCE(${post.message}, posts.message), ' +
				'isEdited = COALESCE(${post.message}, posts.message) <> posts.message ' +
				'FROM threads ' +
				'WHERE threads.id = posts.thread_id AND posts.id = ${id} ' +
				'RETURNING posts.username as author, posts.created, posts.id, posts.isEdited, posts.message, ' +
				'posts.parent_id as parent, posts.thread_id as thread, threads.forum_slug as forum ', {
				id: id, post: post
			});
			updatedPost['isEdited'] = updatedPost.isedited;
			delete updatedPost.isedited;
			return updatedPost;
		} catch (error) {
			throw new errors.NotFoundError();
		}
	}
}

module.exports = postsM;