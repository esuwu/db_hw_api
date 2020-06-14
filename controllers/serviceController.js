const model = require('../models/serviceModel');

class serviceController {
	static async status(req, res) {
		try {
			const status = await model.status();
			res.status(200).json(status);
		} catch(error) {
			res.status(500).json({error: error});
		}
	}

	static async clear(req, res) {
		try {
			await model.clear();
			res.status(200).json({message: 'Successfully cleared.'});
		} catch(error) {
			res.status(500).json({error: error});
		}
	}
}

module.exports = serviceController;