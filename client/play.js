const username = sessionStorage.getItem("username");
const roomID = sessionStorage.getItem("roomID");

if (!username || !roomID) {
  window.location.href = "/?invalid=true";
}

const WS_SERVER_URI = "ws://" + window.location.host;

let GAME_RUNNING = false;

const ws = new WebSocket(WS_SERVER_URI + "/join");

ws.onopen = () => {
  ws.send(JSON.stringify({ username, roomID }));
};

const messages = [];

ws.onmessage = (msg) => {
  parsedMsg = JSON.parse(msg.data);
  if (parsedMsg.type === "chat" || parsedMsg.type === "system") {
    messages.push(parsedMsg.content);
    renderMessages();
  } else if (parsedMsg.type === "start") {
    GAME_RUNNING = true;
    document.getElementById("startGameButton").disabled = true;
  } else if (parsedMsg.type === "end") {
    GAME_RUNNING = false;
    document.getElementById("startGameButton").disabled = false;
    displayError("Game has ended!");
  } else if (parsedMsg.type === "error") {
    displayError(parsedMsg.content.message);
  } else if (parsedMsg.type === "board") {
    renderBoard(parsedMsg.content);
  } else if (parsedMsg.type === "win" || parsedMsg.type === "loss") {
    document.getElementById("gameMessages").innerText =
      parsedMsg.content.message;
  } else if (parsedMsg.type === "move") {
    let move = parsedMsg.content;
    markMove(move);
  } else if (parsedMsg.type === "turn") {
    document.getElementById("gameMessages").innerText =
      parsedMsg.content.message;
  }
};

ws.onclose = (e) => {
  window.location.href = "/?invalid=true";
  console.log(e);
};

const chatInputForm = document.getElementById("chatInputForm");
chatInputForm.onclick = (e) => {
  e.preventDefault();
  if (ws.readyState === WebSocket.CLOSED) {
    return;
  }
  const chatText = document.getElementById("chatTextInput").value;
  if (!chatText) {
    return;
  }
  ws.send(
    JSON.stringify({
      type: "chat",
      content: { message: chatText, author: username },
    })
  );
  document.getElementById("chatTextInput").value = "";
  messages.push({ author: "You", message: chatText });
  renderMessages();
};

function renderMessages() {
  const messageSpace = document.getElementById("message-space");
  messageSpace.scrollIntoView = false;
  messageSpace.innerHTML = "";
  if (messages.length > 20) {
    messages.splice(0, 1);
  }
  for (const msg of messages) {
    const parentDiv = document.createElement("div");
    const authorEl = document.createElement("p");
    const msgEl = document.createElement("p");
    msgEl.innerText = msg.message;
    authorEl.classList.add("author");
    authorEl.innerText = msg.author;
    parentDiv.classList.add("message-entity");
    if (msg.author === "You") {
      parentDiv.classList.add("own-message");
    }
    parentDiv.appendChild(authorEl);
    parentDiv.appendChild(msgEl);
    messageSpace.appendChild(parentDiv);
  }
}

const HUDElement = document.getElementById("hud");
let usernameEl = document.createElement("p");
let roomIDEl = document.createElement("p");
usernameEl.innerText = "Username: " + username;
roomIDEl.innerText = "Room ID: " + roomID;
roomIDEl.style.cursor = "pointer";
roomIDEl.onclick = () => {
  navigator.clipboard.writeText(roomID);
  displayError("copied text to clipboard!");
};
HUDElement.appendChild(usernameEl);
HUDElement.appendChild(roomIDEl);

// ---------------- game logics -------------------
document.getElementById("startGameButton").onclick = function startGame() {
  ws.send(JSON.stringify({ type: "start" }));
};

function renderBoard(board) {
  const bingoBoard = document.getElementById("bingoBoard");
  bingoBoard.innerHTML = "";
  for (let i = 0; i < board.length; i++) {
    for (let j = 0; j < board.length; j++) {
      const bingoBox = document.createElement("div");
      bingoBox.classList.add("bingoBox");
      bingoBox.id = board[i][j];
      bingoBox.innerText = board[i][j];
      bingoBox.addEventListener("click", makeMove);
      bingoBoard.appendChild(bingoBox);
    }
  }
}
renderBoard([
  [1, 2, 3, 4, 5],
  [6, 7, 8, 9, 10],
  [11, 12, 13, 14, 15],
  [16, 17, 18, 19, 20],
  [21, 22, 23, 24, 25],
]);

function makeMove(e) {
  if (!GAME_RUNNING) {
    return;
  }
  ws.send(JSON.stringify({ type: "move", content: parseInt(e.target.id) }));
}

function markMove(move) {
  const bingoBoard = document.getElementById("bingoBoard");
  for (let bingoBox of bingoBoard.children) {
    if (bingoBox.id == move) {
      bingoBox.classList.add("marked");
    }
  }
}
