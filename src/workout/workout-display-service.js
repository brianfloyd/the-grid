import {ErrorCode, ServerError} from "../server/server-error.js";
import {WorkoutView} from "../model/view/workout-view.js";
import {SetView} from "../model/view/set-view.js";

/**
 * Orchestration layer that converts workout information into a display ready format.
 */
export class WorkoutDisplayService {

    workoutService;
    exerciseService

    constructor(workoutService, exerciseService) {
        if (!workoutService) {
            throw new Error('Workout service was not provided.');
        }

        if (!exerciseService) {
            throw new Error('Exercise service was not provided.');
        }

        this.workoutService = workoutService;
        this.exerciseService = exerciseService;
    }

    /**
     * Gets a workout for a given date.
     */
    async getWorkoutViewForDate(dateString) {
        try {
            console.time('WorkoutDisplayService#TOTAL-getWorkoutViewForDate');
            // TODO: Switch to passing the date when the front end is ready for it.
            console.time('WorkoutDisplayService#getWorkoutForDate');
            const workout = await this.workoutService.getWorkoutForDate('08/19/2024');
            console.timeEnd('WorkoutDisplayService#getWorkoutForDate');

            let workoutView = new WorkoutView(workout);
            workoutView.date = dateString;
            console.time('WorkoutDisplayService#populateSets');
            workoutView = await this.populateSets(workoutView, workout);
          
            console.timeEnd('WorkoutDisplayService#populateSets');
            console.timeEnd('WorkoutDisplayService#TOTAL-getWorkoutViewForDate');
            return workoutView;
        } catch (e) {
            if (e instanceof ServerError) {
                throw e;
            }

            console.error('An error occurred while getting workout by date.', e);
            throw new ServerError(ErrorCode.GENERIC_ERROR, e.message);
        }
    }

    async populateSets(workoutView, workout) {
        let setViews = [];
        for (const set of workout.sets) {
            const setView = new SetView(set);
            setViews.push(setView);
        }
        setViews = await this.populateExercises(setViews);

        workoutView.sets = setViews;
        console.log(workoutView)
        return workoutView;
    }

    async populateExercises(setViews) {
        const exerciseIds = setViews.map(sv => sv.exerciseId);

        const exercisesByIds = await this.exerciseService.getExercisesByIds(exerciseIds);
        for (const setView of setViews) {
            console.log(setView)
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
