(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["Reports"],{1708:function(t,e,a){"use strict";a.r(e);var r=function(){var t=this,e=t.$createElement,a=t._self._c||e;return t.loading?a("div",{staticClass:"mt-3"},[a("b-spinner",{attrs:{small:""}})],1):t.certAtt?a("CERTAttendanceForm",t._b({},"CERTAttendanceForm",t.certAtt,!1)):t._e()},o=[],s=a("a34a"),n=a.n(s),i=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("form",{attrs:{id:"cert-att-form"},on:{submit:function(e){return e.preventDefault(),t.onSubmit(e)}}},[a("div",{staticClass:"report-title"},[t._v("CERT Attendance Report")]),a("b-form-group",{attrs:{label:"Report on team","label-cols-sm":"auto","label-class":"cert-att-label"}},[a("b-form-radio-group",{staticClass:"cert-att-radio-group",attrs:{options:t.teamList},model:{value:t.team,callback:function(e){t.team=e},expression:"team"}})],1),a("b-form-group",{attrs:{label:"Date range","label-cols-sm":"auto","label-class":"cert-att-label",state:!t.dateError&&null,"invalid-feedback":t.dateError}},[a("b-form-input",{attrs:{id:"cert-att-date-from",type:"date",state:!t.dateError&&null},model:{value:t.dateFromI,callback:function(e){t.dateFromI=e},expression:"dateFromI"}}),t._v("\nthrough\n"),a("b-form-input",{attrs:{id:"cert-att-date-to",type:"date",state:!t.dateError&&null},model:{value:t.dateToI,callback:function(e){t.dateToI=e},expression:"dateToI"}})],1),a("b-form-group",{attrs:{label:"Statistics by","label-cols-sm":"auto","label-class":"cert-att-label"}},[a("b-form-radio-group",{staticClass:"cert-att-radio-group",attrs:{options:t.statsList},model:{value:t.stats,callback:function(e){t.stats=e},expression:"stats"}})],1),a("b-form-group",{attrs:{label:"Show detail","label-cols-sm":"auto","label-class":"cert-att-label"}},[a("b-form-radio-group",{staticClass:"cert-att-radio-group",attrs:{options:t.detailList},model:{value:t.detail,callback:function(e){t.detail=e},expression:"detail"}})],1),a("div",{staticClass:"mt-3"},[a("b-btn",{attrs:{type:"submit",variant:"primary",disabled:!!t.dateError}},[t._v("Generate Report")])],1)],1)},l=[],c=["Alpha","Bravo","Both"],d=[{value:"count",text:"Number of Events"},{value:"hours",text:"Cumulative Hours"}],u=[{value:"date",text:"Show by date"},{value:"month",text:"Show by month"},{value:"total",text:"Show totals only"}],m=/^20\d\d-(?:0[1-9]|1[012])-(?:0[1-9]|[12][0-9]|3[01])$/,h={props:{dateFrom:String,dateTo:String},data:function(){return{teamList:c,statsList:d,detailList:u,team:"Both",dateFromI:null,dateToI:null,stats:"count",detail:"month",dateError:null}},created:function(){this.dateFromI=this.dateFrom,this.dateToI=this.dateTo},watch:{dateFromI:"checkDates",dateToI:"checkDates"},methods:{checkDates:function(){this.dateFromI&&this.dateToI?this.dateFromI.match(m)&&this.dateToI.match(m)?this.dateFromI>this.dateToI?this.dateError="The starting date must be before the ending date.":this.dateError=null:this.dateError="Valid dates have the form YYYY-MM-DD.":this.dateError="Starting and ending dates are required."},onSubmit:function(){this.dateError||this.$router.push({path:"/reports/cert-attendance",query:{team:this.team,dateFrom:this.dateFromI,dateTo:this.dateToI,stats:this.stats,detail:this.detail}})}}},p=h,b=(a("95cc"),a("2877")),f=Object(b["a"])(p,i,l,!1,null,null,null),v=f.exports;function g(t,e,a,r,o,s,n){try{var i=t[s](n),l=i.value}catch(c){return void a(c)}i.done?e(l):Promise.resolve(l).then(r,o)}function I(t){return function(){var e=this,a=arguments;return new Promise((function(r,o){var s=t.apply(e,a);function n(t){g(s,r,o,n,i,"next",t)}function i(t){g(s,r,o,n,i,"throw",t)}n(void 0)}))}}var E={components:{CERTAttendanceForm:v},data:function(){return{loading:!1,certAtt:null}},created:function(){var t=I(n.a.mark((function t(){var e;return n.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return this.$store.commit("setPage",{title:"Reports"}),this.loading=!0,t.next=4,this.$axios.get("/api/reports");case 4:e=t.sent.data,this.certAtt=e.certAttendance,this.loading=!1;case 7:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}()},T=E,x=(a("1ab2"),Object(b["a"])(T,r,o,!1,null,null,null));e["default"]=x.exports},"1ab2":function(t,e,a){"use strict";var r=a("1ab9"),o=a.n(r);o.a},"1ab9":function(t,e,a){},"33df":function(t,e,a){},"95cc":function(t,e,a){"use strict";var r=a("33df"),o=a.n(r);o.a},a34a:function(t,e,a){t.exports=a("96cf")}}]);
//# sourceMappingURL=Reports.dba7f3a9.js.map