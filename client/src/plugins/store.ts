import { reactive, App } from "vue"

interface PageData {
  title: string
  subtitle?: string
}

const store = {
  state: reactive({
    eventsYear: 0,
    page: { title: '' } as PageData,
  }),
  eventsYear(year: number) { this.state.eventsYear = year },
  setPage(page: PageData) { this.state.page = page },
}

export default function (app: App) {
  app.provide('store', store)
}
