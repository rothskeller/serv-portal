(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["Login"],{"17b0":function(t,e,a){"use strict";var r=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{class:t.$store.state.touch?"touch":"mouse",attrs:{id:"page-top"}},[a("div",{attrs:{id:"page-heading"}},[a("div",{attrs:{id:"page-menu-spacer"}}),a("div",{attrs:{id:"page-titlebox"}},[a("div",{attrs:{id:"page-title"},domProps:{textContent:t._s(t.title)}})]),a("div",{attrs:{id:"page-menu-spacer"}})]),a("div",{staticClass:"page-no-menu",attrs:{id:"page-content"}},[t._t("default")],2)])},i=[],s={props:{title:String}},n=s,o=a("2877"),l=Object(o["a"])(n,r,i,!1,null,null,null);e["a"]=l.exports},1964:function(t,e,a){},6651:function(t,e,a){"use strict";var r=a("1964"),i=a.n(r);i.a},a34a:function(t,e,a){t.exports=a("96cf")},a7fb:function(t,e,a){"use strict";a.r(e);var r=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("PublicPage",{attrs:{title:"Sunnyvale SERV"}},[a("div",{attrs:{id:"login-top"}},[a("div",{attrs:{id:"login-banner"}},[t._v("Please log in.")]),a("div",{attrs:{id:"login-forserv"}},[t._v("This web site is for SERV volunteers only.\nIf you are interested in joining one of the SERV volunteer organizations,\nplease visit Sunnyvale’s "),a("a",{attrs:{href:"https://sunnyvale.ca.gov/government/safety/emergency.htm"}},[t._v("emergency response page")]),t._v(".")]),a("div",{attrs:{id:"login-browserwarn"}},[t._v("Your browser is out of date and lacks features needed by this web site.\nThe site may not look or behave correctly.")]),a("form",{attrs:{id:"login-form"},on:{submit:function(e){return e.preventDefault(),t.onSubmit(e)}}},[a("b-form-group",{attrs:{label:"Email address","label-for":"login-email","label-cols-sm":"4"}},[a("b-input",{attrs:{id:"login-email",autocorrect:"off",autocapitalize:"none",autofocus:"",required:"",trim:""},model:{value:t.email,callback:function(e){t.email=e},expression:"email"}})],1),a("b-form-group",{attrs:{label:"Password","label-for":"login-password","label-cols-sm":"4"}},[a("b-input",{ref:"password",attrs:{id:"login-password",type:"password"},model:{value:t.password,callback:function(e){t.password=e},expression:"password"}})],1),a("div",{attrs:{id:"login-submit-row"}},[a("b-button",{attrs:{type:"submit",variant:"primary"}},[t._v("Log in")])],1),t.failed?a("div",{attrs:{id:"login-failed"}},[t._v("Login incorrect. Please try again.")]):t._e()],1),a("div",{attrs:{id:"login-reset"}},[a("b-btn",{attrs:{to:"/password-reset"}},[t._v("Reset my password")])],1),a("b-link",{attrs:{id:"login-policies",to:"/policies"}},[t._v("Site Policies / Legal Stuff")])],1)])},i=[],s=a("a34a"),n=a.n(s),o=a("17b0");function l(t,e,a,r,i,s,n){try{var o=t[s](n),l=o.value}catch(u){return void a(u)}o.done?e(l):Promise.resolve(l).then(r,i)}function u(t){return function(){var e=this,a=arguments;return new Promise((function(r,i){var s=t.apply(e,a);function n(t){l(s,r,i,n,o,"next",t)}function o(t){l(s,r,i,n,o,"throw",t)}n(void 0)}))}}var c={components:{PublicPage:o["a"]},data:function(){return{email:"",password:"",failed:!1}},methods:{onSubmit:function(){var t=u(n.a.mark((function t(){var e,a;return n.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:if(this.email&&this.password){t.next=2;break}return t.abrupt("return");case 2:return e=new FormData,e.append("username",this.email),e.append("password",this.password),t.prev=5,t.next=8,this.$axios.post("/api/login",e);case 8:a=t.sent.data,this.$store.commit("login",a),this.$router.replace(this.$route.query.redirect||"/"),t.next=19;break;case 13:t.prev=13,t.t0=t["catch"](5),console.error(t.t0),this.failed=!0,this.password="",this.$refs.password.focus();case 19:case"end":return t.stop()}}),t,this,[[5,13]])})));function e(){return t.apply(this,arguments)}return e}()}},d=c,p=(a("6651"),a("2877")),f=Object(p["a"])(d,r,i,!1,null,null,null);e["default"]=f.exports}}]);
//# sourceMappingURL=Login.0ea10ee1.js.map