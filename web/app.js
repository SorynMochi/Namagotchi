const serverStatus = document.querySelector("#server-status");
const databaseStatus = document.querySelector("#database-status");
const uptime = document.querySelector("#uptime");
const version = document.querySelector("#version");
const message = document.querySelector("#message");
const refreshButton = document.querySelector("#refresh-button");
const seedButton = document.querySelector("#seed-button");
const playerPanel = document.querySelector("#player-panel");

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

async function seedDevPlayer() {
  seedButton.disabled = true;
  seedButton.textContent = "Creating...";

  try {
    const response = await fetch("/api/dev/seed-player", {
      method: "POST",
    });

    if (!response.ok) {
      throw new Error(`Seed request failed: ${response.status}`);
    }

    const result = await response.json();
    message.textContent = result.message;
    await loadPlayerStatus();
  } catch (error) {
    console.error(error);
    message.textContent = "Nami-chan tried to create the dev player, but a database goblin bit the clipboard.";
  } finally {
    seedButton.disabled = false;
    seedButton.textContent = "Create Dev Player";
  }
}

async function loadPlayerStatus() {
  try {
    const response = await fetch("/api/player/status");

    if (!response.ok) {
      playerPanel.innerHTML = `<p class="message">No player loaded yet. Create the dev player first.</p>`;
      return;
    }

    const status = await response.json();
    renderPlayerStatus(status);
  } catch (error) {
    console.error(error);
    playerPanel.innerHTML = `<p class="message">Could not load player status.</p>`;
  }
}

function renderPlayerStatus(status) {
  const player = status.player;
  const companion = status.companion;
  const resources = status.resources;

  playerPanel.innerHTML = `
    <div class="mini-grid">
      <div>
        <span class="label">Player</span>
        <strong>${escapeHTML(player.displayName)}</strong>
      </div>
      <div>
        <span class="label">Level</span>
        <strong>${player.level}</strong>
      </div>
      <div>
        <span class="label">Currency</span>
        <strong>${formatCurrency(player.currencyCents)}</strong>
      </div>
      <div>
        <span class="label">Mood</span>
        <strong>${companion.moodScore}</strong>
      </div>
    </div>

    <h3>${escapeHTML(companion.name)}</h3>

    <div class="stat-list">
      ${renderStat("Satiety", companion.satiety)}
      ${renderStat("Connection", companion.connection)}
      ${renderStat("Energy", companion.energy)}
      ${renderStat("Comfort", companion.comfort)}
      ${renderStat("Playfulness", companion.playfulness)}
      ${renderStat("Inspiration", companion.inspiration)}
      ${renderStat("Cleanliness", companion.cleanliness)}
    </div>

    <h3>Resources</h3>

    <div class="resource-grid">
      ${renderResource("Fans", resources.fans)}
      ${renderResource("Memes", resources.memes)}
      ${renderResource("Lost Items", resources.lostItems)}
      ${renderResource("Confidence", resources.confidence)}
      ${renderResource("Receipts", resources.receipts)}
      ${renderResource("Patterns", resources.patterns)}
      ${renderResource("Glitch Drops", resources.glitchDrops)}
    </div>
  `;
}

function renderStat(name, value) {
  return `
    <div class="stat-row">
      <span>${name}</span>
      <div class="stat-bar" aria-label="${name}: ${value}/100">
        <div class="stat-fill" style="width: ${value}%"></div>
      </div>
      <strong>${value}</strong>
    </div>
  `;
}

function renderResource(name, value) {
  return `
    <div>
      <span class="label">${name}</span>
      <strong>${value}</strong>
    </div>
  `;
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

function formatCurrency(cents) {
  return `$${(cents / 100).toFixed(2)}`;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

refreshButton.addEventListener("click", () => {
  loadStatus();
  loadPlayerStatus();
});

seedButton.addEventListener("click", seedDevPlayer);

loadStatus();
loadPlayerStatus();

setInterval(loadStatus, 10000);