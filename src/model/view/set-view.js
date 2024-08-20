
export class SetView {
    id;
    exerciseId;
    name;
    group;
    weight;
    reps;

    constructor(set) {
        this.id = set.id;
        this.exerciseId = set.exerciseId;
        this.weight = set.weight;
        this.reps = set.reps;
    }

}
