# CSS Cleanup Ledger

This ledger tracks the behavior-preserving cleanup of the main Namigotchi Idle stylesheet. The first task is audit-only: no application styling was changed.

## Baseline

- Untouched baseline commit: `05a611335b4d50bafa7240b363840abe8720e31e`
- Main stylesheet: `web/styles.css`
- Initial line count: `10,755`
- Git repository confirmed: `/workspace/Namagotchi`
- Initial working tree status: clean on branch `work`
- Stylesheet load location: `web/index.html` loads `/styles.css` before the theme stylesheet and `auth-landing.css`.

## Tooling and commands

No `package.json`, frontend build config, lint config, browser automation config, screenshot tooling, or visual-regression tooling was present in the repository at audit time. The application is a Go module with static assets under `web/`.

| Purpose | Command | Baseline result |
| --- | --- | --- |
| Git root | `git rev-parse --show-toplevel` | `/workspace/Namagotchi` |
| Git status | `git status --short --branch` | clean, `## work` |
| Baseline commit | `git rev-parse HEAD` | `05a611335b4d50bafa7240b363840abe8720e31e` |
| Stylesheet line count | `wc -l < web/styles.css` | `10755` |
| CSS import search | `rg -n "styles\.css|<link" web/index.html web/*.js internal` | `/styles.css` loaded in `web/index.html` |
| Tooling search | `rg --files -g 'package.json' -g 'Makefile' -g 'Taskfile*' -g 'go.mod' -g 'go.work' -g 'Dockerfile' -g '*.yml' -g '*.yaml'` | `go.mod` only |
| Tests / compile baseline | `go test ./...` | pass |

## Baseline failures

None observed. `go test ./...` passed for the root module, `internal/database`, and `internal/server`.

## Current selector counts

The audit parsed CSS rules with at-rule context and selector-list expansion. Counts below are repeated selector groups within the same cascade context; matching selector text in different media/supports/keyframe contexts was counted separately.

| Risk category | Meaning | Groups |
| --- | --- | ---: |
| 1 | Exact duplicate blocks in the same cascade context | 34 |
| 2 | Earlier declarations unconditionally overridden later by the exact same selector/context | 79 |
| 3 | Duplicate selectors containing both overridden and still-effective declarations | 104 |
| 4 | Duplicate selectors separated by potentially competing selectors | 53 |
| 5 | Duplicate selectors involving states, pseudo-elements, animations, responsive rules, themes, or runtime-generated classes | 78 |
| 6 | Intentional or cannot safely be consolidated without changing cascade behavior | 0 initially classified by script; likely candidates are listed as unresolved until manually proven safe |
| **Total** | **Repeated selector groups** | **348** |

- Provably redundant declarations identified by same-selector/same-context property override accounting: `1,246`.
- High-risk or order-dependent groups to defer initially: `131` (`risk 4 + risk 5`; unresolved intentional groups are treated as high-risk until reviewed).
- Current remaining repeated selector groups: `348`.
- Current remaining provably redundant declarations: `1,246`.

## Visual verification strategy

Preferred available method: use the existing Go application and manually or semi-automatically inspect rendered states in a browser. No existing screenshot or E2E framework was found, and this audit intentionally does not add one.

Recommended verification viewports for future cleanup batches:

- Desktop wide: `1600x1000`
- Desktop threshold: `1380x900`
- Tablet / rail breakpoint: `980x900`
- Mobile: `560x900`
- Narrow mobile: `360x740`
- Short viewport for modal/auth constraints: `900x760`

Major states to verify after any future style cleanup:

