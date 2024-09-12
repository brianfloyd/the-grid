import {ErrorCode, ServerError} from "../server/server-error.js";
import {WorkoutView} from "../model/view/workout-view.js";
import {SetView} from "../model/view/set-view.js";
import {transactional} from "../utils.js";

/**
 * Orchestration layer that converts workout information into a display ready format.
 */
export class WorkoutDisplayService {

    workoutService;
    setService;
    exerciseService;
    dbClientFactory;

    constructor(workoutService, setService, exerciseService, dbClientFactory) {
        if (!workoutService) {
            throw new Error('Workout service was not provided.');
        }

        if (!setService) {
            throw new Error('Set service was not provided.');
        }

        if (!exerciseService) {
            throw new Error('Exercise service was not provided.');
        }

        if (!dbClientFactory) {
            throw new Error('dbClientFactory was not provided.');
        }

        this.workoutService = workoutService;
        this.exerciseService = exerciseService;
        this.setService = setService;
        this.dbClientFactory = dbClientFactory;
    }

    /**
     * Creates a new workout for today. Associates any provided sets with the new workout.
     */
    async saveNewWorkout(dateString, exerciseIds) {
        let client;
        try {
            client = await this.dbClientFactory.obtain();
            await transactional(client, async () => {
                const workoutId = await this.workoutService.createWorkoutForDate(dateString, client);
                if (exerciseIds && exerciseIds.length > 0) {
                    const defaultSets = await this.createDefaultSets(workoutId,exerciseIds, client);
                    for (const set of defaultSets) {
                        await this.setService.createSet(set, client);
                    }
                }
            });
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while creating a workout for today.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        } finally {
            if (client) {
                client.release();
            }
        }
    }

    async addNewSetToWorkout(workoutId, exerciseId) {
        try {
            if (exerciseId <= 0) {
                throw new ServerError('Provided exercise id is not valid.');
            }

            if (workoutId <= 0) {
                throw new ServerError('Provided workout id is not valid.')
            }


            const defaultSets = await this.createDefaultSets(workoutId, [exerciseId]);
            await this.setService.createSet(defaultSets[0]);
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while creating a workout for today.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    /**
     * Gets a workout for a given date.
     */
    async getWorkoutViewForDate(dateString) {
        try {
            const workout = await this.workoutService.getWorkoutForDate(dateString);
            let workoutView = new WorkoutView(workout);
            workoutView.date = dateString;
            workoutView = await this.populateSets(workoutView, workout);
            return workoutView;
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting workout by date.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    async createDefaultSets(workoutId, exerciseIds, client) {
        const exercises = await this.exerciseService.getExercisesByIds(exerciseIds, client);

        const sets = [];
        for (const exerciseId of exerciseIds) {
            const exercise = exercises[exerciseId];
            if (!exercise) {
                throw new ServerError(ErrorCode.INVALID_REQUEST, `An unrecognized exercise id ${exerciseId} was given.`);
            }

            const set = {
                workoutId: workoutId,
                exerciseId: exerciseId,
                weight: exercise.repWeight || 0,
                reps: exercise.repNumber || 0,
                count: 0
            }
            sets.push(set);
        }

        return sets;
    }


    async populateSets(workoutView, workout) {
        let setViews = [];
        for (const set of workout.sets) {
            const setView = new SetView(set);
            setViews.push(setView);
        }
        setViews = await this.populateExercises(setViews);
        workoutView.sets = setViews;
        return workoutView;
    }

    async populateExercises(setViews) {
        const exerciseIds = setViews.map(sv => sv.exerciseId);

        const exercisesByIds = await this.exerciseService.getExercisesByIds(exerciseIds);
        for (const setView of setViews) {
            const exercise = exercisesByIds[setView.exerciseId];
            if (exercise) {
                setView.name = exercise.name;
                setView.group = exercise.group;
            } else {
                setView.name = 'NA';
                setView.group = 'NA';
            }
        }

        return setViews;
    }
}
