const controller = require('../controllers/userC');
const express = require('express');

const router = express.Router();

router.post('/:nickname/create', controller.createUser);
router.get('/:nickname/profile', controller.getProfile);
router.post('/:nickname/profile', controller.updateProfile);

module.exports = router;