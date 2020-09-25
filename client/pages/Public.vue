<!--
Public displays the public page (the root page of the site).  It's the page
shown at / when the user is not logged in.
-->

<template lang="pug">
PublicPage(title='Sunnyvale SERV')
  #public
    #public-head
      #public-head-name Sunnyvale Emergency Response&nbsp;Volunteers
      img#public-head-logo(:src='servLogo')
      #public-head-body.
        SERV is the volunteer arm of the Sunnyvale Department of Public Safety,
        Office of Emergency Services.  SERV volunteers teach disaster
        preparedness classes, assist uniformed officers in emergencies, and
        respond in disasters when professional responders are&nbsp;overloaded.
      #public-head-login
        b-btn(to='/login', variant='primary') Volunteer&nbsp;Login
    #public-folders
      #public-folder-listos.public-folder
        .public-folder-name Disaster Preparedness
        img.public-folder-logo(:src='listosLogo')
        .public-folder-body.
          <b>Listos California</b> disaster preparedness classes, taught by
          SERV volunteers, teach how to prepare your home and family for
          disasters and&nbsp;emergencies.
        .public-folder-links
          .public-folder-link
            a(
              href='https://www.eventbrite.com/o/sunnyvale-emergency-response-volunteers-11393586702',
              target='_blank'
            ) Class schedules and registration
          .public-folder-link
            router-link(to='/disaster-preparedness') Disaster preparedness materials
          .public-folder-link
            | Email <a href="mailto:listos@sunnyvaleserv.org">Listos@SunnyvaleSERV.org</a>
      #public-folder-spanish.public-folder
        .public-folder-name Preparación para desastres
        img.public-folder-logo(:src='spanishLogo')
        .public-folder-body.
          Las clases de preparación para desastres de <b>Listos
          California,</b> impartidas por voluntarios de SERV, enseñan cómo
          preparar su hogar y su familia para desastres y emergencias. Las
          clases están disponibles en&nbsp;español.
        .public-folder-links
          .public-folder-link
            a(
              href='https://www.eventbrite.com/o/sunnyvale-emergency-response-volunteers-11393586702',
              target='_blank'
            ) Horarios de clases y registro
          .public-folder-link
            router-link(to='/preparacion-para-desastres') Materiales de preparación para desastres
          .public-folder-link
            | Email <a href="mailto:listos@sunnyvaleserv.org">Listos@SunnyvaleSERV.org</a>
      #public-folder-snap.public-folder
        .public-folder-name Neighborhood Organization
        img.public-folder-logo(:src='snapLogo')
        .public-folder-body.
          <b>Sunnyvale Neighborhoods Actively Prepare (SNAP)</b> provides
          support for neighborhood organization of disaster preparedness and
          response&nbsp;volunteers.
        .public-folder-links
          .public-folder-link
            router-link(to='/neighborhood-organization') Neighborhood organization materials
          .public-folder-link
            | Email <a href="mailto:snap@sunnyvaleserv.org">SNAP@SunnyvaleSERV.org</a>
      #public-folder-cert.public-folder
        .public-folder-name Disaster Response
        img.public-folder-logo(:src='certLogo')
        .public-folder-body.
          Sunnyvale’s <b>Community Emergency Response Team (CERT)</b>
          supplements uniformed officers in responding to emergencies,
          disasters, and public service&nbsp;events.
        .public-folder-links
          .public-folder-link
            a(
              href='https://www.eventbrite.com/o/sunnyvale-emergency-response-volunteers-11393586702',
              target='_blank'
            ) Class schedules and registration
          .public-folder-link
            router-link(to='/disaster-response') Disaster response and CERT materials
          .public-folder-link
            | Email <a href="mailto:cert@sunnyvaleserv.org">CERT@SunnyvaleSERV.org</a>
      #public-folder-sares.public-folder
        .public-folder-name Emergency Communications
        img.public-folder-logo(:src='saresLogo')
        .public-folder-body.
          The <b>Sunnyvale Amateur Radio Service (SARES)</b> provides radio
          communication services during a disaster when regular methods are
          overloaded or unavailable. SARES provides a wealth of training for
          radio&nbsp;operators.
        .public-folder-links
          .public-folder-link
            a(href='https://saresrg.org/', target='_blank') SARES Repeater Group website
          .public-folder-link
            a(href='https://www.scc-ares-races.org/activities/events.php', target='_blank') Class schedules and registration
          .public-folder-link
            | Email <a href="mailto:sares@sunnyvaleserv.org">SARES@SunnyvaleSERV.org</a>
</template>

<script>
import PublicPage from '@/base/PublicPage'
import servLogo from './public/serv-logo.png'
import listosLogo from './public/listos-logo.png'
import snapLogo from './public/snap-logo.png'
import certLogo from './public/cert-logo.png'
import saresLogo from './public/sares-logo.png'
import spanishLogo from './public/spanish-logo.png'

export default {
  components: { PublicPage },
  data: () => ({
    servLogo, listosLogo, snapLogo, certLogo, saresLogo, spanishLogo,
  }),
}
</script>

<style lang="stylus">
#public-head
  display grid
  grid-template 'name name' 'logo body' 'login login' / 25% auto
  grid-row-gap 0.5rem
  grid-column-gap 1rem
  @media (min-width 504px)
    grid-template 'name name' 'logo body' 'logo login' / 120px auto
  @media (min-width 560px)
    margin 0 auto
    max-width 536px
  @media (min-width 700px)
    max-width 676px
    grid-template 'name name name' 'logo body login' / 120px auto min-content
#public-head-name
  text-align center
  font-weight bold
  font-size 1.25rem
  line-height 1.2
  grid-area name
#public-head-logo
  max-width 120px
  width 100%
  grid-area logo
#public-head-body
  align-self center
  line-height 1.2
  grid-area body
#public-head-login
  align-self center
  text-align right
  grid-area login
#public-folders
  display grid
  margin 0 auto
  margin-top 1.5rem
  max-width 536px
  grid-template 'listos' 'spanish' 'cert' 'sares' 'snap' / auto
  grid-gap 1rem
  @media (min-width 632px)
    max-width 1088px
    grid-template unquote("'listos cert' 'spanish sares' 'snap none' / calc(50% - 8px) calc(50% - 8px)")
#public-folder-listos
  grid-area listos
#public-folder-spanish
  grid-area spanish
#public-folder-cert
  grid-area cert
  .public-folder-logo
    margin-top 0.5rem
#public-folder-sares
  grid-area sares
  .public-folder-logo
    margin-top 1rem
#public-folder-snap
  grid-area snap
.public-folder
  display grid
  grid-template 'name name' 'logo body' 'links links' / 25% auto
  grid-column-gap 1rem
  @media (min-width 450px)
    grid-template 'logo name' 'logo body' 'logo links' / 25% auto
  @media (min-width 504px)
    grid-template 'logo name' 'logo body' 'logo links' / 120px auto
  @media (min-width 632px)
    grid-template 'name name' 'logo body' 'links links' / 25% auto
  @media (min-width 892px)
    grid-template 'logo name' 'logo body' 'logo links' / 25% auto
  @media (min-width 1000px)
    grid-template 'logo name' 'logo body' 'logo links' / 120px auto
.public-folder-name
  color #888
  font-weight bold
  font-size 1.25rem
  grid-area name
.public-folder-logo
  max-width 120px
  max-height 120px
  width 100%
  grid-area logo
.public-folder-body
  line-height 1.2
  grid-area body
.public-folder-links
  grid-area links
</style>
