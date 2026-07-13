import type { App, Directive, DirectiveBinding, VNode } from 'vue'
import { withDirectives } from 'vue'
import { router } from '@/router'
import type { AccessButtonRecord } from '@/types/app'

const access: Directive<HTMLElement, string> = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string>) {
    const currentRoute = unref(router.currentRoute)
    const btns = (currentRoute.meta?.btns as AccessButtonRecord[] | undefined)?.map(item => item.code) || []
    if (!btns.includes(binding.value)) {
      el.remove()
    }
  },
}

export function setupDirectives(app: App): void {
  app.directive('access', access)
}

/**
 * 用于h函数使用自定义权限指令
 *
 * @param {*} vnode 虚拟节点
 * @param {*} code 权限码
 * @returns 返回一个包含权限指令的vnode
 *
 * 使用示例：withAccess(h('button', {class: 'text-red-500'}, '删除'), 'user:delete')
 *
 */
export function withAccess(vnode: VNode, code: string): VNode {
  return withDirectives(vnode, [[access, code]])
}
