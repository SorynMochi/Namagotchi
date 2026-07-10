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
2. **Batch 2: exact global Wardrobe panel shell duplicates.** Handle exact duplicate global grouped panel rules for `.wardrobe-equipment-panel`, `.wardrobe-inventory-panel`, and `.wardrobe-bonuses-panel`. Risk 1, 3 selector groups.
3. **Batch 3: exact global Nami message/need text duplicates.** Handle `.nami-need-card p:nth-child(1)`, `.nami-need-card p:nth-child(2)`, and `.nami-log-message:last-child`. Risk 1, 3 selector groups.
4. **Batch 4: exact global chat hidden-state duplicates.** Handle `.emoji-picker.hidden` and `.chat-form .emoji-picker.hidden`, then verify emoji picker hidden/open states. Risk 1, 2 selector groups.
5. **Batch 5: exact scrollbar block duplicates.** Handle `.main-panel::-webkit-scrollbar` only after desktop scrollbar screenshots are captured. Risk 1 but browser-pseudo-sensitive.
6. **Batch 6: same-context low-risk layout overrides.** Review non-state risk-2 groups with simple global selectors such as `.game-shell`, `.left-rail`, `.right-rail`, and `.center-column`; keep batch small and verify desktop/tablet layout.
7. **Batch 7: navigation/button low-risk overrides.** Review `.main-nav`, `.nav-item`, `.primary-button`, `.secondary-button`, and related non-hover base rules; defer hover/active rules.
8. **Batch 8+: component-specific passes.** Proceed by subsystem: Chat, Home/Nami cards, Wardrobe layout, Wardrobe modal, top rail, responsive rail reset, theme-effect guards, and only then pseudo/state/animation/theme groups.

## Completed cleanup batches

No cleanup batches have been performed yet.

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

- Remaining repeated selector groups: `348`
- Remaining provably redundant declarations: `1,246`
- Remaining high-risk/order-dependent groups: `131`