- Auth landing, prelanding, Google sign-in card, short-height auth layout.
- Main shell with left and right rails expanded/collapsed.
- Home screen: idle stage, sleep stage, care cards, care stat rows, Nami message log, care action queue/progress.
- Chat: normal, hidden/collapsed, resized desktop, emoji picker inline/floating, tabs, unread state, profile modal.
- Work / gathering screen: task cards, active task, Nami work animation, responsive grids.
- Playdeck / combat: dev buttons, combat log rows, top rail player/playdeck meters.
- Wardrobe: equipment cards, inventory category groups, bonus rows, rarity variants, item modal, drag state, read-only/share state, mobile modal header layouts.
- Theme settings and dropdown/select controls.
- Theme effects and theme-specific DOM classes: Sakura, Tokyo Night, Candy, Cafe, Rainy Mood, Pearl Tide, Woodland Sun, Woodland Moon, and reduced-motion mode.
- Scrollbar styling in main panel, wardrobe columns, and scrollable logs.

Smallest useful tooling addition for a later task: a lightweight Playwright script or Go-driven static test harness that starts the existing server, captures the viewports above, toggles known sections/themes/classes, and compares screenshots against a committed baseline. Do not install this during cleanup batch 1 unless specifically approved.

## Prioritized cleanup batch plan

Batch sizes should stay under eight repeated selector groups and about 200 changed lines. Each batch should preserve cascade context and be verified with the states above.

1. **Batch 1: exact global Wardrobe bonus row duplicates.** Remove only exact duplicate global rules for `.wardrobe-bonus-list`, `.wardrobe-bonus-row`, and `.wardrobe-bonus-row strong` around the first duplicated Wardrobe region, preserving the later identical declarations. Risk 1, 3 selector groups, approximately 20 changed lines.
2. **Batch 2: completed 2026-07-10 — exact global Wardrobe panel shell duplicates.** Removed the earlier exact duplicate global grouped panel min-height block for `.wardrobe-equipment-panel`, `.wardrobe-inventory-panel`, and `.wardrobe-bonuses-panel`.
3. **Batch 3: exact global Nami message/need text duplicates.** Handle `.nami-need-card p:nth-child(1)`, `.nami-need-card p:nth-child(2)`, and `.nami-log-message:last-child`. Risk 1, 3 selector groups.
4. **Batch 4: exact global chat hidden-state duplicates.** Handle `.emoji-picker.hidden` and `.chat-form .emoji-picker.hidden`, then verify emoji picker hidden/open states. Risk 1, 2 selector groups.
5. **Batch 5: exact scrollbar block duplicates.** Handle `.main-panel::-webkit-scrollbar` only after desktop scrollbar screenshots are captured. Risk 1 but browser-pseudo-sensitive.
6. **Batch 6: same-context low-risk layout overrides.** Review non-state risk-2 groups with simple global selectors such as `.game-shell`, `.left-rail`, `.right-rail`, and `.center-column`; keep batch small and verify desktop/tablet layout.
7. **Batch 7: navigation/button low-risk overrides.** Review `.main-nav`, `.nav-item`, `.primary-button`, `.secondary-button`, and related non-hover base rules; defer hover/active rules.
8. **Batch 8+: component-specific passes.** Proceed by subsystem: Chat, Home/Nami cards, Wardrobe layout, Wardrobe modal, top rail, responsive rail reset, theme-effect guards, and only then pseudo/state/animation/theme groups.

## Completed cleanup batches

### Batch 1 — 2026-07-10 — exact global Wardrobe bonus row duplicates

