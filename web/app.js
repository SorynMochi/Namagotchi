const sectionButtons = document.querySelectorAll("[data-section], [data-section-link]");
const sections = document.querySelectorAll(".content-section");

const careStats = document.querySelector("#care-stats");
const namiMessage = document.querySelector("#nami-message");
const namiMessageLog = document.querySelector("#nami-message-log");
const namiRoomStage = document.querySelector("#nami-room-stage");
const namiRoomBackground = document.querySelector("#nami-room-background");
const namiIdleVideo = document.querySelector("#nami-idle-video");
const namiLevel = document.querySelector("#nami-level");
const namiXpLabel = document.querySelector("#nami-xp-label");
const namiXpFill = document.querySelector("#nami-xp-fill");
const namiMoodLabel = document.querySelector("#nami-mood-label");
const namiPrimaryNeed = document.querySelector("#nami-primary-need");
const namiSuggestedAction = document.querySelector("#nami-suggested-action");
const careButtons = document.querySelectorAll("[data-care-action]");
const sleepToggleButton = document.querySelector("#sleep-toggle-button");
const serverTime = document.querySelector("#server-time");
const onlineUsers = document.querySelector("#online-users");

const wealthCredits = document.querySelector("#wealth-credits");
const wealthNibbles = document.querySelector("#wealth-nibbles");
const wealthNamiCoin = document.querySelector("#wealth-namicoin");
const topWardrobe = document.querySelector("#top-wardrobe");

const topMood = document.querySelector("#top-mood");
const topNamiStatus = document.querySelector("#top-nami-status");
const personalMoodBonus = document.querySelector("#personal-mood-bonus");

const playdeckTopLevel = document.querySelector("#playdeck-top-level");
const playdeckEquipLevel = document.querySelector("#playdeck-equip-level");
const playdeckIngredients = document.querySelector("#playdeck-ingredients");

const combatStatusTitle = document.querySelector("#combat-status-title");
const combatStatusCopy = document.querySelector("#combat-status-copy");
const combatEnemyName = document.querySelector("#combat-enemy-name");
const combatEnemyLevel = document.querySelector("#combat-enemy-level");
const combatEnemyHpLabel = document.querySelector("#combat-enemy-hp-label");
const combatEnemyHpFill = document.querySelector("#combat-enemy-hp-fill");
const combatPlayerHpLabel = document.querySelector("#combat-player-hp-label");
const combatPlayerAttack = document.querySelector("#combat-player-attack");
const combatPlayerDefense = document.querySelector("#combat-player-defense");
const combatWinLoss = document.querySelector("#combat-win-loss");
const combatLogList = document.querySelector("#combat-log-list");
const wardrobeCapacityLabel = document.querySelector("#wardrobe-capacity-label");
const wardrobeInlineCount = document.querySelector("#wardrobe-inline-count");
const inventoryCountLabel = document.querySelector("#inventory-count-label");
const equipmentSlotList = document.querySelector("#equipment-slot-list");
const wardrobeBonusesList = document.querySelector("#wardrobe-bonuses-list");
const inventoryPreviewList = document.querySelector("#inventory-preview-list");
const devSpawnWardrobeItemButton = document.querySelector("#dev-spawn-wardrobe-item-button");

const wardrobeItemModal = document.querySelector("#wardrobe-item-modal");
const wardrobeItemModalClose = document.querySelector("#wardrobe-item-modal-close");
const wardrobeItemModalTitle = document.querySelector("#wardrobe-item-modal-title");
const wardrobeItemModalSlot = document.querySelector("#wardrobe-item-modal-slot");
const wardrobeItemModalMeta = document.querySelector("#wardrobe-item-modal-meta");
const wardrobeItemStatLines = document.querySelector("#wardrobe-item-stat-lines");
const wardrobeCompareTarget = document.querySelector("#wardrobe-compare-target");
const wardrobeComparisonList = document.querySelector("#wardrobe-comparison-list");
const wardrobeAccessoryCompare = document.querySelector("#wardrobe-accessory-compare");
const wardrobeItemActions = document.querySelector("#wardrobe-item-actions");

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
const tickFill = document.querySelector("#tick-fill");
const playdeckHpLabel = document.querySelector("#playdeck-hp-label");
const playdeckHpFill = document.querySelector("#playdeck-hp-fill");
const playdeckXpLabel = document.querySelector("#playdeck-xp-label");
const playdeckXpFill = document.querySelector("#playdeck-xp-fill");

const activeTask = document.querySelector("#active-task");
const resourceRate = document.querySelector("#resource-rate");
const nextTick = document.querySelector("#next-tick");

const chatPanel = document.querySelector("#chat-panel");
const chatResizeHandle = document.querySelector("#chat-resize-handle");
const chatTabs = document.querySelectorAll("[data-chat-channel]");
const chatForm = document.querySelector("#chat-form");
const chatInput = document.querySelector("#chat-input");
const chatLog = document.querySelector("#chat-log");
const emojiButton = document.querySelector("#emoji-button");
const emojiPicker = document.querySelector("#emoji-picker");
const chatToggleButton = document.querySelector("#chat-toggle-button");
const collapseToggles = document.querySelectorAll(".collapse-toggle");
const railToggleButtons = document.querySelectorAll("[data-rail-toggle]");
const gameShell = document.querySelector(".game-shell");
const themeStylesheet = document.querySelector("#theme-stylesheet");
const themeSelect = document.querySelector("#theme-select");
const authLanding = document.querySelector("#auth-landing");
const authPrelandingCard = document.querySelector("#auth-prelanding-card");
const authLoginCard = document.querySelector("#auth-login-card");
const authSparkleLayer = document.querySelector("#auth-sparkle-layer");
const authLandingMessage = document.querySelector("#auth-landing-message");
const googleLoginButton = document.querySelector("#google-login-button");
const authLandingMusic = document.querySelector("#auth-landing-music");
const authMusicToggle = document.querySelector("#auth-music-toggle");
const authMusicToggleIcon = document.querySelector("#auth-music-toggle-icon");
const logoutButton = document.querySelector("#logout-button");

const MAX_CHAT_MESSAGES = 100;
const MAX_NAMI_MESSAGES = 50;
const CHAT_STORAGE_KEY = "namigotchi_chat_store_v1";
const CHAT_CHANNEL_KEY = "namigotchi_chat_active_channel_v1";
const CHAT_HIDDEN_KEY = "namigotchi_chat_hidden_v1";
const CHAT_PREVIOUS_HEIGHT_KEY = "namigotchi_chat_previous_height_v1";
const ACTIVE_SECTION_KEY = "namigotchi_active_section_v1";
const AUTH_LANDING_MUSIC_MUTED_KEY = "namigotchi_auth_landing_music_muted_v1";
const AUTH_LANDING_SKIP_PRELANDING_KEY = "namigotchi_auth_skip_prelanding_once_v1";
let authLandingMusicAutoplayBlocked = false;
let authPrelandingDismissed = false;
let themeBeforeAuthLanding = null;
const EMOJI_USAGE_KEY = "namigotchi_emoji_usage_v2";
const EMOJI_CATEGORY_KEY = "namigotchi_emoji_category_v1";
const RECENT_EMOJI_LIMIT = 70;
const CHAT_IGNORE_KEY = "namigotchi_chat_ignore_list_v1";
const CHAT_LAST_WHISPER_KEY = "namigotchi_chat_last_whisper_v1";
const CHAT_OFFLINE_WHISPERS_KEY = "namigotchi_offline_whispers_v1";

const THEME_STORAGE_KEY = "namigotchi_theme_v1";
const THEME_FILES = {
  "nami-days": "/themes/nami-days.css",
  "nami-default": "/themes/nami-default.css",
  midnight: "/themes/midnight.css",
  strawberry: "/themes/strawberry.css",
  "sakura-dark": "/themes/sakura-dark.css",
  "sakura-light": "/themes/sakura-light.css",
  "tokyo-night": "/themes/tokyo-night.css",
  "candy": "/themes/candy.css",
  "cafe": "/themes/cafe.css",
  "rainy-mood": "/themes/rainy-mood.css",
};

let activeWardrobeModalItemId = 0;
let activeWardrobeModalCompareSlot = "";
let activeWardrobeModalDetail = null;
let activeWardrobeModalReadOnly = false;

const NAMI_ROOM_BACKGROUND_PATHS = [
  "/images/backgrounds/Living_Room_00.webp",
  "/images/backgrounds/Living_Room_01.webp",
  "/images/backgrounds/Living_Room_02.webp",
  "/images/backgrounds/Living_Room_03.webp",
  "/images/backgrounds/Living_Room_04.webp",
  "/images/backgrounds/Living_Room_05.webp",
];

const NAMI_BEDROOM_BACKGROUND_PATHS = [
  "/images/backgrounds/Bedroom_00.webp",
  "/images/backgrounds/Bedroom_01.webp",
  "/images/backgrounds/Bedroom_02.webp",
  "/images/backgrounds/Bedroom_03.webp",
  "/images/backgrounds/Bedroom_04.webp",
  "/images/backgrounds/Bedroom_05.webp",
];

const NAMI_IDLE_VIDEO_SRC = "/images/animations/Nami_Idle_01.webm";
const NAMI_SLEEP_VIDEO_SRC = "/images/animations/Nami_Sleep_01.webm";

const CURRENT_PLAYER_NAME = "Soryn";
const CSRF_COOKIE_NAME = "namigotchi_csrf";
const CSRF_HEADER_NAME = "X-CSRF-Token";
let csrfTokenPromise = null;

function readCookieValue(name) {
  return document.cookie
    .split(";")
    .map((part) => part.trim())
    .find((part) => part.startsWith(name + "="))
    ?.slice(name.length + 1) || "";
}

async function ensureCSRFToken() {
  const existingToken = readCookieValue(CSRF_COOKIE_NAME);
  if (existingToken) {
    return existingToken;
  }

  if (!csrfTokenPromise) {
    csrfTokenPromise = fetch("/api/auth/csrf")
      .then((response) => {
        if (!response.ok) {
          throw new Error("CSRF token request failed: " + response.status);
        }

        return response.json();
      })
      .then((payload) => payload.csrfToken || readCookieValue(CSRF_COOKIE_NAME))
      .finally(() => {
        csrfTokenPromise = null;
      });
  }

  return csrfTokenPromise;
}

async function csrfFetch(input, options = {}) {
  const requestOptions = { ...options };
  const method = String(requestOptions.method || "GET").toUpperCase();
  const headers = new Headers(requestOptions.headers || {});

  if (!["GET", "HEAD", "OPTIONS"].includes(method)) {
    const token = await ensureCSRFToken();
    headers.set(CSRF_HEADER_NAME, token);
  }

  requestOptions.headers = headers;

  return fetch(input, requestOptions);
}
function setTheme(themeKey) {
  const normalizedThemeKey = themeKey === "sakura" ? "sakura-dark" : themeKey;
  const safeThemeKey = Object.hasOwn(THEME_FILES, normalizedThemeKey) ? normalizedThemeKey : "nami-default";

  if (themeStylesheet) {
    themeStylesheet.href = THEME_FILES[safeThemeKey];
  }

  if (themeSelect) {
    themeSelect.value = safeThemeKey;
  }

  localStorage.setItem(THEME_STORAGE_KEY, safeThemeKey);

  window.NamigotchiThemeEffects?.setActiveThemeEffect(safeThemeKey);
}

function initializeTheme() {
  const storedTheme = localStorage.getItem(THEME_STORAGE_KEY) || "nami-default";

  setTheme(storedTheme);

  if (themeSelect) {
    themeSelect.addEventListener("change", () => {
      setTheme(themeSelect.value);
    });
  }
}

const GATHERING_TASK_CONFIG = {
  streaming: {
    name: "Streaming",
    resource: "Fans",
  },
  doom_scrolling: {
    name: "Doom Scrolling",
    resource: "Memes",
  },
  cleaning: {
    name: "Cleaning",
    resource: "Lost Items",
  },
  exercising: {
    name: "Exercise",
    resource: "Confidence",
  },
  shopping: {
    name: "Shopping",
    resource: "Receipts",
  },
  designing: {
    name: "Designing",
    resource: "Patterns",
  },
};

const WARDROBE_EQUIP_SLOT_ORDER = [
  { slotKey: "top", label: "Top", family: "top" },
  { slotKey: "bottom", label: "Bottom", family: "bottom" },
  { slotKey: "dress", label: "Dress", family: "dress" },
  { slotKey: "footwear", label: "Footwear", family: "footwear" },
  { slotKey: "outerwear", label: "Outerwear", family: "outerwear" },
  { slotKey: "necklace", label: "Necklace", family: "necklace" },
  { slotKey: "accessory_1", label: "Accessory 1", family: "accessory" },
  { slotKey: "accessory_2", label: "Accessory 2", family: "accessory" },
];

const WARDROBE_INVENTORY_GROUPS = [
  { key: "top", label: "Top", family: "top" },
  { key: "bottom", label: "Bottom", family: "bottom" },
  { key: "dress", label: "Dress", family: "dress" },
  { key: "footwear", label: "Footwear", family: "footwear" },
  { key: "outerwear", label: "Outerwear", family: "outerwear" },
  { key: "necklace", label: "Necklace", family: "necklace" },
  { key: "accessory", label: "Accessory", family: "accessory" },
];
const WARDROBE_BONUS_DISPLAY_ORDER = [
  "max_health_percent",
  "beauty",
  "attack_percent",
  "targeting_percent",
  "attack_speed_percent",
  "dodge_percent",
  "crit_rate_percent",
  "crit_damage_percent",
  "glamor",
  "recovery",
  "charm",
  "humor",
  "playdeck_xp_percent",
  "work_xp_percent",
  "global_xp_percent",
  "work_resources_percent",
  "ingredient_quality_percent",
  "credit_rate_percent",
  "drop_rate_percent",
];

const WARDROBE_BONUS_DISPLAY_RANK = new Map(
  WARDROBE_BONUS_DISPLAY_ORDER.map((statKey, index) => [statKey, index])
);

const CARE_STAT_DEFINITIONS = [
  { key: "satiety", label: "Satiety" },
  { key: "connection", label: "Connection" },
  { key: "energy", label: "Energy" },
  { key: "comfort", label: "Comfort" },
  { key: "playfulness", label: "Playfulness" },
  { key: "inspiration", label: "Inspiration" },
  { key: "cleanliness", label: "Cleanliness" },
];

const CARE_ACTION_CONFIG = {
  meal: {
    label: "Meal",
    durationSeconds: 30 * 60,
  },
  snack: {
    label: "Snack",
    durationSeconds: 5 * 60,
  },
  drink: {
    label: "Drink",
    durationSeconds: 2 * 60,
  },
  cuddle: {
    label: "Cuddle",
    durationSeconds: 15 * 60,
  },
  play: {
    label: "Play",
    durationSeconds: 25 * 60,
  },
  write_together: {
    label: "Write",
    durationSeconds: 30 * 60,
  },
  read_together: {
    label: "Read",
    durationSeconds: 30 * 60,
  },
  boop: {
    label: "Boop",
    durationSeconds: 30,
  },
  nap: {
    label: "Nap",
    durationSeconds: 60 * 60,
  },
  bath: {
    label: "Bath",
    durationSeconds: 30 * 60,
  },
  freshen_up: {
    label: "Freshen",
    durationSeconds: 10 * 60,
  },
  put_to_bed: {
    label: "Sleep",
    durationSeconds: 60 * 60,
  },
  wake_up: {
    label: "Wake",
    durationSeconds: 5 * 60,
  },
};

