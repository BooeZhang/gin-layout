import type {
  DialogApi,
  DialogOptions,
  LoadingBarApi,
  MessageApi,
  MessageOptions,
  MessageReactive,
  NotificationApi,
} from 'naive-ui'
import * as NaiveUI from 'naive-ui'
import { useAppStore } from '@/store'
import { isNullOrUndef } from '@/utils'

type MessageType = keyof Pick<MessageApi, 'loading' | 'success' | 'error' | 'info' | 'warning'>
type MessageContent = string | string[]

export interface KeyedMessageOptions extends MessageOptions {
  key?: string
}

interface ConfirmDialogOptions extends DialogOptions {
  type?: keyof Pick<DialogApi, 'warning' | 'info' | 'success' | 'error'>
  confirm?: () => void | Promise<void>
  cancel?: () => void
}

export interface MessageInstance {
  loading: (content: MessageContent, option?: KeyedMessageOptions) => unknown
  success: (content: MessageContent, option?: KeyedMessageOptions) => unknown
  error: (content: MessageContent, option?: KeyedMessageOptions) => unknown
  info: (content: MessageContent, option?: KeyedMessageOptions) => unknown
  warning: (content: MessageContent, option?: KeyedMessageOptions) => unknown
  destroy: (key: string, duration?: number) => void
}

export type DialogInstance = DialogApi & {
  confirm: (option?: ConfirmDialogOptions) => ReturnType<DialogApi['warning']>
}

export function setupMessage(NMessage: MessageApi): MessageInstance {
  class Message {
    static instance: Message | null = null
    private message: Record<string, MessageReactive> = {}
    private removeTimer: Record<string, ReturnType<typeof setTimeout>> = {}

    constructor() {
      // 单例模式
      if (Message.instance)
        return Message.instance
      Message.instance = this
    }

    removeMessage(key: string, duration = 5000): void {
      this.removeTimer[key] && clearTimeout(this.removeTimer[key])
      this.removeTimer[key] = setTimeout(() => {
        this.message[key]?.destroy()
      }, duration)
    }

    destroy(key: string, duration = 200): void {
      setTimeout(() => {
        this.message[key]?.destroy()
      }, duration)
    }

    showMessage(type: MessageType, content: MessageContent, option: KeyedMessageOptions = {}) {
      if (Array.isArray(content)) {
        return content.forEach(msg => NMessage[type](msg, option))
      }

      if (!option.key) {
        return NMessage[type](content, option)
      }

      const key = option.key
      const currentMessage = this.message[key]
      if (currentMessage) {
        currentMessage.type = type
        currentMessage.content = content
      }
      else {
        this.message[key] = NMessage[type](content, {
          ...option,
          duration: 0,
          onAfterLeave: () => {
            delete this.message[key]
          },
        })
      }
      this.removeMessage(key, option.duration)
    }

    loading(content: MessageContent, option?: KeyedMessageOptions) {
      this.showMessage('loading', content, option)
    }

    success(content: MessageContent, option?: KeyedMessageOptions) {
      this.showMessage('success', content, option)
    }

    error(content: MessageContent, option?: KeyedMessageOptions) {
      this.showMessage('error', content, option)
    }

    info(content: MessageContent, option?: KeyedMessageOptions) {
      this.showMessage('info', content, option)
    }

    warning(content: MessageContent, option?: KeyedMessageOptions) {
      this.showMessage('warning', content, option)
    }
  }

  return new Message() as unknown as MessageInstance
}

export function setupDialog(NDialog: DialogApi): DialogInstance {
  const dialog = NDialog as DialogInstance
  dialog.confirm = function (option: ConfirmDialogOptions = {}) {
    const showIcon = !isNullOrUndef(option.title)
    const type = option.type || 'warning'
    return NDialog[type]({
      showIcon,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: option.confirm,
      onNegativeClick: option.cancel,
      onMaskClick: option.cancel,
      ...option,
    })
  }

  return dialog
}

export function setupNaiveDiscreteApi(): void {
  const appStore = useAppStore()
  const configProviderProps = computed(() => ({
    theme: appStore.isDark ? NaiveUI.darkTheme : undefined,
    themeOverrides: useAppStore().naiveThemeOverrides,
  }))
  const { message, dialog, notification, loadingBar } = NaiveUI.createDiscreteApi(
    ['message', 'dialog', 'notification', 'loadingBar'],
    { configProviderProps },
  ) as { message: MessageApi, dialog: DialogApi, notification: NotificationApi, loadingBar: LoadingBarApi }

  window.$loadingBar = loadingBar
  window.$notification = notification
  window.$message = setupMessage(message)
  window.$dialog = setupDialog(dialog)
}
