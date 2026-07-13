import { resolveViewComponent } from './route-component'

function expect(condition: unknown, message: string) {
  if (!condition)
    throw new Error(message)
}

const userComponent = () => Promise.resolve({})
const routeComponents = {
  '/src/views/pms/user/index.vue': userComponent,
}

expect(
  resolveViewComponent('/src/views/pms/user/index.vue', routeComponents) === userComponent,
  'absolute /src view path should resolve',
)
expect(
  resolveViewComponent('@/views/pms/user/index.vue', routeComponents) === userComponent,
  '@/views path should resolve',
)
expect(
  resolveViewComponent('pms/user/index.vue', routeComponents) === userComponent,
  'relative view path should resolve',
)
