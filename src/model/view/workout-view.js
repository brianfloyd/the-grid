export class WorkoutView {
    id;
    date;
    sets;

    constructor(workout) {
        this.sets = [];
        this.id = workout.id;
        this.date = workout.date;
    }

    addSet(set) {
        if (set) {
            this.sets.push(set);
        }
    }
}
