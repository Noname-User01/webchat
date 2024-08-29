"use strict";
//厳密に記述する設定
document.addEventListener("DOMContentLoaded", () => {
  function getWebsockerUri() {
    let loc = window.location;
    let uri = "ws:";
    if (loc.protocol === "https:") {
      //httpsの場合はwebsocketも追加する
      uri = "wss:";
    }
    uri += "//" + loc.host;
    uri += loc.pathname + "ws";
    return uri;
  }

  const chat_list = document.getElementById("chat_list");

  function addChat(msg) {
    console.log(msg);

    //createElement()で要素を作成
    const li = document.createElement("li");
    const div = document.createElement("article");
    div.className = "chat_box";

    const name = document.createElement("div");
    name.className = "chat_name";
    name.textContent = msg.name;

    const text = document.createElement("div");
    text.className = "chat_text";
    text.innerText = msg.text;

    div.appendChild(name);
    div.appendChild(text);
    li.appendChild(div);
    chat_list.appendChild(li);
  }
  //型作成
  /** @type {WebSocket} */
  let ws;

  function connectWebSocket() {
    const uri = getWebsockerUri();
    ws = new WebSocket(uri);
    ws.addEventListener("open", () => {
      console.log("Connected");
      ws.send(JSON.stringify({ type: "get" }));
    });
    ws.addEventListener("close", () => {
      console.log("Closed");

      setTimeout(() => {
        console.log("Reconnect");
        connectWebSocket();
      }, 1000);
    });
    ws.addEventListener("error", (event) => {
      console.log("Error", event);
    });
    ws.addEventListener("message", (event) => {
      const data = JSON.parse(event.data);
      console.log(data);
      switch (data.type) {
        case "messages":
          while (chat_list.firstChild) {
            chat_list.removeChild(chat_list.firstChild);
          }

          data.obj.forEach((msg) => {
            addChat(msg);
            7;
          });
          break;
        case "append":
          addChat(data.obj);
          break;
      }
    });
  }

  connectWebSocket();

  const name = document.getElementById("submit_name");
  const text = document.getElementById("submit_text");
  const button = document.getElementById("submit_button");

  function button_update() {
    if (text.value === "") {
      button.disabled = true;
    } else {
      button.disabled = false;
    }
  }

  button.addEventListener("click", () => {
    if (name.value === "") {
      alert("名前を入力してください");
      return;
    }
    const msg = {
      type: "send",
      obj: {
        name: name.value,
        text: text.value,
      },
    };
    console.log(msg);

    ws.send(JSON.stringify(msg));

    text.value = "";
    button_update();
  });

  text.addEventListener("input", () => {
    button_update();
  });
});
