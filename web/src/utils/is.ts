const toString = Object.prototype.toString

export function is(val: unknown, type: string): boolean {
  return toString.call(val) === `[object ${type}]`
}

export function isDef<T = unknown>(val: T | undefined): val is T {
  return typeof val !== 'undefined'
}

export function isUndef(val: unknown): val is undefined {
  return typeof val === 'undefined'
}

export function isNull(val: unknown): val is null {
  return val === null
}

export function isWhitespace(val: unknown): val is '' {
  return val === ''
}

export function isObject(val: unknown): val is Record<PropertyKey, unknown> {
  return !isNull(val) && is(val, 'Object')
}

export function isArray<T = unknown>(val: unknown): val is T[] {
  return Array.isArray(val)
}

export function isString(val: unknown): val is string {
  return is(val, 'String')
}

export function isNumber(val: unknown): val is number {
  return is(val, 'Number')
}

export function isBoolean(val: unknown): val is boolean {
  return is(val, 'Boolean')
}

export function isDate(val: unknown): val is Date {
  return is(val, 'Date')
}

export function isRegExp(val: unknown): val is RegExp {
  return is(val, 'RegExp')
}

export function isFunction(val: unknown): val is (...args: unknown[]) => unknown {
  return typeof val === 'function'
}

export function isPromise<T = unknown>(val: unknown): val is Promise<T> {
  return is(val, 'Promise') && isObject(val) && isFunction(val.then) && isFunction(val.catch)
}

export function isElement(val: unknown): val is Element {
  return isObject(val) && isString(val.tagName)
}

export function isWindow(val: unknown): val is Window {
  return typeof window !== 'undefined' && isDef(window) && is(val, 'Window')
}

export function isNullOrUndef(val: unknown): val is null | undefined {
  return isNull(val) || isUndef(val)
}

export function isNullOrWhitespace(val: unknown): val is null | undefined | '' {
  return isNullOrUndef(val) || isWhitespace(val)
}

/** 空数组 | 空字符串 | 空对象 | 空Map | 空Set */
export function isEmpty(val: unknown): boolean {
  if (isArray(val) || isString(val)) {
    return val.length === 0
  }

  if (val instanceof Map || val instanceof Set) {
    return val.size === 0
  }

  if (isObject(val)) {
    return Object.keys(val).length === 0
  }

  return false
}

/**
 * 类似mysql的IFNULL函数
 *
 * @param {number | boolean | string} val
 * @param {number | boolean | string} def
 * @returns 第一个参数为null | undefined | '' 则返回第二个参数作为备用值，否则返回第一个参数
 */
export function ifNull<T, D = ''>(val: T | null | undefined | '', def: D = '' as D): NonNullable<Exclude<T, ''>> | D {
  return isNullOrWhitespace(val) ? def : val as NonNullable<Exclude<T, ''>>
}

export function isUrl(path: string): boolean {
  const reg = /^https?:\/\/[-\w+&@#/%?=~|!:,.;]+[-\w+&@#/%=~|]$/
  return reg.test(path)
}

export const isServer = typeof window === 'undefined'

export const isClient = !isServer
