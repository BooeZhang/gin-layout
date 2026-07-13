import dayjs from 'dayjs'

type AnyFn<T = unknown> = (this: unknown, ...args: unknown[]) => T
type TimeInput = Parameters<typeof dayjs>[0]

/**
 * @param {(object | string | number)} time
 * @param {string} format
 * @returns {string | null} 格式化后的时间字符串
 *
 */
export function formatDateTime(time: TimeInput = undefined, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjs(time).format(format)
}

export function formatDate(date: TimeInput = undefined, format = 'YYYY-MM-DD'): string {
  return formatDateTime(date, format)
}

/**
 * @param {Function} fn
 * @param {number} wait
 * @returns {Function}  节流函数
 *
 */
export function throttle<T extends AnyFn>(fn: T, wait: number): (this: ThisParameterType<T>, ...args: Parameters<T>) => void {
  let context: ThisParameterType<T>
  let args: Parameters<T>
  let previous = 0

  return function (this: ThisParameterType<T>, ...argArr: Parameters<T>) {
    const now = +new Date()
    context = this
    args = argArr
    if (now - previous > wait) {
      fn.apply(context, args)
      previous = now
    }
  }
}

/**
 * @param {Function} method
 * @param {number} wait
 * @param {boolean} immediate
 * @return {*} 防抖函数
 */
export function debounce<T extends AnyFn>(
  method: T,
  wait: number,
  immediate = false,
): (this: ThisParameterType<T>, ...args: Parameters<T>) => void {
  let timeout: ReturnType<typeof setTimeout> | null = null
  return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
    const context = this
    if (timeout) {
      clearTimeout(timeout)
    }
    // 立即执行需要两个条件，一是immediate为true，二是timeout未被赋值或被置为null
    if (immediate) {
      /**
       * 如果定时器不存在，则立即执行，并设置一个定时器，wait毫秒后将定时器置为null
       * 这样确保立即执行后wait毫秒内不会被再次触发
       */
      const callNow = !timeout
      timeout = setTimeout(() => {
        timeout = null
      }, wait)
      if (callNow) {
        method.apply(context, args)
      }
    }
    else {
      // 如果immediate为false，则函数wait毫秒后执行
      timeout = setTimeout(() => {
        /**
         * args是一个类数组对象，所以使用fn.apply
         * 也可写作method.call(context, ...args)
         */
        method.apply(context, args)
      }, wait)
    }
  }
}

/**
 * @param {number} time 毫秒数
 * @returns 睡一会儿，让子弹暂停一下
 */
export function sleep(time: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, time))
}

/**
 * @param {HTMLElement} el
 * @param {Function} cb
 * @return {ResizeObserver}
 */
export function useResize(el: HTMLElement, cb: (rect: DOMRectReadOnly) => void): ResizeObserver {
  const observer = new ResizeObserver((entries) => {
    cb(entries[0].contentRect)
  })
  observer.observe(el)
  return observer
}
