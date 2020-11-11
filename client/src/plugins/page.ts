import { inject } from "vue"

export type PageData = {
  title: string,
  browserTitle?: string,
  subtitle?: string,
  padding?: boolean,
  menuItem?: string,
}

let fn: (data: PageData) => void

export default function setPage(data: PageData) {
  if (!fn) fn = inject<(data: PageData) => void>('setPage')!
  if (fn) fn(data)
}
