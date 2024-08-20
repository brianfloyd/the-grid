import express from 'express';
import bodyParser from "body-parser";
import path from 'path';

import {WorkoutViewResource} from "./workout/workout-view-resource.js";
import {WorkoutDisplayService} from "./workout/workout-display-service.js";

// Some backwards compatability hack for module versions of javascript.
import {fileURLToPath} from 'url';
import {dirname} from 'path';
import {GlobalErrorHandler} from "./server/global-error-handler.js";
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.static(path.join(__dirname, '../public')));
app.use(bodyParser.json());

const workoutDisplayService = new WorkoutDisplayService();
const workoutResource = new WorkoutViewResource(workoutDisplayService);

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, '../public', 'index.html'));
});
workoutResource.bind(app);
app.use(GlobalErrorHandler.handle);

app.listen(PORT, '0.0.0.0', () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
