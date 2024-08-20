import {ErrorCode, ServerError} from "../server/server-error.js";
import {convertDateToYYYYMMDD, isValidDateString, transactional} from "../utils.js";
import {Workout} from "../model/api/workout.js";
import {WorkoutSet} from "../model/api/workout-set.js";

export class WorkoutService {

    databaseClientFactory;
    workoutDao;

    constructor(databaseClientFactory, workoutDao) {
        if (!workoutDao) {
            throw new Error('Workout dao was not provided.');
        }

        if (!databaseClientFactory) {
            throw new Error('Database client factory was not provided.');
        }

        this.databaseClientFactory = databaseClientFactory;
        this.workoutDao = workoutDao;
    }

    async getWorkoutFromDatabaseDate(dateStr) {
        try {
            const client = await this.databaseClientFactory.obtain();
            // TODO: Test to make sure transactional is working properly.
            await transactional(client, async () => {
                const data = await this.workoutDao.getWorkoutForDate(client, dateStr);
                console.log(data);
                return data;
            })
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error("An unexpected error occurred while getting workout for date.", e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    async getWorkoutForDate(dateString) {
        try {
            if (!dateString) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Date was not provided.');
            }

            if (!isValidDateString(dateString)) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Provided string is not a valid date.');
            }

            const dateKey = convertDateToYYYYMMDD(new Date(Date.parse(dateString)));
            return this.generateFakeWorkout(dateKey);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error("An unexpected error occurred while getting workout for date.", e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    generateFakeWorkout(dateKey) {
        const workout = new Workout();
        workout.id = 1;
        workout.date = dateKey;

        const set1 = new WorkoutSet();
        set1.id = 1;
        set1.exerciseId = 1;
        set1.weight = 30;
        set1.reps = 10;

        const set2 = new WorkoutSet();
        set2.id = 2;
        set2.exerciseId = 1;
        set2.weight = 20;
        set2.reps = 10;

        const set3 = new WorkoutSet();
        set3.id = 3;
        set3.exerciseId = 2;
        set3.weight = 20;
        set3.reps = 10;

        const set4 = new WorkoutSet();
        set4.id = 3;
        set4.exerciseId = 2;
        set4.weight = 10;
        set4.reps = 15;

        workout.sets = [set1, set2, set3, set4];

        return workout;
    }


}
