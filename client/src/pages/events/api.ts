// Types for events pages.

// API return from GET /api/events/${id}.
export type GetEvent = {
  event: GetEventEvent
}

// API return from GET /api/events/${id}?edit=true
export type GetEventEdit = {
  event: GetEventEvent
  types: Array<string>
  roles: Array<GetEventRole>
  venues: Array<GetEventVenue>
  orgs: Array<string>
}

// API return from GET /api/events/${id}?attendance/true
export type GetEventAttendance = {
  event: GetEventEvent
  people: Array<GetEventPerson>
}

// API return from POST /api/events/${id}
export type PostEvent = {
  // It will be one or the other of:
  id?: number
  nameError?: true
}

export type GetEventEvent = {
  id: number
  name: string
  date: string
  start: string
  end: string
  venue: GetEventVenue
  details: string
  coveredByDSW: boolean
  org: string
  type: string
  roles: Array<number>
  canEdit: boolean
  canAttendance: boolean
  canEditDSWFlags: boolean
}
export type GetEventPerson = {
  id: number
  sortName: string
  attended: false | GetEventPersonAttendance
}
export type GetEventPersonAttendance = {
  type: string
  minutes: number
}
export type GetEventRole = {
  id: number
  name: string
  org: string
}
export type GetEventVenue = {
  id: number
  name: string
  // Included in GetEventEvent.venue but not in GetEvent.venues:
  address?: string
  city?: string
  url?: string
}

export type GuestOption = {
  id: number
  sortName: string
}
