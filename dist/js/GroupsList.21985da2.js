(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["GroupsList"],{"05df":function(t,n,s){},"21b2":function(t,n,s){"use strict";s.r(n);var r=function(){var t=this,n=t.$createElement,s=t._self._c||n;return s("div",{attrs:{id:"groups-list"}},[t.loading?s("div",{attrs:{id:"groups-list-spinner"}},[s("b-spinner",{attrs:{small:""}})],1):s("div",{attrs:{id:"groups-list-table"}},[s("div",{staticClass:"groups-list-name groups-list-heading"},[t._v("Group")]),s("div",{staticClass:"groups-list-roles groups-list-heading"},[t._v("Included in Roles")]),t._l(t.groups,(function(n){return[s("div",{staticClass:"groups-list-name"},[s("router-link",{attrs:{to:"/groups/"+n.id},domProps:{textContent:t._s(n.name)}})],1),s("div",{staticClass:"groups-list-roles"},t._l(n.roles,(function(n){return s("div",{domProps:{textContent:t._s(n)}})})),0)]}))],2)])},i=[],e=s("a34a"),o=s.n(e);function a(t,n,s,r,i,e,o){try{var a=t[e](o),u=a.value}catch(l){return void s(l)}a.done?n(u):Promise.resolve(u).then(r,i)}function u(t){return function(){var n=this,s=arguments;return new Promise((function(r,i){var e=t.apply(n,s);function o(t){a(e,r,i,o,u,"next",t)}function u(t){a(e,r,i,o,u,"throw",t)}o(void 0)}))}}var l={data:function(){return{groups:null,loading:!0}},created:function(){var t=u(o.a.mark((function t(){return o.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return this.loading=!0,t.next=3,this.$axios.get("/api/groups");case 3:this.groups=t.sent.data,this.loading=!1;case 5:case"end":return t.stop()}}),t,this)})));function n(){return t.apply(this,arguments)}return n}()},c=l,p=(s("498c"),s("2877")),d=Object(p["a"])(c,r,i,!1,null,null,null);n["default"]=d.exports},"498c":function(t,n,s){"use strict";var r=s("05df"),i=s.n(r);i.a},a34a:function(t,n,s){t.exports=s("96cf")}}]);
//# sourceMappingURL=GroupsList.21985da2.js.map