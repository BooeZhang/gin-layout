# Admin TS UI Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Convert the Vue Naive Admin frontend toward TypeScript while delivering the approved same-tone admin shell and the full animated login page based on `login-animation` with redesigned cute characters.

**Architecture:** Keep the existing Vue 3 + Vite + Pinia + Vue Router + Naive UI + UnoCSS architecture. Add TS foundations and shared types first, then convert high-risk shared modules, then redesign the login page and shell UI without changing auth, permission, route, tab, and API semantics.

**Tech Stack:** Vue 3, Vite, TypeScript, Pinia, Vue Router, Naive UI, UnoCSS, Axios, ECharts.

---

## File Structure

- Create `tsconfig.json`: TypeScript project settings for Vite/Vue.
- Create `src/types/env.d.ts`: declarations for Vue SFC, assets, virtual imports, and app globals.
- Create `src/types/app.ts`: shared app, route, auth, permission, API, tab, and CRUD types.
- Modify `package.json`: add `typecheck`, `typescript`, and `vue-tsc`.
- Replace `jsconfig.json`: remove JS-only config once `tsconfig.json` exists.
- Rename `src/main.js` to `src/main.ts`: app bootstrap.
- Rename `src/settings.js` to `src/settings.ts`: layout, theme, and base permissions.
- Rename `src/views/login/api.js` to `src/views/login/api.ts`: typed login API.
- Modify `src/views/login/index.vue`: full animated login redesign, `script setup lang="ts"`, unchanged auth flow.
- Create `src/views/login/components/CuteLoginCharacters.vue`: redesigned interactive characters inspired by `login-animation`.
- Modify layout shell files under `src/layouts` and `src/layouts/components`: same-tone sidebar and modern shell polish.
- Modify shared page/card/CRUD components under `src/components`: aligned spacing, stable controls, and TS script blocks where touched.

## Task 1: TypeScript Foundation

**Files:**
- Create: `tsconfig.json`
- Create: `src/types/env.d.ts`
- Create: `src/types/app.ts`
- Modify: `package.json`
- Delete: `jsconfig.json`

- [ ] **Step 1: Add TypeScript dependencies and scripts**

Modify `package.json`:

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview",
    "typecheck": "vue-tsc --noEmit",
    "lint:fix": "eslint --fix",
    "postinstall": "npx simple-git-hooks",
    "up": "taze major -I"
  },
  "devDependencies": {
    "typescript": "^5.9.3",
    "vue-tsc": "^3.1.4"
  }
}
```

- [ ] **Step 2: Add `tsconfig.json`**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "strict": false,
    "jsx": "preserve",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "esModuleInterop": true,
    "allowJs": true,
    "checkJs": false,
    "skipLibCheck": true,
    "noEmit": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "~/*": ["*"]
    },
    "types": ["vite/client"]
  },
  "include": [
    "src/**/*.ts",
    "src/**/*.tsx",
    "src/**/*.vue",
    "vite.config.*",
    "uno.config.*",
    "build/**/*.js",
    "src/types/**/*.d.ts"
  ],
  "exclude": ["dist", "node_modules"]
}
```

- [ ] **Step 3: Add environment declarations**

Create `src/types/env.d.ts`:

```ts
/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'

  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, any>
  export default component
}

declare module '*.webp'
declare module '*.png'
declare module '*.svg'
declare module 'virtual:page-pathes' {
  const pagePathes: Record<string, () => Promise<unknown>>
  export default pagePathes
}
declare module 'virtual:icons' {
  const icons: Record<string, string>
  export default icons
}

interface ImportMetaEnv {
  readonly VITE_TITLE: string
  readonly VITE_AXIOS_BASE_URL: string
  readonly VITE_PUBLIC_PATH: string
  readonly VITE_PROXY_TARGET: string
  readonly VITE_USE_HASH: string
}
```

- [ ] **Step 4: Add shared domain types**

Create `src/types/app.ts` with route, permission, auth, API, and CRUD types used by migrated modules.

- [ ] **Step 5: Remove JS-only config**

Delete `jsconfig.json` after confirming `tsconfig.json` contains the alias paths.

## Task 2: Core TS Conversion

**Files:**
- Rename: `src/main.js` to `src/main.ts`
- Rename: `src/settings.js` to `src/settings.ts`
- Rename: `src/views/login/api.js` to `src/views/login/api.ts`
- Modify imports pointing at these modules only if an explicit extension exists.

- [ ] **Step 1: Rename files**

Use native file rename commands for `main`, `settings`, and login API.

- [ ] **Step 2: Convert `main.ts`**

