import {WorkoutSet} from "../../model/api/workout-set.js";

export class SetDao {

    static NEXT_SEQUENCE_VALUE = `nextval('the_grid.seq_set_id')`;
    static SELECT_SET_FOR_ID = 'select * from the_grid.set where set_id = $1';
    static SELECT_SETS_FOR_WORKOUT = 'select * from the_grid.set where set_wrk_id = $1';
    static INSERT_SET = 'insert into the_grid.set (set_id, set_exr_id, set_wrk_id, set_reps, set_weight) values' +
        ` (${SetDao.NEXT_SEQUENCE_VALUE}, $1, $2, $3, $4) returning set_id`;

    async getSetForId(client, id) {
        const result = await client.query(SetDao.SELECT_SET_FOR_ID, [id]);
        return SetRowMapper.map(result);
    }

    async getSetsForWorkout(client, workoutId) {
        const result = await client.query(SetDao.SELECT_SETS_FOR_WORKOUT, [workoutId]);
        return SetRowMapper.map(result);
    }

    async insertSet(client, set) {
        const result = await client.query(SetDao.INSERT_SET, [
            set.exerciseId,
            set.workoutId,
            set.reps,
            set.weight
        ]);
        return result.rows[0]['set_id'];
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
