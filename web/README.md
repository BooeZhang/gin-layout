# Vue Naive Admin

基于 Vue 3、Vite、Naive UI、Pinia 和 UnoCSS 的后台管理前端模板。当前项目已保留核心后台能力，适合作为中后台系统的二次开发基础。

## 技术栈

- Vue 3 + Vite
- Vue Router
- Pinia + pinia-plugin-persistedstate
- Naive UI
- UnoCSS + Iconify
- Axios
- TypeScript

## 环境要求

- Node.js 20 或更高版本
- pnpm 9 或更高版本

## 快速开始

```bash
pnpm install
pnpm run dev
```

默认开发服务端口在 `vite.config.js` 中配置为 `3200`。

## 常用命令

```bash
pnpm run dev        # 启动开发服务
pnpm run typecheck  # Vue/TypeScript 类型检查
pnpm run build      # 生产构建
pnpm run preview    # 预览生产构建结果
pnpm run lint:fix   # ESLint 自动修复
```

## 目录结构

```text
build/                  自定义 Vite 插件，生成页面路径和动态图标 safelist
docs/                   项目文档
public/                 静态公共资源
src/api/                通用接口聚合
src/assets/             图片、自定义图标等资源
src/components/         通用组件和业务组件
src/composables/        组合式函数
src/directives/         全局指令
src/layouts/            页面布局
src/router/             基础路由与路由守卫
src/store/              Pinia 状态管理
src/styles/             全局样式与重置样式
src/types/              全局类型声明
src/utils/              请求、存储、通用工具
src/views/              页面模块
```

## 环境变量

环境变量位于 `.env`、`.env.development`、`.env.production`。

```env
VITE_TITLE='Vue Naive Admin'
VITE_USE_HASH='true'
VITE_PUBLIC_PATH='/'
VITE_AXIOS_BASE_URL='https://m1.apifoxmock.com/m1/3776410-3408296-default'
VITE_PROXY_TARGET='http://localhost:8085'
```

- `VITE_TITLE`：浏览器标题。
- `VITE_USE_HASH`：是否使用 hash 路由。
- `VITE_PUBLIC_PATH`：静态资源公共路径。
- `VITE_AXIOS_BASE_URL`：Axios 请求基础地址。
- `VITE_PROXY_TARGET`：开发代理目标地址。

## 权限与菜单

项目采用后端权限树驱动菜单和动态路由：

- 用户信息来自 `src/api/index.ts` 中的 `/user/detail`。
- 权限树来自 `/role/permissions/tree`。
- `src/store/helper.ts` 合并 `src/settings.ts` 中的 `basePermissions` 和后端权限。
- `src/store/modules/permission.ts` 将权限记录转换为菜单和动态路由。
- 菜单记录中的 `component` 指向 `src/views/**/*.vue` 页面路径。

## 布局

当前保留两类布局：

- `src/layouts/full`：后台主布局，包含侧边栏、顶部栏、标签页和内容区。
- `src/layouts/empty`：空布局，用于登录页、错误页等独立页面。

默认布局配置在 `src/settings.ts`：

```ts
export const defaultLayout = 'full'
```

路由 `meta.layout = 'empty'` 时使用空布局，其余页面使用默认布局。

## 图标

项目同时使用 Iconify 图标和本地自定义 SVG 图标：

- `src/assets/icons/feather`：本地 feather 图标库。
- `src/assets/icons/isme`：项目自定义图标。
- `src/assets/icons/dynamic-icons.ts`：需要动态 safelist 的图标。
- `build/index.js` 会扫描本地图标目录并生成 UnoCSS safelist。

菜单图标建议使用 `i-fe:*`、`i-me:*` 或 Iconify 类名。

## 构建说明

```bash
pnpm run build
```

构建产物输出到 `dist/`。该目录是生成产物，不需要手动维护，可删除后重新构建。

## 二次开发

二次开发流程、页面新增、菜单权限、后端对接、主题调整和部署说明见：

[docs/secondary-development.md](./docs/secondary-development.md)

## 许可证

本项目基于 MIT 协议授权，详见 [LICENSE](./LICENSE)。
