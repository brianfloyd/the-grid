import {WorkoutData} from "./workout-data.js";

export class WorkoutDao {

    static NEXT_WORKOUT_ID_SEQUENCE = `nextval('the_grid.seq_wrk_id')`
    static SELECT_WORKOUT_FOR_DATE= 'select * from the_grid.workout where wrk_date = $1';
    static INSERT_WORKOUT_FOR_TODAY = `insert into the_grid.workout (wrk_id, wrk_date)` +
        `values (${WorkoutDao.NEXT_WORKOUT_ID_SEQUENCE}, $1) returning wrk_id`;

    async getWorkoutForDate(client, date) {
        const result = await client.query(WorkoutDao.SELECT_WORKOUT_FOR_DATE, [date]);
        return WorkoutDataRowMapper.map(result);
    }

    async insertWorkoutForDate(client, dateString) {
        const result = await client.query(WorkoutDao.INSERT_WORKOUT_FOR_TODAY, [dateString]);
        return result.rows[0].wrk_id;
    }

}

class WorkoutDataRowMapper {
    static map(result) {
        const data = [];
        for (const row of result.rows) {
            const workoutData = new WorkoutData();
            workoutData.id = Number(row['wrk_id']);
            workoutData.date = row['wrk_date'];
            data.push(workoutData);
        }

        return data;
    }

}
