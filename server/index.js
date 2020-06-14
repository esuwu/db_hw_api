const express = require('express');
const cluster = require('express-cluster');
const morgan = require('morgan');
const app = express();
app.use(morgan('dev'));
const port = 5000;

cluster(function(worker) {
	const bodyParser = require('body-parser');
	app.use(bodyParser.json());
	const forumR = require('../router/forumRouter');
	app.use('/api/forum', forumR);
	const postR = require('../router/postRouter');
	app.use('/api/post', postR);
	const serviceR = require('../router/serviceRouter');
	app.use('/api/service', serviceR);
	const userR = require('../router/userRouter');
	app.use('/api/user', userR);
	const threadR = require('../router/threadRouter');
	app.use('/api/thread', threadR);
	return app.listen(port, function () {
		console.log(`Worker ${worker.id} started`);
	});
});