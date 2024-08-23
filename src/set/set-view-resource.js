import {_200} from "../server/server-utils.js";

export class SetViewResource {

    static URL_PREFIX = '/set/view';

    // TODO: Refactor to use a display service.
    setService;

    constructor(setService) {
        if (!setService) {
            throw new Error('Set service was not provided.');
        }

        this.setService = setService;
    }

    /**
     * Saves a new set.
     *
     * Path: /set/view/save
     */
    async saveNewSet(request, response) {
        const result = await this.setService.insertSet(request.body);
        return _200(response, result);
    }

    /**
     * Updates an existing set.
     *
     * Path: /set/view/:setId/save
     */
    async saveExistingSet(request, response) {

    }

    /**
     * Binds the routes defined by the resource to the given express application instance.
     *
     * app: Express Application
     */
    bind(app) {
        console.log('binding set view');
        app.post(`${SetViewResource.URL_PREFIX}/save`, this.saveNewSet.bind(this));
        // app.post(`${ExerciseViewResource.URL_PREFIX}/:setId/save`, this.getAllExercisesForGroup.bind(this));
    }
}
