const express = require('express');
const controller =  require('../controllers/postController');

const router = express.Router();

router.get('/:id/details', controller.getDetails);
router.post('/:id/details', controller.updateDetails);

module.exports = router;