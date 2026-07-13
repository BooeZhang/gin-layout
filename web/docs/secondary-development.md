# 二次开发文档

本文档面向接手当前项目的开发者，说明如何在现有 Vue Naive Admin 基础上新增页面、配置菜单权限、对接后端、调整主题和部署上线。

## 1. 开发前准备

安装依赖：

```bash
pnpm install
```

启动开发服务：

```bash
pnpm run dev
```

类型检查：

```bash
pnpm run typecheck
```

生产构建：

```bash
pnpm run build
```

## 2. 关键目录职责

```text
src/views/              页面模块，动态菜单组件路径来自这里
src/api/                通用接口聚合
src/views/**/api.ts     页面模块自己的接口
src/router/             基础路由和路由守卫
src/store/modules/      Pinia 业务状态
src/layouts/            full 和 empty 两类布局
src/components/common/  通用基础组件
src/components/me/      业务封装组件，如 CRUD、Modal
src/composables/        useForm、useCrud、useModal 等组合式函数
src/settings.ts         默认布局、主题色、基础菜单权限
src/utils/http/         Axios 实例、拦截器和请求辅助函数
build/                  自定义 Vite 插件，不是构建产物
```

## 3. 新增页面

建议按业务模块创建目录：

```text
src/views/order/list/index.vue
src/views/order/list/api.ts
```

页面组件示例：

```vue
<template>
  <CommonPage title="订单列表">
    <AppCard>
      订单列表内容
    </AppCard>
  </CommonPage>
</template>

<script setup lang="ts">
defineOptions({ name: 'OrderList' })
</script>
```

注意事项：

- `defineOptions({ name: '...' })` 建议和后端权限 `code` 保持可识别关系，便于 KeepAlive 和调试。
- 页面文件必须位于 `src/views/**/*.vue`，资源管理中的组件路径才能选到。
- 页面内接口建议放在同目录 `api.ts`，避免所有接口堆到全局 `src/api/index.ts`。

## 4. 配置菜单和权限

项目菜单由权限树驱动。登录后会请求：

```text
GET /role/permissions/tree
```

典型菜单记录：

```json
{
  "code": "OrderList",
  "name": "订单列表",
  "type": "MENU",
  "path": "/order/list",
  "component": "/src/views/order/list/index.vue",
  "icon": "i-fe:list",
  "order": 10,
  "enable": true,
  "show": true,
  "keepAlive": true
}
```

字段说明：

- `code`：路由名称和权限编码，需唯一。
- `name`：菜单标题。
- `type`：`MENU` 表示菜单，`BUTTON` 表示按钮权限。
- `path`：前端访问路径。
- `component`：页面组件路径。
- `icon`：菜单图标类名。
- `order`：排序值，越小越靠前。
- `enable`：是否启用。
- `show`：是否显示在菜单中。
- `keepAlive`：是否缓存页面组件。

如果只需要本地固定菜单，可在 `src/settings.ts` 的 `basePermissions` 中添加菜单记录。线上项目建议由后端统一管理权限树。

## 5. 按钮权限

按钮权限作为菜单的 `children` 返回，`type` 为 `BUTTON`：

```json
{
  "code": "OrderCreate",
  "name": "新增订单",
  "type": "BUTTON",
  "enable": true,
  "show": true
}
```

动态路由生成时，按钮权限会进入路由 `meta.btns`。页面可以根据当前路由的按钮权限决定按钮是否展示。

## 6. 对接后端

请求实例位于：

```text
src/utils/http/index.ts
src/utils/http/interceptors.ts
src/utils/http/helpers.ts
```

通用接口位于：

```text
src/api/index.ts
```

当前关键接口：

```ts
getUser: () => request.get('/user/detail')
refreshToken: () => request.get('/auth/refresh/token')
logout: () => request.post('/auth/logout', {}, { needTip: false })
switchCurrentRole: role => request.post(`/auth/current-role/switch/${role}`)
getRolePermissions: () => request.get('/role/permissions/tree')
validateMenuPath: path => request.get(`/permission/menu/validate?path=${path}`)
```

开发环境可使用代理：

```env
VITE_AXIOS_BASE_URL='/api'
VITE_PROXY_TARGET='http://localhost:8085'
```

也可直接使用 Mock 地址：

```env
VITE_AXIOS_BASE_URL='https://m1.apifoxmock.com/m1/3776410-3408296-default'
```

