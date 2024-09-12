import {Exercise} from "../../model/api/exercise.js";

export class ExerciseDao {

    static SELECT_ALL_EXERCISES_FOR_GROUP = 'select * from the_grid.exercise where exr_group = $1';
    static SELECT_EXERCISES_FOR_IDS = `SELECT e.exr_id, e.exr_name, e.exr_group, 
    r.rep_id, r.rep_exr_id, r.rep_update, r.rep_number, r.rep_weight
    FROM the_grid.exercise e
    LEFT JOIN the_grid.reps r ON e.exr_id = r.rep_exr_id
    where e.exr_id= ANY ($1)`
    
    async getAllExercisesForGroup(client, groupName) {
        const result = await client.query(ExerciseDao.SELECT_ALL_EXERCISES_FOR_GROUP, [groupName]);
        return ExerciseRowMapper.map(result);
    }

    async getExercisesForIds(client, ids) {
        const result = await client.query(ExerciseDao.SELECT_EXERCISES_FOR_IDS, [ids]);
        return ExerciseRowMapper.map(result);
    }
    
}

class ExerciseRowMapper {

    static map(result) {
     
        const exercises = [];
        for (const row of result.rows) {
            // TODO: Make a data model of exercises.
            const exercise = new Exercise();
            exercise.id = Number(row['exr_id']);
            exercise.name = row['exr_name'];
            exercise.group = row['exr_group'];
            if (row['rep_number']) {
                exercise.repNumber = row['rep_number'];
            }
            if (row['rep_weight']) {
                exercise.repWeight = row['rep_weight'];
            }
            exercises.push(exercise);
        }

        return exercises;
    }
}
