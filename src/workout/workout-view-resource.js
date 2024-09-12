import {_200} from "../server/server-utils.js";

/**
 * Provides REST endpoints to manage workouts from the frontend.
 */
export class WorkoutViewResource {

    static URL_PREFIX = '/workout/view';

    workoutDisplayService;

    constructor(workoutDisplayService) {
        if (!workoutDisplayService) {
            throw new Error('Workout display service was not provided.');
        }

        this.workoutDisplayService = workoutDisplayService;
    }

    /**
     * Gets a workout for a given date passed in the url parameter.
     *
     * Path: /workout/view/date/:date
     */
    async getWorkoutViewForDate(request, response) {
        const workoutView = await this.workoutDisplayService.getWorkoutViewForDate(request.params.date);
        return _200(response, workoutView);
    }

    /**
     * Saves a new workout and all passed sets for the date specified in the url parameter.
     *
     * Path: POST /workout/view/save-new/:date
     */
    async saveNewWorkout(request, response) {
        const date = request.params.date;
        const exerciseIds = request.body;
        await this.workoutDisplayService.saveNewWorkout(date, exerciseIds);
        return _200(response, {});
    }

    /**
     * Adds a new default set with the specified exercise id to the workout.
     *
     * Path: POST /workout/view/add-set
     */
    async addSet(request, response) {
        const workoutId = request.body.workoutId;
        const exerciseId = request.body.exerciseId;
        await this.workoutDisplayService.addNewSetToWorkout(workoutId, exerciseId);
        return _200(response, {});
    }

    /**
     * Binds the routes defined by the resource to the given express application instance.
     *
     * app: Express Application
     */
    bind(app) {
        app.get(`${WorkoutViewResource.URL_PREFIX}/date/:date`, this.getWorkoutViewForDate.bind(this));
        app.post(`${WorkoutViewResource.URL_PREFIX}/save-new/:date`, this.saveNewWorkout.bind(this));
        app.post(`${WorkoutViewResource.URL_PREFIX}/add-set`, this.addSet.bind(this));
    }
}
