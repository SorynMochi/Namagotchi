# CSS Cleanup Audit

This audit covers `web/styles.css` at baseline commit `05a611335b4d50bafa7240b363840abe8720e31e`. It is intentionally navigable rather than a raw selector dump.

## Scope and loading

- Main stylesheet: `web/styles.css`, `10,755` lines.
- Load order: `/styles.css` is loaded from `web/index.html` before the active theme stylesheet and `auth-landing.css`. Theme CSS and auth CSS can therefore override the main stylesheet.
- Runtime class generation is substantial. `web/app.js` toggles sections, rail collapse, auth state, sleep/idle stage state, chat state, wardrobe modal state, and top rail/player/playdeck classes. `web/theme-effects.js` dynamically inserts theme-effect DOM and classes for Sakura, Tokyo Night, Candy, Cafe, Rainy Mood, Pearl Tide, Woodland Sun, and Woodland Moon.

## Method

A structural parser expanded selector lists and tracked the enclosing at-rule context for each style rule. Duplicate accounting is keyed by normalized selector and exact context, so a selector repeated in different `@media`, `@supports`, or keyframe contexts is not treated as the same accidental duplicate.

For each repeated same-context selector, the audit considered source order, declaration property/value pairs, `!important`, pseudo/state selectors, custom properties, vendor/fallback patterns, and nearby component comments. Specificity was computed approximately for prioritization; exact cleanup decisions still require browser verification because theme CSS, dynamic classes, and source order are central to this file.

## Summary counts

- Parsed style rules: `1,345`
- Repeated selector groups in the same cascade context: `348`
- Provably redundant same-selector/same-context declarations: `1,246`
- Risk 1 groups: `34`
- Risk 2 groups: `79`
- Risk 3 groups: `104`
- Risk 4 groups: `53`
- Risk 5 groups: `78`
- High-risk/order-dependent groups for initial deferral: `131` (`risk 4 + risk 5`)

## Stylesheet regions

- Lines 1-545: root theme foundation, aliases, and custom property tokens.
- Lines 551-1360: initial global shell, panels, Nami room, tasks, chat, right rail, scrollbars, and first desktop fallback.
- Lines 1371-3158: accumulated UI fix passes for rails, tick engine, compact top stats, chat, emoji picker, profile modal, compact rails, gathering, and emoji performance.
- Lines 3160-3864: compact Home/Nami care layout, message log, fixed home height, and care action queue.
- Lines 3948-4550: cute UI theme pass, chat visual refresh, left menu redesign, and global responsive layout pass.
- Lines 4910-6578: first Wardrobe layout and later compact/corrective Wardrobe layout blocks.
- Lines 6585-7404: consolidated user info, three-zone dashboard, rail chrome, Home sleep mode, Work screen, scrollbars, theme settings, and theme-effect guards.
- Lines 7432-8714: Wardrobe item modal, modal compact passes, action positioning, accessory compare, share links, and auth landing.
- Lines 8861-9051: viewport modal overlay, fallback/opacity repair, and modal dragging corrections.
- Lines 9051-10755: side rail glyph polish, rarity text bridge, migrated top rail base, migrated dropdown base, top rail player/playdeck redesign/refinement, and top player slot customization.

## Risk 1: exact duplicate blocks

These are the lowest-risk candidates, provided source order and neighboring comments are preserved by deleting only one identical block at a time.

Representative exact duplicates:

| Selector/context | Occurrences | Notes |
| --- | --- | --- |
| `.nav-item.premium` / global | lines 744-746 and 2630-2632 | Exact opacity block. |
| `.main-panel::-webkit-scrollbar` / global | lines 1334-1339 and 7257-7260 | Exact-looking pseudo scrollbar dimensions; verify desktop scrollbars before cleanup. |
| `body.is-resizing-chat` / global | lines 1688-1691 and 5797-5800 | Exact cursor/user-select resize state. |
| `.emoji-picker.hidden` / global | lines 1881-1883 and 1972-1974 | Exact hidden display rule. |
| `.chat-form .emoji-picker.hidden` / global | lines 2094-2096 and 2138-2140 | Exact hidden display rule in chat form. |
| `.nami-need-card p:nth-child(1)` / global | lines 3390-3392 and 3617-3619 | Exact typography color rule. |
| `.nami-need-card p:nth-child(2)` / global | lines 3394-3396 and 3621-3623 | Exact typography color rule. |
| `.nami-log-message:last-child` / global | lines 3537-3539 and 3697-3699 | Exact border removal. |
| `.equipment-card-grid` / `@media (max-width: 700px)` | lines 5194-5197 and 5459-5462 | Same media context; verify mobile wardrobe grid. |

