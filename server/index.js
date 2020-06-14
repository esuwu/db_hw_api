const express = require('express');
const cluster = require('express-cluster');
const morgan = require('morgan');
const app = express();
app.use(morgan('dev'));
const port = 5000;

cluster(function(worker) {
	const bodyParser = require('body-parser');
	app.use(bodyParser.json());
	const forumR = require('./routers/forumR');
	app.use('/api/forum', forumR);
	const postR = require('./routers/postR');
	app.use('/api/post', postR);
	const serviceR = require('./routers/serviceR');
	app.use('/api/service', serviceR);
	const userR = require('./routers/userR');
	app.use('/api/user', userR);
	const threadR = require('./routers/threadR');
	app.use('/api/thread', threadR);
	return app.listen(port, function () {
		console.log(`Worker ${worker.id} started`);
	});
});