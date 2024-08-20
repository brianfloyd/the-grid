import express from 'express';
import bodyParser from "body-parser";
import path from 'path';

// Some backwards compatability hack for module versions of javascript.
import {fileURLToPath} from 'url';
import {dirname} from 'path';
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Hack for error handling from async functions.
import 'express-async-errors';


export function setupServer() {
    const app = express();
    app.use(express.static(path.join(__dirname, '../../public')));
    app.use(bodyParser.json());
    app.get('/', (req, res) => {
        res.sendFile(path.join(__dirname, '../../public', 'index.html'));
    });
    return app;
}