Batch 1 completed on 2026-07-10: the earlier global Wardrobe bonus row exact duplicates (`.wardrobe-bonus-list`, `.wardrobe-bonus-row`, `.wardrobe-bonus-row strong`) were removed after same-context cascade equivalence was verified. The surviving later global rules remain the canonical declarations.

Batch 2 completed on 2026-07-10: the earlier exact duplicate global Wardrobe panel shell min-height rules (`.wardrobe-equipment-panel`, `.wardrobe-inventory-panel`, `.wardrobe-bonuses-panel`) were removed after same-context cascade equivalence was verified. The surviving later global grouped rule remains the canonical declaration.

## Risk 2: earlier declarations unconditionally overridden later

These groups have repeated declarations where later same-selector/same-context rules override earlier declarations. They are promising but should be handled only after exact duplicates.

Representative groups:

| Selector/context | Occurrences | Notes |
| --- | --- | --- |
| `body` / global | lines 551-566 and 3950-3962 | Later cute UI pass adjusts global presentation; root-level body cleanup requires full-app screenshots. |
| `.game-shell` / global | lines 589-598 and 2511-2515 | Layout sizing override; verify desktop and tablet. |
| `.left-rail` / global | lines 606-617 and 2517-2519 | Layout density pass; coupled to collapse behavior. |
| `.right-rail` / global | lines 606-617 and 2712-2714 | Similar to left rail but right buffs content differs. |
| `.center-column` / global | lines 606-624, 1548-1550, 1662-1664 | Chat resizing/top stats interactions make this more sensitive. |
| `.panel-title` / global | lines 651-661 and 4014-4038 | Later grouped typography pass changes presentation. |
| `.player-name` / global | lines 663-676 and 4014-4057 | Shared with top/user info cards. |
| `.bar`, `.fill`, `.xp-fill`, `.mood-fill`, `.stat-fill` / global | lines 678-700 and 4122-4151 | Progress bar primitives used throughout app. |
| `.main-nav`, `.nav-item`, `.primary-button` / global | lines 706-760 and 2591-4105 | Broad shared UI primitives; defer until exact duplicates are gone. |

## Risk 3: duplicate selectors with both overridden and still-effective declarations

These require merging effective declarations carefully, usually by component. Do not delete entire earlier blocks unless every declaration is proven redundant.

Representative groups:

| Selector/context | Occurrences | Notes |
| --- | --- | --- |
| `.panel` / global | lines 626-635 and 3982-4000 | Later theme pass supplements/replaces panel look. |
| `.nav-panel` / global | lines 702-704, 2587-2589, 4419-4421 | Rail/menu redesign layers on base panel placement. |
| `.main-panel` / global | lines 815-822, 1327-1332, 7251-7255 | Scroll behavior and visual shell interact. |
| `.nami-card`, `.care-card` / global | many occurrences between 881-901 and 3982-4000 | Multiple Home layout passes and cute theme pass. |
| `.nami-message` / global | lines 988-996, 3201-3206, 3776-3781, 4241-4248 | Home message styling across compact/fixed-height passes. |
| `.chat-panel`, `.chat-tabs`, `.chat-tab`, `.chat-form`, `.chat-form input` / global | many occurrences from 1167 through 4322 | Chat has resize, hidden, emoji, tabs, and theme styling dependencies. |
| Wardrobe inventory/card selectors / global | many occurrences from 4910 through 6578 | Must be grouped by wardrobe subsystem and verified by category/rarity state. |

## Risk 4: duplicate selectors separated by potentially competing selectors

These groups either have no trivially redundant declarations or sit among broad grouped selectors that may be intentionally layered.

Representative groups:

- `button` and `input` global resets around lines 574-581 and 3950-3954.
- Heading selectors such as `h1` and `h3` that participate in global typography and later cute UI passes.
- `.log-list`, `.task-card p`, `.buff-card span`, `.buff-card small`, `.user-card`, `.meter-fill`, `.compact-grid-lines`, and `.chat-message`.
- Many of these should be left until component-specific cleanup can compare computed styles before and after.

