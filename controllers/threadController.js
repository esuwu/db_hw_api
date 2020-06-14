const model = require('../models/threadsModel');
const errors = require('../tools/errors');

class threadController {
	static async createPost(req, res) {
		const slugOrId = req.params.slug_or_id;
		const posts = req.body.map(post => {
			return {
				created: null,
				parent: null,
				...post
			};
		});

		try {
			const savedPosts = await model.createPost(slugOrId, posts);
			res.status(201).json(savedPosts);
		} catch (error) {
			console.log(error);
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			if (error instanceof errors.AlreadyExistsError) {
				res.status(409).json({message: error.message});
				return;
			}
			res.status(500).json({ error: error});
		}
	}

	static async getDetails(req, res) {
		const slugOrId = req.params.slug_or_id;

		try {
			const thread = await model.getDetails(slugOrId);
			res.status(200).json(thread);
		} catch (error) {
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			res.status(500).json({ error: error});
		}
	}

	static async updateDetails(req, res) {
		const slugOrId = req.params.slug_or_id;
		const thread = {
			message: null,
			title: null,
			...req.body
		};

		try {
			const updatedThread = await model.updateDetails(slugOrId, thread);
			res.status(200).json(updatedThread);
		} catch (error) {
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			res.status(500).json({ error: error});
		}
	}

	static async vote(req, res) {
		const slugOrId = req.params.slug_or_id;
		const vote = req.body;

		const numbersPattern = '^[0-9]+$';
		try {
			let thread = {};
			if (slugOrId.match(numbersPattern)) {
				thread = await model.voteById(slugOrId, vote);
			} else {
				thread = await model.voteBySlug(slugOrId, vote);
			}
			res.status(200).json(thread);
		} catch (error) {
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			res.status(500).json(error);
		}
	}

	static async getPosts(req, res) {
		const slugOrId = req.params.slug_or_id;
		const {limit, since, sort, desc } = req.query;
		const params = {limit, since, sort, desc};

		const numbersPattern = '^[0-9]+$';
		try {
			let posts;
			if (slugOrId.match(numbersPattern)) {
				posts = await model.getPostsById(slugOrId, params);
			} else {
				posts = await model.getPostsBySlug(slugOrId, params);
			}
			res.status(200).json(posts);
		} catch(error) {
			console.log(error);
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			res.status(500).json(error);
		}
	}
}

module.exports = threadController;