const DEV_PLAYER_DIRECTORY = {
  soryn: {
    displayName: "Soryn",
    online: true,
    level: 2,
  },
  "nami-chan": {
    displayName: "Nami-chan",
    online: true,
    level: 1,
  },
  pixelpuff: {
    displayName: "PixelPuff",
    online: true,
    level: 12,
  },
  mochimancer: {
    displayName: "MochiMancer",
    online: false,
    level: 8,
  },
};

const CHAT_CHANNELS = ["lobby", "whispers", "club", "trade", "help", "system"];
const CHAT_LABELS = {
  lobby: "LOBBY",
  whispers: "WHISPERS",
  club: "CLUB",
  trade: "TRADE",
  help: "HELP",
  system: "SYSTEM",
};

const EMOJI_OPTIONS = [
  "\u{1F600}", "\u{1F603}", "\u{1F604}", "\u{1F601}", "\u{1F606}", "\u{1F602}", "\u{1F923}", "\u{1F60A}", "\u{1F607}", "\u{1F642}",
  "\u{1F643}", "\u{1F609}", "\u{1F60C}", "\u{1F60D}", "\u{1F970}", "\u{1F618}", "\u{1F617}", "\u{1F619}", "\u{1F61A}", "\u{1F60B}",
  "\u{1F61B}", "\u{1F61C}", "\u{1F92A}", "\u{1F61D}", "\u{1F911}", "\u{1F917}", "\u{1F92D}", "\u{1FAE2}", "\u{1FAE3}", "\u{1F92B}",
  "\u{1F914}", "\u{1FAE1}", "\u{1F924}", "\u{1F634}", "\u{1F92F}", "\u{1F973}", "\u{1F97A}", "\u{1F62D}", "\u{1F624}", "\u{1F608}",
  "\u{1F47F}", "\u{1F60E}", "\u{1F913}", "\u{1F9D0}", "\u{1F633}", "\u{1F975}", "\u{1F976}", "\u{1F631}", "\u{1F635}", "\u{1F480}",
  "\u{2620}\u{FE0F}", "\u{1F47B}", "\u{1F47D}", "\u{1F916}", "\u{1F63A}", "\u{1F638}", "\u{1F639}", "\u{1F63B}", "\u{1F63C}", "\u{1F640}",
  "\u{1F44D}", "\u{1F44E}", "\u{1F44F}", "\u{1F64C}", "\u{1F450}", "\u{1F932}", "\u{1F64F}", "\u{1F4AA}", "\u{1FAF6}", "\u{1F91D}",
  "\u{1F440}", "\u{1F441}\u{FE0F}", "\u{1F9E0}", "\u{1FAC0}", "\u{1F48B}", "\u{1F485}", "\u{2728}", "\u{1F4AB}", "\u{2B50}", "\u{1F31F}",
  "\u{1F525}", "\u{1F4A5}", "\u{1F4A2}", "\u{1F4A6}", "\u{1F4A8}", "\u{1F56F}\u{FE0F}", "\u{1F380}", "\u{1F381}", "\u{1F389}", "\u{1F38A}",
  "\u{1F496}", "\u{1F497}", "\u{1F493}", "\u{1F495}", "\u{1F49E}", "\u{1F498}", "\u{1F49D}", "\u{1F49C}", "\u{1F499}", "\u{1FA75}",
  "\u{1F49A}", "\u{1F49B}", "\u{1F9E1}", "\u{2764}\u{FE0F}", "\u{1FA77}", "\u{1F5A4}", "\u{1F90D}", "\u{1F90E}", "\u{2615}", "\u{1F375}",
  "\u{1F36A}", "\u{1F369}", "\u{1F370}", "\u{1F9C1}", "\u{1F36B}", "\u{1F36C}", "\u{1F36D}", "\u{1F35C}", "\u{1F363}", "\u{1F359}",
  "\u{1F355}", "\u{1F354}", "\u{1F35F}", "\u{1F950}", "\u{1F95E}", "\u{1F953}", "\u{1F353}", "\u{1F352}", "\u{1F351}", "\u{1F34E}",
  "\u{1F3AE}", "\u{1F579}\u{FE0F}", "\u{1F3B2}", "\u{1F3A7}", "\u{1F3A4}", "\u{1F3AC}", "\u{1F3A8}", "\u{1F9F5}", "\u{1FAA1}", "\u{1F457}",
  "\u{1F45A}", "\u{1F460}", "\u{1F462}", "\u{1F484}", "\u{1F4BB}", "\u{2328}\u{FE0F}", "\u{1F5B1}\u{FE0F}", "\u{1F4F1}", "\u{1F4E6}", "\u{1F48E}",
  "\u{1FA99}", "\u{1F4B0}", "\u{1F511}", "\u{1F512}", "\u{1F513}", "\u{1F9F8}", "\u{1F43E}", "\u{1F431}", "\u{1F408}", "\u{1F98A}",
  "\u{1F430}", "\u{1F439}", "\u{1F43B}", "\u{1F43C}", "\u{1F427}", "\u{1F984}", "\u{1F319}", "\u{2600}\u{FE0F}", "\u{1F308}", "\u{2744}\u{FE0F}",
  "\u{1F338}", "\u{1F33A}", "\u{1F337}", "\u{1F344}", "\u{1FAB4}", "\u{1F451}", "\u{2694}\u{FE0F}", "\u{1F6E1}\u{FE0F}", "\u{1F3AF}", "\u{1F3C6}",
  "\u{1F605}", "\u{1F62C}", "\u{1F972}", "\u{1FAE0}", "\u{1FAE5}", "\u{1F636}", "\u{1F636}\u{200D}\u{1F32B}\u{FE0F}", "\u{1F610}", "\u{1F611}", "\u{1F612}",
  "\u{1F644}", "\u{1F60F}", "\u{1F615}", "\u{1FAE4}", "\u{1F641}", "\u{2639}\u{FE0F}", "\u{1F61F}", "\u{1F614}", "\u{1F61E}", "\u{1F623}",
  "\u{1F616}", "\u{1F62B}", "\u{1F629}", "\u{1F971}", "\u{1F62E}", "\u{1F62F}", "\u{1F632}", "\u{1F626}", "\u{1F627}", "\u{1F628}",
  "\u{1F630}", "\u{1F625}", "\u{1F613}", "\u{1F622}", "\u{1F62A}", "\u{1F62E}\u{200D}\u{1F4A8}", "\u{1F635}\u{200D}\u{1F4AB}", "\u{1F974}", "\u{1F920}", "\u{1F978}",
  "\u{1F925}", "\u{1F928}", "\u{1F910}", "\u{1F922}", "\u{1F92E}", "\u{1F927}", "\u{1F912}", "\u{1F915}", "\u{1F637}", "\u{1F92C}",
  "\u{1F63D}", "\u{1F63F}", "\u{1F63E}", "\u{1F648}", "\u{1F649}", "\u{1F64A}", "\u{1F47E}", "\u{1F479}", "\u{1F47A}", "\u{1F31A}",
  "\u{1F44B}", "\u{1F91A}", "\u{1F590}\u{FE0F}", "\u{270B}", "\u{1F596}", "\u{1F44C}", "\u{1F90C}", "\u{1F90F}", "\u{270C}\u{FE0F}", "\u{1F91E}",
  "\u{1FAF0}", "\u{1F91F}", "\u{1F918}", "\u{1F919}", "\u{1F448}", "\u{1F449}", "\u{1F446}", "\u{1F447}", "\u{261D}\u{FE0F}", "\u{270D}\u{FE0F}",
  "\u{1F44A}", "\u{270A}", "\u{1F91B}", "\u{1F91C}", "\u{1FAF7}", "\u{1FAF8}", "\u{1FAF1}", "\u{1FAF2}", "\u{1FAF3}", "\u{1FAF4}",
  "\u{1F9E9}", "\u{265F}\u{FE0F}", "\u{1F0CF}", "\u{1F004}", "\u{1F3B4}", "\u{1F3AD}", "\u{1F3B7}", "\u{1F3B8}", "\u{1F3B9}", "\u{1F3BA}",
  "\u{1F3BB}", "\u{1F941}", "\u{1FA98}", "\u{1FA87}", "\u{1FA88}", "\u{1F3BC}", "\u{1F3B5}", "\u{1F3B6}", "\u{1F4F7}", "\u{1F4F8}",
  "\u{1F4F9}", "\u{1F3A5}", "\u{1F4FA}", "\u{1F4FB}", "\u{23F0}", "\u{231A}", "\u{1F9ED}", "\u{1F5FA}\u{FE0F}", "\u{1F45C}", "\u{1F45B}",
  "\u{1F45D}", "\u{1F392}", "\u{1F9F3}", "\u{1F453}", "\u{1F576}\u{FE0F}", "\u{1F97D}", "\u{1F9E5}", "\u{1F97C}", "\u{1F97F}", "\u{1F461}",
  "\u{1FA70}", "\u{1F45F}", "\u{1F97E}", "\u{1F9E6}", "\u{1F9E4}", "\u{1F9E3}", "\u{1F3A9}", "\u{1F9E2}", "\u{1F452}", "\u{1F393}",
  "\u{26D1}\u{FE0F}", "\u{1FA96}", "\u{1F48D}", "\u{1F30D}", "\u{1F30E}", "\u{1F30F}", "\u{1FA90}", "\u{1F311}", "\u{1F312}", "\u{1F313}",
  "\u{1F314}", "\u{1F315}", "\u{1F316}", "\u{1F317}", "\u{1F318}", "\u{1F31D}", "\u{1F31E}", "\u{1F320}", "\u{2604}\u{FE0F}", "\u{1F4A7}",
  "\u{1F30A}", "\u{1F32B}\u{FE0F}", "\u{1F32A}\u{FE0F}", "\u{1F327}\u{FE0F}", "\u{26C8}\u{FE0F}", "\u{1F329}\u{FE0F}", "\u{26A1}", "\u{2614}", "\u{2601}\u{FE0F}", "\u{26C4}",
  "\u{2603}\u{FE0F}", "\u{1F32C}\u{FE0F}", "\u{1F324}\u{FE0F}", "\u{1F36E}", "\u{1F382}", "\u{1F36F}", "\u{1F9CB}", "\u{1F964}", "\u{1F9CB}", "\u{1FAD6}",
  "\u{2705}", "\u{2611}\u{FE0F}", "\u{2714}\u{FE0F}", "\u{274C}", "\u{274E}", "\u{2795}", "\u{2796}", "\u{2797}", "\u{2716}\u{FE0F}", "\u{1F514}",
  "\u{1F515}", "\u{1F4E3}", "\u{1F4E2}", "\u{1F4AC}", "\u{1F5EF}\u{FE0F}", "\u{1F4AD}", "\u{1F5E8}\u{FE0F}", "\u{1F50A}", "\u{1F507}", "\u{1F508}",
  "\u{1F509}", "\u{27A1}\u{FE0F}", "\u{2B05}\u{FE0F}", "\u{2B06}\u{FE0F}", "\u{2B07}\u{FE0F}", "\u{2197}\u{FE0F}", "\u{2198}\u{FE0F}", "\u{2199}\u{FE0F}", "\u{2196}\u{FE0F}", "\u{1F504}",
  "\u{1F501}", "\u{1F500}"
];

const EMOJI_OPTION_SET = new Set(EMOJI_OPTIONS);
const EMOJI_INDEX = new Map(EMOJI_OPTIONS.map((emoji, index) => [emoji, index]));

