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
     * Binds the routes defined by the resource to the given express application instance.
     *
     * app: Express Application
     */
    bind(app) {
        app.get(`${WorkoutViewResource.URL_PREFIX}/date/:date`, this.getWorkoutViewForDate.bind(this));
    }
}
