
import fs from 'fs';
import { setupServer } from './server/server.js';

const env = process.env.NODE_ENV || 'test';
const configPath = `./config/${env}.json`;


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

const server = setupServer();
server.listen(PORT, '0.0.0.0', () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
