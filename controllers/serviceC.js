const model = require('../models/serviceM');

class serviceC {
	static async status(req, res) {
		try {
			const stat = await model.status();
			res.status(200).json(stat);
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

module.exports = serviceC;