- Selectors examined: `.wardrobe-bonus-list`, `.wardrobe-bonus-row`, `.wardrobe-bonus-row strong` in the global cascade context.
- Cascade context recorded: all cleaned occurrences were top-level global rules in `web/styles.css` with no enclosing `@media`, `@supports`, `@layer`, `@scope`, or container rule. No pseudo-classes, pseudo-elements, animations, transitions, vendor-prefixed declarations, `!important`, or fallback declaration pairs were present in the removed rules.
- Usage/dynamic check: `rg -n "wardrobe-bonus-list|wardrobe-bonus-row" .` found the static bonus list in `web/index.html`, row generation in `web/app.js`, additional scoped `#section-inventory` rules in `web/styles.css`, theme overrides in `web/themes/*.css`, and the audit/ledger references. The runtime row class is generated in `web/app.js`; `.is-empty`, `.wardrobe-bonus-label`, and `.wardrobe-bonus-value.is-hidden` state selectors remain untouched.
- Cascade timeline: the earlier global blocks at former lines 5024-5045 and the later global blocks at former lines 5365-5386 had identical selector specificity, identical declarations, identical custom-property references, and identical global context. Deleting only the earlier copies preserves the later source-order position and leaves all intervening `.inventory-card-grid`, `.inventory-item-card`, equipment-card, scoped inventory, and theme overrides in their previous relationship to the surviving rules.
- Overlapping selectors considered: later `#section-inventory .wardrobe-bonus-*` rules have higher specificity and remain later in source order; theme selectors under `body[data-theme=...]` are loaded after `/styles.css` or are more specific and remain untouched.
- Selectors safely cleaned: `.wardrobe-bonus-list`, `.wardrobe-bonus-row`, `.wardrobe-bonus-row strong`.
- Selectors deferred: none in this batch.
- Declarations and lines removed: 14 CSS declarations and 23 stylesheet lines.
- Line count before/after: `10,755` lines before; `10,732` lines after.
- Commands executed: `rg -n "wardrobe-bonus-list|wardrobe-bonus-row" .`; `python3` selector/count checks; `python3` brace parse for `web/styles.css`; `go test ./...`; `git diff -- web/styles.css docs/css-cleanup-ledger.md docs/css-cleanup-audit.md`; `git status --short`.
- Verification results: CSS brace parse passed; `go test ./...` passed. No browser automation or visual-regression tooling exists in this repository, so Wardrobe visual states still need manual browser inspection at the recommended `1600x1000`, `980x900`, and `560x900` viewports.
- Remaining repeated-selector count: `345` (baseline `348` minus these 3 exact duplicate selector groups).
- Remaining provably redundant declarations: `1,232` (baseline `1,246` minus 14 identical declarations removed).
- Risks/manual checks still needed: manually inspect the Wardrobe bonus panel in the default theme and in themes that override `.wardrobe-bonus-row`/`.wardrobe-bonus-row strong`, especially Pearl Tide, Woodland Moon/Sun, Phantom Rebel, Cafe, and Sakura Light inherited text behavior.

### Batch 2 — 2026-07-10 — exact global Wardrobe panel shell duplicates

- Selectors examined: `.wardrobe-equipment-panel`, `.wardrobe-inventory-panel`, `.wardrobe-bonuses-panel` in the global cascade context.
- Cascade context recorded: both duplicate grouped occurrences were top-level global rules in `web/styles.css` with no enclosing `@media`, `@supports`, `@layer`, `@scope`, or container rule. No pseudo-classes, pseudo-elements, animations, transitions, vendor-prefixed declarations, `!important`, inherited-property dependencies, custom-property values, or fallback declaration pairs were present in the removed rule.
- Usage/dynamic check: `rg -n "wardrobe-equipment-panel|wardrobe-inventory-panel|wardrobe-bonuses-panel" .` found the static panel classes in `web/index.html`, scoped inventory-panel and bonuses-panel refinements in `web/styles.css`, theme overrides in `web/themes/candy.css`, `web/themes/cafe.css`, and `web/themes/rainy-mood.css`, backup stylesheet references, and audit/ledger references. No dynamically generated class names for these three panel shell classes were found; runtime behavior targets their static DOM elements and scoped `#section-inventory` states.
- Cascade timeline: the earlier global grouped block at former lines 4910-4914 and the later global grouped block now at lines 5230-5234 had identical selector specificity after selector-list expansion and identical `min-height: 0` declarations. The selector order differed within the comma list, but each expanded selector had the same single declaration in the same global context. Deleting only the earlier copy preserves the later source-order position and leaves intervening `.equipment-card-grid`, `.equipment-card`, `.inventory-card-grid`, rarity, responsive Wardrobe, equipment slot, and preview rules in their prior relationship to the surviving panel shell rule.
- Overlapping selectors considered: later `#section-inventory .wardrobe-bonuses-panel`, `#section-inventory .wardrobe-inventory-panel`, responsive `#section-inventory` inventory-panel rules, and theme-specific `body[data-theme=...]` panel overrides have higher specificity and/or remain later in source order and were not changed.
- Selectors safely cleaned: `.wardrobe-equipment-panel`, `.wardrobe-inventory-panel`, `.wardrobe-bonuses-panel`.
- Selectors deferred: none in this batch.
- Declarations and lines removed: 3 expanded CSS declarations and 5 stylesheet lines.
- Line count before/after: `10,732` lines before; `10,727` lines after.
- Commands executed: `rg -n "wardrobe-equipment-panel|wardrobe-inventory-panel|wardrobe-bonuses-panel" .`; `python3` occurrence and selector/count checks; `python3` brace parse for `web/styles.css`; `go test ./...`; `git diff -- web/styles.css docs/css-cleanup-ledger.md docs/css-cleanup-audit.md`; `git status --short`.
- Verification results: CSS brace parse passed; `go test ./...` passed for `github.com/SorynMochi/Namagotchi`, `internal/database`, and `internal/server`. No browser automation, visual-regression tooling, frontend build config, type-check config, or lint config exists in this repository, so Wardrobe visual states still need manual browser inspection at the recommended `1600x1000`, `980x900`, and `560x900` viewports.
- Remaining repeated-selector count: `342` (previous `345` minus these 3 exact duplicate selector groups).
- Remaining provably redundant declarations: `1,229` (previous `1,232` minus 3 identical expanded declarations removed).
- Risks/manual checks still needed: manually inspect the Wardrobe equipment, inventory, and bonuses panel shells in the default theme and in themes that override these panels, especially Candy, Cafe, and Rainy Mood.

