import {ErrorCode, ServerError} from "../server/server-error.js";
import {transactional} from "../utils.js";

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
        let dbClient = client;
        try {
            if (!client) {
                dbClient = await this.databaseClientFactory.obtain();
            }
            return await this.setDao.getSetsForWorkout(dbClient, workoutId);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting sets for workout.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        } finally {
            if (!client && dbClient) {
                dbClient.release();
            }
        }
    }

    async getSetForId(id, client) {
        let dbClient = client;
        try {
            if (!client) {
                dbClient = await this.databaseClientFactory.obtain();
            }
            const result = await this.setDao.getSetForId(dbClient, id);
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
        } finally {
            if (!client && dbClient) {
                dbClient.release();
            }
        }
    }

    async createSet(set, client) {
        let dbClient = client;
        try {
            if (!client) {
                dbClient = await this.databaseClientFactory.obtain();
            }
           this.validateSet(set);

            return await transactional(dbClient, async () => {
                const setId = await this.setDao.insertSet(dbClient, set);
                return await this.getSetForId(setId, dbClient);
            });
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while inserting sets.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        } finally {
            if (!client && dbClient) {
                dbClient.release();
            }
        }
    }

    async updateSet(setId, set, client) {
        let dbClient = client;
        try {
            if (!setId || setId < 1) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'Set id must be non null and greater than zero.');
            }
            this.validateSet(set);
            set.id = setId;

            if (!client) {
                dbClient = await this.databaseClientFactory.obtain();
            }

            return await transactional(dbClient, async () => {
                await this.setDao.updateSet(dbClient, set);
                return await this.getSetForId(setId, dbClient);
            });
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while updating set data.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        } finally {
            if (!client && dbClient) {
                dbClient.release();
            }
        }
    }

    validateSet(set) {
        if (!set.workoutId || set.workoutId < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided workout id was not valid.');
        }

        if (!set.exerciseId || set.exerciseId < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided exercise id was not valid.');
        }

        if (set.weight === undefined || set.weight < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided weight was not valid.');
        }

        if (set.reps === undefined || set.reps < 0) {
            throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided reps was not valid.');
        }

        if (set.count === undefined || set.count < 0) {
          set.count = 0
        }
    }

}
