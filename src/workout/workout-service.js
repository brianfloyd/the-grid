import {ErrorCode, ServerError} from "../server/server-error.js";
import {convertDateToYYYYMMDD, isValidDateString} from "../utils.js";

export class WorkoutService {

    databaseClientFactory;
    workoutDao;
    setService;

    constructor(databaseClientFactory, workoutDao, setService) {
        if (!workoutDao) {
            throw new Error('Workout dao was not provided.');
        }

        if (!databaseClientFactory) {
            throw new Error('Database client factory was not provided.');
        }

        if (!setService) {
            throw new Error('Set service was not provided.');
        }

        this.databaseClientFactory = databaseClientFactory;
        this.workoutDao = workoutDao;
        this.setService = setService;
    }

    async getWorkoutForDate(dateString, client) {
        try {
            if (!dateString) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Date was not provided.');
            }

            if (!isValidDateString(dateString)) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Provided string is not a valid date.');
            }

            const dateKey = convertDateToYYYYMMDD(new Date(Date.parse(dateString)));

            if (!client) {
                client = await this.databaseClientFactory.obtain();
            }

            const data = await this.workoutDao.getWorkoutForDate(client, dateKey);
            if (!data || data.length === 0) {
                throw new ServerError(ErrorCode.NOT_FOUND, 'No workouts found for the specified date.');
            }

            const workout = data[0];
            return await this.convertToApi(workout, client);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error("An unexpected error occurred while getting workout for date.", e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    async convertToApi(workout, client) {
        workout.sets = await this.setService.getSetsForWorkout(workout.id, client);
        return workout;
    }
}
