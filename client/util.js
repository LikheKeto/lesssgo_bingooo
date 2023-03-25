const ERROR_DISPLAY_INTERVAL = 3000;

function displayError(errMessage) {
  const errContainer = document.getElementById("error-box");
  const errEl = document.createElement("div");
  errEl.innerText = errMessage;
  errEl.classList.add("error-el");
  errContainer.appendChild(errEl);
  setTimeout(() => {
    for (const errEl of errContainer.children) {
      if (errEl.innerText == errMessage) {
        errContainer.removeChild(errEl);
        break;
      }
    }
  }, ERROR_DISPLAY_INTERVAL);
}
