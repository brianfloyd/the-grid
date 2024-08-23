import {ErrorCode, ServerError} from "../server/server-error.js";
import {convertDateToYYYYMMDD, isValidDateString} from "../utils.js";
import {Workout} from "../model/api/workout.js";

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
        let dbClient = client;
        try {
            if (!dateString) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Date was not provided.');
            }

            if (!isValidDateString(dateString)) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Provided string is not a valid date.');
            }

            const dateKey = convertDateToYYYYMMDD(new Date(Date.parse(dateString)));

            if (!client) {
                dbClient = await this.databaseClientFactory.obtain();
            }

            const data = await this.workoutDao.getWorkoutForDate(dbClient, dateKey);
            if (!data || data.length === 0) {
                throw new ServerError(ErrorCode.NOT_FOUND, 'No workouts found for the specified date.');
            }

            const workoutData = data[0];
            const result = await this.convertToApi(workoutData, dbClient);
            return result;
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error("An unexpected error occurred while getting workout for date.", e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        } finally {
            if (!client && dbClient) {
                dbClient.release();
            }
        }
    }

    async convertToApi(workoutData, client) {
        const workout = new Workout();
        workout.id = workoutData.id;
        workout.date = convertDateToYYYYMMDD(workoutData.date);
        workout.sets = await this.setService.getSetsForWorkout(workout.id, client);
        return workout;
    }
}
