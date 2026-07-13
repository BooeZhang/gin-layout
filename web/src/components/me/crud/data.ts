export function normalizeTableRows(data: unknown): unknown[] {
  if (Array.isArray(data))
    return data

  if (data && typeof data === 'object' && 'items' in data) {
    const { items } = data as { items?: unknown }
    return Array.isArray(items) ? items : []
  }

  return []
}
