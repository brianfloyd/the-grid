import {ErrorCode, ServerError} from "../server/server-error.js";
import {ExerciseGroup} from "../model/api/exercise-group.js";
import {transactional} from "../utils.js";

export class ExerciseService {

    // Hardcoded exercises...for now!
    static EXERCISES = {
        1: {
            id: 1,
            group: ExerciseGroup.BICEP.name,
            name: 'Bicep Curls'
        },
        2: {
            id: 2,
            group: ExerciseGroup.TRICEP.name,
            name: 'Tricep Extensions'
        }
    }

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

    async getAllExercisesForGroup(groupName) {
        try {
            const groupNames = Object.values(ExerciseGroup).map(eg => eg.name);
            if (!groupNames.includes(groupName)) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, 'The provided group name is not valid.');
            }

            const client = await this.dbClientFactory.obtain();
            return await transactional(client, async () => {
                return await this.exerciseDao.getAllExercisesForGroup(client, groupName);
            });
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting all exercises for a given group.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    // TODO: Remove me.
    async getExercisesByIds(ids) {
        try {
            if (!ids || !Array.isArray(ids)) {
                return {};
            }
            const uniqueIds = new Set(ids);

            const exercisesByIds = {};
            for (const id of uniqueIds) {
                const exercise = ExerciseService.EXERCISES[id];
                if (exercise) {
                    exercisesByIds[id] = exercise;
                }
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