const EMOJI_CATEGORIES = {
  recent: {
    label: "Recent",
    emojis: null,
  },
  faces: {
    label: "Faces",
    emojis: [
      "\u{1F600}", "\u{1F603}", "\u{1F604}", "\u{1F601}", "\u{1F606}", "\u{1F602}", "\u{1F923}", "\u{1F60A}", "\u{1F607}", "\u{1F642}",
      "\u{1F643}", "\u{1F609}", "\u{1F60C}", "\u{1F60D}", "\u{1F970}", "\u{1F618}", "\u{1F617}", "\u{1F619}", "\u{1F61A}", "\u{1F60B}",
      "\u{1F61B}", "\u{1F61C}", "\u{1F92A}", "\u{1F61D}", "\u{1F911}", "\u{1F917}", "\u{1F92D}", "\u{1FAE2}", "\u{1FAE3}", "\u{1F92B}",
      "\u{1F914}", "\u{1FAE1}", "\u{1F924}", "\u{1F634}", "\u{1F92F}", "\u{1F973}", "\u{1F97A}", "\u{1F62D}", "\u{1F624}", "\u{1F608}",
      "\u{1F47F}", "\u{1F60E}", "\u{1F913}", "\u{1F9D0}", "\u{1F633}", "\u{1F975}", "\u{1F976}", "\u{1F631}", "\u{1F635}", "\u{1F480}",
      "\u{2620}\u{FE0F}", "\u{1F47B}", "\u{1F47D}", "\u{1F916}", "\u{1F63A}", "\u{1F638}", "\u{1F639}", "\u{1F63B}", "\u{1F63C}", "\u{1F640}"
    ],
  },
  gestures: {
    label: "Gestures",
    emojis: [
      "\u{1F44D}", "\u{1F44E}", "\u{1F44F}", "\u{1F64C}", "\u{1F450}", "\u{1F932}", "\u{1F64F}", "\u{1F4AA}", "\u{1FAF6}", "\u{1F91D}",
      "\u{1F440}", "\u{1F441}\u{FE0F}", "\u{1F9E0}", "\u{1FAC0}", "\u{1F48B}", "\u{1F485}", "\u{1F44B}", "\u{1F91A}", "\u{1F590}\u{FE0F}", "\u{270B}",
      "\u{1F596}", "\u{1F44C}", "\u{1F90C}", "\u{1F90F}", "\u{270C}\u{FE0F}", "\u{1F91E}", "\u{1FAF0}", "\u{1F91F}", "\u{1F918}", "\u{1F919}",
      "\u{1F448}", "\u{1F449}", "\u{1F446}", "\u{1F447}", "\u{261D}\u{FE0F}", "\u{270D}\u{FE0F}", "\u{1F44A}", "\u{270A}", "\u{1F91B}", "\u{1F91C}",
      "\u{1FAF7}", "\u{1FAF8}", "\u{1FAF1}", "\u{1FAF2}", "\u{1FAF3}", "\u{1FAF4}"
    ],
  },
  hearts: {
    label: "Hearts",
    emojis: [
      "\u{1F496}", "\u{1F497}", "\u{1F493}", "\u{1F495}", "\u{1F49E}", "\u{1F498}", "\u{1F49D}", "\u{1F49C}", "\u{1F499}", "\u{1FA75}",
      "\u{1F49A}", "\u{1F49B}", "\u{1F9E1}", "\u{2764}\u{FE0F}", "\u{1FA77}", "\u{1F5A4}", "\u{1F90D}", "\u{1F90E}", "\u{1F48B}", "\u{1F4AF}",
      "\u{1F4A0}", "\u{1F48C}", "\u{1F4A4}", "\u{1F4AB}", "\u{2728}", "\u{2B50}", "\u{1F31F}", "\u{1F525}", "\u{1F4A5}", "\u{1F4A2}",
      "\u{1F4A6}", "\u{1F4A8}", "\u{1F56F}\u{FE0F}", "\u{1F380}", "\u{1F381}", "\u{1F389}", "\u{1F38A}"
    ],
  },
  food: {
    label: "Food",
    emojis: [
      "\u{2615}", "\u{1F375}", "\u{1F36A}", "\u{1F369}", "\u{1F370}", "\u{1F9C1}", "\u{1F36B}", "\u{1F36C}", "\u{1F36D}", "\u{1F35C}",
      "\u{1F363}", "\u{1F359}", "\u{1F355}", "\u{1F354}", "\u{1F35F}", "\u{1F950}", "\u{1F95E}", "\u{1F953}", "\u{1F353}", "\u{1F352}",
      "\u{1F351}", "\u{1F34E}", "\u{1F36E}", "\u{1F382}", "\u{1F36F}", "\u{1F9CB}", "\u{1F964}", "\u{1FAD6}"
    ],
  },
  gaming: {
    label: "Gaming",
    emojis: [
      "\u{1F3AE}", "\u{1F579}\u{FE0F}", "\u{1F3B2}", "\u{1F3A7}", "\u{1F3A4}", "\u{1F3AC}", "\u{1F3A8}", "\u{1F4BB}", "\u{2328}\u{FE0F}", "\u{1F5B1}\u{FE0F}",
      "\u{1F4F1}", "\u{1F47E}", "\u{1F916}", "\u{2694}\u{FE0F}", "\u{1F6E1}\u{FE0F}", "\u{1F3AF}", "\u{1F3C6}", "\u{1F947}", "\u{1F948}", "\u{1F949}",
      "\u{1F9E9}", "\u{265F}\u{FE0F}", "\u{1F0CF}", "\u{1F004}", "\u{1F3B4}"
    ],
  },
  fashion: {
    label: "Fashion",
    emojis: [
      "\u{1F457}", "\u{1F45A}", "\u{1F460}", "\u{1F462}", "\u{1F484}", "\u{1F451}", "\u{1F45C}", "\u{1F45B}", "\u{1F45D}", "\u{1F392}",
      "\u{1F9F3}", "\u{1F453}", "\u{1F576}\u{FE0F}", "\u{1F97D}", "\u{1F9E5}", "\u{1F97C}", "\u{1F97F}", "\u{1F461}", "\u{1FA70}", "\u{1F45F}",
      "\u{1F97E}", "\u{1F9E6}", "\u{1F9E4}", "\u{1F9E3}", "\u{1F3A9}", "\u{1F9E2}", "\u{1F452}", "\u{1F393}", "\u{26D1}\u{FE0F}", "\u{1FA96}",
      "\u{1F48D}"
    ],
  },
  nature: {
    label: "Nature",
    emojis: [
      "\u{1F43E}", "\u{1F431}", "\u{1F408}", "\u{1F98A}", "\u{1F430}", "\u{1F439}", "\u{1F43B}", "\u{1F43C}", "\u{1F427}", "\u{1F984}",
      "\u{1F319}", "\u{2600}\u{FE0F}", "\u{1F308}", "\u{2744}\u{FE0F}", "\u{1F338}", "\u{1F33A}", "\u{1F337}", "\u{1F344}", "\u{1FAB4}", "\u{1F30D}",
      "\u{1F30E}", "\u{1F30F}", "\u{1FA90}", "\u{1F311}", "\u{1F312}", "\u{1F313}", "\u{1F314}", "\u{1F315}", "\u{1F316}", "\u{1F317}",
      "\u{1F318}", "\u{1F31D}", "\u{1F31E}", "\u{1F320}", "\u{2604}\u{FE0F}", "\u{1F4A7}", "\u{1F30A}", "\u{1F32B}\u{FE0F}", "\u{1F32A}\u{FE0F}", "\u{1F327}\u{FE0F}",
      "\u{26C8}\u{FE0F}", "\u{1F329}\u{FE0F}", "\u{26A1}", "\u{2614}", "\u{2601}\u{FE0F}", "\u{26C4}", "\u{2603}\u{FE0F}", "\u{1F32C}\u{FE0F}", "\u{1F324}\u{FE0F}"
    ],
  },
  objects: {
    label: "Objects",
    emojis: [
      "\u{1F4E6}", "\u{1F48E}", "\u{1FA99}", "\u{1F4B0}", "\u{1F511}", "\u{1F512}", "\u{1F513}", "\u{1F9F8}", "\u{1F4DA}", "\u{1F4D6}",
      "\u{1F516}", "\u{1F4DD}", "\u{270F}\u{FE0F}", "\u{1F58A}\u{FE0F}", "\u{1F58B}\u{FE0F}", "\u{1F58C}\u{FE0F}", "\u{1F58D}\u{FE0F}", "\u{1F4CC}", "\u{1F4CD}", "\u{1F4CE}",
      "\u{1F587}\u{FE0F}", "\u{2702}\u{FE0F}", "\u{1F4CF}", "\u{1F4D0}", "\u{1F4D3}", "\u{1F4D4}", "\u{1F4D5}", "\u{1F4D7}", "\u{1F4D8}", "\u{1F4D9}",
      "\u{1F4F7}", "\u{1F4F8}", "\u{1F4F9}", "\u{1F3A5}", "\u{1F4FA}", "\u{1F4FB}", "\u{23F0}", "\u{231A}", "\u{1F9ED}", "\u{1F5FA}\u{FE0F}"
    ],
  },
  symbols: {
    label: "Symbols",
    emojis: [
      "\u{2705}", "\u{2611}\u{FE0F}", "\u{2714}\u{FE0F}", "\u{274C}", "\u{274E}", "\u{2795}", "\u{2796}", "\u{2797}", "\u{2716}\u{FE0F}", "\u{1F514}",
      "\u{1F515}", "\u{1F4E3}", "\u{1F4E2}", "\u{1F4AC}", "\u{1F5EF}\u{FE0F}", "\u{1F4AD}", "\u{1F5E8}\u{FE0F}", "\u{1F50A}", "\u{1F507}", "\u{1F508}",
      "\u{1F509}", "\u{27A1}\u{FE0F}", "\u{2B05}\u{FE0F}", "\u{2B06}\u{FE0F}", "\u{2B07}\u{FE0F}", "\u{2197}\u{FE0F}", "\u{2198}\u{FE0F}", "\u{2199}\u{FE0F}", "\u{2196}\u{FE0F}", "\u{1F504}",
      "\u{1F501}", "\u{1F500}", "\u{1F4AF}", "\u{1F4A0}", "\u{1F506}", "\u{1F505}"
    ],
  },
};
const PLAYDECK_ZONE_MAX_STREAK_STORAGE_KEY = "namigotchi_playdeck_zone_max_streaks";
const TOP_RAIL_TICK_SECONDS = 5;

function createEmptyChatStore() {
  return {
    lobby: [],
    whispers: [],
    club: [],
    trade: [],
    help: [],
    system: [],
  };
}

function loadEmojiUsage() {
  try {
    const parsed = JSON.parse(localStorage.getItem(EMOJI_USAGE_KEY) || "{}");
    return parsed && typeof parsed === "object" && !Array.isArray(parsed)
      ? parsed
      : {};
  } catch {
    return {};
  }
}

function loadIgnoredPlayers() {
  try {
    const parsed = JSON.parse(localStorage.getItem(CHAT_IGNORE_KEY) || "[]");
    return new Set(Array.isArray(parsed) ? parsed.map((name) => String(name).toLowerCase()) : []);
  } catch {
    return new Set();
  }
}

let latestPlayerStatus = null;
let forceTickButton = null;
let resetPlaydeckStreakButton = null;
let currentChatChannel = "lobby";
let chatMessages = createEmptyChatStore();
let namiMessages = [];
let namiMessageBottomScrollPending = false;
let unreadChannels = new Set();
let emojiUsage = loadEmojiUsage();
let activeEmojiCategory = "recent";
let ignoredPlayers = loadIgnoredPlayers();
let lastWhisperName = localStorage.getItem(CHAT_LAST_WHISPER_KEY) || "";
let activeChatUserMenu = null;
let activeChatProfileModal = null;
let isResizingChat = false;
let chatResizeStartY = 0;
let chatResizeStartHeight = 0;
let isChatHidden = false;
let previousChatHeight = 190;
let serverClockOffsetMs = 0;
let hasServerClock = false;
let tickStartMs = 0;
let tickEndMs = 0;
let playerStatusRefreshTimer = null;
let careCountdownTimer = null;
let careCompletionRefreshInFlight = false;
let emojiPickerNeedsRender = true;
let emojiPickerPreloadTimer = null;
let namiRoomBackgroundTimer = null;
let namiHomeStageMode = "";
let wardrobeModalDragState = {
  active: false,
  pointerId: 0,
  startClientX: 0,
  startClientY: 0,
  startOffsetX: 0,
  startOffsetY: 0,
  offsetX: 0,
  offsetY: 0,
};

function setTextIfChanged(element, value) {
  if (!element) {
    return;
  }

  const nextValue = String(value ?? "");

  if (element.textContent !== nextValue) {
    element.textContent = nextValue;
  }
}
function setTopUserMetricLabel(element, label, value) {
  if (!element) {
    return;
  }

  const nextLabel = String(label ?? "");
  const nextValue = String(value ?? "");

  if (
    element.dataset.metricLabel === nextLabel &&
    element.dataset.metricValue === nextValue
  ) {
    return;
  }

  element.dataset.metricLabel = nextLabel;
  element.dataset.metricValue = nextValue;

  const labelSpan = document.createElement("span");
  labelSpan.className = "metric-name";
  labelSpan.textContent = nextLabel;

  const valueStrong = document.createElement("strong");
  valueStrong.className = "metric-value";
  valueStrong.textContent = nextValue;

  element.replaceChildren(labelSpan, valueStrong);
}

function setTitleIfChanged(element, value) {
  if (!element) {
    return;
  }

  const nextValue = String(value ?? "");

  if (element.title !== nextValue) {
    element.title = nextValue;
  }
}

function setAttributeIfChanged(element, attribute, value) {
  if (!element) {
    return;
  }

  const nextValue = String(value ?? "");

  if (element.getAttribute(attribute) !== nextValue) {
    element.setAttribute(attribute, nextValue);
  }
}

function setWidthIfChanged(element, value) {
  if (!element) {
    return;
  }

  const nextValue = String(value ?? "");

  if (element.style.width !== nextValue) {
    element.style.width = nextValue;
  }
}

function setClassNameIfChanged(element, value) {
  if (!element) {
    return;
  }

  const nextValue = String(value ?? "");

  if (element.className !== nextValue) {
    element.className = nextValue;
  }
}

function toggleClassIfChanged(element, className, enabled) {
  if (!element) {
    return;
  }

  if (element.classList.contains(className) !== enabled) {
    element.classList.toggle(className, enabled);
  }
}

function withUniqueRenderKeys(items, getBaseKey) {
  const seen = new Map();

  return items.map((item, index) => {
    const baseKey = String(getBaseKey(item, index));
    const count = seen.get(baseKey) || 0;

    seen.set(baseKey, count + 1);

    return {
      key: `${baseKey}\u001f${count}`,
      item,
      index,
    };
  });
}

function syncKeyedChildren(container, entries, createElement, updateElement) {
  if (!container) {
    return;
  }

  const wantedKeys = new Set(entries.map((entry) => String(entry.key)));
  const existingByKey = new Map();

  Array.from(container.children).forEach((child) => {
    const key = child.dataset.renderKey;

    if (key && wantedKeys.has(key) && !existingByKey.has(key)) {
      existingByKey.set(key, child);
    }
  });

  let cursor = container.firstElementChild;

  entries.forEach((entry) => {
    const key = String(entry.key);
    let element = existingByKey.get(key);

    if (element) {
      existingByKey.delete(key);
    } else {
      element = createElement(entry);
      element.dataset.renderKey = key;
    }

    updateElement(element, entry);

    if (element === cursor) {
      cursor = cursor.nextElementSibling;
    } else {
      container.insertBefore(element, cursor);
    }
  });

  Array.from(container.children).forEach((child) => {
    if (!wantedKeys.has(child.dataset.renderKey)) {
      child.remove();
    }
  });
}

sectionButtons.forEach((button) => {
  button.addEventListener("click", () => {
    const section = button.dataset.section || button.dataset.sectionLink;
    showSection(section);
  });
});

collapseToggles.forEach((button) => {
  button.dataset.label ||= cleanCollapseLabel(button.textContent)

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

railToggleButtons.forEach((button) => {
  button.addEventListener("click", () => {
    toggleRail(button.dataset.railToggle);
  });

  updateRailToggleButton(button, false);
});

document.querySelectorAll(".left-rail .panel > .panel-title:not(.user-info-title):not(.rail-toggle)").forEach((title) => {
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
    const taskCard = button.closest(".task-card");
    const task = taskCard?.dataset.gatheringTask || taskFromButtonText(button.textContent);
    setGatheringTask(task);
  });
});

careButtons.forEach((button) => {
  button.addEventListener("click", () => {
    performCareAction(button.dataset.careAction);
  });
});

chatTabs.forEach((button) => {
  button.addEventListener("click", () => {
    switchChatChannel(button.dataset.chatChannel);
  });
});

emojiButton.addEventListener("click", toggleEmojiPicker);
chatToggleButton.addEventListener("click", () => setChatHidden(!isChatHidden));
chatForm.addEventListener("submit", submitChatMessage);

document.addEventListener("click", (event) => {
  if (!emojiPicker.contains(event.target) && !emojiButton.contains(event.target)) {
    closeEmojiPicker();
  }

  if (activeChatUserMenu && !activeChatUserMenu.contains(event.target)) {
    closeChatUserMenu();
  }
});

createForceTickButton();
initializeEmojiPickerPortal();
initializeChat();
initializeChatResize();
initializeNamiMessages();
initializeCareButtonTimer();
initializeHomeStage();

document.querySelectorAll(".work-nami-video").forEach((video) => {
  video.addEventListener("contextmenu", (event) => {
    event.preventDefault();
  });

  video.play().catch(() => {
    // Muted autoplay should work in modern browsers, but this keeps failures quiet.
  });
});

function initializeHomeStage() {
  updateHomeStage(latestPlayerStatus);
  scheduleNextHomeRoomBackgroundUpdate();
  initializeNamiIdleVideo();
}

function updateHomeStage(status = latestPlayerStatus) {
  updateHomeRoomBackground(status);
  updateHomeVideo(status);
}

function updateHomeRoomBackground(status = latestPlayerStatus) {
  if (!namiRoomBackground) {
    return;
  }

  const backgroundIndex = getLivingRoomBackgroundIndex();
  const backgroundPaths = shouldUseNamiSleepStage(status)
    ? NAMI_BEDROOM_BACKGROUND_PATHS
    : NAMI_ROOM_BACKGROUND_PATHS;

  const nextSrc = backgroundPaths[backgroundIndex] || backgroundPaths[5];
  const currentPath = new URL(namiRoomBackground.getAttribute("src") || "", window.location.href).pathname;

  if (currentPath !== nextSrc) {
    namiRoomBackground.src = nextSrc;
  }
}