## Unresolved or intentionally repeated selector groups

Treat the following as unresolved/high-risk until manually proven safe:

- `:root` token append blocks. They may intentionally layer custom properties and aliases over the initial theme token foundation.
- Responsive top rail repetitions for `.top-stats-compact` and top-player/playdeck rules under `max-width: 1600px`, `1200px`, and `820px`.
- Hover, focus, active, disabled, and pseudo-element selectors, including `.nav-item:hover`, `.nav-item.active`, button hover rules, `::before`, `::after`, and `::-webkit-scrollbar*`.
- Theme-effect and dynamically inserted selectors from `web/theme-effects.js`, especially Tokyo Night signs, Sakura petals/lanterns, Candy pieces, Cafe floating pieces, Rainy Mood video state, and Woodland/Pearl canvases.
- Wardrobe modal corrections near the bottom of the stylesheet. These are repeated in many sequential repair blocks and include viewport, drag, mobile, action-grid, and ID-specific overrides.
- Responsive layout resets under `max-width` media queries. Selector text matches across breakpoints should not be merged across different contexts.
- Fallback declarations and vendor-prefixed declarations; repeated properties are not removed unless fallback behavior is understood.

## Remaining counts

- Remaining repeated selector groups after Batch 2: `342` (superseded by the final audit below)
- Remaining provably redundant declarations after Batch 2: `1,229` (superseded by the final audit below)
- Remaining high-risk/order-dependent groups after Batch 2: `131` (superseded by the final audit below)

### Batch 3 — 2026-07-10 — exact global chat hidden-state duplicates

- Selectors examined: `.emoji-picker.hidden` and `.chat-form .emoji-picker.hidden` in the global cascade context.
- Classification: Category A, safely removable redundancy. The removed earlier rules were exact same-selector/same-context duplicates of later hidden-state rules and used only `display: none` declarations matching the surviving later rules.
- Selectors safely cleaned: `.emoji-picker.hidden`, `.chat-form .emoji-picker.hidden`.
- Declarations and lines removed: 2 CSS declarations and 8 stylesheet lines.
- Line count before/after this internal batch: `10,727` to `10,719`.
- Verification: CSS brace parse passed; `go test ./...` passed. No browser or screenshot harness exists, so affected emoji hidden/open states remain manual visual checks.

