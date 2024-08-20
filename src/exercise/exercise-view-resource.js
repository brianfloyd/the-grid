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
     * Gets a list of all the exercise groups.
     *
     * Path: /exercise/view/groups/all
     */
    async getAllExerciseGroups(request, response) {
        const exerciseGroupsView = await this.exerciseService.getAllExerciseGroups();
        return _200(response, exerciseGroupsView);
    }

    /**
     * Gets all the exercises for the given group.
     *
     * Path: /exercise/view/group/:name
     */
    async getAllExercisesForGroup(request, response) {
        const groupName = request.params.name
        const exercisesForGroupView = await this.exerciseService.getAllExercisesForGroup(groupName);
        return _200(response, exercisesForGroupView);
    }

    /**
     * Binds the routes defined by the resource to the given express application instance.
     *
     * app: Express Application
     */
    bind(app) {
        app.get(`${ExerciseViewResource.URL_PREFIX}/groups/all`, this.getAllExerciseGroups.bind(this));
        app.get(`${ExerciseViewResource.URL_PREFIX}/group/:name`, this.getAllExercisesForGroup.bind(this));
    }
}
