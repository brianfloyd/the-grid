import {WorkoutData} from "./workout-data.js";

export class WorkoutDao {

    static SELECT_WORKOUT_FOR_DATE= 'select * from the_grid.workout where wrk_date = $1';

    async getWorkoutForDate(client, date) {
        const result = await client.query(WorkoutDao.SELECT_WORKOUT_FOR_DATE, [date]);
        return WorkoutDataRowMapper.map(result);
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
