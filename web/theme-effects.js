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

    document.documentElement.classList.remove("theme-effect-sakura");
    document.body?.classList.remove("theme-effect-sakura");
  }

  function setActiveThemeEffect(themeKey) {
    activeThemeKey = themeKey;

    document.documentElement.dataset.theme = themeKey;
    document.body?.setAttribute("data-theme", themeKey);

    clearThemeEffects();

    if (themeKey === "sakura-dark" || themeKey === "sakura-light" || themeKey === "sakura") {
      startSakuraThemeEffect();
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

    layer.append(moon, branch, lanternLeft, lanternRight, petalField);
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
      if (activeThemeKey !== "sakura" || !activeLayer) {
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
    if (activeThemeKey !== "sakura" || !parent?.isConnected) {
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

  window.NamigotchiThemeEffects = {
    setActiveThemeEffect,
    clearThemeEffects,
  };
})();