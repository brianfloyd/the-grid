class TheGridRender {

    constructor() {
        this.userId = localStorage.getItem("user");
        this.date = '12/04/2024'
        this.request = null;
        this.workout = null;
    }

    async fecthData() {
        const response = await fetch(`/workouts/date`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                userId: this.userId,
                date: this.date
            })
        });
        const json = await response.json();
        this.workout = json.workout;
        this.request = null;
    }

    render() {
        if (!this.workout || this.workout.date !== this.date) {
            if (!this.request) {
                const r = this.fecthData().then(globalRender)
                this.request = r;
            }

            const div = document.createElement("div");
            const h1 = document.createElement("h1")
            h1.innerText = `Loading data for id ${this.date}`;
            div.appendChild(h1);

            return div;
        } else {
            const p = document.createElement("p");
            p.innerText = JSON.stringify(this.workout);
            return p;
        }
    }

}