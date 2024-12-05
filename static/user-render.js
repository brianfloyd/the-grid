class UserRender {

    constructor() {
        this.users = [];
        this.initial = true;
    }

    async fetchData() {
        const response = await fetch("/users");
        const json = await response.json()
        this.users = json.users
    }

    render() {
        if (this.initial) {
            this.initial = false;

            const h1 = document.createElement("h1")
            h1.innerText = "Loading users..."
            this.fetchData().then(globalRender)
            return h1;
        } else {
            const userContainer = document.createElement("div");
            for (const user of this.users) {
                const uElement = document.createElement("p");
                uElement.innerText = user.id + " " + user.name;
                userContainer.appendChild(uElement)
            }
            return userContainer;
        }
    }
}