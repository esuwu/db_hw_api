const model = require('../models/forumsM');
const errors = require('../../tools/errors');

class forumC {
	static async createForum(req, res) {
		try {
			const forum = await model.createForum(req.body);
			res.status(201).json(forum);
		} catch (error) {
			console.log(error);
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			if (error instanceof errors.AlreadyExistsError) {
				res.status(409).json(error.data);
				return;
			}
			res.status(500).json({ error: error});
		}
	}

	static async createThread(req, res) {
		const slugforum = req.params.slug;
		const thread = {
			slug: null,
			created: null,
			...req.body,
		};

		try {
			const savedThread = await model.createThread(slugforum, thread);
			res.status(201).json(savedThread);
		} catch (error) {
			console.log(error);
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			if (error instanceof errors.AlreadyExistsError) {
				res.status(409).json(error.data);
				return;
			}
			res.status(500).json({ error: error});
		}
	}

	static async getDetails(req, res) {
		const slugforum = req.params.slug;

		try {
			const details = await model.getDetails(slugforum);
			res.status(200).json(details);
		} catch (error) {
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}

			res.status(500).json({ error: error});
		}
	}

	static async getThreads(req, res) {
		const slugforum = req.params.slug;
		const {limit, since, desc} = req.query;
		const params = {limit, since, desc};

		try {
			const threads = await model.getThreads(slugforum, params);
			res.status(200).json(threads);
		} catch (error) {
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}

			res.status(500).json({ error: error});
		}
	}

	static async getUsers(req, res) {
		const forumSlug = req.params.slug;
		const {limit, since, desc} = req.query;
		const params = {limit, since, desc};

		try {
			const users = await model.getUsers(forumSlug, params);
			res.status(200).json(users);
		} catch (error) {
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}

			res.status(500).json({ error: error});
		}
	}
}

module.exports = forumC;

