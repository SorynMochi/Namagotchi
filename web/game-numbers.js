(() => {
  "use strict";

  const SUFFIXES = [
    "", "K", "M", "B", "T", "Qa", "Qi", "Sx", "Sp", "Oc", "No", "Dc",
    "Ud", "Dd", "Td", "Qad", "Qid", "Sxd", "Spd", "Ocd", "Nod", "Vg"
  ];

  function normalizeGameIntegerInput(value, fallback = 0n) {
    if (typeof value === "bigint") {
      return value;
    }

    if (typeof value === "number") {
      if (!Number.isFinite(value)) {
        return fallback;
      }

      return BigInt(Math.trunc(value));
    }

    if (typeof value === "string") {
      const trimmed = value.trim();
      if (!trimmed) {
        return fallback;
      }

      if (/^[+-]?\d+$/.test(trimmed)) {
        return BigInt(trimmed);
      }

      if (/^[+-]?\d+(\.\d+)?$/.test(trimmed)) {
        return BigInt(trimmed.split(".")[0]);
      }

      return fallback;
    }

    if (value && typeof value === "object" && typeof value.value === "string") {
      return normalizeGameIntegerInput(value.value, fallback);
    }

    return fallback;
  }

  function addCommas(value) {
    const text = String(value);
    const sign = text.startsWith("-") ? "-" : "";
    const digits = sign ? text.slice(1) : text;
    return sign + digits.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  }

  function suffixForTier(tier) {
    if (tier < SUFFIXES.length) {
      return SUFFIXES[tier];
    }

    return `e${tier * 3}`;
  }

  function formatGameInteger(value, options = {}) {
    const precision = Math.max(0, Math.min(3, Number(options.precision ?? 1)));
    const compactAt = normalizeGameIntegerInput(options.compactAt ?? 1000n, 1000n);
    const integer = normalizeGameIntegerInput(value);
    const isNegative = integer < 0n;
    let absolute = isNegative ? -integer : integer;

    if (absolute < compactAt) {
      return `${isNegative ? "-" : ""}${addCommas(absolute.toString())}`;
    }

    const digits = absolute.toString();
    const tier = Math.floor((digits.length - 1) / 3);
    const suffix = suffixForTier(tier);

    if (tier <= 0) {
      return `${isNegative ? "-" : ""}${addCommas(digits)}`;
    }

    const divisor = 10n ** BigInt(tier * 3);
    let whole = absolute / divisor;
    let remainder = absolute % divisor;

    if (precision === 0) {
      if (remainder * 2n >= divisor) {
        whole += 1n;
      }

      if (whole >= 1000n && suffixForTier(tier + 1) !== suffix) {
        return `${isNegative ? "-" : ""}1${suffixForTier(tier + 1)}`;
      }

      return `${isNegative ? "-" : ""}${whole.toString()}${suffix}`;
    }

    const scale = 10n ** BigInt(precision);
    let decimal = (remainder * scale + divisor / 2n) / divisor;

    if (decimal >= scale) {
      whole += 1n;
      decimal -= scale;
    }

    if (whole >= 1000n && suffixForTier(tier + 1) !== suffix) {
      whole = 1n;
      decimal = 0n;
      return `${isNegative ? "-" : ""}${whole.toString()}${precision > 0 ? "." + "0".repeat(precision) : ""}${suffixForTier(tier + 1)}`;
    }

    const decimalText = decimal.toString().padStart(precision, "0").replace(/0+$/, "");
    const formatted = decimalText ? `${whole}.${decimalText}` : whole.toString();

    return `${isNegative ? "-" : ""}${formatted}${suffix}`;
  }

  function formatFullGameInteger(value) {
    return addCommas(normalizeGameIntegerInput(value).toString());
  }

  function formatGameIntegerRatio(current, maximum, options = {}) {
    return `${formatGameInteger(current, options)} / ${formatGameInteger(maximum, options)}`;
  }

  function gameIntegerPercent(current, maximum) {
    const numerator = normalizeGameIntegerInput(current);
    const denominator = normalizeGameIntegerInput(maximum);

    if (denominator <= 0n) {
      return 0;
    }

    if (numerator <= 0n) {
      return 0;
    }

    if (numerator >= denominator) {
      return 100;
    }

    const numeratorDigits = numerator.toString();
    const denominatorDigits = denominator.toString();
    const sampleDigits = 15;
    const numeratorSample = Number(numeratorDigits.slice(0, sampleDigits));
    const denominatorSample = Number(denominatorDigits.slice(0, sampleDigits));
    const exponentDelta = numeratorDigits.length - denominatorDigits.length;
    const ratio = (numeratorSample / denominatorSample) * Math.pow(10, exponentDelta);

    if (!Number.isFinite(ratio)) {
      return 0;
    }

    return Math.max(0, Math.min(100, ratio * 100));
  }

  function formatCreditsFromCents(cents, options = {}) {
    const value = normalizeGameIntegerInput(cents);
    const isNegative = value < 0n;
    const absolute = isNegative ? -value : value;
    const wholeCredits = absolute / 100n;
    const leftoverCents = absolute % 100n;

    if (wholeCredits >= normalizeGameIntegerInput(options.compactAt ?? 100000n, 100000n)) {
      return `${isNegative ? "-" : ""}${formatGameInteger(wholeCredits, options)}`;
    }

    if (leftoverCents === 0n) {
      return `${isNegative ? "-" : ""}${addCommas(wholeCredits.toString())}`;
    }

    return `${isNegative ? "-" : ""}${addCommas(wholeCredits.toString())}.${leftoverCents.toString().padStart(2, "0")}`;
  }

  function coerceSmallNumber(value, fallback = 0) {
    if (typeof value === "number") {
      return Number.isFinite(value) ? value : fallback;
    }

    if (typeof value === "bigint") {
      return Number(value);
    }

    if (typeof value === "string" && /^[+-]?\d+(\.\d+)?$/.test(value.trim())) {
      const parsed = Number(value);
      return Number.isFinite(parsed) ? parsed : fallback;
    }

    return fallback;
  }

  const api = {
    parse: normalizeGameIntegerInput,
    format: formatGameInteger,
    formatFull: formatFullGameInteger,
    formatRatio: formatGameIntegerRatio,
    formatCreditsFromCents,
    percent: gameIntegerPercent,
    smallNumber: coerceSmallNumber,
    suffixes: [...SUFFIXES],
  };

  window.NamigotchiNumbers = api;

  // Backward-compatible hooks for the existing app.js display calls.
  window.formatCompactNumber = (value) => formatGameInteger(value, { precision: 1 });
  window.formatWholeCredits = (cents) => formatCreditsFromCents(cents, { precision: 1 });
  window.formatCredits = (cents) => formatCreditsFromCents(cents, { precision: 1 });
  window.percent = gameIntegerPercent;
})();
