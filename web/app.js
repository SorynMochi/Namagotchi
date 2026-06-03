const sectionButtons = document.querySelectorAll("[data-section], [data-section-link]");
const sections = document.querySelectorAll(".content-section");

const careStats = document.querySelector("#care-stats");
const namiMessage = document.querySelector("#nami-message");
const serverTime = document.querySelector("#server-time");
const onlineUsers = document.querySelector("#online-users");

const wealthCredits = document.querySelector("#wealth-credits");
const wealthNibbles = document.querySelector("#wealth-nibbles");
const wealthNamiCoin = document.querySelector("#wealth-namicoin");
const wealthInventory = document.querySelector("#wealth-inventory");

const topMood = document.querySelector("#top-mood");
const topNamiStatus = document.querySelector("#top-nami-status");
const topMoodBonus = document.querySelector("#top-mood-bonus");
const personalMoodBonus = document.querySelector("#personal-mood-bonus");

const playdeckTopLevel = document.querySelector("#playdeck-top-level");
const playdeckEquipLevel = document.querySelector("#playdeck-equip-level");
const playdeckIngredients = document.querySelector("#playdeck-ingredients");

const resFans = document.querySelector("#res-fans");
const resMemes = document.querySelector("#res-memes");
const resLostItems = document.querySelector("#res-lost-items");
const resConfidence = document.querySelector("#res-confidence");
const resReceipts = document.querySelector("#res-receipts");
const resPatterns = document.querySelector("#res-patterns");

const actStreaming = document.querySelector("#act-streaming");
const actDoomScrolling = document.querySelector("#act-doom-scrolling");
const actCleaning = document.querySelector("#act-cleaning");
const actExercising = document.querySelector("#act-exercising");
const actShopping = document.querySelector("#act-shopping");
const actDesigning = document.querySelector("#act-designing");

const currentActionLabel = document.querySelector("#current-action-label");
const playdeckHpLabel = document.querySelector("#playdeck-hp-label");
const playdeckHpFill = document.querySelector("#playdeck-hp-fill");
const playdeckXpLabel = document.querySelector("#playdeck-xp-label");
const playdeckXpFill = document.querySelector("#playdeck-xp-fill");

const activeTask = document.querySelector("#active-task");
const resourceRate = document.querySelector("#resource-rate");
const nextTick = document.querySelector("#next-tick");

const chatForm = document.querySelector("#chat-form");
const chatInput = document.querySelector("#chat-input");
const chatLog = document.querySelector("#chat-log");
const collapseToggles = document.querySelectorAll(".collapse-toggle");

const MAX_CHAT_MESSAGES = 100;

let latestPlayerStatus = null;
let forceTickButton = null;

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

document.querySelectorAll(".task-card button").forEach((button) => {
  button.addEventListener("click", () => {
    setGatheringTask(taskFromButtonText(button.textContent));
  });
});

chatForm.addEventListener("submit", (event) => {
  event.preventDefault();

  const text = chatInput.value.trim();
  if (!text) {
    return;
  }

  addChatMessage("Soryn", text);
  chatInput.value = "";
});

createForceTickButton();

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

async function loadPlayerStatus() {
  try {
    const response = await fetch("/api/player/status");

    if (!response.ok) {
      careStats.innerHTML = `<p class="muted">Dev player not found. Visit /api/dev/seed-player once.</p>`;
      return;
    }

    const status = await response.json();
    latestPlayerStatus = status;
    renderPlayerStatus(status);
  } catch (error) {
    console.error(error);
    careStats.innerHTML = `<p class="muted">Could not load player status.</p>`;
  }
}

async function forceTick() {
  if (!forceTickButton) {
    return;
  }

  forceTickButton.disabled = true;
  forceTickButton.textContent = "Processing...";

  try {
    const response = await fetch("/api/dev/force-tick", {
      method: "POST",
    });

    if (!response.ok) {
      throw new Error(`Force tick failed: ${response.status}`);
    }

    const result = await response.json();
    addChatMessage("System", tickResultMessage(result));
    await loadPlayerStatus();
  } catch (error) {
    console.error(error);
    addChatMessage("System", "Force tick failed. The tick goblin dropped its tiny clipboard.");
  } finally {
    forceTickButton.disabled = false;
    forceTickButton.textContent = "Force Tick";
  }
}

async function setGatheringTask(task) {
  try {
    const response = await fetch("/api/player/gathering", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ task }),
    });

    if (!response.ok) {
      throw new Error(`Set gathering task failed: ${response.status}`);
    }

    await loadPlayerStatus();
    addChatMessage("System", `Gathering task changed to ${labelForTask(task)}.`);
  } catch (error) {
    console.error(error);
    addChatMessage("System", "Could not change gathering task.");
  }
}