## Risk 5: states, pseudo-elements, animation, responsive, theme, or runtime-generated selectors

These are high-risk and should be deferred until screenshots or browser-computed-style checks exist for the affected states.

Representative groups:

| Selector/context | Occurrences | Notes |
| --- | --- | --- |
| `:root` / global | lines 1-545 plus later token append blocks at 3726, 7926, 9215, 9281, 9615, 9693, 10089, 10144, 10281 | Custom property source order is intentional until proven otherwise. |
| `.panel::before` / global | lines 637-649 and 4002-4012 | Pseudo-element visual treatment. |
| `.nav-item:hover`, `.nav-item.active` / global | lines 729-742, 2615-2624, 4072-4087 | Interactive state styling. |
| Button hover selectors / global | lines 762-767 and 4107-4113 | Shared hover effects for multiple button classes. |
| `.chat-tab.active` / global | lines 1201-1205, 1715-1719, 4294-4297 | Stateful chat tabs. |
| `.main-panel::-webkit-scrollbar-thumb` and track / global | lines 1341-1352 and 7262-7273 | Browser-specific scrollbar styling. |
| `.game-shell` / `@media (max-width: 1380px)` | lines 1356-1360 and 2690-2694 | Same media context but responsive layout-sensitive. |
| `.top-stats-compact` / responsive media contexts | repeated under `max-width: 1600px` and `max-width: 1200px` in several top rail migrations | The repeated responsive fallbacks appear intentionally preserved across top rail redesigns. |
| Wardrobe modal selectors under `@media (max-width: 560px)` and ID-specific viewport-modal selectors | many lines from 7432 through 9051 | Sequential repair blocks; cleanup must be modal-state-specific. |

## Dynamic selector inventory to respect

Before treating any selector as unused, search the full repository. Important dynamic sources include:

- Navigation and sections: `.active`, section IDs, `data-section`, rail `.collapsed` toggles.
- Home stage: `.is-sleep-stage`, `.is-idle-stage` on `.nami-room-stage`.
- Auth: `body.auth-logged-out`, `body.auth-prelanding-active`, `.auth-landing.hidden`, sparkle classes.
- Chat: `.chat-hidden`, `.is-resizing-chat`, emoji picker visibility, profile modal.
- Top rail: `.top-tick-pill`, `.top-tick-pill-animating`, `.top-player-tick-pill`, top player detail rows.
- Playdeck/combat: `.combat-log-*` classes and playdeck status bars.
- Wardrobe: equipment/inventory group classes, rarity classes, modal classes, drag state, read-only preview, share button states.
- Theme effects: all `.theme-effect-*`, `.theme-effects-*`, Tokyo sign classes, Sakura petals/lanterns, Candy and Cafe floating pieces, Woodland/Pearl canvases, Rainy Mood video state.

## Completed Batch 1

Batch 1 completed on 2026-07-10. The earlier exact duplicate global Wardrobe bonus row rules were removed and the later identical global declarations were preserved:

- `.wardrobe-bonus-list` earlier duplicate occurrence.
- `.wardrobe-bonus-row` earlier duplicate occurrence.
- `.wardrobe-bonus-row strong` earlier duplicate occurrence.

Classification change: these three selectors are no longer active Risk 1 candidates in the main stylesheet. Manual visual inspection of the Wardrobe screen remains recommended at `1600x1000`, `980x900`, and `560x900`, with inventory categories populated if possible.

## Completed Batch 2

Batch 2 completed on 2026-07-10. The earlier exact duplicate global Wardrobe panel shell rule was removed and the later identical global declaration was preserved for these expanded selectors:

- `.wardrobe-equipment-panel` earlier duplicate occurrence.
- `.wardrobe-inventory-panel` earlier duplicate occurrence.
- `.wardrobe-bonuses-panel` earlier duplicate occurrence.

Classification change: these three selectors are no longer active Risk 1 candidates in the main stylesheet. Manual visual inspection of the Wardrobe equipment, inventory, and bonuses panel shells remains recommended at `1600x1000`, `980x900`, and `560x900`, including Candy, Cafe, and Rainy Mood theme checks.
