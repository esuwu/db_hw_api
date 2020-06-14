const usersModel = require('../models/usersM');
const errors = require('../../tools/errors');

class userC {
	static async createUser(req, res) {
		const name = req.params.nickname;
		const profile = req.body;

		try {
			const user = await usersModel.createUser(name, profile);
			res.status(201).json(user);
		}
		catch(error) {
			if (error instanceof errors.AlreadyExistsError) {
				res.status(409).json(error.data);
				return;
			}
			res.status(500).json({ error: error});
		}
	}

	static async getProfile(req, res) {
		const name = req.params.nickname;

		try {
			const user = await usersModel.getProfile(name);
			res.status(200).json(user);
		}
		catch (error) {
			if (error instanceof errors.NotFoundError) {
				res.status(404).json({message: error.message});
				return;
			}
			res.status(500).json({ error: error});
		}
	}

	static async updateProfile(req, res) {
		const name = req.params.nickname;
		const profile = {
			fullname: null,
			about: null,
			email: null,
			...req.body,
		};

		try {
			const user = await usersModel.updateProfile(name, profile);
			res.status(200).json(user);
		}
		catch (error) {
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
}

module.exports = userC;