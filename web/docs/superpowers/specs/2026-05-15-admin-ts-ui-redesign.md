# Admin TypeScript Migration and UI Redesign

Date: 2026-05-15
Project: Vue Naive Admin 2.x
Approved direction: Option A, refined with a same-tone sidebar

## Goal

Refactor the existing Vue 3 admin frontend from JavaScript to TypeScript while preserving current behavior. Refresh the interaction and visual system so the application feels like a modern admin workspace, and remove redundant or low-value logic where it does not change user-facing capability.

The refactor must keep these core capabilities intact:

- Login, auth token handling, refresh flow, and logout.
- Dynamic permission loading and route generation.
- Route guards for permission, page title, loading, and tab behavior.
- Multiple layouts: normal, simple, full, and empty.
- Sidebar navigation, header actions, breadcrumb, tabs, KeepAlive, fullscreen, user menu, role switch, layout/theme settings.
- Existing business pages for home, profile, PMS user/role/resource, iframe, demo upload, base examples, and error pages unless explicitly identified as removable during implementation review.
- Naive UI, Pinia, Vue Router, Vite, UnoCSS, Iconify/UnoCSS icons, Axios, ECharts, and the current backend API contract.

## Visual Direction

Use the refined Option A direction shown in the browser preview.

The main layout should be quiet, work-focused, and modern:

- Sidebar and content use the same visual tone. The sidebar must not have a separate dark or heavily branded background.
- Navigation hierarchy is expressed through spacing, border, icons, subtle hover states, and a light selected state with a slim active indicator.
- Content surfaces use restrained cards, clear sectioning, compact headers, and stable dimensions for toolbars, tabs, tables, and metric blocks.
- The theme remains compatible with light and dark mode. Light mode is the primary design target; dark mode should mirror the same structural language without introducing a disconnected sidebar color.
- Avoid decorative gradients, floating orbs, marketing-style hero sections, large ornamental cards, and one-note color palettes.
- Preserve the original project's "simple is good" positioning, but make density, contrast, and interaction polish more mature.

## TypeScript Scope

Replace the JavaScript-first project setup with TypeScript:

- Add TypeScript project configuration and replace `jsconfig.json` with TS-aware config.
- Rename source JavaScript modules under `src` to `.ts` where practical.
- Convert Vue SFC scripts to `script setup lang="ts"`.
- Type the app bootstrap, directives, router, route guards, Pinia stores, settings, HTTP utilities, API modules, composables, and common components.
- Add project-level declarations for `.vue`, image assets, SVG imports, virtual UnoCSS imports, and Vite env variables.
- Add explicit route meta typing for `title`, `icon`, `layout`, `keepAlive`, `order`, `hide`, permission codes, and related dynamic-menu fields.
- Add reusable types for permissions, menu nodes, tabs, user info, roles, API responses, pagination, CRUD forms, modal state, and table query payloads.
- Keep typing pragmatic. Use narrow domain types for shared contracts and avoid over-modeling one-off view-local data.

## UI and Interaction Scope

Update the shared shell first, then the highest-value pages:

- App shell: `App.vue`, normal/simple/full layouts, sidebar, side logo, side menu, header, breadcrumb, tab bar, context menu, user avatar, fullscreen, menu collapse, layout setting, theme setting.
- Login page: fully adopt the two-column animated login page direction from `https://gitee.com/niumg9527/login-animation`, while keeping the current form flow, redirect behavior, captcha/token assumptions, and API calls. The animated characters should be redesigned instead of copied directly: rounder, livelier, cuter, and expressive, while preserving interaction ideas such as eye tracking, blinking, typing reactions, and password-visibility avoidance.
- Home page: replace the current overview with a cleaner dashboard-style workspace that uses existing dependencies and local mock/display data only where the current project already does.
- CRUD system: polish `CommonPage`, `AppPage`, `AppCard`, `MeCrud`, `QueryItem`, modal/form composition, PMS user/role/resource pages, and shared table actions.
- Empty/error/iframe/profile pages: align spacing, typography, and surface style with the new design without changing behavior.

Interaction rules:

