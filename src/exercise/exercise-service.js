import {ErrorCode, ServerError} from "../server/server-error.js";
import {ExerciseGroup} from "../model/api/exercise-group.js";

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

    constructor() {

    }

    async getAllExerciseGroups() {
        return Object.values(ExerciseGroup).sort((a, b) => a.order - b.order);
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
