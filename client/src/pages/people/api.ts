// Return type from GET /api/people/${id}/hours/${month}.
export type GetPersonHours = {
  id: number
  name: string
  needsVolgistics: boolean
  events: Array<GetPersonHoursEvent>
  // Added locally for use by People.vue
  canHours: boolean
}

export type GetPersonHoursEvent = {
  id: number
  date: string // YYYY-MM
  name: string
  minutes: number
  type: string
  placeholder: boolean
  canViewType: boolean
  canEdit: boolean
  renewsDSW: boolean
}

// Utility function, here for convenience:
export function fmtMinutes(m: number): string {
  if (!m) return ''
  const hours = Math.floor(m / 60).toString()
  return m === 30
    ? '½ hour'
    : m === 60
    ? '1 hour'
    : m % 60 !== 0
    ? `${hours}½ hours`
    : `${hours} hours`
}