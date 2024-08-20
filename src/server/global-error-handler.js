import {ErrorCode, ServerError} from "./server-error.js";

export class GlobalErrorHandler {
    static handle(error, request, response, next) {
        if (response.headersSent) {
            return next(error);
        }

        if (error instanceof ServerError) {
            response.status(error.statusCode)
                .json(error);
        } else {
            console.error('An unhandled error occurred.', error);
            response.status(500)
                .json(new ServerError(ErrorCode.GENERIC_ERROR, error.message));
        }
    }
}
