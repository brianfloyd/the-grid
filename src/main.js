import fs from 'fs';
import { setupServer } from './server/server.js';
import {DatabaseClientFactory} from "./config/database-client-factory.js";
import {GlobalErrorHandler} from "./server/global-error-handler.js";
import {configureExerciseObjects} from "./config/exercise-config.js";
import {configureWorkoutObjects} from "./config/workout-config.js";
import {_200} from "./server/server-utils.js";
import {configureSetObjets} from "./config/set-config.js";

const env = process.env.NODE_ENV || 'test';
const configPath = `./app-config/${env}.json`;

let config = {}
if (fs.existsSync(configPath)) {
    config = JSON.parse(fs.readFileSync(configPath).toString());
}

process.env = {
    ...process.env,
    ...config
}

const PORT = process.env.PORT || 3000;
const DATABASE_URL = process.env.DATABASE_URL || 'NO_DATABASE_URL';

// Configure and wire up services.
const dbClientFactory = new DatabaseClientFactory(DATABASE_URL);

const {
    exerciseService,
    exerciseViewResource
} = configureExerciseObjects(dbClientFactory);

const {
    setService
} = configureSetObjets(dbClientFactory);

const {
    workoutViewResource
} = configureWorkoutObjects(dbClientFactory, exerciseService, setService);


// Server Configuration.
const server = setupServer();
workoutViewResource.bind(server);
exerciseViewResource.bind(server);
new GlobalErrorHandler().bind(server);

server.listen(PORT, '0.0.0.0', () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
