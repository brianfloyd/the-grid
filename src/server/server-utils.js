/**
 * Sends a 400 response.
 *
 * response: Express response
 * error: Server Error
 */
export function _400(response, error) {
    return response.status(400)
        .json(error);
}

/**
 * Sends a 500 response.
 *
 * response: Express response
 * error: Server Error
 */
export function _500(response, error) {
    return response.status(500)
        .json(error);
}

/**
 * Sends a 200 response.
 *
 * response: Express response
 * obj: Object to serialize to json.
 */
export function _200(response, obj) {
    return response.status(200)
        .json(obj);
}
