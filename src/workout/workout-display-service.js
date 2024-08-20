import {ErrorCode, ServerError} from "../server/server-error.js";
import {convertDateToYYYYMMDD, isValidDateString} from "../utils.js";

/**
 * Orchestration layer that converts workout information into a display ready format.
 */
export class WorkoutDisplayService {

    constructor() {}

    /**
     * Gets a workout for a given date.
     */
    getWorkoutViewForDate(dateString) {
        try {
            if (!dateString) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Date was not provided.');
            }

            if (!isValidDateString(dateString)) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Provided string is not a valid date.');
            }

            const dateKey = convertDateToYYYYMMDD(new Date(Date.parse(dateString)));
            return this.generateFakeResponse(dateKey);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting workout by date.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    generateFakeResponse(dateKey) {
        return {
            id: 1,
            date: dateKey,
            sets: [
                {
                    id: 1,
                    exerciseId: 1,
                    group: 'BICEP',
                    name: 'Bicep Curls',
                    weight: 30,
                    reps: 10
                },
                {
                    id: 2,
                    exerciseId: 1,
                    group: 'BICEP',
                    name: 'Bicep Curls',
                    weight: 20,
                    reps: 10
                },
                {
                    id: 3,
                    exerciseId: 2,
                    group: 'TRICEP',
                    name: 'Tricep Extensions',
                    weight: 20,
                    reps: 10
                },
                {
                    id: 4,
                    exerciseId: 2,
                    group: 'TRICEP',
                    name: 'Tricep Extensions',
                    weight: 10,
                    reps: 15
                }
            ]
        }
    }
}