## 7. 登录和用户状态

登录页位于：

```text
src/views/login/index.vue
src/views/login/api.ts
```

用户状态位于：

```text
src/store/modules/user.ts
src/store/modules/auth.ts
```

常见改造点：

- 替换登录接口地址和参数。
- 调整 token 字段名。
- 接入验证码或多租户参数。
- 修改刷新 token 的时机和失败后的跳转策略。

修改登录流程后，应同时验证登录、刷新、退出、切换角色和无权限跳转。

## 8. 布局开发

当前主布局：

```text
src/layouts/full/index.vue
src/layouts/full/header/index.vue
src/layouts/full/sidebar/index.vue
```

空布局：

```text
src/layouts/empty/index.vue
```

布局选择逻辑在 `src/App.vue`：

- `meta.layout = 'empty'` 使用空布局。
- 其他页面使用 `src/settings.ts` 中的 `defaultLayout`。

如果需要新增布局，应同步扩展：

- `src/layouts/<layout-name>/index.vue`
- `src/types/app.ts` 中的 `LayoutName`
- `src/settings.ts` 中的默认布局或路由 `meta.layout`

## 9. 主题和样式

主题色默认值位于：

```text
src/settings.ts
```

```ts
export const defaultPrimaryColor = '#316C72'
```

全局样式：

```text
src/styles/global.css
src/styles/reset.css
```

Naive UI 主题覆盖由 `src/store/modules/app.ts` 管理，`setThemeColor` 会根据主色生成 hover、pressed 等色值。

开发建议：

- 通用样式放到 `src/styles/global.css`。
- 页面局部样式保留在对应 `.vue` 文件中。
- 优先使用 UnoCSS 原子类，减少零散 CSS。

## 10. 图标开发

图标来源：

```text
src/assets/icons/feather/  本地 feather 图标
src/assets/icons/isme/     项目自定义图标
src/assets/icons/dynamic-icons.ts
```

使用示例：

```vue
<i class="i-fe:user" />
<i class="i-me:docs" />
<i class="i-mdi:upload" />
```

如果新增本地 SVG 图标：

1. 放入 `src/assets/icons/isme/` 或 `src/assets/icons/feather/`。
2. 使用 `i-me:<文件名>` 或 `i-fe:<文件名>`。
3. 运行 `pnpm run build` 验证图标是否被 UnoCSS 正确识别。

## 11. 已裁剪示例页面

项目已移除基础示例、业务示例和外链内嵌页面。若后端权限树仍返回这些页面或外部链接，前端会在 `src/store/helper.ts` 中过滤掉对应菜单，避免注册无效路由。

确认后可以删除对应目录，并运行：

```bash
pnpm run typecheck
pnpm run build
```

## 12. 构建和部署

构建：

```bash
pnpm run build
```

产物目录：

```text
dist/
```

部署时根据服务器路径设置：

```env
VITE_PUBLIC_PATH='/'
```

如果应用部署在子路径，例如 `/admin/`：

```env
VITE_PUBLIC_PATH='/admin/'
```

同时确认服务器已正确处理前端路由。若使用 hash 路由，可设置：

```env
VITE_USE_HASH='true'
```

## 13. 验证清单

每次二开提交前建议执行：

```bash
pnpm run typecheck
pnpm run build
```

涉及 UI 或权限时，手动验证：

- 登录和退出。
- 菜单是否正常显示。
- 动态页面是否能打开。
- 403、404 是否正常。
- 刷新页面后路由是否恢复。
- 主题切换、侧边栏折叠和标签页是否正常。

## 14. 常见问题

### 菜单显示了但页面打不开

检查后端返回的 `component` 是否是有效的 `src/views/**/*.vue` 路径。

### 菜单图标不显示

检查 `icon` 字段是否是有效类名，例如 `i-fe:user`、`i-me:docs`、`i-mdi:upload`。

### 本地代理没有生效

确认 `.env.development`：

```env
VITE_AXIOS_BASE_URL='/api'
VITE_PROXY_TARGET='http://localhost:8085'
```

同时确认后端服务已启动。

### 构建后资源路径错误

检查 `.env.production` 中的 `VITE_PUBLIC_PATH` 是否和部署路径一致。

### 修改权限后前端仍旧显示旧菜单

清理浏览器缓存和登录态后重新登录。项目中部分状态使用 Pinia 持久化，旧权限可能仍留在会话中。
