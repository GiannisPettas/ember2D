const statusEl = document.getElementById("ws-status");
const pingBtn = document.getElementById("ping-btn");
const logEl = document.getElementById("log");

function appendLog(line) {
    const ts = new Date().toLocaleTimeString();
    logEl.textContent += `[${ts}] ${line}\n`;
    logEl.scrollTop = logEl.scrollHeight;
}

function connectWS() {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsURL = `${protocol}//${window.location.host}/ws`;
    const socket = new WebSocket(wsURL);

    socket.addEventListener("open", () => {
        statusEl.textContent = "WS: connected";
        pingBtn.disabled = false;
        appendLog(`Connected to ${wsURL}`);
    });

    socket.addEventListener("message", (event) => {
        appendLog(`<= ${event.data}`);
    });

    socket.addEventListener("close", () => {
        statusEl.textContent = "WS: disconnected (reconnecting...)";
        pingBtn.disabled = true;
        appendLog("Connection closed");
        setTimeout(connectWS, 1000);
    });

    socket.addEventListener("error", () => {
        appendLog("WebSocket error");
    });

    pingBtn.onclick = () => {
        const message = {
            action: "ping",
            data: { client_time: new Date().toISOString() },
        };
        const payload = JSON.stringify(message);
        socket.send(payload);
        appendLog(`=> ${payload}`);
    };
}

pingBtn.disabled = true;
connectWS();
