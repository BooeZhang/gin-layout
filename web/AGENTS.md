# web/ — Vue 3 Admin Frontend

**Generated:** 2026-07-06

## OVERVIEW

Vue 3 + TypeScript admin panel for gin-layout backend. Uses NaiveUI, UnoCSS (Wind3 preset), Pinia stores. Separate project with own dependencies and build config.

## TECH STACK

Vue 3 (Composition API, `<script setup lang="ts">`), TypeScript, Pinia (6 stores), Vue Router (Hash/History mode), Axios (auto token refresh on 401), NaiveUI, UnoCSS, ECharts.

## STRUCTURE

```
src/
├── api/            # Axios API layer: auth, login, profile, sys/{user,role,resource}
├── assets/         # Icons (Feather + custom), images
├── components/     # me/crud (table/form/modal system) + common (AppCard, AppPage, ToggleTheme)
├── composables/    # useCrud, useForm, useModal, useAliveData
├── directives/     # v-access (button-level permission check)
├── layouts/        # full (SideBar + Header + AppTab + content), empty
├── router/         # Routes + 4 guards (page-loading → access → page-title → tab)
├── store/          # Pinia: auth, user, access, app, tab, router
├── styles/         # Global styles (reset, CSS variables)
├── types/          # TypeScript: app, auth, user, role, access, route, common, env
├── utils/          # http (axios + interceptors), storage, common helpers
└── views/          # login, home, profile, sys/{user,role,menu}, error-page (404/403)
```

## WHERE TO LOOK

| Task                   | Location                            | Notes                                                              |
| ---------------------- | ----------------------------------- | ------------------------------------------------------------------ |
| Add backend-facing API | `src/api/sys/` or `src/api/auth.ts` | Axios with auto token refresh                                      |
| New CRUD page          | `src/views/sys/`                    | Use `MeCrud` + `useCrud()` composable                              |
| Auth state             | `src/store/modules/auth.ts`         | Token pair stored in Pinia + localStorage                          |
| Route guards           | `src/router/guards/`                | `access-guard.ts` is the key guard (token→user→menus→routes)       |
| Button permission      | `v-access="'btnCode'"` directive    | Checks `route.meta.btns` array                                     |
| UI theme               | `src/store/modules/app.ts`          | Dark mode, primary color, sidebar collapse                         |
| CRUD table             | `src/components/me/crud/`           | Reusable MeCrud + MeQueryItem + MeModal                            |
| Backend API routes     | `internal/server/admin_router.go`   | `/api/auth/`, `/api/v1/users/`, `/api/v1/roles/`, `/api/v1/menus/` |

## CONVENTIONS

- Composition API throughout (`<script setup lang="ts">`).
- Pinia stores for state management (actions for mutations, getters for computed).
- Axios interceptors for token injection + silent 401 refresh (mutex pattern via `refreshTokenOnce`).
- Router guards: 4 guards in `guards/index.ts` — order matters.
- `v-access` directive for button-level permission (code compare against route meta).
- UnoCSS atomic classes (e.g., `wh-full`, `f-c-c`, `flex-col`) for styling — no inline styles.

## ANTI-PATTERNS

- No `as any` or `@ts-ignore` type suppressions.
- No direct `localStorage` access — use `utils/storage/` abstraction.
- No inline styles — use scoped SFC styles or UnoCSS utility classes.

## NOTES

- Token handling: `authStore` stores `accessToken` + `refreshToken`, synced to localStorage via `utils/storage`.
- `accessStore.setAccessMenus()` → generates dynamic routes + side menu from backend menu tree.
- `web/` is an independent project — run `pnpm install` in this directory.
- Build output: `web/dist/`. Dev proxy: `/api` → backend in `vite.config.js`.
- 前端二次开发指南: [`docs/secondary-development.md`](docs/secondary-development.md).