function updateHomeVideo(status = latestPlayerStatus) {
  if (!namiIdleVideo) {
    return;
  }

  const isSleepStage = shouldUseNamiSleepStage(status);
  const nextMode = isSleepStage ? "sleep" : "idle";
  const nextSrc = isSleepStage ? NAMI_SLEEP_VIDEO_SRC : NAMI_IDLE_VIDEO_SRC;
  const currentPath = new URL(namiIdleVideo.getAttribute("src") || "", window.location.href).pathname;

  namiRoomStage?.classList.toggle("is-sleep-stage", isSleepStage);
  namiRoomStage?.classList.toggle("is-idle-stage", !isSleepStage);

  if (namiHomeStageMode === nextMode && currentPath === nextSrc) {
    return;
  }

  namiHomeStageMode = nextMode;

  if (currentPath !== nextSrc) {
    namiIdleVideo.src = nextSrc;
    namiIdleVideo.load();
  }

  namiIdleVideo.play().catch(() => {
    // Muted autoplay should work in modern browsers, but this keeps failures quiet.
  });
}

function shouldUseNamiSleepStage(status = latestPlayerStatus) {
  const companionStatus = String(status?.companion?.status || "").toLowerCase();
  const activeAction = getActiveCareAction(status);
  const activeActionKey = activeAction?.action;
  const activeRemainingSeconds = getCareActionRemainingSeconds(activeAction);

  if (companionStatus === "sleeping") {
    return true;
  }

  return activeRemainingSeconds > 0 && (
    activeActionKey === "nap" ||
    activeActionKey === "put_to_bed" ||
    activeActionKey === "wake_up"
  );
}

function getLivingRoomBackgroundIndex(date = new Date()) {
  const seconds =
    date.getHours() * 60 * 60 +
    date.getMinutes() * 60 +
    date.getSeconds();

  if (seconds >= 5 * 60 * 60 + 30 * 60 && seconds < 9 * 60 * 60) {
    return 0;
  }

  if (seconds >= 9 * 60 * 60 && seconds < 12 * 60 * 60) {
    return 1;
  }

  if (seconds >= 12 * 60 * 60 && seconds < 15 * 60 * 60) {
    return 2;
  }

  if (seconds >= 15 * 60 * 60 && seconds < 18 * 60 * 60) {
    return 3;
  }

  if (seconds >= 18 * 60 * 60 && seconds < 21 * 60 * 60) {
    return 4;
  }

  return 5;
}

function scheduleNextHomeRoomBackgroundUpdate() {
  if (namiRoomBackgroundTimer) {
    clearTimeout(namiRoomBackgroundTimer);
  }

  const now = new Date();
  const nextBoundary = getNextLivingRoomBoundary(now);
  const delay = Math.max(1000, nextBoundary.getTime() - now.getTime() + 250);

  namiRoomBackgroundTimer = setTimeout(() => {
    updateHomeStage(latestPlayerStatus);
    scheduleNextHomeRoomBackgroundUpdate();
  }, delay);
}

function getNextLivingRoomBoundary(date = new Date()) {
  const boundarySeconds = [
    5 * 60 * 60 + 30 * 60,
    9 * 60 * 60,
    12 * 60 * 60,
    15 * 60 * 60,
    18 * 60 * 60,
    21 * 60 * 60,
  ];

  const seconds =
    date.getHours() * 60 * 60 +
    date.getMinutes() * 60 +
    date.getSeconds();

  const nextSecond = boundarySeconds.find((boundary) => boundary > seconds);
  const nextDate = new Date(date);

  nextDate.setHours(0, 0, 0, 0);

  if (nextSecond === undefined) {
    nextDate.setDate(nextDate.getDate() + 1);
    nextDate.setSeconds(boundarySeconds[0]);
    return nextDate;
  }

  nextDate.setSeconds(nextSecond);
  return nextDate;
}

function initializeNamiIdleVideo() {
  if (!namiIdleVideo) {
    return;
  }

  namiIdleVideo.addEventListener("contextmenu", (event) => {
    event.preventDefault();
  });

  namiIdleVideo.play().catch(() => {
    // Muted autoplay should work in modern browsers, but this keeps failures quiet.
  });
}

function showSection(sectionName, options = {}) {
  const targetSection = document.querySelector(`#section-${sectionName}`);
  const safeSectionName = targetSection ? sectionName : "home";

  document.querySelectorAll(".nav-item").forEach((button) => {
    button.classList.toggle("active", button.dataset.section === safeSectionName);
  });

  sections.forEach((section) => {
    section.classList.toggle("active", section.id === `section-${safeSectionName}`);
  });

  if (safeSectionName === "home" && namiMessageBottomScrollPending) {
    scrollNamiMessagesToBottomWhenVisible();
  }

  if (options.save ?? true) {
    localStorage.setItem(ACTIVE_SECTION_KEY, safeSectionName);
  }
}

showSection(localStorage.getItem(ACTIVE_SECTION_KEY) || "home", { save: false });

function initializeAuthSparkles() {
  if (!authSparkleLayer || authSparkleLayer.dataset.ready === "true") {
    return;
  }

  authSparkleLayer.dataset.ready = "true";
  authSparkleLayer.innerHTML = "";

  const sparkleShapes = ["diamond", "plus", "star", "soft"];
  const sparkleColors = ["#fff8fc", "#ffd4ea", "#ff9fd0", "#fff0b6", "#bfefff", "#d8c4ff"];
  const sparkleCount = getAuthSparkleCount(500);

  for (let index = 0; index < sparkleCount; index += 1) {
    const sparkle = document.createElement("span");
    const shape = sparkleShapes[Math.floor(Math.random() * sparkleShapes.length)];
    const color = sparkleColors[Math.floor(Math.random() * sparkleColors.length)];
    const size = randomNumber(6, 15);
    const x = randomNumber(4, 96);
    const y = randomNumber(5, 94);
    const driftX = randomNumber(-34, 34);
    const driftY = randomNumber(-26, 26);
    const duration = randomNumber(9, 18);
    const delay = randomNumber(-18, 0);
    const rotation = randomNumber(-90, 120);

    sparkle.className = `auth-sparkle auth-sparkle--${shape}`;
    sparkle.style.setProperty("--sparkle-size", `${size.toFixed(2)}px`);
    sparkle.style.setProperty("--sparkle-x", `${x.toFixed(2)}%`);
    sparkle.style.setProperty("--sparkle-y", `${y.toFixed(2)}%`);
    sparkle.style.setProperty("--sparkle-drift-x", `${driftX.toFixed(2)}px`);
    sparkle.style.setProperty("--sparkle-drift-y", `${driftY.toFixed(2)}px`);
    sparkle.style.setProperty("--sparkle-duration", `${duration.toFixed(2)}s`);
    sparkle.style.setProperty("--sparkle-delay", `${delay.toFixed(2)}s`);
    sparkle.style.setProperty("--sparkle-rotation", `${rotation.toFixed(2)}deg`);
    sparkle.style.setProperty("--sparkle-color", color);

    authSparkleLayer.appendChild(sparkle);
  }
}

function getAuthSparkleCount(desktopCount) {
  const reducedMotion = window.matchMedia?.("(prefers-reduced-motion: reduce)")?.matches ?? false;
  const compactViewport = window.matchMedia?.("(max-width: 640px), (max-height: 760px)")?.matches ?? false;
  const mediumViewport = window.matchMedia?.("(max-width: 1024px), (max-height: 900px)")?.matches ?? false;
  const coarsePointer = window.matchMedia?.("(pointer: coarse)")?.matches ?? false;
  const deviceMemory = Number(navigator.deviceMemory || 0);

  if (reducedMotion) {
    return Math.min(desktopCount, 16);
  }

  if (compactViewport || coarsePointer) {
    return Math.min(desktopCount, 52);
  }

  if (mediumViewport) {
    return Math.min(desktopCount, 78);
  }

  if (deviceMemory > 0 && deviceMemory <= 4) {
    return Math.min(desktopCount, 64);
  }

  return desktopCount;
}
function randomNumber(min, max) {
  return min + Math.random() * (max - min);
}

function initializeAuthLogoutMusicBridge() {
  if (!logoutButton) {
    return;
  }

  logoutButton.addEventListener("click", prepareAuthLandingAfterLogout, { capture: true });
}

function prepareAuthLandingAfterLogout() {
  sessionStorage.setItem(AUTH_LANDING_SKIP_PRELANDING_KEY, "true");
  authPrelandingDismissed = true;
  document.body.classList.remove("auth-prelanding-active");

  if (!isAuthLandingMusicMuted()) {
    startAuthLandingMusic();
  }
}

function shouldSkipAuthPrelandingOnce() {
  const shouldSkip = sessionStorage.getItem(AUTH_LANDING_SKIP_PRELANDING_KEY) === "true";

  if (shouldSkip) {
    sessionStorage.removeItem(AUTH_LANDING_SKIP_PRELANDING_KEY);
  }

  return shouldSkip;
}

function showAuthLoginCardWithoutPrelanding() {
  authPrelandingDismissed = true;
  document.body.classList.remove("auth-prelanding-active");

  if (authPrelandingCard) {
    authPrelandingCard.setAttribute("aria-hidden", "true");
  }

  if (authLoginCard) {
    authLoginCard.setAttribute("aria-hidden", "false");
  }
}

function initializeAuthPrelanding() {
  if (!authPrelandingCard) {
    return;
  }

  authPrelandingCard.addEventListener("click", dismissAuthPrelanding);
  authPrelandingCard.addEventListener("keydown", handleAuthPrelandingKeydown);
}

function showAuthPrelanding() {
  if (!authPrelandingCard || !authLoginCard) {
    authPrelandingDismissed = true;
    return;
  }

  authPrelandingDismissed = false;
  document.body.classList.add("auth-prelanding-active");
  authPrelandingCard.setAttribute("aria-hidden", "false");
  authLoginCard.setAttribute("aria-hidden", "true");
}

async function dismissAuthPrelanding() {
  if (authPrelandingDismissed) {
    return;
  }

  authPrelandingDismissed = true;
  document.body.classList.remove("auth-prelanding-active");

  if (authPrelandingCard) {
    authPrelandingCard.setAttribute("aria-hidden", "true");
  }

  if (authLoginCard) {
    authLoginCard.setAttribute("aria-hidden", "false");
  }

  await startAuthLandingMusic();
}

function handleAuthPrelandingKeydown(event) {
  if (event.key !== "Enter" && event.key !== " ") {
    return;
  }

  event.preventDefault();
  dismissAuthPrelanding();
}

function initializeAuthMusic() {
  if (!authLandingMusic || !authMusicToggle) {
    return;
  }

  authLandingMusic.loop = true;
  authLandingMusic.volume = 0.36;
  authLandingMusic.muted = isAuthLandingMusicMuted();
  updateAuthMusicToggle();

  authMusicToggle.addEventListener("click", toggleAuthLandingMusic);
}

function isAuthLandingMusicMuted() {
  return localStorage.getItem(AUTH_LANDING_MUSIC_MUTED_KEY) === "true";
}

function setAuthLandingMusicMuted(muted) {
  localStorage.setItem(AUTH_LANDING_MUSIC_MUTED_KEY, muted ? "true" : "false");

  if (authLandingMusic) {
    authLandingMusic.muted = muted;
  }

  updateAuthMusicToggle();
}

async function startAuthLandingMusic() {
  if (!authLandingMusic) {
    return;
  }

  authLandingMusic.volume = 0.36;
  authLandingMusic.loop = true;
  authLandingMusic.muted = isAuthLandingMusicMuted();

  try {
    await authLandingMusic.play();
    authLandingMusicAutoplayBlocked = false;
  } catch {
    authLandingMusicAutoplayBlocked = true;
  }

  updateAuthMusicToggle();
}

function stopAuthLandingMusic() {
  if (!authLandingMusic) {
    return;
  }

  authLandingMusic.pause();
  authLandingMusic.currentTime = 0;
  authLandingMusicAutoplayBlocked = false;
  updateAuthMusicToggle();
}

async function toggleAuthLandingMusic() {
  if (!authLandingMusic) {
    return;
  }

  const muted = isAuthLandingMusicMuted();

  if (muted || authLandingMusic.paused || authLandingMusicAutoplayBlocked) {
    setAuthLandingMusicMuted(false);

    try {
      await authLandingMusic.play();
      authLandingMusicAutoplayBlocked = false;
    } catch {
      authLandingMusicAutoplayBlocked = true;
    }
  } else {
    setAuthLandingMusicMuted(true);
    authLandingMusic.pause();
    authLandingMusicAutoplayBlocked = false;
  }

  updateAuthMusicToggle();
}

function updateAuthMusicToggle() {
  if (!authMusicToggle || !authMusicToggleIcon || !authLandingMusic) {
    return;
  }

  const mutedByPreference = isAuthLandingMusicMuted();
  const waitingForClick = !mutedByPreference && (authLandingMusic.paused || authLandingMusicAutoplayBlocked);
  const label = mutedByPreference
    ? "Play landing music"
    : waitingForClick
      ? "Click to start landing music"
      : "Mute landing music";

  setTextIfChanged(authMusicToggleIcon, mutedByPreference ? "volume_off" : "volume_up");
  authMusicToggle.setAttribute("aria-label", label);
  authMusicToggle.title = label;
  authMusicToggle.classList.toggle("muted", mutedByPreference);
  authMusicToggle.classList.toggle("blocked", waitingForClick);
}

function initializeAuthLanding() {
  initializeAuthLogoutMusicBridge();
  initializeAuthPrelanding();
  initializeAuthMusic();
  initializeAuthSparkles();
  if (!googleLoginButton) {
    return;
  }

  googleLoginButton.addEventListener("click", () => {
    window.location.href = "/api/auth/google/start";
  });
}

async function initializeAuthGate() {
  try {
    const response = await fetch("/api/auth/me", {
      cache: "no-store",
    });

    if (!response.ok) {
      throw new Error(`Auth check failed: ${response.status}`);
    }

    const auth = await response.json();

    if (!auth.loggedIn) {
      showAuthLanding("Player accounts only use Google sign-in.\nNew accounts will choose a display name.");
      return;
    }

    hideAuthLanding();
    await loadPlayerStatus();
  } catch (error) {
    console.error(error);
    showAuthLanding("Could not check your login yet. Try refreshing, or continue with Google.");
  }
}

function showAuthLanding(message) {
  if (themeBeforeAuthLanding === null) {
    themeBeforeAuthLanding = document.body.dataset.theme || "";
  }

  document.body.classList.add("auth-logged-out");
  document.body.dataset.theme = "auth-landing";

  if (shouldSkipAuthPrelandingOnce()) {
    showAuthLoginCardWithoutPrelanding();
    startAuthLandingMusic();
  } else {
    showAuthPrelanding();

    if (authPrelandingDismissed) {
      startAuthLandingMusic();
    }
  }

  if (themeStylesheet) {
    themeStylesheet.disabled = true;
  }

  if (authLanding) {
    authLanding.classList.remove("hidden");
    authLanding.setAttribute("aria-hidden", "false");
  }

  if (authLandingMessage && message) {
    setTextIfChanged(authLandingMessage, message);
  }
}

