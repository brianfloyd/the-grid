import {ErrorCode, ServerError} from "../server/server-error.js";
import {ExerciseGroup} from "../model/api/exercise-group.js";

export class ExerciseService {

    dbClientFactory;
    exerciseDao;

    constructor(dbClientFactory, exerciseDao) {
        if (!dbClientFactory) {
            throw new Error('Db Client Factory was not provided.');
        }

        if (!exerciseDao) {
            throw new Error('Exercise dao was not provided.');
        }

        this.dbClientFactory = dbClientFactory;
        this.exerciseDao = exerciseDao;
    }

    async getAllExerciseGroups() {
        return Object.values(ExerciseGroup).sort((a, b) => a.order - b.order);
    }

    async getAllExercisesForGroup(groupName, client) {
        try {
            const groupNames = Object.values(ExerciseGroup).map(eg => eg.name);
            if (!groupNames.includes(groupName)) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided group name is not valid.');
            }

            if (!client) {
                client = await this.dbClientFactory.obtain();
            }

            return await this.exerciseDao.getAllExercisesForGroup(client, groupName);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting all exercises for a given group.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    async getExercisesByIds(ids, client) {
        try {
            if (!ids || !Array.isArray(ids)) {
                return {};
            }

            const uniqueIds = Array.from(new Set(ids));
            if (uniqueIds.length === 0) {
                return {};
            }

            if (!client) {
                client = await this.dbClientFactory.obtain();
            }

            const exercisesByIds = {};
            const exercises = await this.exerciseDao.getExercisesForIds(client, uniqueIds);
            for (const exercise of exercises) {
                exercisesByIds[exercise.id] = exercise;
            }
            return exercisesByIds;
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting exercises by ids.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }
}
