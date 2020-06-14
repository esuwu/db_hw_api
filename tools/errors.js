class BaseError extends Error {
	constructor (message) {
		super();
		Error.captureStackTrace( this, this.constructor );
		this.name = this.constructor.name;
		if (message) {
			this.message = message;
		}
	}
}

class AlreadyExistsError extends BaseError {
	constructor(data = {}) {
		super('Error: Already exists');
		this.data = data;
	}
}

class NotFoundError extends BaseError {
	constructor(message = 'Error: Not Found') {
		super(message);
	}
}

exports.AlreadyExistsError = AlreadyExistsError;
exports.NotFoundError = NotFoundError;