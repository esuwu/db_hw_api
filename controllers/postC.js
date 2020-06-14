const model = require('../models/postsM');
const errors = require('../tools/errors');

class postC {
	static async getDetails(req, res) {
		const id = req.params.id;
		const connected = req.query.related ? req.query.related.split(',') : [];
		let params = {};
		connected.forEach(param => params[param] = true);

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
		const postDetail = {
			message: null,
			...req.body
		};

		try {
			const updatedPost = await model.updateDetails(id, postDetail);
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

module.exports = postC;