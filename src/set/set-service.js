
export class SetService {

    databaseClientFactory;
    setDao;

    constructor(databaseClientFactory, setDao) {
        if (!databaseClientFactory) {
            throw new Error('Database client factory was not provided.');
        }

        if (!setDao) {
            throw new Error('Set dao was not provided.');
        }


        this.databaseClientFactory = databaseClientFactory;
        this.setDao = setDao;
    }

    async getSetsForWorkout(workoutId) {

    }

}
