export type GetClearance = {
  parameters: GetClearanceParameters
  rows: Array<GetClearanceRow>
  // Set and used locally
  justLoaded: boolean
}
export type GetClearanceParameters = {
  role: number
  with: string
  without: string
  allowedRoles: Array<{
    id: number
    name: string
  }>
  allowedRestrictions: Array<{
    value: string
    label: string
  }>
}
export type GetClearanceRow = {
  id: number
  sortName: string
  orgs: Record<
    string,
    {
      privLevel: string
      title: string
    }
  >
  dswCERT: boolean
  dswComm: boolean
  cardKey: boolean
  certShirtLS: boolean
  certShirtSS: boolean
  idPhoto: boolean
  servShirt: boolean
  volgistics: boolean
  bgCheck?: boolean
  bgCheckDOJ?: string
  bgCheckFBI?: string
  bgCheckPHS?: string
}
