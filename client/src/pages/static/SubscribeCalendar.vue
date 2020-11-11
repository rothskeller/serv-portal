<!--
SubscribeCalendar displays a page with information about how to subscribe to the
events calendar.
-->

<template lang="pug">
#subscribe-calendar
  h2 SERV Calendar Subscription
  p.
    You can subscribe to the SERV calendar so that SERV events will
    automatically appear in the calendar app on your phone, or in your desktop
    calendar software. Please see the instructions for your phone or software
    below.
  h3 iPhone or iPad Calendar App
  ol
    li Open the Settings app.
    li Go to “Passwords &amp; Accounts”.
    li Tap on “Other”.
    li Tap on “Add Subscribed Calendar”.
    li In the “Server” field, enter <tt>https://sunnyvaleserv.org/calendar.ics</tt>.
    li Tap “Next”.
    li Optional: change the “Description” field to a name that’s meaningful to you, such as “SERV Calendar”.
    li Tap “Save”.
  h3 Google Calendar (including Android Phones)
  ol
    li In a web browser, go to Google Calendar (<tt>https://calendar.google.com</tt>). Log in if necessary.
    li In the left sidebar, click the large “+” sign next to “Other Calendars”.
    li Click “From URL”.
    li In the “URL of calendar” field, enter <tt>https://sunnyvaleserv.org/calendar.ics</tt>.
    li Click “Add calendar”.
  h3 Microsoft Outlook
  ol
    li Open Microsoft Outlook.
    li Go to the calendar page.
    li In the Home ribbon, click on “Open Calendar”, then “From Internet”.
    li Enter <tt>https://sunnyvaleserv.org/calendar.ics</tt>.
    li Click “Yes”.
    li In the left sidebar, under “Other Calendars”, right-click on “Untitled” and choose “Rename Calendar”.
    li Give the calendar a name meaningful to you, such as “SERV Calendar”.
  h3 Mac Calendar App
  ol
    li Open the Calendar app.
    li From the menu, choose File → New Calendar Subscription.
    li Enter <tt>https://sunnyvaleserv.org/calendar.ics</tt>.
    li Click “Subscribe”.
    li Set the options to suit your preferences and click “OK”.
  h3 Other Software
  p.
    Most calendar software has the ability to subscribe to Internet calendars.
    Consult the documentation for your software to find out how.  The address
    of the SERV calendar is <tt>https://sunnyvaleserv.org/calendar.ics</tt>.
  div(style='margin: 1.5rem 0')
    SButton(@click='$router.go(-1)', variant='primary') Back
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import setPage from '../../plugins/page'
import SButton from '../../base/SButton.vue'

export default defineComponent({
  components: { SButton },
  setup() {
    setPage({ title: 'SERV Calendar Subscription' })
  },
})
</script>

<style lang="postcss">
#subscribe-calendar {
  margin: 1.5rem 0.75rem;
  max-width: 40rem;
}
</style>
