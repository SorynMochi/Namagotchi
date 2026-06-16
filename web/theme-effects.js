(() => {
  const reduceMotionQuery = window.matchMedia("(prefers-reduced-motion: reduce)");

  let activeThemeKey = "";
  let activeLayer = null;
  let petalTimer = null;
  let lanternTimer = null;

  function randomBetween(min, max) {
    return Math.random() * (max - min) + min;
  }

  function randomItem(items) {
    return items[Math.floor(Math.random() * items.length)];
  }

  function isSakuraTheme(themeKey) {
    return themeKey === "sakura" || themeKey === "sakura-dark" || themeKey === "sakura-light";
  }

  function isTokyoNightTheme(themeKey) {
    return themeKey === "tokyo-night";
  }

  function clearThemeEffects() {
    if (petalTimer) {
      window.clearInterval(petalTimer);
      petalTimer = null;
    }

    if (lanternTimer) {
      window.clearInterval(lanternTimer);
      lanternTimer = null;
    }

    if (activeLayer) {
      activeLayer.remove();
      activeLayer = null;
    }

    document.documentElement.classList.remove("theme-effect-sakura", "theme-effect-tokyo-night", "theme-effect-candy");
    document.body?.classList.remove("theme-effect-sakura", "theme-effect-tokyo-night", "theme-effect-candy");
    document.querySelectorAll(".tokyo-night-rail-signs, .tokyo-night-left-signs").forEach((node) => node.remove());
  }

  function setActiveThemeEffect(themeKey) {
    activeThemeKey = themeKey;

    document.documentElement.dataset.theme = themeKey;
    document.body?.setAttribute("data-theme", themeKey);

    clearThemeEffects();

    if (isSakuraTheme(themeKey)) {
      startSakuraThemeEffect();
    }

    if (isTokyoNightTheme(themeKey)) {
      startTokyoNightThemeEffect();
    }
  }

  function startSakuraThemeEffect() {
    document.documentElement.classList.add("theme-effect-sakura");
    document.body?.classList.add("theme-effect-sakura");

    const layer = document.createElement("div");
    layer.className = "theme-effects-layer theme-effects-sakura";
    layer.setAttribute("aria-hidden", "true");

    const moon = document.createElement("div");
    moon.className = "sakura-moon";

    const branch = document.createElement("div");
    branch.className = "sakura-branch";

    const lanternLeft = document.createElement("div");
    lanternLeft.className = "sakura-lantern sakura-lantern-left";

    const lanternRight = document.createElement("div");
    lanternRight.className = "sakura-lantern sakura-lantern-right";

    const petalField = document.createElement("div");
    petalField.className = "sakura-petal-field";

    if (activeThemeKey === "sakura-light") {
      const ryokanLayer = document.createElement("div");
      ryokanLayer.className = "sakura-light-ryokan-layer";

      const noren = document.createElement("div");
      noren.className = "sakura-light-noren";

      const tokonoma = document.createElement("div");
      tokonoma.className = "sakura-light-tokonoma";

      const teaService = document.createElement("div");
      teaService.className = "sakura-light-tea-service";

      const steam = document.createElement("div");
      steam.className = "sakura-light-steam";

      const tatamiMark = document.createElement("div");
      tatamiMark.className = "sakura-light-tatami-mark";

      ryokanLayer.append(noren, tokonoma, teaService, steam, tatamiMark);
      layer.append(ryokanLayer, petalField);
    } else {
      layer.append(moon, branch, lanternLeft, lanternRight, petalField);
    }

    const leftRail = document.querySelector(".left-rail");
    const navPanel = document.querySelector(".nav-panel");

    if (leftRail && navPanel) {
      leftRail.querySelectorAll(".tokyo-night-left-signs").forEach((node) => node.remove());

      const leftSigns = document.createElement("div");
      leftSigns.className = "tokyo-night-left-signs";

      [
        { label: "\u6E0B\u8C37\u99C5", caption: "SHIBUYA ST.", className: "jp-station" },
        { label: "\u30B3\u30F3\u30D3\u30CB", caption: "KONBINI", className: "jp-konbini" },
        { label: "\u30E9\u30FC\u30E1\u30F3", caption: "RAMEN", className: "jp-ramen" },
        { label: "\u30CA\u30DF", caption: "NAMI", className: "jp-nami" }
      ].forEach((config) => {
        const sign = document.createElement("span");
        sign.className = `tokyo-night-left-sign ${config.className}`;
        sign.textContent = config.label;
        sign.dataset.caption = config.caption;
        leftSigns.append(sign);
      });

      navPanel.insertAdjacentElement("afterend", leftSigns);
    }
    document.body.append(layer);
    activeLayer = layer;

    if (reduceMotionQuery.matches) {
      layer.classList.add("theme-effects-reduced-motion");

      for (let index = 0; index < 26; index++) {
        createStaticSakuraPetal(petalField);
      }

      return;
    }

    for (let index = 0; index < 42; index++) {
      spawnSakuraPetal(petalField, true);
    }

    petalTimer = window.setInterval(() => {
      const burstCount = Math.random() > 0.72 ? 3 : 1;

      for (let index = 0; index < burstCount; index++) {
        spawnSakuraPetal(petalField, false);
      }
    }, 260);

    lanternTimer = window.setInterval(() => {
      if (!isSakuraTheme(activeThemeKey) || !activeLayer) {
        return;
      }

      activeLayer.classList.toggle("sakura-lantern-breathe");
    }, 2400);
  }

  function createStaticSakuraPetal(parent) {
    const petal = document.createElement("span");
    petal.className = `sakura-petal ${randomItem(["petal-a", "petal-b", "petal-c"])}`;

    petal.style.setProperty("--petal-size", `${randomBetween(8, 18).toFixed(1)}px`);
    petal.style.setProperty("--petal-left", `${randomBetween(0, 100).toFixed(2)}vw`);
    petal.style.setProperty("--petal-top", `${randomBetween(0, 100).toFixed(2)}vh`);
    petal.style.setProperty("--petal-opacity", randomBetween(0.4, 0.95).toFixed(2));
    petal.style.setProperty("--petal-rotate", `${randomBetween(-70, 70).toFixed(1)}deg`);

    parent.append(petal);
  }

  function spawnSakuraPetal(parent, isInitialBurst) {
    if (!isSakuraTheme(activeThemeKey) || !parent?.isConnected) {
      return;
    }

    const petal = document.createElement("span");
    const petalType = randomItem(["petal-a", "petal-b", "petal-c", "petal-d"]);

    petal.className = `sakura-petal ${petalType}`;

    const duration = randomBetween(12, 26);
    const delay = isInitialBurst ? randomBetween(-duration, 0) : randomBetween(0, 1.2);
    const size = randomBetween(8, 22);
    const startLeft = randomBetween(-8, 108);
    const drift = randomBetween(-18, 18);
    const sway = randomBetween(18, 54);
    const fall = randomBetween(112, 142);
    const spin = randomBetween(-720, 720);
    const opacity = randomBetween(0.52, 1);

    petal.style.setProperty("--petal-size", `${size.toFixed(1)}px`);
    petal.style.setProperty("--petal-left", `${startLeft.toFixed(2)}vw`);
    petal.style.setProperty("--petal-drift", `${drift.toFixed(2)}vw`);
    petal.style.setProperty("--petal-sway", `${sway.toFixed(2)}px`);
    petal.style.setProperty("--petal-fall", `${fall.toFixed(2)}vh`);
    petal.style.setProperty("--petal-spin", `${spin.toFixed(1)}deg`);
    petal.style.setProperty("--petal-duration", `${duration.toFixed(2)}s`);
    petal.style.setProperty("--petal-delay", `${delay.toFixed(2)}s`);
    petal.style.setProperty("--petal-opacity", opacity.toFixed(2));
    petal.style.setProperty("--petal-blur", `${randomBetween(0, 0.8).toFixed(2)}px`);

    parent.append(petal);

    const cleanupDelay = Math.max(1200, (duration + Math.max(delay, 0) + 1) * 1000);
    window.setTimeout(() => petal.remove(), cleanupDelay);
  }


  function startTokyoNightThemeEffect() {
    document.documentElement.classList.add("theme-effect-tokyo-night");
    document.body?.classList.add("theme-effect-tokyo-night");

    const layer = document.createElement("div");
    layer.className = "theme-effects-layer theme-effects-tokyo-night";
    layer.setAttribute("aria-hidden", "true");

    const moon = document.createElement("div");
    moon.className = "tokyo-night-moon";

    const skyline = document.createElement("div");
    skyline.className = "tokyo-night-skyline";

    const signs = document.createElement("div");
    signs.className = "tokyo-night-signs";

    const vapor = document.createElement("div");
    vapor.className = "tokyo-night-vapor";

    const cityProps = document.createElement("div");
    cityProps.className = "tokyo-night-city-props";

        const antennaField = document.createElement("div");
    antennaField.className = "tokyo-night-antenna-field";

    [
      { x: "12", h: "70", d: "-0.6s" },
      { x: "26", h: "94", d: "-2.1s" },
      { x: "43", h: "78", d: "-1.2s" },
      { x: "61", h: "108", d: "-3.4s" },
      { x: "79", h: "86", d: "-1.8s" }
    ].forEach((config) => {
      const antenna = document.createElement("span");
      antenna.className = "tokyo-night-antenna";
      antenna.style.setProperty("--antenna-x", config.x);
      antenna.style.setProperty("--antenna-h", config.h);
      antenna.style.setProperty("--antenna-delay", config.d);
      antennaField.append(antenna);
    });

    const crosswalk = document.createElement("div");
    crosswalk.className = "tokyo-night-crosswalk";
    const commuterTrain = document.createElement("div");
    commuterTrain.className = "tokyo-night-commuter-train";

    const trainBody = document.createElement("div");
    trainBody.className = "tokyo-night-train-body";

    for (let index = 0; index < 18; index += 1) {
      const windowLight = document.createElement("span");
      windowLight.className = "tokyo-night-train-window";
      trainBody.append(windowLight);
    }

    commuterTrain.append(trainBody);

        const jpSignCluster = document.createElement("div");
    jpSignCluster.className = "tokyo-night-jp-sign-cluster";

    [
      { label: "\u6E0B\u8C37\u99C5", caption: "SHIBUYA ST.", className: "jp-station" },
      { label: "\u30B3\u30F3\u30D3\u30CB", caption: "KONBINI", className: "jp-konbini" },
      { label: "\u30E9\u30FC\u30E1\u30F3", caption: "RAMEN", className: "jp-ramen" },
      { label: "\u30CA\u30DF", caption: "NAMI", className: "jp-nami" }
    ].forEach((config) => {
      const sign = document.createElement("span");
      sign.className = `tokyo-night-jp-sign ${config.className}`;
      sign.textContent = config.label;
      sign.dataset.caption = config.caption;
      jpSignCluster.append(sign);
    });

    cityProps.append(antennaField, crosswalk, commuterTrain, jpSignCluster);

    const rainField = document.createElement("div");
    rainField.className = "tokyo-night-rain-field";

    if (!reduceMotionQuery.matches) {
      for (let index = 0; index < 90; index++) {
        const drop = document.createElement("span");
        drop.className = "tokyo-rain-drop";
        drop.style.setProperty("--rain-left", `${randomBetween(-5, 105).toFixed(2)}vw`);
        drop.style.setProperty("--rain-top", `${randomBetween(-20, 100).toFixed(2)}vh`);
        drop.style.setProperty("--rain-length", `${randomBetween(28, 82).toFixed(1)}px`);
        drop.style.setProperty("--rain-speed", `${randomBetween(0.9, 2.4).toFixed(2)}s`);
        drop.style.setProperty("--rain-delay", `${randomBetween(-3.5, 0).toFixed(2)}s`);
        drop.style.setProperty("--rain-opacity", randomBetween(0.12, 0.46).toFixed(2));
        rainField.append(drop);
      }
    } else {
      layer.classList.add("theme-effects-reduced-motion");
    }

        layer.append(moon, skyline, signs, vapor, cityProps, rainField);

            const rightRail = document.querySelector(".right-rail");
    const buffsPanel = document.querySelector(".right-buffs-panel");
    if (rightRail && buffsPanel) {
      rightRail.querySelectorAll(".tokyo-night-rail-signs").forEach((node) => node.remove());

      const railSigns = document.createElement("div");
      railSigns.className = "tokyo-night-rail-signs";

      [
        { label: "\u6771\u4EAC", caption: "TOKYO", className: "sign-cyan" },
        { label: "\u30CA\u30DF", caption: "NAMI", className: "sign-pink" },
        { label: "\u6DF1\u591C", caption: "MIDNIGHT", className: "sign-violet" },
        { label: "24\u6642", caption: "24H", className: "sign-gold" }
      ].forEach((config) => {
        const sign = document.createElement("span");
        sign.className = `tokyo-night-rail-sign ${config.className}`;
        sign.textContent = config.label;

        if (config.caption) {
          sign.dataset.caption = config.caption;
        }

        railSigns.append(sign);
      });

      buffsPanel.insertAdjacentElement("afterend", railSigns);
    }

    const leftRail = document.querySelector(".left-rail");
    const navPanel = document.querySelector(".nav-panel");

    if (leftRail && navPanel) {
      leftRail.querySelectorAll(".tokyo-night-left-signs").forEach((node) => node.remove());

      const leftSigns = document.createElement("div");
      leftSigns.className = "tokyo-night-left-signs";

      [
        { label: "\u6E0B\u8C37\u99C5", caption: "SHIBUYA ST.", className: "jp-station" },
        { label: "\u30B3\u30F3\u30D3\u30CB", caption: "KONBINI", className: "jp-konbini" },
        { label: "\u30E9\u30FC\u30E1\u30F3", caption: "RAMEN", className: "jp-ramen" },
        { label: "\u30CA\u30DF", caption: "NAMI", className: "jp-nami" }
      ].forEach((config) => {
        const sign = document.createElement("span");
        sign.className = `tokyo-night-left-sign ${config.className}`;
        sign.textContent = config.label;
        sign.dataset.caption = config.caption;
        leftSigns.append(sign);
      });

      navPanel.insertAdjacentElement("afterend", leftSigns);
    }
    document.body.append(layer);
    activeLayer = layer;
  }

  function startCandyThemeEffect() {
    document.documentElement.classList.add("theme-effect-candy");
    document.body?.classList.add("theme-effect-candy");

    const layer = document.createElement("div");
    layer.className = "theme-effects-layer theme-effects-candy";
    layer.setAttribute("aria-hidden", "true");

    const candyCloud = document.createElement("div");
    candyCloud.className = "candy-cloud";

    const candyShopAwning = document.createElement("div");
    candyShopAwning.className = "candy-shop-awning";

    const lollipopLeft = document.createElement("div");
    lollipopLeft.className = "candy-lollipop candy-lollipop-left";

    const lollipopRight = document.createElement("div");
    lollipopRight.className = "candy-lollipop candy-lollipop-right";

    const chocolateDrip = document.createElement("div");
    chocolateDrip.className = "candy-chocolate-drip";

    const gumdropRow = document.createElement("div");
    gumdropRow.className = "candy-gumdrop-row";

    for (let index = 0; index < 18; index += 1) {
      const gumdrop = document.createElement("span");
      gumdrop.className = `candy-gumdrop gumdrop-${(index % 6) + 1}`;
      gumdropRow.append(gumdrop);
    }

    const sprinkleField = document.createElement("div");
    sprinkleField.className = "candy-sprinkle-field";

    if (!reduceMotionQuery.matches) {
      for (let index = 0; index < 95; index += 1) {
        const sprinkle = document.createElement("span");
        sprinkle.className = `candy-sprinkle sprinkle-${(index % 7) + 1}`;
        sprinkle.style.setProperty("--sprinkle-left", `${randomBetween(-4, 104).toFixed(2)}vw`);
        sprinkle.style.setProperty("--sprinkle-top", `${randomBetween(-12, 104).toFixed(2)}vh`);
        sprinkle.style.setProperty("--sprinkle-size", `${randomBetween(5, 13).toFixed(1)}px`);
        sprinkle.style.setProperty("--sprinkle-speed", `${randomBetween(5.5, 13).toFixed(2)}s`);
        sprinkle.style.setProperty("--sprinkle-delay", `${randomBetween(-12, 0).toFixed(2)}s`);
        sprinkle.style.setProperty("--sprinkle-rotate", `${randomBetween(-120, 120).toFixed(1)}deg`);
        sprinkle.style.setProperty("--sprinkle-opacity", randomBetween(0.26, 0.72).toFixed(2));
        sprinkleField.append(sprinkle);
      }
    } else {
      layer.classList.add("theme-effects-reduced-motion");
    }

    layer.append(candyCloud, candyShopAwning, lollipopLeft, lollipopRight, chocolateDrip, gumdropRow, sprinkleField);
    document.body.append(layer);
    activeLayer = layer;
  }
  window.NamigotchiThemeEffects = {
    setActiveThemeEffect,
    clearThemeEffects,
  };
})();
// Candy pass 4: floating candy background
(() => {
  const LAYER_ID = "candy-float-layer";
  const TYPES = ["wrapped", "lollipop", "peppermint", "gumdrop", "truffle"];
  const PIECE_COUNT = 16;

  function getCurrentThemeName() {
    return document.body?.dataset?.theme || "";
  }

  function isCandyTheme() {
    return getCurrentThemeName() === "candy";
  }

  function removeCandyLayer() {
    const existing = document.getElementById(LAYER_ID);
    if (existing) existing.remove();
  }

  function buildCandyPiece(type) {
    const piece = document.createElement("div");
    piece.className = `candy-piece candy-${type}`;

    const size = 56 + Math.floor(Math.random() * 54);
    const left = Math.random() * 94;
    const top = Math.random() * 88;
    const opacity = 0.16 + Math.random() * 0.20;
    const driftX = -18 + Math.random() * 36;
    const driftY = -24 + Math.random() * 48;
    const spin = 34 + Math.random() * 34;
    const float = 24 + Math.random() * 18;
    const rot = Math.floor(Math.random() * 360);

    piece.style.setProperty("--candy-size", `${size}px`);
    piece.style.setProperty("--candy-opacity", opacity.toFixed(2));
    piece.style.setProperty("--candy-drift-x", `${driftX}px`);
    piece.style.setProperty("--candy-drift-y", `${driftY}px`);
    piece.style.setProperty("--candy-spin", `${spin.toFixed(1)}s`);
    piece.style.setProperty("--candy-float", `${float.toFixed(1)}s`);
    piece.style.setProperty("--candy-rot", `${rot}deg`);
    piece.style.left = `${left}%`;
    piece.style.top = `${top}%`;
    piece.style.animationDelay = `${(Math.random() * -30).toFixed(1)}s, ${(Math.random() * -30).toFixed(1)}s`;

    if (type === "wrapped") {
      piece.innerHTML = `
        <span class="candy-tail candy-tail-left"></span>
        <span class="candy-center"></span>
        <span class="candy-tail candy-tail-right"></span>
      `;
    } else if (type === "lollipop") {
      piece.innerHTML = `
        <span class="candy-disc"></span>
        <span class="candy-stick"></span>
      `;
    } else if (type === "peppermint") {
      piece.innerHTML = `<span class="candy-disc"></span>`;
    } else if (type === "gumdrop") {
      piece.innerHTML = `<span class="candy-body"></span>`;
    } else if (type === "truffle") {
      piece.innerHTML = `
        <span class="candy-base"></span>
        <span class="candy-top"></span>
      `;
    }

    return piece;
  }

  function createCandyLayer() {
    if (!document.body || document.getElementById(LAYER_ID) || !isCandyTheme()) return;

    const layer = document.createElement("div");
    layer.id = LAYER_ID;
    layer.className = "candy-float-layer";

    for (let i = 0; i < PIECE_COUNT; i += 1) {
      const type = TYPES[Math.floor(Math.random() * TYPES.length)];
      layer.appendChild(buildCandyPiece(type));
    }

    document.body.appendChild(layer);
  }

  function syncCandyLayer() {
    if (isCandyTheme()) {
      createCandyLayer();
    } else {
      removeCandyLayer();
    }
  }

  function installCandyWatcher() {
    if (!document.body) return;

    syncCandyLayer();

    const observer = new MutationObserver(() => {
      syncCandyLayer();
    });

    observer.observe(document.body, {
      attributes: true,
      attributeFilter: ["data-theme", "class"]
    });

    window.addEventListener("pageshow", syncCandyLayer);
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", installCandyWatcher, { once: true });
  } else {
    installCandyWatcher();
  }
})();