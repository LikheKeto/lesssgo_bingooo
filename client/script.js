// let ws = new WebSocket("ws://localhost:8000/create");

// ws.onopen = () => {
//   ws.send(JSON.stringify({ username: "jivan" }));
// };

// ws.onmessage = (msg) => {
//   console.log(msg);
// };

// let sendButton = document.getElementById("sendButton");

// sendButton.onclick = () => {
//   let testString = document.getElementById("testString").value;
//   console.log("sending: ", testString);
//   ws.send(JSON.stringify({ move: testString }));
// };

const SERVER_URI = "http://localhost:8000";

const createRoomButton = document.getElementById("createRoomButton");
createRoomButton.onclick = async () => {
  console.log("sup");
  const username = document.getElementById("usernameInput").value;
  if (username === "") {
    displayError("Username must be provided!");
    return;
  }
  const res = await fetch(`${SERVER_URI}/create`);
  const parsedRes = await res.json();
  if (!parsedRes.roomID) {
    return;
  }
  const roomID = parsedRes.roomID;
  sessionStorage.setItem("username", username);
  sessionStorage.setItem("roomID", roomID);
  location.href = "/play.html";
};

const joinRoomButton = document.getElementById("joinRoomButton");
joinRoomButton.onclick = async () => {
  const username = document.getElementById("usernameInput").value;
  if (username === "") {
    displayError("Username must be provided!");
    return;
  }
  const roomID = document.getElementById("roomIDInput").value;
  if (roomID === "") {
    displayError("Room ID must be provided!");
    return;
  }
  sessionStorage.setItem("username", username);
  sessionStorage.setItem("roomID", roomID);
  location.href = "/play.html";
};

// check for redirect errors
const params = new Proxy(new URLSearchParams(window.location.search), {
  get: (searchParams, prop) => searchParams.get(prop),
});

if (params.invalid == "true") {
  displayError("Invalid request!");
}
