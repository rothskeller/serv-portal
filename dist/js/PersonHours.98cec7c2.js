(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["PersonHours"],{"0296":function(t,e,n){"use strict";n.r(e);var r=function(){var t=this,e=t.$createElement,n=t._self._c||e;return t.unregistered?n("div",{staticClass:"mt-3 ml-2",staticStyle:{"max-width":"600px"}},[t._v("You are not currently registered as a City of Sunnyvale volunteer.  We\nappreciate your volunteer efforts, but we cannot record your hours until you\nare registered.  To do so, please fill out\n"),n("a",{attrs:{href:"https://www.volgistics.com/ex/portal.dll/ap?AP=929478828",target:"_blank"}},[t._v("this form")]),t._v(".  In the “City employee status or referral” box, please enter"),n("pre",{staticClass:"ml-4 mt-3"},[t._v("Rebecca Elizondo\nDepartment of Public Safety")]),t._v("and the names of the organizations you're volunteering for (CERT, LISTOS,\nSNAP, and/or SARES).  Come back a week or so later and we should have your\nregistration on file.  If you have any difficulties with this, please\ncontact Rebecca at RElizondo@sunnyvale.ca.gov.")]):t.months?n("form",{attrs:{id:"person-hours"},on:{submit:function(e){return e.preventDefault(),t.onSubmit(e)}}},[t._l(t.months,(function(e){return n("div",{staticClass:"person-hours"},[n("div",{staticClass:"person-hours-heading",domProps:{textContent:t._s("Volunteer Hours for "+e.month)}}),n("table",{staticClass:"person-hours-table"},[t._l(e.events,(function(e){return n("tr",[n("td",{staticClass:"person-hours-event",domProps:{textContent:t._s(t.eventText(e))}}),n("td",[n("input",{staticClass:"person-hours-time",attrs:{type:"number",min:"0",step:"0.5"},domProps:{value:t.eventTime(e)},on:{change:function(n){return t.setEventTime(e,n)}}})])])})),n("tr",[n("td",{staticClass:"person-hours-total-label"},[t._v("TOTAL")]),n("td",[n("div",{staticClass:"person-hours-total-time",domProps:{textContent:t._s(t.totalHours(e))}})])])],2)])})),n("div",{staticClass:"mt-3"},[n("b-btn",{attrs:{type:"submit",variant:"primary"}},[t._v("Save Hours")]),n("b-btn",{staticClass:"ml-2",on:{click:t.onCancel}},[t._v("Cancel")])],1)],2):n("div",{staticClass:"mt-3 ml-2"},[n("b-spinner",{attrs:{small:""}})],1)},o=[],s=n("a34a"),a=n.n(s);function i(t,e,n,r,o,s,a){try{var i=t[s](a),u=i.value}catch(c){return void n(c)}i.done?e(u):Promise.resolve(u).then(r,o)}function u(t){return function(){var e=this,n=arguments;return new Promise((function(r,o){var s=t.apply(e,n);function a(t){i(s,r,o,a,u,"next",t)}function u(t){i(s,r,o,a,u,"throw",t)}a(void 0)}))}}var c={props:{onLoadPerson:Function},data:function(){return{months:null,unregistered:!1}},created:function(){var t=u(a.a.mark((function t(){var e;return a.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,this.$axios.get("/api/people/".concat(this.$route.params.id,"/hours"));case 2:e=t.sent.data,e?this.months=e:this.unregistered=!0;case 4:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),methods:{eventText:function(t){return t.placeholder?t.name:"".concat(t.date," ").concat(t.name)},eventTime:function(t){return 0===t.minutes?"":Math.floor(t.minutes/30)/2},onCancel:function(){this.$router.go(-1)},onSubmit:function(){var t=u(a.a.mark((function t(){var e;return a.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return e=new FormData,this.months.forEach((function(t){t.events.forEach((function(t){e.append("e".concat(t.id),t.minutes)}))})),t.next=4,this.$axios.post("/api/people/".concat(this.$route.params.id,"/hours"),e);case 4:this.$router.push("/people/".concat(this.$route.params.id));case 5:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),setEventTime:function(t,e){t.minutes=60*e.target.value},totalHours:function(t){var e=t.events.reduce((function(t,e){return t+e.minutes}),0);return Math.floor(e/30)/2}}},l=c,p=(n("5cba"),n("2877")),h=Object(p["a"])(l,r,o,!1,null,null,null);e["default"]=h.exports},"48ce":function(t,e,n){},"5cba":function(t,e,n){"use strict";var r=n("48ce"),o=n.n(r);o.a},a34a:function(t,e,n){t.exports=n("96cf")}}]);
//# sourceMappingURL=PersonHours.98cec7c2.js.map