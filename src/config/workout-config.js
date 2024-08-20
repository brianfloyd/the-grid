import {WorkoutService} from "../workout/workout-service.js";
import {WorkoutDisplayService} from "../workout/workout-display-service.js";
import {WorkoutViewResource} from "../workout/workout-view-resource.js";
import {WorkoutDao} from "../workout/data/workout-dao.js";

export function configureWorkoutObjects(dbClientFactory, exerciseService, setService) {
    const workoutDao = new WorkoutDao();
    const workoutService = new WorkoutService(dbClientFactory, workoutDao, setService);
    const workoutDisplayService = new WorkoutDisplayService(workoutService, exerciseService);
    const workoutViewResource = new WorkoutViewResource(workoutDisplayService);

    return {
        workoutDao,
        workoutService,
        workoutDisplayService,
        workoutViewResource
    };
}