function hideAuthLanding() {
  document.body.classList.remove("auth-logged-out");

  stopAuthLandingMusic();
  authPrelandingDismissed = false;
  document.body.classList.remove("auth-prelanding-active");

  if (themeStylesheet) {
    themeStylesheet.disabled = false;
  }

  if (themeBeforeAuthLanding !== null) {
    if (themeBeforeAuthLanding) {
      document.body.dataset.theme = themeBeforeAuthLanding;
    } else {
      delete document.body.dataset.theme;
    }

    themeBeforeAuthLanding = null;
  }

  if (authLanding) {
    authLanding.classList.add("hidden");
    authLanding.setAttribute("aria-hidden", "true");
  }
}

async function loadStatus() {
  try {
    const response = await fetch("/api/status");

    if (!response.ok) {
      throw new Error(`Status request failed: ${response.status}`);
    }

    const status = await response.json();

    syncServerClock(status.timestamp);
onlineUsers.textContent = status.onlineUsers ?? 1;
  } catch (error) {
    console.error(error);
    serverTime.textContent = "Offline";
    onlineUsers.textContent = "?";
  }
}

async function loadPlayerStatus() {
  try {
    const response = await csrfFetch("/api/player/sync", {
      method: "POST",
    });

    if (!response.ok) {
      if (response.status === 401) {
        showAuthLanding("Please sign in with Google to keep taking care of Nami-chan.");
        return;
      }

      careStats.innerHTML = `<p class="muted">Could not sync player status. Please log in again.</p>`;
      return;
    }

    const status = await response.json();
    latestPlayerStatus = status;
    renderPlayerStatus(status);
    await loadNamiMessagesFromServer();
  } catch (error) {
    console.error(error);
    careStats.innerHTML = `<p class="muted">Could not load player status.</p>`;
  }
}

async function forceTick() {
  // Dev controls are server-gated at /dev.
}

function getPlaydeckZoneMaxStreaks() {
  try {
    const parsed = JSON.parse(localStorage.getItem(PLAYDECK_ZONE_MAX_STREAK_STORAGE_KEY) || "{}");
    return parsed && typeof parsed === "object" && !Array.isArray(parsed)
      ? parsed
      : {};
  } catch {
    return {};
  }
}

function savePlaydeckZoneMaxStreaks(streaks) {
  localStorage.setItem(PLAYDECK_ZONE_MAX_STREAK_STORAGE_KEY, JSON.stringify(streaks));
}

function getPlaydeckZoneStreakKey(status) {
  const zoneID = status?.tick?.playdeckZoneId ?? status?.playdeck?.zone?.id ?? "unknown";
  return String(zoneID || "unknown");
}

function getPlaydeckCurrentAndMaxStreak(status) {
  const current = Math.max(0, Number(status?.tick?.playdeckStreak ?? 0));
  const zoneKey = getPlaydeckZoneStreakKey(status);
  const maxStreaks = getPlaydeckZoneMaxStreaks();
  const previousMax = Math.max(0, Number(maxStreaks[zoneKey] ?? 0));
  const nextMax = Math.max(previousMax, current);

  if (nextMax !== previousMax) {
    maxStreaks[zoneKey] = nextMax;
    savePlaydeckZoneMaxStreaks(maxStreaks);
  }

  return {
    current,
    max: nextMax,
  };
}

function updateTopPlayerName(player) {
  const playerPill = document.querySelector(".top-player-name-button");
  if (!playerPill) {
    return;
  }

  const displayName = String(player?.displayName || CURRENT_PLAYER_NAME || "Player").trim() || "Player";
  setTextIfChanged(playerPill, `\u{1F48E} ${displayName}`);
}
function syncTopRailTickPill(tick) {
  const playerPill = document.querySelector(".top-player-name-button");
  if (!playerPill) {
    return;
  }

  const secondsUntilNextTick = Math.max(
    0,
    Math.min(TOP_RAIL_TICK_SECONDS, Number(tick?.secondsUntilNextTick ?? TOP_RAIL_TICK_SECONDS))
  );

  const progressPercent = Math.max(
    0,
    Math.min(100, ((TOP_RAIL_TICK_SECONDS - secondsUntilNextTick) / TOP_RAIL_TICK_SECONDS) * 100)
  );

  playerPill.classList.add("top-tick-pill");
  playerPill.style.setProperty("--top-tick-progress", `${progressPercent}%`);
  playerPill.style.setProperty("--top-tick-duration", `${Math.max(0.1, secondsUntilNextTick)}s`);
  playerPill.title = `Next tick in ${Math.ceil(secondsUntilNextTick).toLocaleString()}s`;

  playerPill.classList.remove("top-tick-pill-animating");
  void playerPill.offsetWidth;
  playerPill.classList.add("top-tick-pill-animating");
}

function syncTopPlayerTickPill(tick) {
  const playerPill = document.querySelector(".top-player-name-button");

  if (!playerPill) {
    return;
  }

  const nextStartMs = Date.parse(tick?.lastTickAt);
  const nextEndMs = Date.parse(tick?.nextTickAt);

  if (!Number.isNaN(nextStartMs) && !Number.isNaN(nextEndMs) && nextEndMs > nextStartMs) {
    tickStartMs = nextStartMs;
    tickEndMs = nextEndMs;
  }

  playerPill.classList.add("top-player-tick-pill");
  playerPill.classList.remove("top-player-tick-pill-animating");

  updateTickProgressBar();
}
function renderPlayerStatus(status) {
  const player = status.player;
  const companion = status.companion;
  const tick = status.tick;
  const bonus = getMoodBonus(companion.moodScore);
  const namiXpPercent = percent(companion.xpIntoLevel, companion.xpToNext);

  updateTopPlayerName(player);

  setTextIfChanged(namiLevel, Number(companion.level ?? 1).toLocaleString());
  setTextIfChanged(
    namiXpLabel,
    `${Number(companion.xpIntoLevel ?? 0).toLocaleString()} / ${Number(companion.xpToNext ?? 120).toLocaleString()}`
  );
  setWidthIfChanged(namiXpFill, `${namiXpPercent}%`);
  setTextIfChanged(namiMoodLabel, companion.moodLabel || "Okay");
  setTextIfChanged(namiPrimaryNeed, companion.primaryNeed || "Waiting");
  setTextIfChanged(namiSuggestedAction, companion.suggestedAction || "Any care action");

  setTextIfChanged(wealthCredits, formatWholeCredits(player.creditsCents ?? player.currencyCents));
  setTextIfChanged(wealthNibbles, formatCompactNumber(player.nibbles ?? 0));
  setTextIfChanged(wealthNamiCoin, formatCompactNumber(player.namiCoin ?? 0));

  const wardrobe = status.wardrobe || { used: 0, capacity: 100 };
  const wardrobeCountText = `${Number(wardrobe.used ?? 0).toLocaleString()} / ${Number(wardrobe.capacity ?? 100).toLocaleString()}`;

  setTextIfChanged(topWardrobe, wardrobeCountText);
  setTextIfChanged(topMood, Math.round(Number(companion.moodScore ?? 0)));
  setTextIfChanged(topNamiStatus, capitalize(companion.status));
  setTextIfChanged(personalMoodBonus, `+${bonus}% Resource Gain`);

  setTextIfChanged(playdeckTopLevel, Number(player.level ?? 1).toLocaleString());
  setTextIfChanged(playdeckEquipLevel, formatCompactNumber(status.playdeck?.equipmentPower ?? 0));
  setTextIfChanged(playdeckIngredients, formatCompactNumber(0));
  setTitleIfChanged(playdeckIngredients, "Ingredients are not implemented yet.");

  setTextIfChanged(resFans, formatCompactNumber(status.resources?.fans ?? 0));
  setTextIfChanged(resMemes, formatCompactNumber(status.resources?.memes ?? 0));
  setTextIfChanged(resLostItems, formatCompactNumber(status.resources?.lostItems ?? 0));
  setTextIfChanged(resConfidence, formatCompactNumber(status.resources?.confidence ?? 0));
  setTextIfChanged(resReceipts, formatCompactNumber(status.resources?.receipts ?? 0));
  setTextIfChanged(resPatterns, formatCompactNumber(status.resources?.patterns ?? 0));

  setTextIfChanged(actStreaming, activityLevel(status.activities?.streaming));
  setTextIfChanged(actDoomScrolling, activityLevel(status.activities?.doomScrolling));
  setTextIfChanged(actCleaning, activityLevel(status.activities?.cleaning));
  setTextIfChanged(actExercising, activityLevel(status.activities?.exercising));
  setTextIfChanged(actShopping, activityLevel(status.activities?.shopping));
  setTextIfChanged(actDesigning, activityLevel(status.activities?.designing));

  const xpPercent = percent(player.xpIntoLevel, player.xpToNext);
  setTopUserMetricLabel(
    playdeckXpLabel,
    "XP:",
    `${Number(player.xpIntoLevel ?? 0).toLocaleString()} / ${Number(player.xpToNext ?? 720).toLocaleString()}`
  );
  setWidthIfChanged(playdeckXpFill, `${xpPercent}%`);

  setTopUserMetricLabel(playdeckHpLabel, "Sparkles:", "100 / 100");
  setWidthIfChanged(playdeckHpFill, "100%");
  const currentStreak = Math.max(0, Number(tick.playdeckStreak ?? 0));
  const maxStreak = Math.max(currentStreak, Number(tick.playdeckMaxStreak ?? currentStreak));
  const streakPercent = percent(currentStreak, maxStreak);

  setTopUserMetricLabel(
    currentActionLabel,
    "Streak:",
    `${currentStreak.toLocaleString()} / ${maxStreak.toLocaleString()}`
  );
  setWidthIfChanged(tickFill, `${streakPercent}%`);
  syncTopPlayerTickPill(tick);
  scheduleNextPlayerStatusRefresh(tick);

  renderCareStats(companion);
  updateGatheringCards(status);
  renderCareButtons(status);
  renderPlaydeckStatus(status);
  renderWardrobeStatus(status);
  updateHomeStage(status);

  setTextIfChanged(namiMessage, companion.caption || "Nami-chan is waiting sweetly.");
}

function activityLevel(activity) {
  return Number(activity?.level ?? 1).toLocaleString();
}

function normalizeCareStatValue(value) {
  const numberValue = Number(value ?? 0);

  if (!Number.isFinite(numberValue)) {
    return 0;
  }

  return Math.max(0, Math.min(100, numberValue));
}

function formatCareStatValue(value) {
  const normalized = normalizeCareStatValue(value);

  return Number.isInteger(normalized)
    ? normalized.toLocaleString()
    : normalized.toFixed(1);
}

function renderCareStats(companion) {
  if (!careStats) {
    return;
  }

  const entries = CARE_STAT_DEFINITIONS.map((definition) => {
    const value = normalizeCareStatValue(companion?.[definition.key]);

    return {
      key: definition.key,
      label: definition.label,
      value,
      displayValue: formatCareStatValue(value),
    };
  });

  syncKeyedChildren(careStats, entries, createCareStatRow, updateCareStatRow);
}

function createCareStatRow() {
  const row = document.createElement("div");
  row.className = "stat-row";

  const label = document.createElement("span");
  label.className = "care-stat-name";

  const bar = document.createElement("div");
  bar.className = "bar stat-bar";

  const fill = document.createElement("div");
  fill.className = "fill stat-fill";

  const value = document.createElement("strong");
  value.className = "care-stat-value";

  bar.appendChild(fill);
  row.append(label, bar, value);

  return row;
}

function updateCareStatRow(row, entry) {
  const label = row.querySelector(".care-stat-name");
  const bar = row.querySelector(".stat-bar");
  const fill = row.querySelector(".stat-fill");
  const value = row.querySelector(".care-stat-value");

  setTextIfChanged(label, entry.label);
  setAttributeIfChanged(bar, "aria-label", `${entry.label}: ${entry.displayValue}/100`);
  setWidthIfChanged(fill, `${entry.value}%`);
  setTextIfChanged(value, entry.displayValue);
}

function renderPlaydeckStatus(status) {
  const playdeck = status.playdeck || {};
  const enemy = playdeck.enemy || {};

  const playerHp = Number(playdeck.playerHp ?? 100);
  const playerMaxHp = Number(playdeck.playerMaxHp ?? 100);
  const enemyHp = Number(enemy.hp ?? 0);
  const enemyMaxHp = Number(enemy.maxHp ?? 0);
  const timeoutTicks = Number(playdeck.timeoutTicks ?? 0);

  if (combatStatusTitle) {
    combatStatusTitle.textContent = timeoutTicks > 0
      ? `> PLAYDECK RECOVERING`
      : `> ${String(playdeck.lastOutcome || "ready").toUpperCase()}`;
  }

  if (combatStatusCopy) {
    combatStatusCopy.textContent = timeoutTicks > 0
      ? `Recovering for ${timeoutTicks.toLocaleString()} tick(s). Nami-chan has applied a tiny digital bandage.`
      : `${playdeck.zone?.name || "Starter Deck"} is ready. Combat state is now backend-backed.`;
  }

  if (combatEnemyName) {
    combatEnemyName.textContent = enemy.name || "No enemy loaded";
  }

  if (combatEnemyLevel) {
    combatEnemyLevel.textContent = `Lv${Number(enemy.level ?? 1).toLocaleString()}`;
  }

  if (combatEnemyHpLabel) {
    combatEnemyHpLabel.textContent = `${enemyHp.toLocaleString()} / ${enemyMaxHp.toLocaleString()}`;
  }

  if (combatEnemyHpFill) {
    combatEnemyHpFill.style.width = `${percent(enemyHp, enemyMaxHp)}%`;
  }

  if (combatPlayerHpLabel) {
    combatPlayerHpLabel.textContent = `${playerHp.toLocaleString()} / ${playerMaxHp.toLocaleString()}`;
  }

  if (combatPlayerAttack) {
    combatPlayerAttack.textContent = Number(playdeck.attack ?? 0).toLocaleString();
  }

  if (combatPlayerDefense) {
    combatPlayerDefense.textContent = Number(playdeck.defense ?? 0).toLocaleString();
  }

  if (combatWinLoss) {
    combatWinLoss.textContent = `${Number(playdeck.wins ?? 0).toLocaleString()} / ${Number(playdeck.losses ?? 0).toLocaleString()}`;
  }

  renderCombatLog(playdeck.combatLog || []);
}

function renderCombatLog(logs) {
  if (!combatLogList) {
    return;
  }

  const entries = Array.isArray(logs) && logs.length
    ? withUniqueRenderKeys(logs, combatLogEntryKey)
    : [{ key: "__empty__", empty: true }];

  syncKeyedChildren(combatLogList, entries, createCombatLogRow, updateCombatLogRow);
}

function combatLogEntryKey(entry) {
  return [
    entry.outcome,
    entry.enemyName,
    entry.enemyLevel,
    entry.playerDamage,
    entry.enemyDamage,
    entry.xpGained,
    entry.creditsCentsGained,
    entry.nibblesGained,
    entry.itemName,
    entry.itemQuantity,
  ].join("\u001f");
}

function createCombatLogRow(entry) {
  const row = document.createElement("p");

  if (entry.empty) {
    const tag = document.createElement("span");
    tag.className = "combat-log-tag";

    const text = document.createElement("span");
    text.className = "combat-log-empty";

    row.append(tag, " ", text);
    return row;
  }

  row.className = "combat-log-row";

  const tag = document.createElement("span");
  tag.className = "combat-log-tag";

  const outcome = document.createElement("span");
  outcome.className = "combat-log-outcome";

  const enemy = document.createElement("span");
  enemy.className = "combat-log-enemy";

  const damage = document.createElement("em");
  damage.className = "combat-log-damage";

  const xp = document.createElement("span");
  xp.className = "combat-log-xp";

  const credits = document.createElement("span");
  credits.className = "combat-log-credits";

  const nibbles = document.createElement("span");
  nibbles.className = "combat-log-nibbles";

  row.append(tag, " ", outcome, " ", enemy, " ", damage, " ", xp, " ", credits, " ", nibbles);

  return row;
}

