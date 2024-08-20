import {_200} from "../server/server-utils.js";

export class ExerciseViewResource {

    static URL_PREFIX = '/exercise/view';

    // TODO: Refactor to use a display service.
    exerciseService;

    constructor(exerciseService) {
        if (!exerciseService) {
            throw new Error('Exercise service was not provided.');
        }

        this.exerciseService = exerciseService;
    }

    /**
     * Gets a workout for a given date passed in the url parameter.
     *
     * Path: /exercise/view/groups/all
     */
    async getAllExerciseGroups(request, response) {
        const workoutView = await this.exerciseService.getAllExerciseGroups();
        return _200(response, workoutView);
    }

    /**
     * Binds the routes defined by the resource to the given express application instance.
     *
     * app: Express Application
     */
    bind(app) {
        app.get(`${ExerciseViewResource.URL_PREFIX}/groups/all`, this.getAllExerciseGroups.bind(this));
    }
}
