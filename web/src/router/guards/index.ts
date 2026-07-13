import type { Router } from 'vue-router'
import { createAccessGuard } from './access-guard'
import { createPageLoadingGuard } from './page-loading-guard'
import { createPageTitleGuard } from './page-title-guard'
import { createTabGuard } from './tab-guard'

export function setupRouterGuards(router: Router): void {
  createPageLoadingGuard(router)
  createAccessGuard(router)
  createPageTitleGuard(router)
  createTabGuard(router)
}
