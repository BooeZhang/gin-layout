import type { Router } from 'vue-router'

export function createPageLoadingGuard(router: Router): void {
  router.beforeEach(() => {
    $loadingBar.start()
  })

  router.afterEach(() => {
    setTimeout(() => {
      $loadingBar.finish()
    }, 200)
  })

  router.onError(() => {
    $loadingBar.error()
  })
}