function updateCombatLogRow(row, entry) {
  if (entry.empty) {
    setClassNameIfChanged(row, "");
    setTextIfChanged(row.querySelector(".combat-log-tag"), "[000]");
    setTextIfChanged(row.querySelector(".combat-log-empty"), "[LOG] No combat logs yet.");
    return;
  }

  const log = entry.item;
  const tag = String(entry.index + 1).padStart(3, "0");
  const itemText = log.itemName
    ? ` [${log.itemName}${Number(log.itemQuantity ?? 0) > 1 ? ` x${Number(log.itemQuantity).toLocaleString()}` : ""}]`
    : "";

  setClassNameIfChanged(row, "combat-log-row");
  setTextIfChanged(row.querySelector(".combat-log-tag"), `[${tag}]`);
  setTextIfChanged(row.querySelector(".combat-log-outcome"), `[${formatOutcome(log.outcome)}]`);
  setTextIfChanged(
    row.querySelector(".combat-log-enemy"),
    `${log.enemyName || "Enemy"} Lv${Number(log.enemyLevel ?? 1).toLocaleString()}`
  );
  setTextIfChanged(
    row.querySelector(".combat-log-damage"),
    `-${Number(log.playerDamage ?? 0).toLocaleString()} / -${Number(log.enemyDamage ?? 0).toLocaleString()}`
  );
  setTextIfChanged(row.querySelector(".combat-log-xp"), `+${Number(log.xpGained ?? 0).toLocaleString()} XP`);
  setTextIfChanged(row.querySelector(".combat-log-credits"), `+${formatCredits(log.creditsCentsGained ?? 0)} Credits`);
  setTextIfChanged(
    row.querySelector(".combat-log-nibbles"),
    `[+${Number(log.nibblesGained ?? 0).toLocaleString()} Nibbles]${itemText}`
  );
}

function renderWardrobeStatus(status) {
  const wardrobe = status.wardrobe || { used: 0, capacity: 100 };
  const playdeck = status.playdeck || {};
  const equipment = Array.isArray(playdeck.equipment) ? playdeck.equipment : [];
  const rawInventoryItems = Array.isArray(playdeck.inventoryPreview) ? playdeck.inventoryPreview : [];
  const items = rawInventoryItems.filter((item) => !String(item?.equippedSlot || "").trim());

  const equipmentByKey = new Map(
    equipment.map((slot) => [String(slot.slotKey || "").toLowerCase(), slot])
  );

  const filledSlots = WARDROBE_EQUIP_SLOT_ORDER.filter((definition) => {
    const slot = equipmentByKey.get(definition.slotKey);
    return Number(slot?.itemId ?? 0) > 0;
  }).length;

  const wardrobeCountText = `${Number(wardrobe.used ?? 0).toLocaleString()} / ${Number(wardrobe.capacity ?? 100).toLocaleString()}`;

  setTextIfChanged(topWardrobe, wardrobeCountText);
  setTextIfChanged(wardrobeInlineCount, wardrobeCountText);
  setTextIfChanged(wardrobeCapacityLabel, wardrobeCountText);
  setTextIfChanged(inventoryCountLabel, `${items.length.toLocaleString()} shown, ${Number(wardrobe.used ?? 0).toLocaleString()} total`);

  if (equipmentSlotList) {
    const slotEntries = WARDROBE_EQUIP_SLOT_ORDER.map((definition) => {
      const slot = equipmentByKey.get(definition.slotKey) || {
        slotKey: definition.slotKey,
        displayName: definition.label,
        acceptsSlot: definition.family,
        itemId: 0,
        itemName: "",
        rarity: "basic",
        powerLevel: 0,
      };

      return {
        key: definition.slotKey,
        definition,
        slot,
      };
    });

    syncKeyedChildren(equipmentSlotList, slotEntries, createEquipmentSlotCard, updateEquipmentSlotCard);
  }

  if (wardrobeBonusesList) {
    const bonusRows = buildEquippedWardrobeBonusRows(equipment);
    syncKeyedChildren(wardrobeBonusesList, bonusRows, createWardrobeBonusRow, updateWardrobeBonusRow);
  }

  if (inventoryPreviewList) {
    const groupedItems = groupWardrobeInventoryItems(items);
    const groupEntries = WARDROBE_INVENTORY_GROUPS.map((group) => ({
      key: group.key,
      group,
      items: groupedItems[group.key] || [],
    }));

    syncKeyedChildren(inventoryPreviewList, groupEntries, createInventoryGroupSection, updateInventoryGroupSection);
  }
}

function createEquipmentSlotCard() {
  const card = document.createElement("article");

  card.innerHTML = `
    <div class="gear-card-topline">
      <span class="equipment-slot-label"></span>
      <span class="gear-level-tag"></span>
    </div>

    <div class="gear-card-name-wrap">
      <strong class="gear-card-name"></strong>
    </div>

    <div class="gear-card-bottomline">
      <span class="gear-tailoring-tag"></span>
    </div>
  `;

  return card;
}

function updateEquipmentSlotCard(card, entry) {
  const { definition, slot } = entry;
  const hasItem = Number(slot.itemId ?? 0) > 0;
  const itemID = Number(slot.itemId ?? 0);
  const rarity = normalizeWardrobeRarity(slot.rarity || "basic");

  setClassNameIfChanged(card, `equipment-card ${hasItem ? wardrobeRarityClass(rarity) : "is-empty"}`);
  toggleClassIfChanged(card, "wardrobe-item-clickable", hasItem);

if (hasItem) {
  card.dataset.wardrobeItemId = String(itemID);
  card.dataset.wardrobeCompareSlot = definition.slotKey;
  card.setAttribute("role", "button");
  card.tabIndex = 0;
  setTitleIfChanged(card, `View ${slot.itemName || "item"} details`);
} else {
  delete card.dataset.wardrobeItemId;
  delete card.dataset.wardrobeCompareSlot;
  card.removeAttribute("role");
  card.removeAttribute("tabindex");
  card.removeAttribute("title");
}
  setTextIfChanged(card.querySelector(".equipment-slot-label"), definition.label);

  const levelTag = card.querySelector(".gear-level-tag");
  setTextIfChanged(levelTag, hasItem ? Number(slot.powerLevel ?? 0).toLocaleString() : "");
  toggleClassIfChanged(levelTag, "is-hidden", !hasItem);

  const name = card.querySelector(".gear-card-name");
  setClassNameIfChanged(name, hasItem ? "gear-card-name" : "gear-card-name gear-card-empty-name");
  setTextIfChanged(name, hasItem ? slot.itemName || "Unknown Item" : "EMPTY");

  const tailoring = card.querySelector(".gear-tailoring-tag");
  setTextIfChanged(tailoring, hasItem ? `T: ${formatTailoringPoints(slot)}` : "");
  toggleClassIfChanged(tailoring, "is-hidden", !hasItem);
}



function wardrobeBonusDisplayRank(statKey) {
  return WARDROBE_BONUS_DISPLAY_RANK.get(String(statKey || "")) ?? 9999;
}

function buildEquippedWardrobeBonusRows(equipment) {
  const totals = new Map();

  (Array.isArray(equipment) ? equipment : []).forEach((slot) => {
    if (Number(slot?.itemId ?? 0) <= 0 || !Array.isArray(slot?.statLines)) {
      return;
    }

    slot.statLines.forEach((line) => {
      const statKey = String(line?.statKey || "").trim();

      if (!statKey) {
        return;
      }

      const value = Number(line?.value ?? 0);

      if (!Number.isFinite(value) || value === 0) {
        return;
      }

      const current = totals.get(statKey) || {
        key: statKey,
        label: line.displayName || statKey,
        valueKind: line.valueKind || "flat",
        rawValue: 0,
        tooltip: line.tooltip || "",
        sortOrder: Number(line.sortOrder ?? 9999),
      };

      current.rawValue += value;
      current.sortOrder = Math.min(current.sortOrder, Number(line.sortOrder ?? 9999));
      totals.set(statKey, current);
    });
  });

  const rows = [...totals.values()]
    .filter((row) => Math.abs(Number(row.rawValue ?? 0)) > 0.0001)
    .sort((a, b) => {
      const rankDelta = wardrobeBonusDisplayRank(a.key) - wardrobeBonusDisplayRank(b.key);

      if (rankDelta !== 0) {
        return rankDelta;
      }

      return String(a.label || "").localeCompare(String(b.label || ""));
    })
    .map((row) => ({
      key: row.key,
      label: row.label,
      value: formatWardrobeStatValue(row.rawValue, row.valueKind),
      tooltip: row.tooltip,
    }));

  if (rows.length === 0) {
    return [{
      key: "__empty__",
      empty: true,
      label: "No equipped gear bonuses yet.",
      value: "",
      tooltip: "",
    }];
  }

  return rows;
}

function createWardrobeBonusRow() {
  const row = document.createElement("div");
  row.className = "wardrobe-bonus-row";

  const label = document.createElement("span");
  label.className = "wardrobe-bonus-label";

  const value = document.createElement("strong");
  value.className = "wardrobe-bonus-value";

  row.append(label, value);

  return row;
}

function updateWardrobeBonusRow(row, entry) {
  setTextIfChanged(row.querySelector(".wardrobe-bonus-label"), entry.label);
  setTextIfChanged(row.querySelector(".wardrobe-bonus-value"), entry.value);
}

function createInventoryGroupSection() {
  const section = document.createElement("section");
  section.className = "inventory-group-section";

  const header = document.createElement("div");
  header.className = "inventory-group-header";

  const label = document.createElement("span");
  label.className = "inventory-group-label";

  const count = document.createElement("strong");
  count.className = "inventory-group-count";

  const grid = document.createElement("div");
  grid.className = "inventory-group-grid";

  header.append(label, count);
  section.append(header, grid);

  return section;
}

function updateInventoryGroupSection(section, entry) {
  const grid = section.querySelector(".inventory-group-grid");

  setTextIfChanged(section.querySelector(".inventory-group-label"), entry.group.label);
  setTextIfChanged(section.querySelector(".inventory-group-count"), `(${entry.items.length.toLocaleString()})`);

  const itemEntries = entry.items.length
    ? withUniqueRenderKeys(entry.items, wardrobeInventoryItemKey)
    : [{ key: "__empty__", empty: true }];

  syncKeyedChildren(grid, itemEntries, createInventoryItemElement, updateInventoryItemElement);
}

function wardrobeInventoryItemKey(item) {
  const id = Number(item.id ?? item.itemId ?? item.inventoryId ?? 0);

  if (id > 0) {
    return `item:${id}`;
  }

  return [
    item.name,
    item.equipmentSlot,
    item.itemType,
    item.rarity,
    item.powerLevel,
  ].join("\u001f");
}

function createInventoryItemElement(entry) {
  if (entry.empty) {
    const empty = document.createElement("div");
    empty.className = "inventory-empty-row";
    return empty;
  }

  const card = document.createElement("article");

  card.innerHTML = `
    <div class="gear-card-topline">
      <span class="gear-card-filler"></span>
      <span class="gear-level-tag"></span>
    </div>

    <div class="gear-card-name-wrap">
      <strong class="gear-card-name"></strong>
    </div>

    <div class="gear-card-bottomline">
      <span class="gear-tailoring-tag"></span>
    </div>
  `;

  return card;
}

function updateInventoryItemElement(element, entry) {
  if (entry.empty) {
    setClassNameIfChanged(element, "inventory-empty-row");
    setTextIfChanged(element, "Empty");
    return;
  }

  const item = entry.item;
  const itemID = Number(item.id ?? item.itemId ?? item.inventoryId ?? 0);
  const rarity = normalizeWardrobeRarity(item.rarity || "basic");

  setClassNameIfChanged(element, `inventory-item-card ${wardrobeRarityClass(rarity)}`);
  toggleClassIfChanged(element, "wardrobe-item-clickable", itemID > 0);

if (itemID > 0) {
  element.dataset.wardrobeItemId = String(itemID);
  element.dataset.wardrobeCompareSlot = defaultCompareSlotForWardrobeItem(item);
  element.setAttribute("role", "button");
  element.tabIndex = 0;
  setTitleIfChanged(element, `View ${item.name || "item"} details`);
} else {
  delete element.dataset.wardrobeItemId;
  delete element.dataset.wardrobeCompareSlot;
  element.removeAttribute("role");
  element.removeAttribute("tabindex");
  element.removeAttribute("title");
}
  setTextIfChanged(element.querySelector(".gear-level-tag"), Number(item.powerLevel ?? 1).toLocaleString());
  setTextIfChanged(element.querySelector(".gear-card-name"), item.name || "Unknown Item");
  setTextIfChanged(element.querySelector(".gear-tailoring-tag"), `T: ${formatTailoringPoints(item)}`);
}

function groupWardrobeInventoryItems(items) {
  const groups = {};

  WARDROBE_INVENTORY_GROUPS.forEach((group) => {
    groups[group.key] = [];
  });

  items.forEach((item) => {
    const key = normalizeWardrobeInventoryGroupKey(item.equipmentSlot || item.itemType || "");

    if (groups[key]) {
      groups[key].push(item);
    }
  });

  Object.keys(groups).forEach((key) => {
    groups[key].sort((a, b) => {
      const rarityDelta = wardrobeRarityWeight(b.rarity) - wardrobeRarityWeight(a.rarity);
      if (rarityDelta !== 0) {
        return rarityDelta;
      }

      const powerDelta = Number(b.powerLevel ?? 0) - Number(a.powerLevel ?? 0);
      if (powerDelta !== 0) {
        return powerDelta;
      }

      return String(a.name || "").localeCompare(String(b.name || ""));
    });
  });

  return groups;
}

function normalizeWardrobeInventoryGroupKey(value) {
  const key = String(value || "").trim().toLowerCase();

  switch (key) {
    case "top":
    case "bottom":
    case "dress":
    case "footwear":
    case "outerwear":
    case "necklace":
      return key;
    case "accessory":
    case "accessory_1":
    case "accessory_2":
      return "accessory";
    default:
      return "";
  }
}

function wardrobeRarityWeight(rarity) {
  switch (normalizeWardrobeRarity(rarity)) {
    case "devastating":
      return 7;
    case "iconic":
      return 6;
    case "glam":
      return 5;
    case "trendy":
      return 4;
    case "chic":
      return 3;
    case "cute":
      return 2;
    case "basic":
      return 1;
    default:
      return 0;
  }
}

function formatTailoringPoints(source) {
  const current = Number(
    source?.tailoringCurrent ??
    source?.tailoringPointsCurrent ??
    source?.tailoringPoints ??
    0
  );

  const max = Number(
    source?.tailoringMax ??
    source?.tailoringPointsMax ??
    0
  );

  return `${current.toLocaleString()}/${max.toLocaleString()}`;
}

function normalizeWardrobeRarity(value) {
  return String(value || "basic").trim().toLowerCase();
}

