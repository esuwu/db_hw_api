const express = require('express');
const controller = require('../controllers/serviceC');

const router = express.Router();

router.get('/status', controller.status);
router.post('/clear', controller.clear);

module.exports = router;