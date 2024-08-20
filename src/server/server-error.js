
export const ErrorCode = Object.freeze({
    INVALID_REQUEST: {
        statusCode: 400,
        name: 'INVALID_REQUEST'
    },
    NOT_FOUND: {
        statusCode: 404,
        name: 'NOT_FOUND'
    },
    GENERIC_ERROR: {
        responseCode: 500,
        name: 'GENERIC_ERROR'
    }
});

export class ServerError extends Error {
    statusCode; // number
    errorCode;  // string
    reason;     // string

    /**
     *  errorCode: ErrorCode;
     *  reason: string;
     */
    constructor(errorCode, reason) {
        super(reason || 'Generic error response.');
        if (errorCode && typeof errorCode === 'object' && errorCode.statusCode && errorCode.name) {
            this.statusCode = errorCode.statusCode;
            this.errorCode = errorCode.name;
        } else {
            this.statusCode = 500;
            this.errorCode = 'GENERIC_ERROR';
        }
        if (typeof reason === 'string' && reason.length > 0) {
            this.reason = reason;
        } else {
            this.reason = 'Generic error response.'
        }
    }
}