### Batch 4 — 2026-07-10 — exact global Home/Nami row duplicates

- Selectors examined: `.nami-need-card p:nth-child(1)`, `.nami-need-card p:nth-child(2)`, and `.nami-log-message:last-child` in the global cascade context.
- Classification: Category A, safely removable redundancy. The removed earlier rules were exact same-selector/same-context duplicates of later Home compact layout rules; the surviving later rules preserve the final grid-row and last-message margin behavior.
- Selectors safely cleaned: `.nami-need-card p:nth-child(1)`, `.nami-need-card p:nth-child(2)`, `.nami-log-message:last-child`.
- Declarations and lines removed: 3 CSS declarations and 12 stylesheet lines.
- Line count before/after this internal batch: `10,719` to `10,707`.
- Verification: CSS brace parse passed; `go test ./...` passed. No browser or screenshot harness exists, so affected Home/Nami care summary and message log states remain manual visual checks.

## Final audit — 2026-07-10

- Final main stylesheet line count: `10,707` lines, down from the original `10,755` baseline.
- Final repeated selector groups in the same cascade context: `337`, down from the original `348` baseline. Reconciliation: completed batches removed 11 same-context repeated selector groups total: 3 Wardrobe bonus row groups, 3 Wardrobe panel shell groups, 2 chat hidden-state groups, and 3 Home/Nami groups.
- Final overridden declarations by same-selector/same-context property accounting: `1,357` occurrences by the fresh parser used for this final audit. This differs from the original audit's `1,246` because the final parser counts every earlier property occurrence overridden by a later same-property rule, including repeated fallback-style properties; the comparable repeated-selector group count reconciles exactly with the batch removals.
- Final identical declaration repeat occurrences: `650` by fresh same-selector/same-context accounting.
- Final exact duplicate block groups: `41` by fresh same-selector/same-context accounting. These are not all safe to delete because many are responsive, order-dependent, browser-pseudo, top-rail repair, modal repair, or intentionally repeated layout clamps.
- Final empty blocks: `0`.
- Potentially unused selectors: unresolved. Runtime class generation in `web/app.js` and theme-effect DOM generation in `web/theme-effects.js` make static deletion unsafe without browser coverage.
- Final complete deterministic verification: `go test ./...` passed. No `package.json`, frontend build, type-check, lint, or visual-regression tooling exists in this repository, so no additional deterministic frontend suite could be run.
- Final visual verification: not automated in-repo. Required manual coverage remains the route/screen/viewport list above, with special attention to Chat emoji hidden/open states and Home/Nami compact rows touched in the final batches.
- Theme stylesheet behavior: `/styles.css` still loads before active theme CSS and `auth-landing.css`; no theme stylesheet was changed in these final batches.
- Unrelated files: final changed files are limited to `web/styles.css`, `docs/css-cleanup-ledger.md`, and `docs/css-cleanup-audit.md`.

### Final category classification of remaining repeated groups

- Category A — safely removable redundancy: the obvious exact global duplicates targeted in Batches 1-4 were removed. Remaining Category A candidates are isolated but not implemented here because they need browser-sensitive checks, e.g. `.nav-item.premium`, `body.is-resizing-chat`, and selected one-line grid clamps.
- Category B — safely consolidatable with small local restructuring: small component-local candidates remain, including selected Home compact stat sizing and simple Wardrobe mobile grid repeats. They should be handled only in new small batches with screenshots or computed-style checks.
- Category C — consolidatable only if all competing selectors are reorganized together: broad base layout and component primitives, especially `body`, `.game-shell`, `.left-rail`, `.right-rail`, `.center-column`, `.panel`, `.main-panel`, `.nav-panel`, `.main-nav`, `.nav-item`, button primitives, `.chat-panel`, `.chat-form`, Wardrobe inventory/card rules, and top rail stat/player/playdeck rules.
- Category D — intentional progressive enhancement or browser fallback: vendor scrollbar pseudo-elements, fallback/repair declarations, opacity/fallback modal repairs, and browser-specific UI guard blocks.
- Category E — intentional cascade or source-order behavior: `:root` token append blocks, theme-effect guard layers, responsive rail/top-stat repetitions across breakpoints, and late modal/top-rail repair blocks that intentionally win by source order.
- Category F — uncertain due to insufficient runtime or visual coverage: potentially unused dynamic selectors, interaction states, theme-effect generated classes, drag/read-only/share Wardrobe states, auth/prelanding states, and modal variants. These were retained.

