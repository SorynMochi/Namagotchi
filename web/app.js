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

const MAX_CHAT_MESSAGES = 100;
const CHAT_STORAGE_KEY = "namagotchi_chat_store_v1";
const CHAT_CHANNEL_KEY = "namagotchi_chat_active_channel_v1";
const CHAT_HIDDEN_KEY = "namagotchi_chat_hidden_v1";
const CHAT_PREVIOUS_HEIGHT_KEY = "namagotchi_chat_previous_height_v1";
const EMOJI_USAGE_KEY = "namagotchi_emoji_usage_v1";
const EMOJI_CATEGORY_KEY = "namagotchi_emoji_category_v1";
const CHAT_IGNORE_KEY = "namagotchi_chat_ignore_list_v1";
const CHAT_LAST_WHISPER_KEY = "namagotchi_chat_last_whisper_v1";
const CHAT_OFFLINE_WHISPERS_KEY = "namagotchi_offline_whispers_v1";

const CURRENT_PLAYER_NAME = "Soryn";

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
  "😀", "😃", "😄", "😁", "😆", "😂", "🤣", "😊", "😇", "🙂",
  "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚", "😋",
  "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🫢", "🫣", "🤫",
  "🤔", "🫡", "🤤", "😴", "🤯", "🥳", "🥺", "😭", "😤", "😈",
  "👿", "😎", "🤓", "🧐", "😳", "🥵", "🥶", "😱", "😵", "💀",
  "☠️", "👻", "👽", "🤖", "😺", "😸", "😹", "😻", "😼", "🙀",
  "👍", "👎", "👏", "🙌", "👐", "🤲", "🙏", "💪", "🫶", "🤝",
  "👀", "👁️", "🧠", "🫀", "💋", "💅", "✨", "💫", "⭐", "🌟",
  "🔥", "💥", "💢", "💦", "💨", "🕯️", "🎀", "🎁", "🎉", "🎊",
  "💖", "💗", "💓", "💕", "💞", "💘", "💝", "💜", "💙", "🩵",
  "💚", "💛", "🧡", "❤️", "🩷", "🖤", "🤍", "🤎", "☕", "🍵",
  "🍪", "🍩", "🍰", "🧁", "🍫", "🍬", "🍭", "🍜", "🍣", "🍙",
  "🍕", "🍔", "🍟", "🥐", "🥞", "🥓", "🍓", "🍒", "🍑", "🍎",
  "🎮", "🕹️", "🎲", "🎧", "🎤", "🎬", "🎨", "🧵", "🪡", "👗",
  "👚", "👠", "👢", "💄", "💻", "⌨️", "🖱️", "📱", "📦", "💎",
  "🪙", "💰", "🔑", "🔒", "🔓", "🧸", "🐾", "🐱", "🐈", "🦊",
  "🐰", "🐹", "🐻", "🐼", "🐧", "🦄", "🌙", "☀️", "🌈", "❄️",
  "🌸", "🌺", "🌷", "🍄", "🪴", "👑", "⚔️", "🛡️", "🎯", "🏆",
  "😅", "😬", "🥲", "🫠", "🫥", "😶", "😶‍🌫️", "😐", "😑", "😒",
  "🙄", "😏", "😕", "🫤", "🙁", "☹️", "😟", "😔", "😞", "😣",
  "😖", "😫", "😩", "🥱", "😮", "😯", "😲", "😦", "😧", "😨",
  "😰", "😥", "😓", "😢", "😪", "😮‍💨", "😵‍💫", "🥴", "🤠", "🥸",
  "🤥", "🤨", "🤐", "🤢", "🤮", "🤧", "🤒", "🤕", "😷", "🤬",
  "😽", "😿", "😾", "🙈", "🙉", "🙊", "👾", "👹", "👺", "🌚",
  "👋", "🤚", "🖐️", "✋", "🖖", "👌", "🤌", "🤏", "✌️", "🤞",
  "🫰", "🤟", "🤘", "🤙", "👈", "👉", "👆", "👇", "☝️", "✍️",
  "👊", "✊", "🤛", "🤜", "🫷", "🫸", "🫱", "🫲", "🫳", "🫴",
  "🦾", "🦿", "🦷", "🦴", "👂", "👃", "👄", "👅", "🫦", "🫂",
  "❤️‍🔥", "❤️‍🩹", "❣️", "💔", "💟", "💌", "💤", "💯", "💠", "🔆",
  "🔅", "💡", "🔦", "🏮", "🪔", "🫧", "🔮", "🪄", "🧿", "🪬",
  "📚", "📖", "🔖", "📝", "✏️", "🖊️", "🖋️", "🖌️", "🖍️", "📌",
  "📍", "📎", "🖇️", "✂️", "📏", "📐", "📓", "📔", "📕", "📗",
  "📘", "📙", "📰", "🗞️", "🧩", "♟️", "🃏", "🀄", "🎴", "🎭",
  "🎷", "🎸", "🎹", "🎺", "🎻", "🥁", "🪘", "🪇", "🪈", "🎼",
  "🎵", "🎶", "📷", "📸", "📹", "🎥", "📺", "📻", "⏰", "⌚",
  "🧭", "🗺️", "👜", "👛", "👝", "🎒", "🧳", "👓", "🕶️", "🥽",
  "🧥", "🥼", "🥿", "👡", "🩰", "👟", "🥾", "🧦", "🧤", "🧣",
  "🎩", "🧢", "👒", "🎓", "⛑️", "🪖", "💍", "🌍", "🌎", "🌏",
  "🪐", "🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘", "🌝",
  "🌞", "🌠", "☄️", "💧", "🌊", "🌫️", "🌪️", "🌧️", "⛈️", "🌩️",
  "⚡", "☔", "☁️", "⛄", "☃️", "🌬️", "🌤️", "⛅", "🌥️", "🌦️",
  "☂️", "🌡️", "🍀", "☘️", "🌿", "🌱", "🌲", "🌳", "🌴", "🌵",
  "🌾", "💐", "🥀", "🌹", "🪷", "🪻", "🌼", "🌻", "🍁", "🍂",
  "🍃", "🪹", "🪺", "🪵", "🪨", "🍏", "🍐", "🍊", "🍋", "🍌",
  "🍉", "🍇", "🫐", "🥝", "🥭", "🍍", "🥥", "🥑", "🍆", "🥕",
  "🌽", "🌶️", "🥒", "🥬", "🥦", "🧄", "🧅", "🥔", "🍠", "🥯",
  "🍞", "🥖", "🧀", "🥚", "🍳", "🧇", "🥨", "🥩", "🍗", "🍖",
  "🌭", "🥪", "🌮", "🌯", "🫔", "🥙", "🧆", "🥘", "🍲", "🥣",
  "🥗", "🍿", "🧈", "🧂", "🥫", "🍱", "🍘", "🍥", "🥮", "🍡",
  "🥠", "🥡", "🦪", "🍤", "🍛", "🍚", "🍝", "🥟", "🍧", "🍨",
  "🍦", "🥧", "🍮", "🍯", "🥛", "🍼", "🫖", "🧋", "🧃", "🥤",
  "🐶", "🐕", "🐩", "🐈‍⬛", "🦁", "🐯", "🐅", "🐆", "🐴", "🫎",
  "🫏", "🐎", "🦌", "🦬", "🐮", "🐂", "🐃", "🐄", "🐷", "🐖",
  "🐗", "🐽", "🐏", "🐑", "🐐", "🐪", "🐫", "🦙", "🦒", "🐘",
  "🦣", "🦏", "🦛", "🐭", "🐁", "🐀", "🦔", "🦇", "🐿️", "🦫",
  "🦥", "🦦", "🦨", "🦘", "🦡", "🦃", "🐔", "🐓", "🐣", "🐤",
  "🐥", "🐦", "🐦‍⬛", "🪿", "🦆", "🦅", "🦉", "🦤", "🪶", "🦩",
  "🦚", "🦜", "🪽", "🐸", "🐊", "🐢", "🦎", "🐍", "🐲", "🐉",
  "🦕", "🦖", "🐳", "🐋", "🐬", "🦭", "🐟", "🐠", "🐡", "🦈",
  "🐙", "🐚", "🪸", "🪼", "🐌", "🦋", "🐛", "🐜", "🐝", "🪲",
  "🐞", "🦗", "🪳", "🕷️", "🕸️", "🦂", "🦟", "🪰", "🪱", "🦠",
  "🏠", "🏡", "🏘️", "🏰", "🏯", "⛩️", "🗼", "🌉", "🎠", "🎡",
  "🎢", "🚗", "🚕", "🚙", "🚌", "🚎", "🏎️", "🚓", "🚑", "🚒",
  "🚚", "🚛", "🛻", "🚲", "🛴", "🛵", "🏍️", "🚂", "🚆", "🚇",
  "🚀", "🛸", "✈️", "🛶", "⛵", "🚢", "⚓", "🗿", "🧱", "🪜",
  "🪑", "🛋️", "🛏️", "🪞", "🪟", "🚪", "🧺", "🧻", "🧼", "🪥",
  "🧽", "🪣", "🛁", "🚿", "✅", "☑️", "✔️", "❌", "❎", "➕",
  "➖", "➗", "✖️", "🔔", "🔕", "📣", "📢", "💬", "🗯️", "💭",
  "🗨️", "🔊", "🔇", "🔈", "🔉", "➡️", "⬅️", "⬆️", "⬇️", "↗️",
  "↘️", "↙️", "↖️", "🔄", "🔁", "🔀", "⚽", "🏀", "🏈", "⚾",
  "🥎", "🎾", "🏐", "🏉", "🥏", "🎳", "🏏", "🏑", "🏒", "🥍",
  "🏓", "🏸", "🥊", "🥋", "🥅", "⛳", "⛸️", "🎣", "🤿", "🎿",
  "🛷", "🥌", "🏹", "🪃", "🪁", "🛼", "🛹", "🩹", "🩺", "💊",
  "💉", "🧬", "🔬", "🔭", "💸", "💳", "🧾", "🏦", "🏧", "🔨",
  "🪓", "⛏️", "⚒️", "🛠️", "🗡️", "🔧", "🪛", "🔩", "⚙️", "🧲",
  "🧰", "🪚", "🧪", "⚗️", "🧫", "🧯", "🧨", "🚧", "⛓️", "🪝",
  "🧷", "🪫", "🔋", "🔌", "💽", "💾", "💿", "📀", "🖥️", "🖨️",
  "🧮", "🎞️", "📽️", "🧑‍💻", "👩‍💻", "👨‍💻", "🧑‍🎨", "👩‍🎨", "👨‍🎨", "🧑‍🍳",
  "👩‍🍳", "👨‍🍳", "🧑‍🔧", "👩‍🔧", "👨‍🔧", "🧑‍🔬", "👩‍🔬", "👨‍🔬", "🧑‍🚀", "👩‍🚀",
  "👨‍🚀", "🧙", "🧙‍♀️", "🧙‍♂️", "🧚", "🧚‍♀️", "🧚‍♂️", "🧛", "🧛‍♀️", "🧛‍♂️",
  "🧜", "🧜‍♀️", "🧜‍♂️", "🧝", "🧝‍♀️", "🧝‍♂️", "🧞", "🧞‍♀️", "🧞‍♂️", "🧟",
  "🧟‍♀️", "🧟‍♂️", "🥷", "🦸", "🦸‍♀️", "🦸‍♂️", "🦹", "🦹‍♀️", "🦹‍♂️", "🧌",
  "🎪", "🎟️", "🎫", "🎰", "🪅", "🪩", "🛍️", "🏷️", "📫", "📬",
  "📭", "📮", "🗳️", "✉️", "📩", "📨", "📧", "📤", "📥", "📜",
  "📃", "📄", "📑", "📊", "📈", "📉", "🗃️", "🗄️", "🗑️"
];

