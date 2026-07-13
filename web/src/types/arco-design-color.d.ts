declare module '@arco-design/color' {
  export function generate(color: string, options?: { list?: boolean, dark?: boolean }): string[]
  export function getRgbStr(color: string): string
}
