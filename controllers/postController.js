const model = require('../models/postsModel');
const errors = require('../tools/errors');

class postController {
	static async getDetails(req, res) {
		const id = req.params.id;
		const related = req.query.related ? req.query.related.split(',') : [];
		let params = {};
		related.forEach(param => params[param] = true);

		try {
			const post = await model.getDetails(id, params);
			res.status(200).json(post);
		} catch (error) {
			if (error instanceof  errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			res.status(500).json({error: error});
		}
	}

	static async updateDetails(req, res) {
		const id = req.params.id;
		const post = {
			message: null,
			...req.body
		};

		try {
			const updatedPost = await model.updateDetails(id, post);
			res.status(200).json(updatedPost);
		} catch (error) {
			if (error instanceof  errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			res.status(500).json({error: error});
		}
	}
}

module.exports = postController;