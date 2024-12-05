const u = new UserRender();
const g = new TheGridRender();

const bindingElement = document.getElementById("root")

function globalRender() {
    let render;
    if (!localStorage.getItem("user")) {
        render = u.render()
    } else {
        render = g.render();
    }

    if (render) {
        bindingElement.replaceChildren(render);
    }
}

localStorage.setItem("user", "46c6e118-8da6-425b-a06e-105a027ce668")
globalRender()
