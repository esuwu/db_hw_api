const express = require('express');
const controller = require('../controllers/serviceController');

const router = express.Router();

router.get('/status', controller.status);
router.post('/clear', controller.clear);

module.exports = router;