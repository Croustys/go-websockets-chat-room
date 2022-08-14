const send = document.querySelector("#send")
const textWrapperDiv = document.querySelector('#text-wrapper')
textWrapperDiv.style.display = 'none';

let username;

const name = document.querySelector("#name-input");
const okBtn = document.querySelector("#name-input-btn")
okBtn.addEventListener("click", () => {
  username = name.value;
  document.querySelector('#modal').hidden = true;
  textWrapperDiv.style.display = "block";
})

const socket = new WebSocket("ws://localhost:8080/websocket")

socket.addEventListener("message", (s) => {
  const { data } = s
  const message = JSON.parse(data)

  const parentDiv = document.querySelector("#messages")
  const wrapper = document.createElement("div")
  const newMessageAuthor = document.createElement("p")
  const newMessage = document.createElement("h1")

  newMessage.textContent = message.text;
  newMessageAuthor.textContent = message.username;

  wrapper.appendChild(newMessageAuthor)
  wrapper.appendChild(newMessage)
  wrapper.style.border = "1px solid black";
  parentDiv.appendChild(wrapper)
})

send.addEventListener('click', () => {
  const text = document.querySelector("#msg").value;
  const data = {
    username,
    text,
  }
  socket.send(JSON.stringify(data))

  document.querySelector("#msg").value = "";
})