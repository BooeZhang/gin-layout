import { normalizeTableRows } from './data'

function expect(condition: unknown, message: string) {
  if (!condition)
    throw new Error(message)
}

const rows = normalizeTableRows({
  items: null,
  total: 0,
  page: 1,
  pageSize: 20,
})

expect(Array.isArray(rows), 'items: null should normalize to an array')
expect(rows.length === 0, 'items: null should clear the table rows')