### Category C restructuring plan for later work

1. **Global shell and rails**
   - Move together: `body`, `button`, `input`, `.game-shell`, `.left-rail`, `.right-rail`, `.center-column`, `.main-panel`, `.panel`, `.nav-panel`.
   - Current cascade behavior: early base shell rules are partially overwritten by later cute-UI, rail-density, and scroll repair passes.
   - Proposed behavior: create one base shell section plus one clearly marked later responsive/repair section, preserving final computed values.
   - Equivalence tests: computed-style snapshots and screenshots at `1600x1000`, `1380x900`, `980x900`, `560x900`, and `360x740`, with rails expanded/collapsed and chat hidden/open.
   - Rollback boundary: one commit touching only shell/rail selectors.
   - Estimated changed-line count: 250-450.
2. **Chat and emoji subsystem**
   - Move together: `.chat-panel`, `.chat-tabs`, `.chat-tab`, `.chat-tab.active`, `.chat-form`, `.chat-form input`, `.chat-emoji-button`, `.emoji-picker`, `.emoji-option`, floating picker selectors, and hidden/collapsed chat selectors.
   - Current cascade behavior: early chat layout, resize fixes, hidden-state hard overrides, and emoji category pass layer by source order and `!important`.
   - Proposed behavior: consolidate base chat layout, then explicit states for hidden, floating, active tab, hover, and responsive behavior.
   - Equivalence tests: screenshots and computed styles for normal chat, hidden chat, resized chat, emoji picker open/hidden, tab active/unread, and profile modal at desktop/mobile viewports.
   - Rollback boundary: one commit touching only chat/emoji selectors.
   - Estimated changed-line count: 300-600.
3. **Wardrobe modal and inventory**
   - Move together: `#section-inventory` layout rules, `.equipment-card-grid`, `.inventory-card-grid`, `.inventory-item-card`, `.wardrobe-*panel`, `.wardrobe-item-modal*`, action buttons, drag, read-only, share, rarity, and mobile modal selectors.
   - Current cascade behavior: first Wardrobe pass is corrected by compact layout, scoped inventory rules, modal viewport repairs, mobile header repairs, and late action-grid fixes.
   - Proposed behavior: split stable base inventory styles from modal state styles and viewport-specific repairs while preserving final scoped overrides.
   - Equivalence tests: screenshots for equipment/inventory/bonuses panels, populated and empty inventory groups, item modal, compare target, action grid, drag state, read-only/share state, and mobile modal header at `1600x1000`, `980x900`, `560x900`, and `360x740`.
   - Rollback boundary: one commit touching only Wardrobe selectors.
   - Estimated changed-line count: 500-900.
4. **Top rail migration cleanup**
   - Move together: `.top-stats-compact`, `.stat-group.compact`, `.top-user-info`, top-player/top-playdeck slot selectors, and responsive `max-width: 1600px`, `1200px`, and `820px` repeats.
   - Current cascade behavior: late migration/refinement blocks use `!important` and repeated nth-child placement to override older top-stat assumptions.
   - Proposed behavior: one canonical top rail grid and one breakpoint ladder, with player/playdeck cards integrated into the same grid definition.
   - Equivalence tests: screenshots and computed grid styles for all top stat groups and player/playdeck slots at wide, 1600, 1200, 980, 820, and mobile widths.
   - Rollback boundary: one commit touching only top rail selectors.
   - Estimated changed-line count: 350-700.

## Other screens / shared components structural batch — 2026-07-10

