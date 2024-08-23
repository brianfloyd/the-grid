
export class SetView {
    id;
    exerciseId;
    name;
    group;
    weight;
    reps;
    count;

    constructor(set) {
        this.id = set.id;
        this.exerciseId = set.exerciseId;
        this.weight = set.weight;
        this.reps = set.reps;
        this.count = set.count;
    }

}