function wardrobeRarityClass(rarity) {
  return `rarity-${normalizeWardrobeRarity(rarity).replace(/[^a-z0-9]+/g, "-")}`;
}

function formatWardrobeRarity(rarity) {
  const value = normalizeWardrobeRarity(rarity);
  return value.charAt(0).toUpperCase() + value.slice(1);
}

function formatWardrobeSlot(value) {
  const key = String(value || "").trim().toLowerCase();

  switch (key) {
    case "top":
      return "Top";
    case "bottom":
      return "Bottom";
    case "dress":
      return "Dress";
    case "footwear":
      return "Footwear";
    case "outerwear":
      return "Outerwear";
    case "necklace":
      return "Necklace";
    case "bag":
      return "Bag";
    case "accessory":
      return "Accessory";
    case "accessory_1":
      return "Accessory 1";
    case "accessory_2":
      return "Accessory 2";
    default:
      return capitalize(String(value || "Item").replace(/_/g, " "));
  }
}

function defaultCompareSlotForWardrobeItem(item) {
  const slot = normalizeWardrobeInventoryGroupKey(item?.equipmentSlot || item?.itemType || "");

  if (slot === "accessory") {
    return "accessory_1";
  }

  return slot;
}

function updateWardrobeScrollState() {
  const wardrobeSection = document.querySelector("#section-inventory");
  const equippedScroll = document.querySelector("#section-inventory .wardrobe-equipped-scroll");

  if (!wardrobeSection || !equippedScroll) {
    return;
  }

  const wardrobeIsActive = wardrobeSection.classList.contains("active");
  document.body.classList.toggle("is-wardrobe-active", wardrobeIsActive);

  if (!wardrobeIsActive) {
    equippedScroll.classList.remove("can-scroll");
    return;
  }

  requestAnimationFrame(() => {
    const hasOverflow = equippedScroll.scrollHeight > equippedScroll.clientHeight + 2;
    equippedScroll.classList.toggle("can-scroll", hasOverflow);
  });
}

window.addEventListener("resize", updateWardrobeScrollState);

function initializeWardrobeItemModal() {
  [equipmentSlotList, inventoryPreviewList].forEach((container) => {
    if (!container) {
      return;
    }

    container.addEventListener("click", handleWardrobeItemClick);
    container.addEventListener("keydown", handleWardrobeItemKeydown);
  });

  wardrobeItemModalClose?.addEventListener("click", closeWardrobeItemModal);
  initializeWardrobeModalDragging();

  wardrobeItemModal?.addEventListener("click", (event) => {
    if (event.target === wardrobeItemModal) {
      closeWardrobeItemModal();
    }
  });

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && wardrobeItemModal && !wardrobeItemModal.classList.contains("hidden")) {
      closeWardrobeItemModal();
    }
  });
}

function getWardrobeModalCard() {
  return wardrobeItemModal?.querySelector(".wardrobe-item-modal-card") || null;
}

function setWardrobeModalDragOffset(x, y) {
  const card = getWardrobeModalCard();

  wardrobeModalDragState.offsetX = Number(x) || 0;
  wardrobeModalDragState.offsetY = Number(y) || 0;

  if (!card) {
    return;
  }

  card.style.setProperty("--wardrobe-modal-drag-x", `${wardrobeModalDragState.offsetX}px`);
  card.style.setProperty("--wardrobe-modal-drag-y", `${wardrobeModalDragState.offsetY}px`);
}

function resetWardrobeModalDragPosition() {
  wardrobeModalDragState.active = false;
  wardrobeModalDragState.pointerId = 0;
  document.body.classList.remove("wardrobe-modal-dragging");
  setWardrobeModalDragOffset(0, 0);
}

function clampWardrobeModalDragOffset(x, y) {
  const card = getWardrobeModalCard();

  if (!card) {
    return { x, y };
  }

  const rect = card.getBoundingClientRect();
  const currentX = wardrobeModalDragState.offsetX;
  const currentY = wardrobeModalDragState.offsetY;
  const baseLeft = rect.left - currentX;
  const baseTop = rect.top - currentY;
  const margin = 8;

  const minX = margin - baseLeft;
  const maxX = window.innerWidth - margin - baseLeft - rect.width;
  const minY = margin - baseTop;
  const maxY = window.innerHeight - margin - baseTop - rect.height;

  const clamp = (value, min, max) => {
    if (max < min) {
      return Math.max(max, Math.min(min, value));
    }

    return Math.max(min, Math.min(max, value));
  };

  return {
    x: clamp(x, minX, maxX),
    y: clamp(y, minY, maxY),
  };
}

function initializeWardrobeModalDragging() {
  const card = getWardrobeModalCard();

  if (!wardrobeItemModal || !card || card.dataset.dragReady === "true") {
    return;
  }

  card.dataset.dragReady = "true";

  const dragHandle = card.querySelector(".wardrobe-item-modal-header") || card;

  dragHandle.addEventListener("pointerdown", (event) => {
    if (event.button !== 0) {
      return;
    }

    if (event.target.closest("button, a, input, select, textarea, [contenteditable='true']")) {
      return;
    }

    wardrobeModalDragState.active = true;
    wardrobeModalDragState.pointerId = event.pointerId;
    wardrobeModalDragState.startClientX = event.clientX;
    wardrobeModalDragState.startClientY = event.clientY;
    wardrobeModalDragState.startOffsetX = wardrobeModalDragState.offsetX;
    wardrobeModalDragState.startOffsetY = wardrobeModalDragState.offsetY;

    dragHandle.setPointerCapture?.(event.pointerId);
    document.body.classList.add("wardrobe-modal-dragging");
    event.preventDefault();
  });

  dragHandle.addEventListener("pointermove", (event) => {
    if (!wardrobeModalDragState.active || wardrobeModalDragState.pointerId !== event.pointerId) {
      return;
    }

    const nextX =
      wardrobeModalDragState.startOffsetX +
      event.clientX -
      wardrobeModalDragState.startClientX;

    const nextY =
      wardrobeModalDragState.startOffsetY +
      event.clientY -
      wardrobeModalDragState.startClientY;

    const clamped = clampWardrobeModalDragOffset(nextX, nextY);
    setWardrobeModalDragOffset(clamped.x, clamped.y);
  });

  const stopDragging = (event) => {
    if (!wardrobeModalDragState.active || wardrobeModalDragState.pointerId !== event.pointerId) {
      return;
    }

    wardrobeModalDragState.active = false;
    wardrobeModalDragState.pointerId = 0;
    document.body.classList.remove("wardrobe-modal-dragging");
  };

  dragHandle.addEventListener("pointerup", stopDragging);
  dragHandle.addEventListener("pointercancel", stopDragging);

  window.addEventListener("resize", () => {
    if (!wardrobeItemModal || wardrobeItemModal.classList.contains("hidden")) {
      return;
    }

    const clamped = clampWardrobeModalDragOffset(
      wardrobeModalDragState.offsetX,
      wardrobeModalDragState.offsetY
    );

    setWardrobeModalDragOffset(clamped.x, clamped.y);
  });
}
function initializeDevWardrobeSpawner() {
  // Dev controls are server-gated at /dev.
}

function handleWardrobeItemClick(event) {
  const target = event.target.closest(".wardrobe-item-clickable");

  if (!target) {
    return;
  }

  openWardrobeItemDetail(target.dataset.wardrobeItemId, target.dataset.wardrobeCompareSlot || "");
}

function handleWardrobeItemKeydown(event) {
  if (event.key !== "Enter" && event.key !== " ") {
    return;
  }

  const target = event.target.closest(".wardrobe-item-clickable");
  if (!target) {
    return;
  }

  event.preventDefault();
  openWardrobeItemDetail(target.dataset.wardrobeItemId, target.dataset.wardrobeCompareSlot || "");
}

async function fetchWardrobeItemDetailForCompareV4(itemID, compareSlot = "") {
  const params = new URLSearchParams({ id: String(itemID) });
  if (compareSlot) { params.set("compareSlot", compareSlot); }
  const response = await fetch("/api/player/wardrobe/item?" + params.toString());
  if (!response.ok) { throw new Error("Item detail failed: " + response.status); }
  return response.json();
}
function getWardrobeEquippedSlotForCompareV4(item) {
  return String(item?.equippedSlot || "").trim().toLowerCase();
}
function isWardrobeAccessoryItemForCompareV4(item) {
  return normalizeWardrobeInventoryGroupKey(item?.equipmentSlot || item?.itemType || "") === "accessory";
}
function getOtherWardrobeAccessorySlotForCompareV4(slotKey) {
  const key = String(slotKey || "").trim().toLowerCase();
  if (key === "accessory_1") { return "accessory_2"; }
  if (key === "accessory_2") { return "accessory_1"; }
  return "";
}
function hideWardrobeComparisonForCompareV4(detail) {
  return { ...detail, hideComparison: true, accessoryCompareSlots: [] };
}
async function resolveWardrobeEquippedItemComparisonV4(detail) {
  const item = detail?.item || {};
  const equippedSlot = getWardrobeEquippedSlotForCompareV4(item);
  if (!equippedSlot) { return { ...detail, hideComparison: false }; }
  if (!isWardrobeAccessoryItemForCompareV4(item)) { return hideWardrobeComparisonForCompareV4(detail); }
  const otherSlot = getOtherWardrobeAccessorySlotForCompareV4(equippedSlot);
  if (!otherSlot) { return hideWardrobeComparisonForCompareV4(detail); }
  const itemID = Number(item.id ?? item.itemId ?? activeWardrobeModalItemId ?? 0);
  const currentCompareSlot = String(detail?.compareSlot || "").trim().toLowerCase();
  let resolvedDetail = detail;
  if (itemID > 0 && currentCompareSlot !== otherSlot) {
    resolvedDetail = await fetchWardrobeItemDetailForCompareV4(itemID, otherSlot);
  }
  if (!resolvedDetail?.compareItem) { return hideWardrobeComparisonForCompareV4(resolvedDetail); }
  const filteredAccessorySlots = (Array.isArray(resolvedDetail?.accessoryCompareSlots) ? resolvedDetail.accessoryCompareSlots : []).filter((slot) => String(slot.slotKey || "").trim().toLowerCase() === otherSlot);
  return { ...resolvedDetail, hideComparison: false, accessoryCompareSlots: filteredAccessorySlots };
}

function getWardrobeEquippedSlotsForCompareV5() {
  return Array.isArray(latestPlayerStatus?.playdeck?.equipment)
    ? latestPlayerStatus.playdeck.equipment
    : [];
}
function getWardrobeEquippedSlotForItemV5(item) {
  const apiSlot = String(item?.equippedSlot || "").trim().toLowerCase();
  if (apiSlot) { return apiSlot; }
  const itemID = Number(item?.id ?? item?.itemId ?? activeWardrobeModalItemId ?? 0);
  if (itemID <= 0) { return ""; }
  const equippedSlot = getWardrobeEquippedSlotsForCompareV5().find((slot) => Number(slot?.itemId ?? 0) === itemID);
  return String(equippedSlot?.slotKey || "").trim().toLowerCase();
}
function getWardrobeEquippedItemIdInSlotV5(slotKey) {
  const key = String(slotKey || "").trim().toLowerCase();
  const equippedSlot = getWardrobeEquippedSlotsForCompareV5().find((slot) => String(slot?.slotKey || "").trim().toLowerCase() === key);
  return Number(equippedSlot?.itemId ?? 0);
}
function isWardrobeAccessoryForCompareV5(item, equippedSlot) {
  const itemFamily = normalizeWardrobeInventoryGroupKey(item?.equipmentSlot || item?.itemType || "");
  const slotKey = String(equippedSlot || "").trim().toLowerCase();
  return itemFamily === "accessory" || slotKey === "accessory_1" || slotKey === "accessory_2";
}
function getOtherWardrobeAccessorySlotForCompareV5(slotKey) {
  const key = String(slotKey || "").trim().toLowerCase();
  if (key === "accessory_1") { return "accessory_2"; }
  if (key === "accessory_2") { return "accessory_1"; }
  return "";
}
function hideWardrobeComparisonForCompareV5(detail, equippedSlot = "") {
  return {
    ...detail,
    item: { ...(detail?.item || {}), equippedSlot },
    hideComparison: true,
    accessoryCompareSlots: [],
  };
}
async function resolveWardrobeEquippedItemComparisonV5(detail) {
  const item = detail?.item || {};
  const equippedSlot = getWardrobeEquippedSlotForItemV5(item);
  if (!equippedSlot) { return { ...detail, hideComparison: false }; }
  if (!isWardrobeAccessoryForCompareV5(item, equippedSlot)) {
    return hideWardrobeComparisonForCompareV5(detail, equippedSlot);
  }
  const otherSlot = getOtherWardrobeAccessorySlotForCompareV5(equippedSlot);
  if (!otherSlot || getWardrobeEquippedItemIdInSlotV5(otherSlot) <= 0) {
    return hideWardrobeComparisonForCompareV5(detail, equippedSlot);
  }
  const itemID = Number(item.id ?? item.itemId ?? activeWardrobeModalItemId ?? 0);
  const currentCompareSlot = String(detail?.compareSlot || "").trim().toLowerCase();
  let resolvedDetail = detail;
  if (itemID > 0 && currentCompareSlot !== otherSlot) {
    resolvedDetail = await fetchWardrobeItemDetailForCompareV4(itemID, otherSlot);
  }
  if (!resolvedDetail?.compareItem) {
    return hideWardrobeComparisonForCompareV5(resolvedDetail, equippedSlot);
  }
  const filteredAccessorySlots = (Array.isArray(resolvedDetail?.accessoryCompareSlots) ? resolvedDetail.accessoryCompareSlots : []).filter((slot) => String(slot.slotKey || "").trim().toLowerCase() === otherSlot);
  const accessorySlots = filteredAccessorySlots.length > 0
    ? filteredAccessorySlots
    : [{ slotKey: otherSlot, displayName: formatWardrobeSlot(otherSlot), selected: true }];
  return {
    ...resolvedDetail,
    item: { ...(resolvedDetail?.item || {}), equippedSlot },
    hideComparison: false,
    accessoryCompareSlots: accessorySlots,
  };
}
function syncWardrobeCompareSectionVisibilityV5(detail) {
  const compareSection = wardrobeComparisonList?.closest(".wardrobe-compare-section") || null;
  const shouldHideComparison = Boolean(detail?.hideComparison);
  if (compareSection) {
    compareSection.classList.toggle("hidden", shouldHideComparison);
    compareSection.toggleAttribute("hidden", shouldHideComparison);
    compareSection.style.display = shouldHideComparison ? "none" : "";
  }
  if (!shouldHideComparison) { return; }
  wardrobeAccessoryCompare?.classList.add("hidden");
  setTextIfChanged(wardrobeCompareTarget, "");
  wardrobeComparisonList?.replaceChildren();
}

async function openWardrobeItemDetail(itemID, compareSlot = "") {
  const safeItemID = Number(itemID);

  if (!safeItemID || safeItemID < 1) {
    return;
  }

  activeWardrobeModalItemId = safeItemID;

  try {
    const params = new URLSearchParams({
      id: String(safeItemID),
    });

    if (compareSlot) {
      params.set("compareSlot", compareSlot);
    }

    const response = await fetch(`/api/player/wardrobe/item?${params.toString()}`);
    if (!response.ok) {
      throw new Error(`Item detail failed: ${response.status}`);
    }

    const detail = await response.json();
    const resolvedDetail = await resolveWardrobeEquippedItemComparisonV5(detail);
    renderWardrobeItemModal(resolvedDetail);
  } catch (error) {
    console.error(error);
    addChatMessage("System", "Could not load item details. The wardrobe gremlin misplaced the tag.", "system");
  }
}

