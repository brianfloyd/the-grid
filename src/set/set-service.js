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

    async getSetForId(id, client) {
        try {
            if (!client) {
                client = await this.databaseClientFactory.obtain();
            }
            const result = await this.setDao.getSetForId(client, id);
            if (!result || result.length === 0) {
                throw new ServerError(ErrorCode.NOT_FOUND, `Could not find set for id ${id}.`);
            }
            return result;
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting set by id.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    async insertSet(set, client) {
        try {
            this.validateSet(set);
            if (!client) {
                client = await this.databaseClientFactory.obtain();
            }
            const setId = await this.setDao.insertSet(client, set);
            return await this.getSetForId(setId, client);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while inserting sets.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    async updateSet(setId, set, client) {
        try {
            if (!setId || setId < 1) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Set id must be non null and greater than zero.');
            }
            this.validateSet(set);
            set.id = setId;

            if (!client) {
                client = await this.databaseClientFactory.obtain();
            }
            await this.setDao.updateSet(client, set);
            return await this.getSetForId(setId, client);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while updating set data.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }

    }

    validateSet(set) {
        if (!set.workoutId || set.workoutId < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided workout id was not valid.');
        }

        if (!set.exerciseId || set.exerciseId < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided exercise id was not valid.');
        }

        if (!set.weight || set.weight < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided weight was not valid.');
        }

        if (!set.reps || set.reps < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided reps was not valid.');
        }

        if (!set.count || set.count < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided count was not valid.');
        }
    }

}
