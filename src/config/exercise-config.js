import {ExerciseService} from "../exercise/exercise-service.js";
import {ExerciseViewResource} from "../exercise/exercise-view-resource.js";
import {ExerciseDao} from "../exercise/data/exercise-dao.js";

export function configureExerciseObjects(dbClientFactory) {
    const exerciseDao = new ExerciseDao();
    const exerciseService = new ExerciseService(dbClientFactory, exerciseDao);
    const exerciseViewResource = new ExerciseViewResource(exerciseService);
    return {exerciseService, exerciseViewResource};
}