function closeWardrobeItemModal() {
  if (!wardrobeItemModal) {
    return;
  }

  wardrobeItemModal.classList.add("hidden");
  wardrobeItemModal.setAttribute("aria-hidden", "true");
  activeWardrobeModalItemId = 0;
  resetWardrobeModalDragPosition();
}

function syncWardrobeCompareSectionVisibilityV4(detail) {
  const compareSection = wardrobeComparisonList?.closest(".wardrobe-compare-section") || null;
  const shouldHideComparison = Boolean(detail?.hideComparison);
  compareSection?.classList.toggle("hidden", shouldHideComparison);
  compareSection?.toggleAttribute("hidden", shouldHideComparison);
  if (!shouldHideComparison) { return; }
  wardrobeAccessoryCompare?.classList.add("hidden");
  setTextIfChanged(wardrobeCompareTarget, "");
  wardrobeComparisonList?.replaceChildren();
}

function renderWardrobeItemModal(detail) {
  if (!wardrobeItemModal) {
    return;
  }

  const item = detail?.item || {};
  activeWardrobeModalItemId = Number(item.id ?? activeWardrobeModalItemId ?? 0);
  activeWardrobeModalCompareSlot = detail?.compareSlot || defaultCompareSlotForWardrobeItem(item);
  const rarity = formatWardrobeRarity(item.rarity || "basic");
  const slot = formatWardrobeSlot(item.equipmentSlot || "item");
  const level = Number(item.powerLevel ?? 1).toLocaleString();

  setTextIfChanged(wardrobeItemModalTitle, item.name || "Unknown Item");
  setTextIfChanged(wardrobeItemModalSlot, `${slot} | ${rarity}`);
  setTextIfChanged(
    wardrobeItemModalMeta,
    `Item Level ${level} | T: ${formatTailoringPoints(item)}`
  );

  renderWardrobeAccessoryCompare(detail);
  renderWardrobeItemActions(detail);
  renderWardrobeStatLines(detail?.item?.statLines || []);
  renderWardrobeComparison(detail);
  syncWardrobeCompareSectionVisibilityV5(detail);

  wardrobeItemModal.classList.remove("hidden");
  wardrobeItemModal.setAttribute("aria-hidden", "false");
}

function renderWardrobeAccessoryCompare(detail) {
  if (!wardrobeAccessoryCompare) {
    return;
  }

  wardrobeAccessoryCompare.replaceChildren();

  if (detail?.hideComparison) {
    wardrobeAccessoryCompare.classList.add("hidden");
    return;
  }

  const slots = Array.isArray(detail?.accessoryCompareSlots) ? detail.accessoryCompareSlots : [];
  wardrobeAccessoryCompare.classList.toggle("hidden", slots.length === 0);

  slots.forEach((slot) => {
    const button = document.createElement("button");
    button.type = "button";
    button.className = `wardrobe-accessory-compare-button${slot.selected ? " active" : ""}`;
    button.textContent = slot.displayName;
    button.addEventListener("click", () => {
      openWardrobeItemDetail(activeWardrobeModalItemId, slot.slotKey);
    });

    wardrobeAccessoryCompare.appendChild(button);
  });
}


function renderWardrobeItemActions(detail, options = {}) {
  if (!wardrobeItemActions) {
    return;
  }

  wardrobeItemActions.replaceChildren();
  wardrobeItemActions.classList.toggle("hidden", Boolean(options.readOnly));

  if (options.readOnly) {
    return;
  }

  const item = detail?.item || {};
  const itemID = Number(item.id ?? item.itemId ?? 0);
  const itemSlot = normalizeWardrobeInventoryGroupKey(item.equipmentSlot || item.itemType || "");
  const equippedSlot = String(item.equippedSlot || "").trim().toLowerCase();

  const primaryActions = document.createElement("div");
  primaryActions.className = "wardrobe-item-action-row wardrobe-item-primary-actions";

  const utilityActions = document.createElement("div");
  utilityActions.className = "wardrobe-item-action-row wardrobe-item-utility-actions";

  if (!itemID || !itemSlot) {
    const empty = document.createElement("p");
    empty.className = "muted";
    empty.textContent = "This item cannot be worn.";
    wardrobeItemActions.appendChild(empty);
    return;
  }

  if (itemSlot === "accessory") {
    [
      { slotKey: "accessory_1", label: "A1" },
      { slotKey: "accessory_2", label: "A2" },
    ].forEach((slot) => {
      const isWearingHere = equippedSlot === slot.slotKey;
      primaryActions.appendChild(
        createWardrobeActionButton(
          isWearingHere ? `Wearing ${slot.label}` : `Wear ${slot.label}`,
          () => equipWardrobeItem(itemID, slot.slotKey),
          isWearingHere,
          `${isWearingHere ? "Already worn in" : "Wear in"} ${formatWardrobeSlot(slot.slotKey)}`
        )
      );
    });
  } else {
    const targetSlot = detail?.compareSlot || defaultCompareSlotForWardrobeItem(item);
    const isWearingHere = equippedSlot === targetSlot;

    primaryActions.appendChild(
      createWardrobeActionButton(
        isWearingHere ? "Wearing" : "Wear",
        () => equipWardrobeItem(itemID, targetSlot),
        isWearingHere,
        isWearingHere ? "Nami-Chan is already wearing this." : `Wear in ${formatWardrobeSlot(targetSlot)}`
      )
    );
  }

  if (equippedSlot) {
    primaryActions.appendChild(
      createWardrobeActionButton(
        "Take Off",
        () => unequipWardrobeItem(itemID, equippedSlot),
        false,
        `Remove from ${formatWardrobeSlot(equippedSlot)}`,
        "danger"
      )
    );
  }

  ["Recycle", "Sell", "Beautify"].forEach((label) => {
    utilityActions.appendChild(
      createWardrobeActionButton(
        label,
        null,
        true,
        `${label} will be added later.`,
        "utility"
      )
    );
  });

  wardrobeItemActions.append(primaryActions, utilityActions);
}

function appendFormattedChatText(parent, text) {
  const source = String(text);
  const pattern = /(\[\[item:\d+\|[a-z0-9_-]+\|[a-z0-9_-]*\|[^\]]+\]\]|\*\*\*[^*]+\*\*\*|\*\*[^*]+\*\*|\*[^*]+\*|@[A-Za-z0-9_-]{2,32})/g;

  let lastIndex = 0;

  source.replace(pattern, (match, _ignored, offset) => {
    appendPlainChatText(parent, source.slice(lastIndex, offset));

    if (match.startsWith("[[item:")) {
      const tokenInfo = parseChatItemToken(match);

      if (tokenInfo) {
        parent.appendChild(createChatItemLink(tokenInfo));
      } else {
        appendPlainChatText(parent, match);
      }
    } else if (match.startsWith("@")) {
      appendMention(parent, match);
    } else if (match.startsWith("***") && match.endsWith("***")) {
      const strong = document.createElement("strong");
      const em = document.createElement("em");
      em.textContent = match.slice(3, -3);
      strong.appendChild(em);
      parent.appendChild(strong);
    } else if (match.startsWith("**") && match.endsWith("**")) {
      const strong = document.createElement("strong");
      strong.textContent = match.slice(2, -2);
      parent.appendChild(strong);
    } else if (match.startsWith("*") && match.endsWith("*")) {
      const em = document.createElement("em");
      em.textContent = match.slice(1, -1);
      parent.appendChild(em);
    }

    lastIndex = offset + match.length;
    return match;
  });

  appendPlainChatText(parent, source.slice(lastIndex));
}

function shareActiveWardrobeItemToChat() {
  const item = activeWardrobeModalDetail?.item;

  if (!item || !Number(item.id ?? 0)) {
    return;
  }

  insertChatItemChip(item, activeWardrobeModalCompareSlot || defaultCompareSlotForWardrobeItem(item));
}

function createChatItemToken(item, compareSlot = "") {
  const id = Number(item?.id ?? item?.itemId ?? 0);
  const rarity = normalizeWardrobeRarity(item?.rarity || "basic");
  const slot = String(compareSlot || "").trim().toLowerCase();
  const name = encodeURIComponent(String(item?.name || "Unknown Item"));

  return `[[item:${id}|${rarity}|${slot}|${name}]]`;
}

function parseChatItemToken(token) {
  const match = String(token || "").match(/^\[\[item:(\d+)\|([a-z0-9_-]+)\|([a-z0-9_-]*)\|([^\]]+)\]\]$/);

  if (!match) {
    return null;
  }

  return {
    itemId: Number(match[1]),
    rarity: normalizeWardrobeRarity(match[2]),
    compareSlot: match[3] || "",
    name: decodeURIComponent(match[4] || "Unknown Item"),
  };
}

function insertChatItemChip(item, compareSlot = "") {
  if (!chatInput) {
    return;
  }

  const chip = createChatItemInputChip(item, compareSlot);

  if (!chatInput.isContentEditable) {
    insertIntoChatInput(`[${item?.name || "Unknown Item"}] `);
    return;
  }

  chatInput.appendChild(chip);
  chatInput.appendChild(document.createTextNode(" "));
  chatInput.focus();
  setChatCursorToEnd();
}

function createChatItemInputChip(item, compareSlot = "") {
  const chip = document.createElement("span");
  chip.className = `chat-item-token rarity-${normalizeWardrobeRarity(item?.rarity || "basic")}`;
  chip.contentEditable = "false";
  chip.dataset.itemId = String(Number(item?.id ?? item?.itemId ?? 0));
  chip.dataset.itemName = item?.name || "Unknown Item";
  chip.dataset.rarity = normalizeWardrobeRarity(item?.rarity || "basic");
  chip.dataset.compareSlot = compareSlot || "";

  const open = document.createElement("span");
  open.className = "chat-item-bracket";
  open.textContent = "[";

  const name = document.createElement("span");
  name.className = "chat-item-name";
  name.textContent = item?.name || "Unknown Item";

  const close = document.createElement("span");
  close.className = "chat-item-bracket";
  close.textContent = "]";

  chip.append(open, name, close);
  return chip;
}

function serializeChatInput() {
  if (!chatInput) {
    return "";
  }

  if (!chatInput.isContentEditable) {
    return chatInput.value || "";
  }

  let output = "";

  chatInput.childNodes.forEach((node) => {
    if (node.nodeType === Node.TEXT_NODE) {
      output += node.textContent || "";
      return;
    }

    if (node.nodeType !== Node.ELEMENT_NODE) {
      return;
    }

    if (node.classList.contains("chat-item-token")) {
      output += createChatItemToken({
        id: node.dataset.itemId,
        name: node.dataset.itemName,
        rarity: node.dataset.rarity,
      }, node.dataset.compareSlot || "");
      return;
    }

    output += node.textContent || "";
  });

  return output;
}

function getChatInputPlainText() {
  if (!chatInput) {
    return "";
  }

  return chatInput.isContentEditable ? chatInput.textContent || "" : chatInput.value || "";
}

function setChatInputPlainText(text) {
  if (!chatInput) {
    return;
  }

  if (chatInput.isContentEditable) {
    chatInput.replaceChildren(document.createTextNode(String(text || "")));
  } else {
    chatInput.value = String(text || "");
  }
}

function clearChatInput() {
  if (!chatInput) {
    return;
  }

  if (chatInput.isContentEditable) {
    chatInput.replaceChildren();
  } else {
    chatInput.value = "";
  }
}

function setChatInputPlaceholder(text) {
  if (!chatInput) {
    return;
  }

  if (chatInput.isContentEditable) {
    chatInput.dataset.placeholder = text;
  } else {
    chatInput.placeholder = text;
  }
}

function setChatInputDisabled(disabled) {
  if (!chatInput) {
    return;
  }

  if (chatInput.isContentEditable || chatInput.classList.contains("chat-input")) {
    chatInput.contentEditable = disabled ? "false" : "true";
    chatInput.classList.toggle("is-disabled", disabled);
    chatInput.setAttribute("aria-disabled", String(disabled));
  } else {
    chatInput.disabled = disabled;
  }
}

function insertIntoChatInput(text) {
  if (chatInput.isContentEditable) {
    chatInput.appendChild(document.createTextNode(String(text || "")));
  } else {
    chatInput.value = `${chatInput.value}${text}`;
  }

  chatInput.focus();
  setChatCursorToEnd();
}

function setChatCursorToEnd() {
  if (!chatInput) {
    return;
  }

  if (!chatInput.isContentEditable) {
    const end = chatInput.value.length;
    chatInput.setSelectionRange(end, end);
    return;
  }

  const range = document.createRange();
  range.selectNodeContents(chatInput);
  range.collapse(false);

  const selection = window.getSelection();
  selection.removeAllRanges();
  selection.addRange(range);
}

function createChatItemLink(tokenInfo) {
  const button = document.createElement("button");
  button.type = "button";
  button.className = `chat-item-link rarity-${tokenInfo.rarity}`;
  button.title = `View ${tokenInfo.name}`;

  const open = document.createElement("span");
  open.className = "chat-item-bracket";
  open.textContent = "[";

  const name = document.createElement("span");
  name.className = "chat-item-name";
  name.textContent = tokenInfo.name;

  const close = document.createElement("span");
  close.className = "chat-item-bracket";
  close.textContent = "]";

  button.append(open, name, close);
  button.addEventListener("click", (event) => {
    event.preventDefault();
    event.stopPropagation();
    openWardrobeItemDetail(tokenInfo.itemId, tokenInfo.compareSlot || "", { readOnly: true });
  });

  return button;
}

function handleChatInputEnterKey(event) {
  if (event.key !== "Enter" || event.shiftKey) {
    return;
  }

  event.preventDefault();
  chatForm.requestSubmit();
}

document.querySelector("#wardrobe-item-share-button")?.addEventListener("click", shareActiveWardrobeItemToChat);
chatInput?.addEventListener("keydown", handleChatInputEnterKey);

/* Wardrobe modal close after wear or share */

async function equipWardrobeItem(itemID, slotKey = "") {
  try {
    const response = await csrfFetch("/api/player/wardrobe/equip", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        itemId: Number(itemID),
        slotKey,
      }),
    });

    if (!response.ok) {
      throw new Error(`Equip wardrobe item failed: ${response.status}`);
    }

    const result = await response.json();

    await loadPlayerStatus();

    const itemName = result?.detail?.item?.name || "Item";
    addChatMessage("System", `${itemName} equipped. Nami-Chan is now officially more stylish.`, "system");

    closeWardrobeItemModal();
  } catch (error) {
    console.error(error);
    addChatMessage("System", "Could not equip that item. The wardrobe clasp fought back.", "system");
  }
}

function shareActiveWardrobeItemToChat() {
  const item = activeWardrobeModalDetail?.item;

  if (!item || !Number(item.id ?? 0)) {
    return;
  }

  insertChatItemChip(item, activeWardrobeModalCompareSlot || defaultCompareSlotForWardrobeItem(item));
  closeWardrobeItemModal();

  chatInput?.focus();
  setChatCursorToEnd();
}