function renderPlayerStatus(status) {
  const player = status.player;
  const companion = status.companion;
  const tick = status.tick;
  const bonus = getMoodBonus(companion.moodScore);

  wealthCredits.textContent = formatWholeCredits(player.creditsCents ?? player.currencyCents);
wealthNibbles.textContent = Number(player.nibbles ?? 0).toLocaleString();
wealthNamiCoin.textContent = Number(player.namiCoin ?? 0).toLocaleString();
wealthInventory.textContent = "0 / 40";

topMood.textContent = Math.round(Number(companion.moodScore));
topNamiStatus.textContent = capitalize(companion.status);
topMoodBonus.textContent = `+${bonus}%`;
personalMoodBonus.textContent = `+${bonus}% Resource Gain`;

playdeckTopLevel.textContent = Number(player.level).toLocaleString();
playdeckEquipLevel.textContent = "—";
playdeckIngredients.textContent = "0";
playdeckIngredients.title = "Ingredients are not implemented yet.";

resFans.textContent = Number(status.resources.fans ?? 0).toLocaleString();
resMemes.textContent = Number(status.resources.memes ?? 0).toLocaleString();
resLostItems.textContent = Number(status.resources.lostItems ?? 0).toLocaleString();
resConfidence.textContent = Number(status.resources.confidence ?? 0).toLocaleString();
resReceipts.textContent = Number(status.resources.receipts ?? 0).toLocaleString();
resPatterns.textContent = Number(status.resources.patterns ?? 0).toLocaleString();

actStreaming.textContent = "1";
actDoomScrolling.textContent = "1";
actCleaning.textContent = "1";
actExercising.textContent = "1";
actShopping.textContent = "1";
actDesigning.textContent = "1";

  const xpPercent = percent(player.xpIntoLevel, player.xpToNext);
  playdeckXpLabel.textContent = `XP: ${player.xpIntoLevel.toLocaleString()} / ${player.xpToNext.toLocaleString()}`;
  playdeckXpFill.style.width = `${xpPercent}%`;

  playdeckHpLabel.textContent = "HP: 100 / 100";
  playdeckHpFill.style.width = "100%";

  currentActionLabel.textContent = `Playdeck + ${tick.activeGatheringName} [x${tick.playdeckStreak.toLocaleString()}]`;

  careStats.innerHTML = `
    ${renderStat("Satiety", companion.satiety)}
    ${renderStat("Connection", companion.connection)}
    ${renderStat("Energy", companion.energy)}
    ${renderStat("Comfort", companion.comfort)}
    ${renderStat("Playfulness", companion.playfulness)}
    ${renderStat("Inspiration", companion.inspiration)}
    ${renderStat("Cleanliness", companion.cleanliness)}
  `;

  updateGatheringCards(tick.activeGatheringTask);

  namiMessage.textContent = "Nami-chan’s tick engine is running. She is now professionally obligated to be productive every five seconds.";
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

function updateGatheringCards(activeGatheringTask) {
  document.querySelectorAll(".task-card").forEach((card) => {
    const button = card.querySelector("button");
    if (!button) {
      return;
    }

    const task = taskFromButtonText(button.textContent);
    const isActive = task === activeGatheringTask;

    card.classList.toggle("active", isActive);
    const status = card.querySelector("span");
    if (status) {
      status.textContent = isActive ? "Active" : "Idle";
    }
  });
}

function createForceTickButton() {
  const header = document.querySelector("#section-playdeck .section-header");
  if (!header || document.querySelector("#force-tick-button")) {
    return;
  }

  forceTickButton = document.createElement("button");
  forceTickButton.id = "force-tick-button";
  forceTickButton.className = "secondary-button";
  forceTickButton.type = "button";
  forceTickButton.textContent = "Force Tick";
  forceTickButton.addEventListener("click", forceTick);

  header.appendChild(forceTickButton);
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

function tickResultMessage(result) {
  if (!result || !result.ok) {
    return "No tick result received.";
  }

  if (result.ticksProcessed === 0) {
    return result.message || "No ticks ready yet.";
  }

  const levelText = result.levelUps > 0 ? `, ${result.levelUps} level-up(s)` : "";

  return `Processed ${result.ticksProcessed} tick(s): +${result.syncXpGained.toLocaleString()} Sync XP, +${formatCredits(result.creditsCentsGained)} Credits, +${Number(result.nibblesGained).toLocaleString()} Nibbles, +${Number(result.resourceAmountGained).toLocaleString()} ${result.resourceName}${levelText}.`;
}

function getMoodBonus(mood) {
  return Math.round((Number(mood) / 200) * 100);
}

function formatCredits(cents) {
  return (Number(cents) / 100).toLocaleString(undefined, {
    maximumFractionDigits: 2,
  });
}

function formatWholeCredits(cents) {
  return Math.round(Number(cents) / 100).toLocaleString();
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

function taskFromButtonText(text) {
  const normalized = String(text).toLowerCase();

  if (normalized.includes("stream")) {
    return "streaming";
  }

  if (normalized.includes("scroll")) {
    return "doom_scrolling";
  }

  if (normalized.includes("clean")) {
    return "cleaning";
  }

  if (normalized.includes("exercis")) {
    return "exercising";
  }

  if (normalized.includes("shop")) {
    return "shopping";
  }

  if (normalized.includes("design")) {
    return "designing";
  }

  return "streaming";
}

function labelForTask(task) {
  switch (task) {
    case "streaming":
      return "Streaming";
    case "doom_scrolling":
      return "Doom Scrolling";
    case "cleaning":
      return "Cleaning";
    case "exercising":
      return "Exercising";
    case "shopping":
      return "Shopping";
    case "designing":
      return "Designing";
    default:
      return "Streaming";
  }
}

function percent(value, max) {
  if (!max || max <= 0) {
    return 0;
  }

  return Math.max(0, Math.min(100, (Number(value) / Number(max)) * 100));
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
setInterval(loadPlayerStatus, 5000);