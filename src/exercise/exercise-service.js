import {ErrorCode, ServerError} from "../server/server-error.js";

export class ExerciseService {

    // Hardcoded exercises...for now!
    static EXERCISES = {
        1: {
            id: 1,
            group: 'BICEP',
            name: 'Bicep Curls'
        },
        2: {
            id: 2,
            group: 'TRICEP',
            name: 'Tricep Extensions'
        }
    }

    constructor() {

    }

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
