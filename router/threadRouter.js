const controller = require('../controllers/threadController');
const express = require('express');

const router = express.Router();

router.post('/:slug_or_id/create', controller.createPost);
router.get('/:slug_or_id/details', controller.getDetails);
router.post('/:slug_or_id/details', controller.updateDetails);
router.post('/:slug_or_id/vote', controller.vote);
router.get('/:slug_or_id/posts', controller.getPosts);

module.exports = router;
