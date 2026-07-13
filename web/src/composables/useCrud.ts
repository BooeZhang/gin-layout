import type { DialogOptions } from 'naive-ui'
import { cloneDeep } from 'lodash-es'
import type { Ref } from 'vue'
import { useForm, useModal } from '.'
import type { Id } from '@/types/app'

const ACTIONS = {
  view: '查看',
  edit: '编辑',
  add: '新增',
} as const

type CrudAction = keyof typeof ACTIONS
type CrudForm = Record<string, unknown>
type CrudApi<T extends CrudForm> = (data: T) => Promise<unknown>

interface CrudOptions<T extends CrudForm> {
  name: string
  initForm?: T
  doCreate: CrudApi<T>
  doDelete: (id: Id) => Promise<unknown>
  doUpdate: CrudApi<T>
  refresh: (data?: unknown, deleted?: boolean) => void
}

interface OpenOptions<T extends CrudForm> {
  action?: CrudAction
  row?: Partial<T>
  title?: string
  onOk?: () => Promise<unknown> | unknown
  [key: string]: unknown
}

interface SaveAction {
  api: () => Promise<unknown>
  cb: () => void
}

export function useCrud<T extends CrudForm = CrudForm>({
  name,
  initForm = {} as T,
  doCreate,
  doDelete,
  doUpdate,
  refresh,
}: CrudOptions<T>) {
  const modalAction = ref<CrudAction | ''>('')
  const [modalRef, okLoading] = useModal()
  const [modalFormRef, modalForm, validation] = useForm(initForm)

  /** 新增 */
  function handleAdd(row: Partial<T> = {}, title?: string) {
    handleOpen({ action: 'add', title, row: Object.assign({}, cloneDeep(initForm), cloneDeep(row)) })
  }

  /** 修改 */
  function handleEdit(row: Partial<T>, title?: string) {
    handleOpen({ action: 'edit', title, row })
  }

  /** 查看 */
  function handleView(row: Partial<T>, title?: string) {
    handleOpen({ action: 'view', title, row })
  }

  /** 打开modal */
  function handleOpen(options: OpenOptions<T> = {}) {
    const { action, row, title, onOk } = options
    modalAction.value = action || ''
    modalForm.value = { ...row } as T
    modalRef.value?.open({
      ...options,
      async onOk() {
        if (typeof onOk === 'function') {
          return await onOk()
        }
        else {
          return await handleSave()
        }
      },
      title: title ?? `${modalAction.value ? ACTIONS[modalAction.value] : ''}${name}`,
    })
  }

  /** 保存 */
  async function handleSave(action?: SaveAction) {
    if (!action && !['edit', 'add'].includes(modalAction.value)) {
      return false
    }
    await validation()
    const actions: Record<'add' | 'edit', SaveAction> = {
      add: {
        api: () => doCreate(modalForm.value),
        cb: () => $message.success('新增成功'),
      },
      edit: {
        api: () => doUpdate(modalForm.value),
        cb: () => $message.success('保存成功'),
      },
    }

    action = action || actions[modalAction.value as 'add' | 'edit']

    try {
      okLoading.value = true
      const data = await action.api()
      action.cb()
      okLoading.value = false
      data && refresh(data)
    }
    catch (error) {
      console.error(error)
      okLoading.value = false
      return false
    }
  }

  /** 删除 */
  function handleDelete(id: Id | undefined | null, confirmOptions?: DialogOptions) {
    if (!id && id !== 0)
      return
    const d = $dialog.warning({
      content: '确定删除？',
      title: '提示',
      positiveText: '确定',
      negativeText: '取消',
      async onPositiveClick() {
        try {
          d.loading = true
          const data = await doDelete(id)
          $message.success('删除成功')
          d.loading = false
          refresh(data, true)
        }
        catch (error) {
          console.error(error)
          d.loading = false
        }
      },
      ...confirmOptions,
    })
  }

  return {
    modalRef,
    modalFormRef,
    modalAction,
    modalForm,
    okLoading,
    validation,
    handleAdd,
    handleDelete,
    handleEdit,
    handleView,
    handleOpen,
    handleSave,
  }
}
