import {WorkoutSet} from "../../model/api/workout-set.js";

export class SetDao {

    static NEXT_SEQUENCE_VALUE = `nextval('the_grid.seq_set_id')`;
    static SELECT_SET_FOR_ID = `SELECT * from the_grid.exercise e where e.exr_id=($1)`;

    static SELECT_SETS_FOR_WORKOUT = `SELECT 
    s.*, -- Select all columns from the set table
    r.rep_id, r.rep_exr_id, r.rep_update, r.rep_number, r.rep_weight -- Select specific columns from the reps table
    FROM the_grid.set s
    LEFT JOIN the_grid.reps r ON s.set_exr_id = r.rep_exr_id -- Join reps table based on the connecting key
    WHERE s.set_wrk_id = $1;
    `;
    static SELECT_DEFAULT_REPS_FOR_EXCERCISE_ID = `SELECT rep_weight, rep_number, rep_update
    FROM the_grid.reps
    WHERE rep_exr_id = $1
    ORDER BY rep_update DESC
    LIMIT 1;
    `
    static INSERT_SET = 'insert into the_grid.set (set_id,set_exr_id, set_wrk_id, set_reps, set_weight, set_count)' +
        ` values (${SetDao.NEXT_SEQUENCE_VALUE}, $1, $2, $3, $4, $5) returning set_id`;
    static UPDATE_SET = 'update the_grid.set set set_exr_id = $1, set_wrk_id = $2, set_reps = $3, set_weight = $4,' +
        ' set_count = $5 where set_id = $6';




    async getSetForId(client, id) {
        const result = await client.query(SetDao.SELECT_SET_FOR_ID, [id]);
        return SetRowMapper.map(result);
    }

    async getSetsForWorkout(client, workoutId) {
        const result = await client.query(SetDao.SELECT_SETS_FOR_WORKOUT, [workoutId]);
        return SetRowMapper.map(result);
    }

    async insertSet(client, set) {
        const result = await client.query(SetDao.INSERT_SET, this.getParameters(set));
        return result.rows[0]['set_id'];
    }

    async updateSet(client, set) {
        await client.query(SetDao.UPDATE_SET, [...this.getParameters(set), set.id]);
    }

    async getDefaultRepsForExcerciseId(client,id,set){
        let result = await client.query(SetDao.SELECT_DEFAULT_REPS_FOR_EXCERCISE_ID,[id]);
        let rep=result.rows[0];
        console.log(rep)
        set.weight=rep.rep_weight;
        set.reps=rep.rep_number;
        set.count=0;
       return set;
    }

    getParameters(set) {
        return [
            set.exerciseId,
            set.workoutId,
            set.reps,
            set.weight,
            set.count
        ];
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
            set.reps = Number(row['rep_number']);
            set.weight = Number(row['rep_weight']);
            set.count = Number(row['set_count']);
            sets.push(set);
        }
        return sets;
    }
}