const EMOJI_CATEGORIES = {
  all: {
    label: "All",
    emojis: null,
  },
  faces: {
    label: "Faces",
    emojis: [
      "😀", "😃", "😄", "😁", "😆", "😂", "🤣", "😊", "😇", "🙂",
      "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚", "😋",
      "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🫢", "🫣", "🤫",
      "🤔", "🫡", "🤤", "😴", "🤯", "🥳", "🥺", "😭", "😤", "😈",
      "👿", "😎", "🤓", "🧐", "😳", "🥵", "🥶", "😱", "😵", "💀",
      "☠️", "👻", "👽", "🤖", "😺", "😸", "😹", "😻", "😼", "🙀",
      "😅", "😬", "🥲", "🫠", "🫥", "😶", "😶‍🌫️", "😐", "😑", "😒",
      "🙄", "😏", "😕", "🫤", "🙁", "☹️", "😟", "😔", "😞", "😣",
      "😖", "😫", "😩", "🥱", "😮", "😯", "😲", "😦", "😧", "😨",
      "😰", "😥", "😓", "😢", "😪", "😮‍💨", "😵‍💫", "🥴", "🤠", "🥸",
      "🤥", "🤨", "🤐", "🤢", "🤮", "🤧", "🤒", "🤕", "😷", "🤬",
      "😽", "😿", "😾", "🙈", "🙉", "🙊", "👾", "👹", "👺", "🌚"
    ],
  },
  hands: {
    label: "Hands",
    emojis: [
      "👍", "👎", "👏", "🙌", "👐", "🤲", "🙏", "💪", "🫶", "🤝",
      "👋", "🤚", "🖐️", "✋", "🖖", "👌", "🤌", "🤏", "✌️", "🤞",
      "🫰", "🤟", "🤘", "🤙", "👈", "👉", "👆", "👇", "☝️", "✍️",
      "👊", "✊", "🤛", "🤜", "🫷", "🫸", "🫱", "🫲", "🫳", "🫴"
    ],
  },
  hearts: {
    label: "Hearts",
    emojis: [
      "💖", "💗", "💓", "💕", "💞", "💘", "💝", "💜", "💙", "🩵",
      "💚", "💛", "🧡", "❤️", "🩷", "🖤", "🤍", "🤎", "❤️‍🔥", "❤️‍🩹",
      "❣️", "💔", "💟", "💌", "💋", "🫦", "✨", "💫", "⭐", "🌟"
    ],
  },
  food: {
    label: "Food",
    emojis: [
      "☕", "🍵", "🍪", "🍩", "🍰", "🧁", "🍫", "🍬", "🍭", "🍜",
      "🍣", "🍙", "🍕", "🍔", "🍟", "🥐", "🥞", "🥓", "🍓", "🍒",
      "🍑", "🍎", "🍏", "🍐", "🍊", "🍋", "🍌", "🍉", "🍇", "🫐",
      "🥝", "🥭", "🍍", "🥥", "🥑", "🍆", "🥕", "🌽", "🌶️", "🥒",
      "🥬", "🥦", "🧄", "🧅", "🥔", "🍠", "🥯", "🍞", "🥖", "🧀",
      "🥚", "🍳", "🧇", "🥨", "🥩", "🍗", "🍖", "🌭", "🥪", "🌮",
      "🌯", "🫔", "🥙", "🧆", "🥘", "🍲", "🥣", "🥗", "🍿", "🧈",
      "🧂", "🥫", "🍱", "🍘", "🍥", "🥮", "🍡", "🥠", "🥡", "🦪",
      "🍤", "🍛", "🍚", "🍝", "🥟", "🍧", "🍨", "🍦", "🥧", "🍮",
      "🍯", "🥛", "🍼", "🫖", "🧋", "🧃", "🥤"
    ],
  },
  animals: {
    label: "Animals",
    emojis: [
      "🐾", "🐱", "🐈", "🐈‍⬛", "🦊", "🐰", "🐹", "🐻", "🐼", "🐧",
      "🦄", "🐶", "🐕", "🐩", "🦁", "🐯", "🐅", "🐆", "🐴", "🫎",
      "🫏", "🐎", "🦌", "🦬", "🐮", "🐂", "🐃", "🐄", "🐷", "🐖",
      "🐗", "🐽", "🐏", "🐑", "🐐", "🐪", "🐫", "🦙", "🦒", "🐘",
      "🦣", "🦏", "🦛", "🐭", "🐁", "🐀", "🦔", "🦇", "🐿️", "🦫",
      "🦥", "🦦", "🦨", "🦘", "🦡", "🦃", "🐔", "🐓", "🐣", "🐤",
      "🐥", "🐦", "🐦‍⬛", "🪿", "🦆", "🦅", "🦉", "🦤", "🪶", "🦩",
      "🦚", "🦜", "🪽", "🐸", "🐊", "🐢", "🦎", "🐍", "🐲", "🐉",
      "🦕", "🦖", "🐳", "🐋", "🐬", "🦭", "🐟", "🐠", "🐡", "🦈",
      "🐙", "🐚", "🪸", "🪼", "🐌", "🦋", "🐛", "🐜", "🐝", "🪲",
      "🐞", "🦗", "🪳", "🕷️", "🕸️", "🦂", "🦟", "🪰", "🪱", "🦠"
    ],
  },
  nature: {
    label: "Nature",
    emojis: [
      "🌙", "☀️", "🌈", "❄️", "🌸", "🌺", "🌷", "🍄", "🪴", "🌍",
      "🌎", "🌏", "🪐", "🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗",
      "🌘", "🌝", "🌞", "🌠", "☄️", "💧", "🌊", "🌫️", "🌪️", "🌧️",
      "⛈️", "🌩️", "⚡", "☔", "☁️", "⛄", "☃️", "🌬️", "🌤️", "⛅",
      "🌥️", "🌦️", "☂️", "🌡️", "🍀", "☘️", "🌿", "🌱", "🌲", "🌳",
      "🌴", "🌵", "🌾", "💐", "🥀", "🌹", "🪷", "🪻", "🌼", "🌻",
      "🍁", "🍂", "🍃", "🪹", "🪺", "🪵", "🪨"
    ],
  },
  gaming: {
    label: "Gaming",
    emojis: [
      "🎮", "🕹️", "🎲", "🎧", "🎤", "🎬", "🎨", "💻", "⌨️", "🖱️",
      "📱", "🖥️", "💽", "💾", "💿", "📀", "👾", "🤖", "🧑‍💻", "👩‍💻",
      "👨‍💻", "⚔️", "🛡️", "🎯", "🏆", "🥇", "🥈", "🥉"
    ],
  },
  objects: {
    label: "Objects",
    emojis: [
      "📦", "💎", "🪙", "💰", "🔑", "🔒", "🔓", "🧸", "👑", "🎀",
      "🎁", "🎉", "🎊", "📚", "📖", "🔖", "📝", "✏️", "🖊️", "🖋️",
      "🖌️", "🖍️", "📌", "📍", "📎", "🖇️", "✂️", "📏", "📐", "📓",
      "📔", "📕", "📗", "📘", "📙", "🧩", "🎭", "🎵", "🎶", "📷",
      "📸", "📹", "🎥", "📺", "📻", "⏰", "⌚", "🧭", "🗺️", "👜",
      "👛", "👝", "🎒", "🧳", "👓", "🕶️", "👗", "👚", "👠", "👢",
      "💄", "🧵", "🪡", "💍", "🏠", "🏡", "🛋️", "🛏️", "🪞", "🧺",
      "🧼", "🧽", "🪣", "🛁", "🚿", "🧾", "🔨", "🛠️", "⚙️", "🔋"
    ],
  },
  symbols: {
    label: "Symbols",
    emojis: [
      "✅", "☑️", "✔️", "❌", "❎", "➕", "➖", "➗", "✖️", "🔔",
      "🔕", "📣", "📢", "💬", "🗯️", "💭", "🗨️", "🔊", "🔇", "🔈",
      "🔉", "➡️", "⬅️", "⬆️", "⬇️", "↗️", "↘️", "↙️", "↖️", "🔄",
      "🔁", "🔀", "💯", "💠", "🔆", "🔅"
    ],
  },
};

