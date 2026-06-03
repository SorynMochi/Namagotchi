const sectionButtons = document.querySelectorAll("[data-section], [data-section-link]");
const sections = document.querySelectorAll(".content-section");
const seedButton = document.querySelector("#seed-button");
const careStats = document.querySelector("#care-stats");
const namiMessage = document.querySelector("#nami-message");
const serverTime = document.querySelector("#server-time");
const onlineUsers = document.querySelector("#online-users");
const level = document.querySelector("#level");
const syncXp = document.querySelector("#sync-xp");
const credits = document.querySelector("#credits");
const moodScore = document.querySelector("#mood-score");
const namiStatus = document.querySelector("#nami-status");
const moodBonus = document.querySelector("#mood-bonus");
const personalMoodBonus = document.querySelector("#personal-mood-bonus");
const miniMoodFill = document.querySelector("#mini-mood-fill");
const chatForm = document.querySelector("#chat-form");
const chatInput = document.querySelector("#chat-input");
const chatLog = document.querySelector("#chat-log");
const collapseToggles = document.querySelectorAll(".collapse-toggle");

const MAX_CHAT_MESSAGES = 100;

sectionButtons.forEach((button) => {
  button.addEventListener("click", () => {
    const section = button.dataset.section || button.dataset.sectionLink;
    showSection(section);
  });
});

collapseToggles.forEach((button) => {
  button.dataset.label = cleanCollapseLabel(button.textContent);

  button.addEventListener("click", () => {
    const target = document.querySelector(`#${button.dataset.collapse}`);
    if (!target) {
      return;
    }

    target.classList.toggle("collapsed");
    updateCollapseButton(button, target.classList.contains("collapsed"));
  });

  updateCollapseButton(button, false);
});

document.querySelectorAll(".left-rail .panel > .panel-title").forEach((title) => {
  title.dataset.label = cleanCollapseLabel(title.textContent);
  title.setAttribute("role", "button");
  title.setAttribute("tabindex", "0");

  title.addEventListener("click", () => toggleLeftPanel(title));
  title.addEventListener("keydown", (event) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      toggleLeftPanel(title);
    }
  });

  updateLeftPanelTitle(title, false);
});

seedButton?.addEventListener("click", seedDevPlayer);

chatForm.addEventListener("submit", (event) => {
  event.preventDefault();

  const text = chatInput.value.trim();
  if (!text) {
    return;
  }

  addChatMessage("Soryn", text);
  chatInput.value = "";
});

function showSection(sectionName) {
  document.querySelectorAll(".nav-item").forEach((button) => {
    button.classList.toggle("active", button.dataset.section === sectionName);
  });

  sections.forEach((section) => {
    section.classList.toggle("active", section.id === `section-${sectionName}`);
  });
}

async function loadStatus() {
  try {
    const response = await fetch("/api/status");

    if (!response.ok) {
      throw new Error(`Status request failed: ${response.status}`);
    }

    const status = await response.json();

    serverTime.textContent = formatDateTime(status.timestamp);
    onlineUsers.textContent = status.onlineUsers ?? 1;
  } catch (error) {
    console.error(error);
    serverTime.textContent = "Offline";
    onlineUsers.textContent = "?";
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
    namiMessage.textContent = result.message;
    await loadPlayerStatus();
  } catch (error) {
    console.error(error);
    namiMessage.textContent = "Nami-chan tried to create the dev player, but a database goblin bit the clipboard.";
  } finally {
    seedButton.disabled = false;
    seedButton.textContent = "Create Dev Player";
  }
}

async function loadPlayerStatus() {
  try {
    const response = await fetch("/api/player/status");

    if (!response.ok) {
      careStats.innerHTML = `<p class="muted">Create the dev player to load care stats.</p>`;
      return;
    }

    const status = await response.json();
    renderPlayerStatus(status);
  } catch (error) {
    console.error(error);
    careStats.innerHTML = `<p class="muted">Could not load player status.</p>`;
  }
}

function renderPlayerStatus(status) {
  const player = status.player;
  const companion = status.companion;
  const bonus = getMoodBonus(companion.moodScore);

  level.textContent = player.level;
  syncXp.textContent = player.totalXp.toLocaleString();
  credits.textContent = formatCurrency(player.currencyCents);
  moodScore.textContent = companion.moodScore;
  namiStatus.textContent = capitalize(companion.status);
  moodBonus.textContent = `+${bonus}%`;
  personalMoodBonus.textContent = `+${bonus}% Resource Gain`;
  miniMoodFill.style.width = `${Math.max(0, Math.min(100, companion.moodScore))}%`;

  careStats.innerHTML = `
    ${renderStat("Satiety", companion.satiety)}
    ${renderStat("Connection", companion.connection)}
    ${renderStat("Energy", companion.energy)}
    ${renderStat("Comfort", companion.comfort)}
    ${renderStat("Playfulness", companion.playfulness)}
    ${renderStat("Inspiration", companion.inspiration)}
    ${renderStat("Cleanliness", companion.cleanliness)}
  `;

  namiMessage.textContent = "Nami-chan is loaded into the command center and trying very hard not to press every button at once.";
}

function renderStat(name, value) {
  return `
    <div class="stat-row">
      <span>${escapeHTML(name)}</span>
      <div class="bar stat-bar" aria-label="${escapeHTML(name)}: ${value}/100">
        <div class="fill stat-fill" style="width: ${value}%"></div>
      </div>
      <strong>${value}</strong>
    </div>
  `;
}

function addChatMessage(username, text) {
  const message = document.createElement("p");
  message.innerHTML = `<span>[${escapeHTML(username)}]</span> ${escapeHTML(text)}`;
  chatLog.appendChild(message);

  while (chatLog.children.length > MAX_CHAT_MESSAGES) {
    chatLog.removeChild(chatLog.firstElementChild);
  }

  chatLog.scrollTop = chatLog.scrollHeight;
}

function getMoodBonus(mood) {
  return Math.round((Number(mood) / 200) * 100);
}

function formatCurrency(cents) {
  return `$${(Number(cents) / 100).toFixed(2)}`;
}

function formatDateTime(value) {
  const date = new Date(value);

  const year = date.getFullYear();
  const month = pad2(date.getMonth() + 1);
  const day = pad2(date.getDate());
  const hours = pad2(date.getHours());
  const minutes = pad2(date.getMinutes());
  const seconds = pad2(date.getSeconds());

  return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`;
}

function pad2(value) {
  return String(value).padStart(2, "0");
}

function cleanCollapseLabel(value) {
  return String(value).replace(/^(\[[+-]\]\s*)+/, "").trim();
}

function updateCollapseButton(button, isCollapsed) {
  button.textContent = `${isCollapsed ? "[+]" : "[-]"} ${button.dataset.label}`;
}

function toggleLeftPanel(title) {
  const panel = title.closest(".panel");
  if (!panel) {
    return;
  }

  panel.classList.toggle("is-collapsed");
  updateLeftPanelTitle(title, panel.classList.contains("is-collapsed"));
}

function updateLeftPanelTitle(title, isCollapsed) {
  title.textContent = `${isCollapsed ? "[+]" : "[-]"} ${title.dataset.label}`;
}

function capitalize(value) {
  const text = String(value);
  return text.charAt(0).toUpperCase() + text.slice(1);
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

loadStatus();
loadPlayerStatus();

setInterval(loadStatus, 10000);