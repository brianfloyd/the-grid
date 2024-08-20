import {WorkoutSet} from "../../model/api/workout-set.js";

export class SetDao {

    static SELECT_SETS_FOR_WORKOUT = 'select * from the_grid.set where set_wrk_id = $1';

    async getSetsForWorkout(client, workoutId) {
        const result = await client.query(SetDao.SELECT_SETS_FOR_WORKOUT, [workoutId]);
        return SetRowMapper.map(result);
    }
}

class SetRowMapper {

    static map(result) {
        const sets = [];
        for (const row of result.rows) {
            // TODO: Make a data model of sets.
            const set = new WorkoutSet();
            set.id = Number(row['set_id']);
            set.exerciseId = Number(row['set_exr_id']);
            set.workoutId = Number(row['set_wrk_id']);
            set.reps = Number(row['set_reps']);
            set.weight = Number(row['set_weight']);
            sets.push(set);
        }
        return sets;
    }
}
