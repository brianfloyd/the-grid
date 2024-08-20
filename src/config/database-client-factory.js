import pg from 'pg';
const { Pool } = pg;

export class DatabaseClientFactory {

    pool;

    constructor(connectionString) {
        this.pool = new Pool({
            connectionString: connectionString,
            connectionTimeoutMillis: 5000
        });
    }

    async obtain() {
        return this.pool.connect();
    }

}
