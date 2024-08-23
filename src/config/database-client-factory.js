import pg from 'pg';
const { Pool } = pg;

export class DatabaseClientFactory {

    pool;

    constructor(connectionString) {
        this.pool = new Pool({
            max: 2,
            connectionString: connectionString,
            connectionTimeoutMillis: 5000,
            idleTimeoutMillis: 0
        });
    }

    async obtain() {
        return this.pool.connect();
    }

}
