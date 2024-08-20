import {Exercise} from "../../model/api/exercise.js";

export class ExerciseDao {

    static SELECT_ALL_EXERCISES_FOR_GROUP = 'select * from the_grid.exercise where exr_group = $1';
    static SELECT_EXERCISES_FOR_IDS = 'select * from the_grid.exercise where exr_id = ANY ($1)';

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
            exercises.push(exercise);
        }
        return exercises;
    }
}
