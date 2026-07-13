import type { RouteRecordRaw } from 'vue-router'

export function resolveViewComponent(
  component: string,
  routeComponents: Record<string, RouteRecordRaw['component']>,
): RouteRecordRaw['component'] | undefined {
  const componentPath = component.startsWith('/src/')
    ? component
    : `/src/views/${component.replace(/^@\/views\/|^\/?src\/views\/|^\//, '')}`
  return routeComponents[componentPath]
}