- Scope: one small structural consolidation batch inside the `Other screens / shared components` family. Touched selectors: `.log-list`, `.log-list::-webkit-scrollbar`, `.log-list::-webkit-scrollbar-thumb`, `.log-list::-webkit-scrollbar-track`, `.task-card p`, `.buff-card span`, `.buff-card small`, `.meter-fill`, `.compact-grid-lines`, `.menu-logo-wrap`, `.menu-status-card`, `#section-gathering .task-card h2`, and `#section-gathering .task-card p`.
- Blocks before/after for duplicate candidates: `.log-list` 2→1 normal blocks, `.task-card p` 2→1 standalone declarations with the base grouped margin retained, `.buff-card span` 2→1 family-level declaration plus scoped right-rail/right-buff overrides retained, `.buff-card small` 2→1 family-level declaration plus scoped overrides retained, `.meter-fill` 2→1 shared meter declaration plus scoped top-rail/top-playdeck overrides retained, `.compact-grid-lines` 2→1 base grid declaration plus top-stat scoped variants retained, `.menu-logo-wrap` 2→1, `.menu-status-card` 2→1, `#section-gathering .task-card h2` 2→1, and `#section-gathering .task-card p` 2→1.
- Old line ranges: `.log-list` base at 1090-1099 and scrollbar membership at 1327-1353; `.task-card p` at 1128-1140; `.buff-card span`/`small` at 1317-1320 and 1518-1524; `.meter-fill` at 1430-1434 and 4112-4120; `.compact-grid-lines` at 1604-1614; `.menu-logo-wrap`/`.menu-status-card` at 6725-6756; `#section-gathering .task-card h2`/`p` at 7081-7100.
- New line ranges: `.log-list` canonical block and scrollbar pseudo-elements at 1090-1130; `.task-card p` at 1154-1157; `.buff-card span`/`small` at 1334-1338; `.meter-fill` at 1440-1447; `.compact-grid-lines` at 1609-1613; `.menu-logo-wrap` at 6729-6736; `.menu-status-card` at 6751-6758; `#section-gathering .task-card h2` at 7083-7088; `#section-gathering .task-card p` at 7095-7099.
- Source-order relationships preserved: inherited/base selectors remain before scoped right-rail, right-buff-list, top-rail, top-playdeck, and `#section-gathering` refinements; only same-specificity duplicate declarations were moved into the nearest canonical component block or removed from later grouped patches after being reproduced earlier. The later theme pass still supplies the same meter shadow to `.fill`, `.tick-fill`, `.mood-fill`, and `.stat-fill`; `.meter-fill` now carries that value directly in its canonical block.
- Intervening competitors reviewed: grouped scrollbar rules for `.main-panel` and `.chat-log`, cute-UI fill shadow group, right-rail/right-buff-list `.buff-card` descendants, top-user/top-playdeck `.meter-fill`, top-stat `.compact-grid-lines`, collapsed-left-rail `.menu-logo-wrap`/`.menu-status-card`, and the scoped gathering task-card rules.
- Line count: `web/styles.css` went from 10,707 to 10,706 lines; net 1 line removed.
- Duplicate blocks removed or flattened: 10 candidate duplicate memberships/blocks were reduced, with browser scrollbar pseudo-element behavior retained as explicit `.log-list` pseudo-element rules.
- Baseline/equivalence checks: used static computed-cascade reasoning for equal-specificity same-property declarations and repository selector searches; no browser automation harness is available in this checkout, so visual state coverage remains manual.
- Commands and results: `find /workspace -name AGENTS.md -print` found no AGENTS.md; `rg -n` selector searches found stylesheet, markup, script, theme, and doc references; `python3` stylesheet brace and occurrence check passed with brace balance 0; `go test ./...` passed.
- Unresolved risks/manual inspection: manually inspect combat log scrolling, gathering task cards, buff cards, left menu logo/status cards, player meters, top-stat compact grids, and gathering task cards at desktop, 1380px, 980px, 560px, and collapsed-left-rail states.
