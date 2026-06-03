const serverStatus = document.querySelector("#server-status");
const databaseStatus = document.querySelector("#database-status");
const uptime = document.querySelector("#uptime");
const version = document.querySelector("#version");
const message = document.querySelector("#message");
const refreshButton = document.querySelector("#refresh-button");

async function loadStatus() {
  setPending();

  try {
    const response = await fetch("/api/status");

    if (!response.ok) {
      throw new Error(`Status request failed: ${response.status}`);
    }

    const status = await response.json();

    setStatusText(serverStatus, status.server);
    setStatusText(databaseStatus, status.database);
    uptime.textContent = status.uptime;
    version.textContent = status.version;

    if (status.server === "online" && status.database === "online") {
      message.textContent = "Nami-chan says the server lights are glowing and the database gremlin is behaving.";
    } else if (status.server === "online") {
      message.textContent = "Nami-chan is online, but she is side-eyeing the database very intensely.";
    } else {
      message.textContent = "Nami-chan cannot find the server heartbeat. Tiny panic blanket deployed.";
    }
  } catch (error) {
    console.error(error);
    setStatusText(serverStatus, "offline");
    setStatusText(databaseStatus, "unknown");
    uptime.textContent = "...";
    version.textContent = "...";
    message.textContent = "Nami-chan tried to check the server, but the status lantern flickered out.";
  }
}

function setPending() {
  serverStatus.textContent = "Checking...";
  databaseStatus.textContent = "Checking...";
  serverStatus.className = "";
  databaseStatus.className = "";
  uptime.textContent = "...";
  version.textContent = "...";
}

function setStatusText(element, value) {
  element.textContent = value;
  element.className = value === "online" ? "online" : "offline";
}

refreshButton.addEventListener("click", loadStatus);

loadStatus();
setInterval(loadStatus, 10000);