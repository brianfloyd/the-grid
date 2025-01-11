
const globalErrorPlaceholder = document.createElement("div")
globalErrorPlaceholder.id = "global-error"

function removeElement(event) {
    event.target.replaceWith(globalErrorPlaceholder)
}

function setupAnimation() {
    const globalError = document.getElementById('global-error');
    if (globalError) {
        globalError.addEventListener('animationend', removeElement);
    }
}
document.addEventListener('htmx:afterRequest', setupAnimation)