Add `import type { App as VueApp } from 'vue'` only if a helper needs it. Keep bootstrap behavior unchanged.

- [ ] **Step 3: Convert `settings.ts`**

Type `basePermissions` as `PermissionRecord[]` from `src/types/app.ts`.

- [ ] **Step 4: Convert `login/api.ts`**

Type login payload and response with `LoginPayload` and `TokenPayload`.

## Task 3: Animated Login Page

**Files:**
- Create: `src/views/login/components/CuteLoginCharacters.vue`
- Modify: `src/views/login/index.vue`
- Modify: `src/views/login/api.ts`

- [ ] **Step 1: Preserve auth behavior**

Before changing layout, identify the existing behavior to preserve:

```ts
quickLogin()
handleLogin(isQuick?: boolean)
initCaptcha()
onLoginSuccess(data?: TokenPayload)
```

- [ ] **Step 2: Add `CuteLoginCharacters.vue`**

Implement a local Vue component with props:

```ts
interface Props {
  focusField?: 'username' | 'password' | 'captcha' | ''
  passwordVisible?: boolean
  typing?: boolean
}
```

The component owns pointer tracking, blinking, cute character markup, and CSS-only animation. It must not depend on external packages.

- [ ] **Step 3: Replace login template**

Use the `login-animation`-style two-column layout:

- left animation/brand panel
- right Naive UI form panel
- same fields as current project: username, password, captcha, remember checkbox, quick login, login

- [ ] **Step 4: Convert login script to TS**

Use `script setup lang="ts"`, preserve local storage behavior, login messages, captcha refresh on error code `10003`, and redirect behavior.

- [ ] **Step 5: Validate manually**

Run dev server and check: username focus, password focus, password visibility state, captcha refresh, quick login, normal login validation, mobile layout.

## Task 4: Same-Tone Admin Shell

**Files:**
- Modify: `src/layouts/normal/index.vue`
- Modify: `src/layouts/simple/index.vue`
- Modify: `src/layouts/full/index.vue`
- Modify: `src/layouts/components/SideLogo.vue`
- Modify: `src/layouts/components/SideMenu.vue`
- Modify: `src/layouts/normal/header/index.vue`
- Modify: `src/layouts/components/tab/index.vue`

- [ ] **Step 1: Remove independent sidebar color treatment**

Set sidebar surfaces to the same visual tone as the content area. Use border and active menu state for hierarchy.

- [ ] **Step 2: Stabilize layout dimensions**

Keep existing collapsed widths and tab/header heights. Avoid hover scale effects that shift layout.

- [ ] **Step 3: Modernize header and tabs**

Keep breadcrumb, collapse, fullscreen, guide, avatar, and tabs. Improve spacing, borders, active state, and focus/hover states.

## Task 5: Shared UI and Cleanup

**Files:**
- Modify: `src/components/common/AppPage.vue`
- Modify: `src/components/common/AppCard.vue`
- Modify: `src/components/common/CommonPage.vue`
- Modify: `src/components/me/crud/index.vue`
- Modify: `src/components/me/crud/QueryItem.vue`
- Modify touched PMS pages only where needed for style/type compatibility.

- [ ] **Step 1: Align shared surfaces**

Use consistent radius, borders, padding, and table action spacing.

- [ ] **Step 2: Remove stale comments and unused imports**

Only remove code proven unused by local search or compiler/linter feedback.

- [ ] **Step 3: Keep demo pages unless proven unreachable**

Do not delete demo/base pages in this pass.

## Task 6: Verification

**Files:**
- Verify all changed files.

- [ ] **Step 1: Install dependencies if missing**

Run: `pnpm install`
Expected: dependencies installed and lockfile consistent.

- [ ] **Step 2: Run typecheck**

Run: `pnpm run typecheck`
Expected: exit code 0.

- [ ] **Step 3: Run build**

Run: `pnpm run build`
Expected: exit code 0.

- [ ] **Step 4: Browser check**

Run: `pnpm run dev`
Open: `http://localhost:3200`
Check login page, shell layout, sidebar collapse, tabs, theme switch, layout setting, and one CRUD route.

## Self-Review

- Spec coverage: TS foundation, login-animation full-page direction, same-tone sidebar, UI polish, and cleanup all map to tasks.
- Placeholder scan: no `TBD` or `TODO` markers are used as implementation steps.
- Type consistency: shared `LoginPayload`, `TokenPayload`, and `PermissionRecord` are introduced before use.

## Execution Choice

This plan will be executed inline in the current session because the user explicitly asked to start execution and did not request subagents.
