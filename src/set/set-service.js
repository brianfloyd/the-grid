import {ErrorCode, ServerError} from "../server/server-error.js";

export class SetService {

    databaseClientFactory;
    setDao;

    constructor(databaseClientFactory, setDao) {
        if (!databaseClientFactory) {
            throw new Error('Database client factory was not provided.');
        }

        if (!setDao) {
            throw new Error('Set dao was not provided.');
        }


        this.databaseClientFactory = databaseClientFactory;
        this.setDao = setDao;
    }

    async getSetsForWorkout(workoutId, client) {
        try {
            if (!client) {
                client = await this.databaseClientFactory.obtain();
            }
            return await this.setDao.getSetsForWorkout(client, workoutId);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting sets for workout.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

}