- Clickable elements must have clear hover and focus states.
- Table rows, action buttons, menu items, tab controls, and context menus should feel stable and not shift layout on hover.
- Toolbar controls should use existing icon conventions and Naive UI controls.
- Avoid visible instructional copy that explains how the UI works unless the existing page already requires it for business context.

## Logic Cleanup Scope

Cleanup is allowed when it is directly tied to the migration or UI refactor:

- Remove unused imports, exports, variables, stale comments, and redundant wrapper functions.
- Consolidate duplicated layout styles across normal/simple/full where doing so reduces maintenance without hiding behavior.
- Simplify theme token handling while preserving persisted layout, collapsed state, primary color, and Naive UI overrides.
- Keep core route, permission, auth, tab, KeepAlive, and HTTP semantics unchanged.
- Treat demo pages carefully. They may be typed and visually aligned, but deletion requires a specific implementation-time finding that they are unreachable or intentionally out of scope.

## Architecture

Use the current architecture instead of introducing a new framework layer:

- Vue 3 Composition API remains the default.
- Pinia remains the state layer.
- Vue Router dynamic routes remain driven by permissions.
- Axios wrapper and interceptors remain the HTTP layer.
- UnoCSS remains the utility styling layer.
- Naive UI remains the component foundation.

Add type boundaries near existing module boundaries:

- `src/types` for shared app, route, auth, permission, API, CRUD, and component helper types.
- Module-local types inside view folders when they are not reused.
- Typed exports from API modules so composables and pages do not infer `any` across network boundaries.

## Data Flow

The current data flow stays intact:

1. `main` creates the Vue app, installs Pinia/directives/router, mounts, and sets up Naive UI discrete APIs.
2. Router guards fetch user info and permissions when needed.
3. Permission store converts backend permission nodes into menu options and access routes.
4. Layout components read app, tab, user, and permission stores.
5. Views call typed API modules through the existing HTTP wrapper.
6. CRUD views use shared composables and components for query, table, modal, and form operations.

## Error Handling

Preserve current error handling behavior:

- HTTP interceptors continue to normalize response, auth, and notification behavior.
- Permission guard keeps the distinction between unauthorized `403` and missing `404`.
- Login and logout continue to route through the existing auth store behavior.
- UI refresh should improve empty/loading/error presentation where the existing component already exposes those states.

## Testing and Verification

Minimum verification for implementation:

- Install or reuse dependencies without changing package manager assumptions.
- Run type checking once a TypeScript script exists.
- Run lint or build checks available in `package.json`.
- Run `npm run build` or the project-equivalent build command.
- Open the app in a browser and manually verify login page, main shell, sidebar collapse, tabs, route navigation, theme switch, layout setting, and at least one CRUD page.
- Check responsive behavior at mobile, tablet, desktop, and wide desktop widths.

## Implementation Phases

1. TypeScript foundation: config, env declarations, source renames, import fixes, and build script updates.
2. Shared types: route meta, permissions, user/auth, tabs, API response, pagination, CRUD, modal/form helpers.
3. Core modules: bootstrap, settings, router, guards, store, directives, HTTP utilities, API modules.
4. Shared shell UI: layouts, sidebar/header/tabs/breadcrumb/user actions/theme/layout settings using same-tone sidebar design.
5. Shared UI primitives and composables: page/card/crud/modal/form hooks and common components.
6. Business views: login, home, profile, PMS pages, iframe, base/demo pages, error pages.
7. Cleanup and verification: unused code removal, formatting, typecheck/build/browser checks.

## Non-Goals

- Do not replace Naive UI, Pinia, Vue Router, Vite, or UnoCSS.
- Do not redesign the app as a landing page.
- Do not change backend API URLs, payload semantics, or permission model.
- Do not remove core functionality for speed.
- Do not introduce a new state machine, schema validation library, UI framework, or route generation framework unless a concrete blocker appears during implementation.
- Do not copy the external animated characters verbatim. Use the external project as the login-page interaction/design reference and create project-local character styling.

## Open Constraint

This workspace is not currently a Git repository, so the design document cannot be committed from this environment. The file is still written in the requested spec location for review and later version control.
