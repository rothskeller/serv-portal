(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["TextsSend"],{"1ef7":function(e,t,s){"use strict";s.r(t);var r=function(){var e=this,t=e.$createElement,s=e._self._c||t;return e.groups?s("b-form",{attrs:{id:"texts-send"},on:{submit:function(t){return t.preventDefault(),e.onSubmit(t)}}},[s("b-form-group",{attrs:{label:"Message","label-for":"texts-send-message","label-cols-sm":"auto","label-class":"texts-send-label",state:!e.messageError&&null,invalidFeedback:e.messageError}},[s("b-textarea",{attrs:{id:"texts-send-message",rows:"5",autofocus:""},model:{value:e.message,callback:function(t){e.message=t},expression:"message"}}),e.countMessage?s("b-form-text",{class:e.countClass},[e._v(e._s(e.countMessage))]):e._e()],1),s("b-form-group",{attrs:{label:"Recipients","label-for":"texts-send-groups","label-cols-sm":"auto","label-class":"texts-send-label pt-0",state:!e.groupsError&&null,invalidFeedback:e.groupsError}},[s("b-form-checkbox-group",{attrs:{id:"texts-send-groups",options:e.groups,stacked:"","text-field":"name","value-field":"id"},model:{value:e.recipients,callback:function(t){e.recipients=t},expression:"recipients"}})],1),s("div",{staticClass:"mt-3"},[s("b-button",{attrs:{type:"submit",variant:"primary",disabled:e.sending||!e.valid}},[e.sending?s("b-spinner",{attrs:{small:""}}):s("span",[e._v("Send Message")])],1)],1)],1):s("div",{staticClass:"mt-3 ml-2"},[s("b-spinner",{attrs:{small:""}})],1)},n=[],a=s("a34a"),i=s.n(a);function o(e,t,s,r,n,a,i){try{var o=e[a](i),u=o.value}catch(l){return void s(l)}o.done?t(u):Promise.resolve(u).then(r,n)}function u(e){return function(){var t=this,s=arguments;return new Promise((function(r,n){var a=e.apply(t,s);function i(e){o(a,r,n,i,u,"next",e)}function u(e){o(a,r,n,i,u,"throw",e)}i(void 0)}))}}var l={data:function(){return{groups:null,message:"",recipients:[],countClass:"",countMessage:"0/160",messageError:null,groupsError:null,submitted:!1,sending:!1}},created:function(){var e=u(i.a.mark((function e(){var t;return i.a.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.next=2,this.$axios.get("/api/sms/NEW");case 2:t=e.sent.data,this.groups=t.groups;case 4:case"end":return e.stop()}}),e,this)})));function t(){return e.apply(this,arguments)}return t}(),watch:{message:"validate",recipients:"validate"},computed:{valid:function(){return!this.messageError&&!this.groupsError}},methods:{onSubmit:function(){var e=u(i.a.mark((function e(){var t,s;return i.a.wrap((function(e){while(1)switch(e.prev=e.next){case 0:if(this.submitted=!0,this.validate(),this.valid){e.next=4;break}return e.abrupt("return");case 4:return t=new FormData,t.append("message",this.message),this.recipients.forEach((function(e){return t.append("group",e)})),this.sending=!0,e.next=10,this.$axios.post("/api/sms",t);case 10:s=e.sent.data,this.sending=!1,this.$router.push("/texts/".concat(s.id));case 13:case"end":return e.stop()}}),e,this)})));function t(){return e.apply(this,arguments)}return t}(),validate:function(){var e=!1,t=0,s=0,r=!0,n=!1,a=void 0;try{for(var i,o=this.message[Symbol.iterator]();!(r=(i=o.next()).done);r=!0){var u=i.value;t++,"£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà".includes(u)?s++:"\f\n^{}\\[~]|€".includes(u)?s+=2:e=!0}}catch(l){n=!0,a=l}finally{try{r||null==o.return||o.return()}finally{if(n)throw a}}this.submitted&&!t?(this.messageError="Please enter the text of your message.",this.countMessage="",this.countClass=""):e?(this.messageError="",this.countMessage="".concat(t,"/70"),this.countClass=t>70?"texts-send-long":""):(this.messageError="",this.countMessage="".concat(s,"/160"),this.countClass=s>160?"texts-send-long":""),this.submitted&&!this.recipients.length?this.groupsError="Please select the recipients of your message.":this.groupsError=null}}},c=l,d=(s("525e"),s("2877")),p=Object(d["a"])(c,r,n,!1,null,null,null);t["default"]=p.exports},"525e":function(e,t,s){"use strict";var r=s("eab0"),n=s.n(r);n.a},a34a:function(e,t,s){e.exports=s("96cf")},eab0:function(e,t,s){}}]);
//# sourceMappingURL=TextsSend.39b33d4d.js.map