let latestPlayerStatus = null;
let forceTickButton = null;
let currentChatChannel = "lobby";
let chatMessages = createEmptyChatStore();
let unreadChannels = new Set();
let emojiUsage = loadEmojiUsage();
let activeEmojiCategory = getSavedEmojiCategory();
let ignoredPlayers = loadIgnoredPlayers();
let lastWhisperName = localStorage.getItem(CHAT_LAST_WHISPER_KEY) || "";
let activeChatUserMenu = null;
let activeChatProfileModal = null;
let isResizingChat = false;
let isChatHidden = false;
let previousChatHeight = 190;
let serverClockOffsetMs = 0;
let hasServerClock = false;

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
    addChatMessage("System", tickResultMessage(result), "system");
    await loadPlayerStatus();
  } catch (error) {
    console.error(error);
    addChatMessage("System", "Force tick failed. The tick goblin dropped its tiny clipboard.", "system");
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
    addChatMessage("System", `Gathering task changed to ${labelForTask(task)}.`, "system");
  } catch (error) {
    console.error(error);
    addChatMessage("System", "Could not change gathering task.", "system");
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

function createEmptyChatStore() {
  return CHAT_CHANNELS.reduce((store, channel) => {
    store[channel] = [];
    return store;
  }, {});
}

function initializeChat() {
  chatMessages = loadChatStore();
  currentChatChannel = getSavedChatChannel();

  const hasStoredMessages = CHAT_CHANNELS.some((channel) => chatMessages[channel].length > 0);

  if (!hasStoredMessages) {
    pushChatMessage("system", "System", "Welcome to Namagotchi Phase 3B.", false);
    pushChatMessage("lobby", "Nami-chan", "Lobby chat is online. I am absolutely not testing the buttons with my tiny chaos paws.", false);
    pushChatMessage("help", "System", "Help chat will be used for player questions once multiplayer chat is wired in.", false);
    saveChatStore();
  }

  unreadChannels.clear();
  renderEmojiPicker();
  initializeChatVisibility();
  switchChatChannel(currentChatChannel);
}

function initializeChatResize() {
  const savedHeight = Number(localStorage.getItem("namagotchi_chat_height"));
  if (savedHeight) {
    setChatHeight(savedHeight);
  }

  chatResizeHandle.addEventListener("pointerdown", (event) => {
  if (isChatHidden) {
    return;
  }

  event.preventDefault();
  isResizingChat = true;
  document.body.classList.add("is-resizing-chat");
  chatResizeHandle.setPointerCapture(event.pointerId);
});

  window.addEventListener("pointermove", (event) => {
    if (!isResizingChat) {
      return;
    }

    const centerColumn = document.querySelector(".center-column");
    const rect = centerColumn.getBoundingClientRect();
    const newHeight = rect.bottom - event.clientY;

    setChatHeight(newHeight);
  });

  window.addEventListener("pointerup", () => {
    if (!isResizingChat) {
      return;
    }

    isResizingChat = false;
    document.body.classList.remove("is-resizing-chat");

    const chatHeight = Math.round(chatPanel.getBoundingClientRect().height);
    previousChatHeight = chatHeight;
    localStorage.setItem("namagotchi_chat_height", String(chatHeight));
    localStorage.setItem(CHAT_PREVIOUS_HEIGHT_KEY, String(chatHeight));
  });
}

function setChatHeight(height) {
  const clampedHeight = Math.max(120, Math.min(430, Number(height)));
  document.documentElement.style.setProperty("--chat-height", `${clampedHeight}px`);
}

function submitChatMessage(event) {
  event.preventDefault();

  const text = normalizeChatText(chatInput.value);
  if (!text) {
    return;
  }

  if (currentChatChannel === "system") {
    addChatMessage("System", "System chat is read-only. Tiny velvet rope deployed.", "system");
    chatInput.value = "";
    return;
  }

  handleChatInput(text);
  chatInput.value = "";
}

function switchChatChannel(channel) {
  if (!CHAT_CHANNELS.includes(channel)) {
    return;
  }

  currentChatChannel = channel;
  localStorage.setItem(CHAT_CHANNEL_KEY, channel);
  unreadChannels.delete(channel);

  chatTabs.forEach((button) => {
    const isActive = button.dataset.chatChannel === channel;
    button.classList.toggle("active", isActive);
    button.classList.toggle("has-unread", unreadChannels.has(button.dataset.chatChannel));
  });

  chatInput.placeholder = `Message [${CHAT_LABELS[channel]}]...`;
  chatInput.disabled = channel === "system";

  if (channel === "system") {
    chatInput.placeholder = "System messages are read-only.";
  }

  if (channel === "whispers" && lastWhisperName && !chatInput.value.trim()) {
    chatInput.value = `/w ${lastWhisperName} `;
    setChatCursorToEnd();
  }

  renderChatChannel();
}

function addChatMessage(username, text, channel = currentChatChannel, options = {}) {
  const destination = username === "System" && channel === currentChatChannel ? "system" : channel;

  pushChatMessage(destination, username, text, options);

  if (destination === currentChatChannel) {
    renderChatChannel();
  } else {
    unreadChannels.add(destination);
    updateChatUnreadTabs();
  }
}

function pushChatMessage(channel, username, text, optionsOrMarkUnread = {}) {
  if (!CHAT_CHANNELS.includes(channel)) {
    channel = "lobby";
  }

  const options =
    typeof optionsOrMarkUnread === "boolean"
      ? { markUnread: optionsOrMarkUnread }
      : optionsOrMarkUnread;

  const message = {
    username: normalizeChatText(username, 40),
    text: normalizeChatText(text),
    timestamp: getChatTimestamp(),
    kind: options.kind || "normal",
    target: options.target || "",
    whisperTo: options.whisperTo || "",
    whisperFrom: options.whisperFrom || "",
  };

  chatMessages[channel].push(message);

  while (chatMessages[channel].length > MAX_CHAT_MESSAGES) {
    chatMessages[channel].shift();
  }

  saveChatStore();

  if ((options.markUnread ?? true) && channel !== currentChatChannel) {
    unreadChannels.add(channel);
  }
}

function renderChatChannel() {
  const messages = chatMessages[currentChatChannel];

  chatLog.replaceChildren();

  messages
    .filter((message) => !isIgnored(message.username))
    .forEach((message) => {
      const row = document.createElement("p");
      row.className = "chat-message";
      row.classList.toggle("chat-emote", message.kind === "emote");
      row.classList.toggle("chat-whisper", message.kind === "whisper");
      row.classList.toggle("chat-mentioned", messageMentionsCurrentPlayer(message));

      const time = document.createElement("span");
      time.className = "chat-time";
      time.textContent = `[${message.timestamp}]`;

      row.append(time, " ");

      if (message.kind === "emote") {
        const emote = document.createElement("span");
        emote.className = "chat-emote-text";

        const nameButton = createPlayerNameButton(message.username);
        emote.append(nameButton, " ");
        appendFormattedChatText(emote, message.text);

        row.append(emote);
      } else {
        const nameButton = createPlayerNameButton(message.username);
        row.append(nameButton);

        if (message.kind !== "system") {
          row.append(": ");
        } else {
          row.append(" ");
        }

        const text = document.createElement("span");
        text.className = "chat-text";
        appendFormattedChatText(text, message.text);
        row.append(text);
      }

      chatLog.appendChild(row);
    });

  chatLog.scrollTop = chatLog.scrollHeight;
  updateChatUnreadTabs();
}

function updateChatUnreadTabs() {
  chatTabs.forEach((button) => {
    const channel = button.dataset.chatChannel;
    button.classList.toggle("has-unread", unreadChannels.has(channel) && channel !== currentChatChannel);
  });
}

function getChatTimestamp() {
  const date = new Date();
  return `${pad2(date.getHours())}:${pad2(date.getMinutes())}:${pad2(date.getSeconds())}`;
}

function handleChatInput(text) {
  if (text.startsWith("/")) {
    handleChatCommand(text);
    return;
  }

  addChatMessage(CURRENT_PLAYER_NAME, text, currentChatChannel);
}

function handleChatCommand(text) {
  const [commandRaw, ...parts] = text.split(/\s+/);
  const command = commandRaw.toLowerCase();
  const rest = text.slice(commandRaw.length).trim();

  switch (command) {
    case "/me":
      handleMeCommand(rest);
      return;

    case "/w":
    case "/whisper":
      handleWhisperCommand(rest);
      return;

    case "/ignore":
      handleIgnoreCommand(parts[0]);
      return;

    case "/unignore":
      handleUnignoreCommand(parts[0]);
      return;

    case "/ignored":
      handleIgnoredCommand();
      return;

    case "/profile":
      handleProfileCommand(parts[0]);
      return;

    default:
      addChatMessage("System", `Unknown command: ${commandRaw}`, "system");
  }
}

function handleMeCommand(message) {
  if (!message) {
    addChatMessage("System", "Usage: /me action text", "system");
    return;
  }

  addChatMessage(CURRENT_PLAYER_NAME, message, currentChatChannel, {
    kind: "emote",
  });
}

function handleWhisperCommand(rest) {
  const match = rest.match(/^(\S+)\s+([\s\S]+)$/);

  if (!match) {
    addChatMessage("System", "Usage: /w playername message", "whispers");
    switchChatChannel("whispers");
    return;
  }

  const requestedName = match[1];
  const message = normalizeChatText(match[2]);
  const player = getPlayerByName(requestedName);

  switchChatChannel("whispers");

  if (!player) {
    addChatMessage("System", `Player "${requestedName}" does not exist.`, "whispers");
    return;
  }

  if (isIgnored(player.displayName)) {
    addChatMessage("System", `You are ignoring ${player.displayName}. Use /unignore ${player.displayName} first.`, "whispers");
    return;
  }

  setLastWhisperName(player.displayName);

  addChatMessage(CURRENT_PLAYER_NAME, `to ${player.displayName}: ${message}`, "whispers", {
    kind: "whisper",
    whisperTo: player.displayName,
    whisperFrom: CURRENT_PLAYER_NAME,
  });

  if (!player.online) {
    queueOfflineWhisper(player.displayName, message);
    addChatMessage("System", `${player.displayName} is offline. Your whisper will be delivered when they come online again.`, "whispers");
  }
}

function handleIgnoreCommand(name) {
  const player = getPlayerByName(name);

  if (!name) {
    addChatMessage("System", "Usage: /ignore playername", "system");
    return;
  }

  if (!player) {
    addChatMessage("System", `Player "${name}" does not exist.`, "system");
    return;
  }

  if (normalizePlayerName(player.displayName) === normalizePlayerName(CURRENT_PLAYER_NAME)) {
    addChatMessage("System", "You cannot ignore yourself. Even if past-you deserved it.", "system");
    return;
  }

  ignoredPlayers.add(normalizePlayerName(player.displayName));
  saveIgnoredPlayers();
  addChatMessage("System", `You are now ignoring ${player.displayName}.`, "system");
  renderChatChannel();
}

function handleUnignoreCommand(name) {
  const player = getPlayerByName(name);

  if (!name) {
    addChatMessage("System", "Usage: /unignore playername", "system");
    return;
  }

  const normalizedName = normalizePlayerName(player?.displayName || name);

  if (!ignoredPlayers.has(normalizedName)) {
    addChatMessage("System", `${name} is not on your ignore list.`, "system");
    return;
  }

  ignoredPlayers.delete(normalizedName);
  saveIgnoredPlayers();
  addChatMessage("System", `You are no longer ignoring ${player?.displayName || name}.`, "system");
  renderChatChannel();
}

function handleIgnoredCommand() {
  const ignored = [...ignoredPlayers];

  if (!ignored.length) {
    addChatMessage("System", "Your ignore list is empty.", "system");
    return;
  }

  addChatMessage("System", `Ignored players: ${ignored.join(", ")}`, "system");
}

function handleProfileCommand(name) {
  if (!name) {
    addChatMessage("System", "Usage: /profile playername", "system");
    return;
  }

  const player = getPlayerByName(name);

  if (!player) {
    addChatMessage("System", `Player "${name}" does not exist.`, "system");
    return;
  }

  openProfileModal(player.displayName);
}

function queueOfflineWhisper(playerName, message) {
  const queued = loadOfflineWhispers();

  queued.push({
    to: playerName,
    from: CURRENT_PLAYER_NAME,
    text: message,
    queuedAt: new Date().toISOString(),
  });

  localStorage.setItem(CHAT_OFFLINE_WHISPERS_KEY, JSON.stringify(queued.slice(-100)));
}

function loadOfflineWhispers() {
  try {
    const parsed = JSON.parse(localStorage.getItem(CHAT_OFFLINE_WHISPERS_KEY));
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

function setLastWhisperName(name) {
  lastWhisperName = name;
  localStorage.setItem(CHAT_LAST_WHISPER_KEY, name);
}

function getPlayerByName(name) {
  if (!name) {
    return null;
  }

  return DEV_PLAYER_DIRECTORY[normalizePlayerName(name)] || null;
}

function normalizePlayerName(name) {
  return String(name || "").trim().toLowerCase();
}

function isIgnored(name) {
  return ignoredPlayers.has(normalizePlayerName(name));
}

function loadIgnoredPlayers() {
  try {
    const parsed = JSON.parse(localStorage.getItem(CHAT_IGNORE_KEY));
    return new Set(Array.isArray(parsed) ? parsed.map(normalizePlayerName) : []);
  } catch {
    return new Set();
  }
}

function saveIgnoredPlayers() {
  localStorage.setItem(CHAT_IGNORE_KEY, JSON.stringify([...ignoredPlayers]));
}

function messageMentionsCurrentPlayer(message) {
  return getMentionedNames(message.text).some((name) => normalizePlayerName(name) === normalizePlayerName(CURRENT_PLAYER_NAME));
}

function getMentionedNames(text) {
  return [...String(text).matchAll(/@([A-Za-z0-9_-]{2,32})/g)].map((match) => match[1]);
}

function appendFormattedChatText(parent, text) {
  const source = String(text);
  const pattern = /(\*\*\*[^*]+\*\*\*|\*\*[^*]+\*\*|\*[^*]+\*|@[A-Za-z0-9_-]{2,32})/g;

  let lastIndex = 0;

  source.replace(pattern, (match, _ignored, offset) => {
    appendPlainChatText(parent, source.slice(lastIndex, offset));

    if (match.startsWith("@")) {
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

function appendPlainChatText(parent, text) {
  if (text) {
    parent.appendChild(document.createTextNode(text));
  }
}

function appendMention(parent, mentionText) {
  const playerName = mentionText.slice(1);
  const button = document.createElement("button");
  button.type = "button";
  button.className = "chat-mention";
  button.textContent = mentionText;

  button.addEventListener("click", (event) => {
    event.stopPropagation();
    openChatUserMenu(playerName, button);
  });

  parent.appendChild(button);
}

function createPlayerNameButton(username) {
  if (username === "System") {
    const span = document.createElement("span");
    span.className = "chat-name system-name";
    span.textContent = "System";
    return span;
  }

  const button = document.createElement("button");
  button.type = "button";
  button.className = "chat-name chat-player-name";
  button.textContent = username;

  button.addEventListener("click", (event) => {
    event.stopPropagation();
    openChatUserMenu(username, button);
  });

  return button;
}

function openChatUserMenu(playerName, anchorElement) {
  closeChatUserMenu();

  const player = getPlayerByName(playerName);
  const displayName = player?.displayName || playerName;

  activeChatUserMenu = document.createElement("div");
  activeChatUserMenu.className = "chat-user-menu";

  activeChatUserMenu.append(
    createChatUserMenuButton("@Mention", () => {
      insertIntoChatInput(`@${displayName} `);
      closeChatUserMenu();
    }),
    createChatUserMenuButton("Message", () => {
      setLastWhisperName(displayName);
      switchChatChannel("whispers");
      chatInput.value = `/w ${displayName} `;
      setChatCursorToEnd();
      closeChatUserMenu();
    }),
    createChatUserMenuButton("Profile", () => {
      openProfileModal(displayName);
      closeChatUserMenu();
    }),
    createChatUserMenuButton("View Market", () => {
      showSection("market");
      addChatMessage("System", `Showing market listings for ${displayName}. Full market filtering comes later.`, "system");
      closeChatUserMenu();
    }),
    createChatUserMenuButton(isIgnored(displayName) ? "Unignore User" : "Ignore User", () => {
      const playerRecord = getPlayerByName(displayName);

      if (!playerRecord) {
        addChatMessage("System", `Player "${displayName}" does not exist.`, "system");
        closeChatUserMenu();
        return;
      }

      if (isIgnored(displayName)) {
        ignoredPlayers.delete(normalizePlayerName(displayName));
        addChatMessage("System", `You are no longer ignoring ${displayName}.`, "system");
      } else {
        ignoredPlayers.add(normalizePlayerName(displayName));
        addChatMessage("System", `You are now ignoring ${displayName}.`, "system");
      }

      saveIgnoredPlayers();
      renderChatChannel();
      closeChatUserMenu();
    })
  );

  document.body.appendChild(activeChatUserMenu);
  positionFloatingElement(activeChatUserMenu, anchorElement);
}

function createChatUserMenuButton(label, action) {
  const button = document.createElement("button");
  button.type = "button";
  button.textContent = label;
  button.addEventListener("click", (event) => {
    event.stopPropagation();
    action();
  });
  return button;
}

function closeChatUserMenu() {
  if (activeChatUserMenu) {
    activeChatUserMenu.remove();
    activeChatUserMenu = null;
  }
}

function openProfileModal(playerName) {
  const player = getPlayerByName(playerName);

  if (!player) {
    addChatMessage("System", `Player "${playerName}" does not exist.`, "system");
    return;
  }

  closeProfileModal();

  activeChatProfileModal = document.createElement("div");
  activeChatProfileModal.className = "chat-profile-modal";
  activeChatProfileModal.innerHTML = "";

  const title = document.createElement("h2");
  title.textContent = player.displayName;

  const status = document.createElement("p");
  status.textContent = `Status: ${player.online ? "Online" : "Offline"}`;

  const level = document.createElement("p");
  level.textContent = `Level: ${player.level.toLocaleString()}`;

  const note = document.createElement("p");
  note.className = "muted";
  note.textContent = "Full player profiles will be server-backed later.";

  const close = document.createElement("button");
  close.type = "button";
  close.textContent = "Close";
  close.addEventListener("click", closeProfileModal);

  activeChatProfileModal.append(title, status, level, note, close);
  document.body.appendChild(activeChatProfileModal);
}

function closeProfileModal() {
  if (activeChatProfileModal) {
    activeChatProfileModal.remove();
    activeChatProfileModal = null;
  }
}

function positionFloatingElement(element, anchorElement) {
  const rect = anchorElement.getBoundingClientRect();
  const gap = 8;
  const padding = 8;

  const width = element.offsetWidth || 180;
  const height = element.offsetHeight || 180;

  let left = rect.left;
  let top = rect.bottom + gap;

  left = Math.max(padding, Math.min(left, window.innerWidth - width - padding));

  if (top + height > window.innerHeight - padding) {
    top = rect.top - height - gap;
  }

  top = Math.max(padding, Math.min(top, window.innerHeight - height - padding));

  element.style.left = `${left}px`;
  element.style.top = `${top}px`;
}

function insertIntoChatInput(text) {
  chatInput.value = `${chatInput.value}${text}`;
  chatInput.focus();
  setChatCursorToEnd();
}

function setChatCursorToEnd() {
  const end = chatInput.value.length;
  chatInput.setSelectionRange(end, end);
}

function normalizeChatText(value, maxLength = 255) {
  return Array.from(String(value).trim()).slice(0, maxLength).join("");
}

function loadChatStore() {
  try {
    const parsed = JSON.parse(localStorage.getItem(CHAT_STORAGE_KEY));
    return normalizeChatStore(parsed);
  } catch {
    return createEmptyChatStore();
  }
}

function saveChatStore() {
  localStorage.setItem(CHAT_STORAGE_KEY, JSON.stringify(chatMessages));
}

function normalizeChatStore(value) {
  const store = createEmptyChatStore();

  if (!value || typeof value !== "object") {
    return store;
  }

  CHAT_CHANNELS.forEach((channel) => {
    if (!Array.isArray(value[channel])) {
      return;
    }

    store[channel] = value[channel]
      .slice(-MAX_CHAT_MESSAGES)
      .map((message) => {
        return {
          username: normalizeChatText(message.username ?? "Unknown", 40),
          text: normalizeChatText(message.text ?? ""),
          timestamp: normalizeChatText(message.timestamp ?? getChatTimestamp(), 12),
          kind: normalizeChatText(message.kind ?? "normal", 20),
          target: normalizeChatText(message.target ?? "", 40),
          whisperTo: normalizeChatText(message.whisperTo ?? "", 40),
          whisperFrom: normalizeChatText(message.whisperFrom ?? "", 40),
        };
      });
  });

  return store;
}

function getSavedChatChannel() {
  const savedChannel = localStorage.getItem(CHAT_CHANNEL_KEY);
  return CHAT_CHANNELS.includes(savedChannel) ? savedChannel : "lobby";
}

function loadEmojiUsage() {
  try {
    const parsed = JSON.parse(localStorage.getItem(EMOJI_USAGE_KEY));
    return parsed && typeof parsed === "object" ? parsed : {};
  } catch {
    return {};
  }
}

function saveEmojiUsage() {
  localStorage.setItem(EMOJI_USAGE_KEY, JSON.stringify(emojiUsage));
}

function renderEmojiPicker() {
  emojiPicker.replaceChildren();

  const tabs = document.createElement("div");
  tabs.className = "emoji-picker-tabs";

  Object.entries(EMOJI_CATEGORIES).forEach(([categoryKey, category]) => {
    const tab = document.createElement("button");
    tab.type = "button";
    tab.className = "emoji-category-tab";
    tab.classList.toggle("active", categoryKey === activeEmojiCategory);
    tab.textContent = category.label;

    tab.addEventListener("click", (event) => {
      event.preventDefault();
      event.stopPropagation();
      setEmojiCategory(categoryKey);
    });

    tabs.appendChild(tab);
  });

  const grid = document.createElement("div");
  grid.className = "emoji-grid";

  getSortedEmojis(activeEmojiCategory).forEach((emoji) => {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "emoji-option";
    button.textContent = emoji;
    button.title = `${emojiUsage[emoji] ?? 0} use(s)`;

    button.addEventListener("click", (event) => {
      event.preventDefault();
      event.stopPropagation();
      addEmojiToChat(emoji);
    });

    grid.appendChild(button);
  });

  emojiPicker.append(tabs, grid);
}

function getSortedEmojis(categoryKey = "all") {
  const category = EMOJI_CATEGORIES[categoryKey] ?? EMOJI_CATEGORIES.all;
  const emojiSource = category.emojis ?? EMOJI_OPTIONS;

  return [...emojiSource]
    .filter((emoji, index, array) => array.indexOf(emoji) === index)
    .filter((emoji) => EMOJI_OPTIONS.includes(emoji))
    .sort((a, b) => {
      const countDifference = (emojiUsage[b] ?? 0) - (emojiUsage[a] ?? 0);

      if (countDifference !== 0) {
        return countDifference;
      }

      return EMOJI_OPTIONS.indexOf(a) - EMOJI_OPTIONS.indexOf(b);
    });
}

function setEmojiCategory(categoryKey) {
  if (!EMOJI_CATEGORIES[categoryKey]) {
    return;
  }

  activeEmojiCategory = categoryKey;
  localStorage.setItem(EMOJI_CATEGORY_KEY, categoryKey);
  renderEmojiPicker();
  positionEmojiPicker();
}

function getSavedEmojiCategory() {
  const savedCategory = localStorage.getItem(EMOJI_CATEGORY_KEY);
  return EMOJI_CATEGORIES[savedCategory] ? savedCategory : "all";
}

function initializeEmojiPickerPortal() {
  if (emojiPicker.parentElement !== document.body) {
    document.body.appendChild(emojiPicker);
  }

  emojiPicker.classList.add("emoji-picker-floating");

  emojiPicker.addEventListener("click", (event) => {
    event.stopPropagation();
  });

  window.addEventListener("resize", positionVisibleEmojiPicker);
  window.addEventListener("scroll", positionVisibleEmojiPicker, true);
}

function positionVisibleEmojiPicker() {
  if (!emojiPicker.classList.contains("hidden")) {
    positionEmojiPicker();
  }
}

function positionEmojiPicker() {
  const buttonRect = emojiButton.getBoundingClientRect();
  const gap = 8;
  const pagePadding = 8;

  const pickerWidth = emojiPicker.offsetWidth || 324;
  const pickerHeight = emojiPicker.offsetHeight || 248;

  let left = buttonRect.left;
  let top = buttonRect.top - pickerHeight - gap;

  left = Math.max(pagePadding, Math.min(left, window.innerWidth - pickerWidth - pagePadding));

  if (top < pagePadding) {
    top = buttonRect.bottom + gap;
  }

  top = Math.max(pagePadding, Math.min(top, window.innerHeight - pickerHeight - pagePadding));

  emojiPicker.style.left = `${left}px`;
  emojiPicker.style.top = `${top}px`;
}

function toggleEmojiPicker(event) {
  event?.preventDefault();
  event?.stopPropagation();

  const shouldOpen = emojiPicker.classList.contains("hidden");

  if (shouldOpen) {
    renderEmojiPicker();
    emojiPicker.classList.remove("hidden");
    emojiButton.setAttribute("aria-expanded", "true");
    positionEmojiPicker();
    return;
  }

  closeEmojiPicker();
}

function closeEmojiPicker() {
  emojiPicker.classList.add("hidden");
  emojiButton.setAttribute("aria-expanded", "false");
}

function addEmojiToChat(emoji) {
  const currentValue = chatInput.value;
  chatInput.value = Array.from(`${currentValue}${emoji}`).slice(0, 255).join("");
  chatInput.focus();

  emojiUsage[emoji] = (emojiUsage[emoji] ?? 0) + 1;
  saveEmojiUsage();
  renderEmojiPicker();
  closeEmojiPicker();
}

function toggleEmojiPicker() {
  const shouldOpen = emojiPicker.classList.contains("hidden");

  if (shouldOpen) {
    renderEmojiPicker();
    emojiPicker.classList.remove("hidden");
    emojiButton.setAttribute("aria-expanded", "true");
    positionEmojiPicker();
    return;
  }

  closeEmojiPicker();
}

function closeEmojiPicker() {
  emojiPicker.classList.add("hidden");
  emojiButton.setAttribute("aria-expanded", "false");
}

function addEmojiToChat(emoji) {
  const currentValue = chatInput.value;
  chatInput.value = Array.from(`${currentValue}${emoji}`).slice(0, 255).join("");
  chatInput.focus();

  emojiUsage[emoji] = (emojiUsage[emoji] ?? 0) + 1;
  saveEmojiUsage();
  renderEmojiPicker();
  closeEmojiPicker();
}

function initializeChatVisibility() {
  previousChatHeight =
    Number(localStorage.getItem(CHAT_PREVIOUS_HEIGHT_KEY)) ||
    Number(localStorage.getItem("namagotchi_chat_height")) ||
    190;

  setChatHidden(localStorage.getItem(CHAT_HIDDEN_KEY) === "true", false);
}

function setChatHidden(hidden, shouldSave = true) {
  isChatHidden = hidden;

  if (hidden) {
    previousChatHeight = Math.round(chatPanel.getBoundingClientRect().height) || previousChatHeight || 190;
    localStorage.setItem(CHAT_PREVIOUS_HEIGHT_KEY, String(previousChatHeight));

    chatPanel.classList.add("chat-hidden");
    chatToggleButton.textContent = "Show Chat";
    chatToggleButton.setAttribute("aria-pressed", "true");
    closeEmojiPicker();
    document.documentElement.style.setProperty("--chat-height", "46px");
  } else {
    chatPanel.classList.remove("chat-hidden");
    chatToggleButton.textContent = "Hide Chat";
    chatToggleButton.setAttribute("aria-pressed", "false");
    setChatHeight(previousChatHeight || 190);
  }

  if (shouldSave) {
    localStorage.setItem(CHAT_HIDDEN_KEY, String(hidden));
  }
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

function syncServerClock(timestamp) {
  const serverDate = new Date(timestamp);

  if (Number.isNaN(serverDate.getTime())) {
    return;
  }

  serverClockOffsetMs = serverDate.getTime() - Date.now();
  hasServerClock = true;
  updateLiveServerClock();
}

function updateLiveServerClock() {
  if (!hasServerClock) {
    return;
  }

  serverTime.textContent = formatDateTime(Date.now() + serverClockOffsetMs);
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

setInterval(updateLiveServerClock, 1000);
setInterval(loadStatus, 10000);
setInterval(loadPlayerStatus, 5000);