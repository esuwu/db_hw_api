const usersModel = require('../models/usersModel');
const errors = require('../tools/errors');

class userController {
	static async createUser(req, res) {
		const nickname = req.params.nickname;
		const profile = req.body;

		try {
			const user = await usersModel.createUser(nickname, profile);
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
		const nickname = req.params.nickname;

		try {
			const user = await usersModel.getProfile(nickname);
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
		const nickname = req.params.nickname;
		const profile = {
			fullname: null,
			about: null,
			email: null,
			...req.body,
		};

		try {
			const user = await usersModel.updateProfile(nickname, profile);
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

module.exports = userController;