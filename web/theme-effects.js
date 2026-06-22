(() => {
  const reduceMotionQuery = window.matchMedia("(prefers-reduced-motion: reduce)");

  let activeThemeKey = "";
  let activeLayer = null;
  let petalTimer = null;
  let lanternTimer = null;
  let pearlTideBubbleAnimationFrame = 0;
  let pearlTideBubbleResizeListener = null;

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
  function isCafeTheme(themeKey) {
    return themeKey === "cafe";
  }

  function isRainyMoodTheme(themeKey) {
    return themeKey === "rainy-mood";
  }

  function isPearlTideTheme(themeKey) {
    return themeKey === "pearl-tide";
  }

  function isWoodlandSunTheme(themeKey) {
    return themeKey === "woodland-sun";
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

    if (pearlTideBubbleAnimationFrame) {
      window.cancelAnimationFrame(pearlTideBubbleAnimationFrame);
      pearlTideBubbleAnimationFrame = 0;
    }

    if (pearlTideBubbleResizeListener) {
      window.removeEventListener("resize", pearlTideBubbleResizeListener);
      pearlTideBubbleResizeListener = null;
    }

    if (activeLayer) {
      activeLayer.remove();
      activeLayer = null;
    }

    document.documentElement.classList.remove("theme-effect-sakura", "theme-effect-tokyo-night", "theme-effect-candy", "theme-effect-cafe", "theme-effect-rainy-mood", "theme-effect-pearl-tide", "theme-effect-woodland-sun");
    document.body?.classList.remove("theme-effect-sakura", "theme-effect-tokyo-night", "theme-effect-candy", "theme-effect-cafe", "theme-effect-rainy-mood", "theme-effect-pearl-tide", "theme-effect-woodland-sun");
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
    if (isCafeTheme(themeKey)) {
      startCafeThemeEffect();
    }

    if (isRainyMoodTheme(themeKey)) {
      startRainyMoodThemeEffect();
    }

    if (isPearlTideTheme(themeKey)) {
      startPearlTideThemeEffect();
    }

    if (isWoodlandSunTheme(themeKey)) {
      startWoodlandSunThemeEffect();
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

  function startCafeThemeEffect() {
    document.documentElement.classList.add("theme-effect-cafe");
    document.body?.classList.add("theme-effect-cafe");

    const layer = document.createElement("div");
    layer.className = "theme-effects-layer theme-effects-cafe cafe-float-layer";
    layer.setAttribute("aria-hidden", "true");

    const types = ["cup", "croissant", "macaron", "donut", "beans"];
    const piecesPerZone = 8;

    function buildCafePiece(type) {
      const piece = document.createElement("div");
      piece.className = `cafe-piece cafe-${type}`;

      const size = randomBetween(58, 102);
      const left = randomBetween(5, 78);
      const top = randomBetween(4, 82);
      const opacity = randomBetween(0.34, 0.52);
      const driftX = randomBetween(-10, 10);
      const driftY = randomBetween(-16, 16);
      const spin = randomBetween(58, 106);
      const float = randomBetween(30, 48);
      const rot = randomBetween(0, 360);

      piece.style.setProperty("--cafe-size", `${size.toFixed(1)}px`);
      piece.style.setProperty("--cafe-opacity", opacity.toFixed(2));
      piece.style.setProperty("--cafe-drift-x", `${driftX.toFixed(1)}px`);
      piece.style.setProperty("--cafe-drift-y", `${driftY.toFixed(1)}px`);
      piece.style.setProperty("--cafe-spin", `${spin.toFixed(1)}s`);
      piece.style.setProperty("--cafe-float", `${float.toFixed(1)}s`);
      piece.style.setProperty("--cafe-rot", `${rot.toFixed(1)}deg`);
      piece.style.left = `${left.toFixed(1)}%`;
      piece.style.top = `${top.toFixed(1)}%`;
      piece.style.animationDelay = `${randomBetween(-42, 0).toFixed(1)}s, ${randomBetween(-42, 0).toFixed(1)}s`;

      if (type === "cup") {
        piece.innerHTML = `
          <span class="cup-steam"></span>
          <span class="cup-body"></span>
          <span class="cup-handle"></span>
          <span class="cup-saucer"></span>
        `;
      } else if (type === "croissant") {
        piece.innerHTML = `
          <span class="croissant-base"></span>
          <span class="croissant-lines"></span>
        `;
      } else if (type === "macaron") {
        piece.innerHTML = `
          <span class="macaron-top"></span>
          <span class="macaron-cream"></span>
          <span class="macaron-bottom"></span>
        `;
      } else if (type === "donut") {
        piece.innerHTML = `
          <span class="donut-ring"></span>
          <span class="donut-glaze"></span>
        `;
      } else if (type === "beans") {
        piece.innerHTML = `
          <span class="bean bean-1"></span>
          <span class="bean bean-2"></span>
        `;
      }

      return piece;
    }

    function buildCafeZone(zoneName, railElement, panelElement) {
      const railRect = railElement.getBoundingClientRect();
      const panelRect = panelElement.getBoundingClientRect();

      const zone = document.createElement("div");
      zone.className = `cafe-float-zone cafe-float-zone-${zoneName}`;

      const top = Math.max(panelRect.bottom + 12, railRect.top + 48);
      const height = Math.max(window.innerHeight - top - 8, 150);

      zone.style.left = `${Math.max(0, railRect.left).toFixed(1)}px`;
      zone.style.top = `${top.toFixed(1)}px`;
      zone.style.width = `${Math.max(railRect.width, 38).toFixed(1)}px`;
      zone.style.height = `${height.toFixed(1)}px`;

      for (let index = 0; index < piecesPerZone; index += 1) {
        const type = types[(index + Math.floor(randomBetween(0, types.length))) % types.length];
        zone.append(buildCafePiece(type));
      }

      return zone;
    }

    const leftRail = document.querySelector(".left-rail");
    const rightRail = document.querySelector(".right-rail");
    const leftPanel = document.querySelector(".left-rail > .nav-panel") || leftRail;
    const rightPanel = document.querySelector(".right-rail > .right-buffs-panel") || rightRail;

    if (leftRail && rightRail && leftPanel && rightPanel) {
      layer.append(
        buildCafeZone("left", leftRail, leftPanel),
        buildCafeZone("right", rightRail, rightPanel)
      );
    }

    document.body.append(layer);
    activeLayer = layer;
  }





  /* Woodland Sun forest light effect START */
  function startWoodlandSunThemeEffect() {
    document.documentElement.classList.add("theme-effect-woodland-sun");
    document.body?.classList.add("theme-effect-woodland-sun");

    const layer = document.createElement("div");
    layer.className = "theme-effects-layer theme-effects-woodland-sun";
    layer.setAttribute("aria-hidden", "true");

    const atmosphere = document.createElement("div");
    atmosphere.className = "woodland-sun-atmosphere";

    const shaftField = document.createElement("div");
    shaftField.className = "woodland-sun-shaft-field";

    const bokehField = document.createElement("div");
    bokehField.className = "woodland-sun-bokeh-field";

    layer.append(atmosphere, shaftField, bokehField);
    document.body.append(layer);
    activeLayer = layer;

    const reducedMotion = reduceMotionQuery.matches;
    const viewportWidth = Math.max(1, window.innerWidth || document.documentElement.clientWidth || 1);
    const viewportHeight = Math.max(1, window.innerHeight || document.documentElement.clientHeight || 1);
    const viewportArea = viewportWidth * viewportHeight;

    if (reducedMotion) {
      layer.classList.add("theme-effects-reduced-motion");
    }

    function configureLightShaft(shaft, isInitial) {
      const duration = randomBetween(26, 52);
      const width = randomBetween(12, 31);
      const height = randomBetween(92, 172);
      const left = randomBetween(-22, 94);
      const top = randomBetween(-24, 18);
      const rotation = randomBetween(-19, 15);
      const drift = randomBetween(-2.8, 3.6);
      const strength = randomBetween(0.34, 0.62);
      const blur = randomBetween(6, 15);

      shaft.style.setProperty("--shaft-left", `${left.toFixed(2)}vw`);
      shaft.style.setProperty("--shaft-top", `${top.toFixed(2)}vh`);
      shaft.style.setProperty("--shaft-width", `${width.toFixed(2)}vw`);
      shaft.style.setProperty("--shaft-height", `${height.toFixed(2)}vh`);
      shaft.style.setProperty("--shaft-rotation", `${rotation.toFixed(2)}deg`);
      shaft.style.setProperty("--shaft-drift", `${drift.toFixed(2)}vw`);
      shaft.style.setProperty("--shaft-strength", strength.toFixed(3));
      shaft.style.setProperty("--shaft-blur", `${blur.toFixed(2)}px`);
      shaft.style.setProperty("--shaft-warmth", randomBetween(0.22, 0.38).toFixed(3));
      shaft.style.animationDuration = `${duration.toFixed(2)}s`;
      shaft.style.animationDelay = isInitial
        ? `${randomBetween(-duration * 0.82, 0).toFixed(2)}s`
        : "0s";
    }

    function configureBokehLight(bokeh, isInitial) {
      const duration = randomBetween(18, 42);
      const size = randomBetween(34, 128);
      const left = randomBetween(-8, 104);
      const top = randomBetween(4, 98);
      const driftX = randomBetween(-34, 34);
      const driftY = randomBetween(-26, 20);
      const alpha = randomBetween(0.26, 0.52);
      const core = Math.min(alpha + randomBetween(0.16, 0.32), 0.78);
      const blur = randomBetween(0.8, 4.8);

      bokeh.style.setProperty("--bokeh-size", `${size.toFixed(2)}px`);
      bokeh.style.setProperty("--bokeh-left", `${left.toFixed(2)}vw`);
      bokeh.style.setProperty("--bokeh-top", `${top.toFixed(2)}vh`);
      bokeh.style.setProperty("--bokeh-drift-x", `${driftX.toFixed(2)}px`);
      bokeh.style.setProperty("--bokeh-drift-y", `${driftY.toFixed(2)}px`);
      bokeh.style.setProperty("--bokeh-alpha", alpha.toFixed(3));
      bokeh.style.setProperty("--bokeh-core", core.toFixed(3));
      bokeh.style.setProperty("--bokeh-blur", `${blur.toFixed(2)}px`);
      bokeh.style.setProperty("--bokeh-scale", randomBetween(0.74, 1.36).toFixed(3));
      bokeh.style.animationDuration = `${duration.toFixed(2)}s`;
      bokeh.style.animationDelay = isInitial
        ? `${randomBetween(-duration, 0).toFixed(2)}s`
        : "0s";
    }

    const shaftCount = reducedMotion
      ? 7
      : Math.max(10, Math.min(18, Math.round(viewportWidth / 250)));

    for (let index = 0; index < shaftCount; index += 1) {
      const shaft = document.createElement("span");
      shaft.className = "woodland-sun-light-shaft";
      configureLightShaft(shaft, true);

      if (!reducedMotion) {
        shaft.addEventListener("animationiteration", () => {
          if (isWoodlandSunTheme(activeThemeKey) && shaft.isConnected) {
            configureLightShaft(shaft, false);
          }
        });
      }

      shaftField.append(shaft);
    }

    const bokehCount = reducedMotion
      ? Math.max(22, Math.min(40, Math.round(viewportArea / 160000)))
      : Math.max(54, Math.min(118, Math.round(viewportArea / 62000)));

    for (let index = 0; index < bokehCount; index += 1) {
      const bokeh = document.createElement("span");
      bokeh.className = `woodland-sun-bokeh bokeh-${(index % 5) + 1}`;
      configureBokehLight(bokeh, true);

      if (!reducedMotion) {
        bokeh.addEventListener("animationiteration", () => {
          if (isWoodlandSunTheme(activeThemeKey) && bokeh.isConnected) {
            configureBokehLight(bokeh, false);
          }
        });
      }

      bokehField.append(bokeh);
    }
  }
  /* Woodland Sun forest light effect END */

  /* Pearl Tide bubble canvas effect START */
  function startPearlTideThemeEffect() {
    document.documentElement.classList.add("theme-effect-pearl-tide");
    document.body?.classList.add("theme-effect-pearl-tide");

    const layer = document.createElement("div");
    layer.className = "theme-effects-layer theme-effects-pearl-tide";
    layer.setAttribute("aria-hidden", "true");

    const canvas = document.createElement("canvas");
    canvas.className = "pearl-tide-bubble-canvas";
    layer.append(canvas);

    document.body.append(layer);
    activeLayer = layer;

    const ctx = canvas.getContext("2d", { alpha: true });
    const backgroundCanvas = document.createElement("canvas");
    const backgroundCtx = backgroundCanvas.getContext("2d", { alpha: false });

    if (!ctx || !backgroundCtx) {
      return;
    }

    const backgroundImage = new Image();
    backgroundImage.src = "/images/Nami_Splash.webp";

    const reducedMotion = reduceMotionQuery.matches;
    const bubbles = [];
    let viewportWidth = 0;
    let viewportHeight = 0;
    let pixelRatio = 1;
    let backgroundReady = false;
    let coverRect = { dx: 0, dy: 0, dw: 1, dh: 1 };

    function calculateCoverRect(imageWidth, imageHeight, targetWidth, targetHeight) {
      const scale = Math.max(targetWidth / imageWidth, targetHeight / imageHeight);
      const drawWidth = imageWidth * scale;
      const drawHeight = imageHeight * scale;

      return {
        dx: (targetWidth - drawWidth) / 2,
        dy: (targetHeight - drawHeight) / 2,
        dw: drawWidth,
        dh: drawHeight,
      };
    }

    function redrawBackgroundBuffer() {
      backgroundCtx.setTransform(1, 0, 0, 1, 0, 0);
      backgroundCtx.clearRect(0, 0, viewportWidth, viewportHeight);

      if (!backgroundReady) {
        const fallback = backgroundCtx.createLinearGradient(0, 0, 0, viewportHeight);
        fallback.addColorStop(0, "#06182c");
        fallback.addColorStop(1, "#08243a");
        backgroundCtx.fillStyle = fallback;
        backgroundCtx.fillRect(0, 0, viewportWidth, viewportHeight);
        return;
      }

      backgroundCtx.drawImage(
        backgroundImage,
        coverRect.dx,
        coverRect.dy,
        coverRect.dw,
        coverRect.dh
      );
    }

    function drawBackgroundToVisibleCanvas() {
      ctx.drawImage(backgroundCanvas, 0, 0, viewportWidth, viewportHeight);
    }

    function makeBubble(isInitial) {
      const radiusScale = Math.max(0.85, Math.min(1.22, viewportWidth / 1500));
      const radius = randomBetween(9, 32) * radiusScale;
      const originX = randomBetween(-radius, viewportWidth + radius);

      return {
        originX,
        x: originX,
        y: isInitial
          ? randomBetween(-radius, viewportHeight + radius)
          : viewportHeight + radius + randomBetween(0, viewportHeight * 0.18),
        radius,
        heightScale: randomBetween(0.74, 1.04),
        lift: randomBetween(0.16, 0.58) * radiusScale,
        drift: randomBetween(10, 38) * radiusScale,
        wobbleSpeed: randomBetween(0.0008, 0.0023),
        wobbleOffset: randomBetween(0, Math.PI * 2),
        rotation: randomBetween(0, Math.PI * 2),
        rotationSpeed: randomBetween(-0.003, 0.003),
        opacity: randomBetween(0.28, 0.58),
        lensPower: randomBetween(0.16, 0.30),
      };
    }

    function replaceBubble(targetBubble) {
      const replacement = makeBubble(false);

      targetBubble.originX = replacement.originX;
      targetBubble.x = replacement.x;
      targetBubble.y = replacement.y;
      targetBubble.radius = replacement.radius;
      targetBubble.heightScale = replacement.heightScale;
      targetBubble.lift = replacement.lift;
      targetBubble.drift = replacement.drift;
      targetBubble.wobbleSpeed = replacement.wobbleSpeed;
      targetBubble.wobbleOffset = replacement.wobbleOffset;
      targetBubble.rotation = replacement.rotation;
      targetBubble.rotationSpeed = replacement.rotationSpeed;
      targetBubble.opacity = replacement.opacity;
      targetBubble.lensPower = replacement.lensPower;
    }

    function resizeCanvas() {
      viewportWidth = Math.max(1, window.innerWidth || document.documentElement.clientWidth || 1);
      viewportHeight = Math.max(1, window.innerHeight || document.documentElement.clientHeight || 1);
      pixelRatio = Math.min(window.devicePixelRatio || 1, 2);

      canvas.width = Math.round(viewportWidth * pixelRatio);
      canvas.height = Math.round(viewportHeight * pixelRatio);
      canvas.style.width = `${viewportWidth}px`;
      canvas.style.height = `${viewportHeight}px`;

      backgroundCanvas.width = viewportWidth;
      backgroundCanvas.height = viewportHeight;

      ctx.setTransform(pixelRatio, 0, 0, pixelRatio, 0, 0);

      if (backgroundReady) {
        coverRect = calculateCoverRect(
          backgroundImage.naturalWidth || backgroundImage.width,
          backgroundImage.naturalHeight || backgroundImage.height,
          viewportWidth,
          viewportHeight
        );
      }

      redrawBackgroundBuffer();

      bubbles.length = 0;

      const viewportArea = viewportWidth * viewportHeight;
      const targetCount = reducedMotion
        ? Math.max(30, Math.min(60, Math.round(viewportArea / 140000)))
        : Math.max(50, Math.min(180, Math.round(viewportArea / 52000)));

      for (let index = 0; index < targetCount; index += 1) {
        bubbles.push(makeBubble(true));
      }

      drawFrame(performance.now(), true);
    }

    function drawClippedBackgroundSlice(sourceX, sourceY, sourceWidth, sourceHeight, destX, destY, destWidth, destHeight) {
      const sourceRight = sourceX + sourceWidth;
      const sourceBottom = sourceY + sourceHeight;

      const clippedSourceX = Math.max(0, sourceX);
      const clippedSourceY = Math.max(0, sourceY);
      const clippedSourceRight = Math.min(viewportWidth, sourceRight);
      const clippedSourceBottom = Math.min(viewportHeight, sourceBottom);

      const clippedSourceWidth = clippedSourceRight - clippedSourceX;
      const clippedSourceHeight = clippedSourceBottom - clippedSourceY;

      if (clippedSourceWidth <= 0 || clippedSourceHeight <= 0 || sourceWidth <= 0 || sourceHeight <= 0) {
        return;
      }

      const leftRatio = (clippedSourceX - sourceX) / sourceWidth;
      const topRatio = (clippedSourceY - sourceY) / sourceHeight;
      const rightRatio = (sourceRight - clippedSourceRight) / sourceWidth;
      const bottomRatio = (sourceBottom - clippedSourceBottom) / sourceHeight;

      const clippedDestX = destX + destWidth * leftRatio;
      const clippedDestY = destY + destHeight * topRatio;
      const clippedDestWidth = destWidth * Math.max(0, 1 - leftRatio - rightRatio);
      const clippedDestHeight = destHeight * Math.max(0, 1 - topRatio - bottomRatio);

      if (clippedDestWidth <= 0 || clippedDestHeight <= 0) {
        return;
      }

      ctx.drawImage(
        backgroundCanvas,
        clippedSourceX,
        clippedSourceY,
        clippedSourceWidth,
        clippedSourceHeight,
        clippedDestX,
        clippedDestY,
        clippedDestWidth,
        clippedDestHeight
      );
    }

    function drawLensBackground(bubble, metrics) {
      if (!backgroundReady || !metrics) {
        return;
      }

      const { x, y, radiusX, radiusY } = metrics;
      const sliceHeight = 2.25;
      const lensPower = bubble.lensPower;

      ctx.save();

      ctx.translate(x, y);
      ctx.rotate(bubble.rotation);
      ctx.beginPath();
      ctx.ellipse(0, 0, radiusX, radiusY, 0, 0, Math.PI * 2);
      ctx.clip();

      ctx.rotate(-bubble.rotation);
      ctx.translate(-x, -y);

      for (let localY = -radiusY; localY <= radiusY; localY += sliceHeight) {
        const normalizedY = localY / radiusY;
        const chord = Math.sqrt(Math.max(0, 1 - normalizedY * normalizedY));
        const destHalfWidth = radiusX * chord;

        if (destHalfWidth <= 0.5) {
          continue;
        }

        const centerBoost = 1 - Math.abs(normalizedY);
        const horizontalScale = 1 + lensPower * (0.45 + centerBoost * 1.10);
        const verticalScale = 1 + lensPower * 0.46;

        const sourceHalfWidth = destHalfWidth / horizontalScale;
        const sourceHeight = sliceHeight / verticalScale;

        const sourceX = x - sourceHalfWidth;
        const sourceY = y + localY / verticalScale;
        const sourceWidth = sourceHalfWidth * 2;

        drawClippedBackgroundSlice(
          sourceX,
          sourceY,
          sourceWidth,
          sourceHeight,
          x - destHalfWidth,
          y + localY,
          destHalfWidth * 2,
          sliceHeight
        );
      }

      ctx.restore();
    }

    function drawBubbleShell(bubble, metrics) {
      if (!metrics) {
        return;
      }

      const { x, y, radiusX, radiusY } = metrics;

      const shellGradient = ctx.createRadialGradient(
        x - radiusX * 0.30,
        y - radiusY * 0.34,
        0,
        x,
        y,
        Math.max(radiusX, radiusY) * 1.22
      );

      shellGradient.addColorStop(0, `rgba(255, 255, 255, ${Math.min(0.22, bubble.opacity + 0.06).toFixed(3)})`);
      shellGradient.addColorStop(0.42, `rgba(151, 236, 246, ${(bubble.opacity * 0.22).toFixed(3)})`);
      shellGradient.addColorStop(1, "rgba(151, 236, 246, 0)");

      ctx.save();
      ctx.translate(x, y);
      ctx.rotate(bubble.rotation);

      ctx.beginPath();
      ctx.ellipse(0, 0, radiusX, radiusY, 0, 0, Math.PI * 2);
      ctx.fillStyle = shellGradient;
      ctx.fill();

      ctx.lineWidth = Math.max(1, bubble.radius * 0.065);
      ctx.strokeStyle = `rgba(234, 252, 255, ${Math.min(0.48, bubble.opacity + 0.10).toFixed(3)})`;
      ctx.stroke();

      ctx.beginPath();
      ctx.ellipse(-radiusX * 0.30, -radiusY * 0.34, radiusX * 0.18, radiusY * 0.10, -0.52, 0, Math.PI * 2);
      ctx.fillStyle = `rgba(255, 255, 255, ${Math.min(0.48, bubble.opacity + 0.18).toFixed(3)})`;
      ctx.fill();

      ctx.beginPath();
      ctx.ellipse(radiusX * 0.28, radiusY * 0.32, radiusX * 0.20, radiusY * 0.12, -0.34, 0, Math.PI * 2);
      ctx.fillStyle = "rgba(110, 233, 245, 0.055)";
      ctx.fill();

      ctx.restore();
    }

    function getBubbleMetrics(bubble, now) {
      const wobble = Math.sin(now * bubble.wobbleSpeed + bubble.wobbleOffset);
      const pulse = Math.cos(now * bubble.wobbleSpeed * 1.8 + bubble.wobbleOffset);
      const x = bubble.originX + wobble * bubble.drift;
      const y = bubble.y;
      const radiusX = bubble.radius * (0.88 + pulse * 0.055);
      const radiusY = bubble.radius * bubble.heightScale * (1.04 - pulse * 0.04);

      bubble.x = x;

      return { x, y, radiusX, radiusY };
    }

    function drawFrame(now, drawOnce = false) {
      if (!isPearlTideTheme(activeThemeKey) || !canvas.isConnected) {
        return;
      }

      ctx.clearRect(0, 0, viewportWidth, viewportHeight);
      drawBackgroundToVisibleCanvas();

      for (const bubble of bubbles) {
        if (!drawOnce && !reducedMotion) {
          bubble.y -= bubble.lift;
          bubble.rotation += bubble.rotationSpeed;

          if (bubble.y < -bubble.radius * 2.6) {
            replaceBubble(bubble);
          }
        }

        const metrics = getBubbleMetrics(bubble, now);
        drawLensBackground(bubble, metrics);
        drawBubbleShell(bubble, metrics);
      }

      if (!drawOnce && !reducedMotion) {
        pearlTideBubbleAnimationFrame = window.requestAnimationFrame(drawFrame);
      }
    }

    pearlTideBubbleResizeListener = () => resizeCanvas();
    window.addEventListener("resize", pearlTideBubbleResizeListener, { passive: true });

    backgroundImage.addEventListener("load", () => {
      backgroundReady = true;
      coverRect = calculateCoverRect(
        backgroundImage.naturalWidth || backgroundImage.width,
        backgroundImage.naturalHeight || backgroundImage.height,
        viewportWidth,
        viewportHeight
      );
      redrawBackgroundBuffer();
      drawFrame(performance.now(), true);
    }, { once: true });

    backgroundImage.addEventListener("error", () => {
      backgroundReady = false;
      redrawBackgroundBuffer();
      drawFrame(performance.now(), true);
    }, { once: true });

    resizeCanvas();

    if (!reducedMotion) {
      pearlTideBubbleAnimationFrame = window.requestAnimationFrame(drawFrame);
    } else {
      layer.classList.add("theme-effects-reduced-motion");
    }
  }
  /* Pearl Tide bubble canvas effect END */

  function startRainyMoodThemeEffect() {
    document.documentElement.classList.add("theme-effect-rainy-mood");
    document.body?.classList.add("theme-effect-rainy-mood");

    const layer = document.createElement("div");
    layer.className = "theme-effects-layer theme-effects-rainy-mood";
    layer.setAttribute("aria-hidden", "true");

    const video = document.createElement("video");
    video.className = "rainy-mood-video-bg";
    video.src = "/images/rainy_mood.webm";
    video.autoplay = true;
    video.muted = true;
    video.loop = true;
    video.playsInline = true;
    video.preload = "auto";

    const veil = document.createElement("div");
    veil.className = "rainy-mood-video-veil";

    layer.append(video, veil);
    document.body.append(layer);
    activeLayer = layer;

    const playResult = video.play?.();

    if (playResult?.catch) {
      playResult.catch(() => {
        layer.classList.add("rainy-mood-video-paused");
      });
    }
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

