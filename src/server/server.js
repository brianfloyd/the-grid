import express from 'express';
import bodyParser from "body-parser";
import path from 'path';

import {WorkoutViewResource} from "../workout/workout-view-resource.js";
import {WorkoutDisplayService} from "../workout/workout-display-service.js";

// Some backwards compatability hack for module versions of javascript.
import {fileURLToPath} from 'url';
import {dirname} from 'path';
import {GlobalErrorHandler} from "./global-error-handler.js";
import {WorkoutService} from "../workout/workout-service.js";
import {ExerciseService} from "../exercise/exercise-service.js";
import {ExerciseViewResource} from "../exercise/exercise-view-resource.js";
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);


export function setupServer() {
    const app = express();

    app.use(express.static(path.join(__dirname, '../../public')));
    app.use(bodyParser.json());

    const workoutService = new WorkoutService();
    const exerciseService = new ExerciseService();

    const workoutDisplayService = new WorkoutDisplayService(workoutService, exerciseService);

    const workoutViewResource = new WorkoutViewResource(workoutDisplayService);
    const exerciseViewResource = new ExerciseViewResource(exerciseService);

    app.get('/', (req, res) => {
        res.sendFile(path.join(__dirname, '../../public', 'index.html'));
    });

    workoutViewResource.bind(app);
    exerciseViewResource.bind(app);

    app.use(GlobalErrorHandler.handle);

    return app;
}
