let e=document.createElement("style");e.innerHTML="#list{margin:1.5rem .75rem}",document.head.appendChild(e);import{d as a,s as t,S as l,j as s,k as i,l as n,e as o,f as d,_ as u,m,n as r,u as p,r as c,c as b,p as v,w as S,o as y,a as h,b as f}from"./index.be4dc8e3.js";import{s as g}from"./page.46d1c173.js";const L=[{value:"email",label:"Email"},{value:"sms",label:"SMS"}];var w=a({components:{MessageBox:t,SButton:l,SForm:s,SFInput:i,SFRadioGroup:n,SSpinner:o},setup(){const e=r(),a=p();g({title:"NEW"===e.params.lid?"Add List":"Edit List"});const t="NEW"===e.params.lid?"Add List":"Save List",l=d(!0),s=d({});u.get(`/api/lists/${e.params.lid}`).then((e=>{s.value=e.data,l.value=!1}));const i=d(""),n=m((()=>"email"===s.value.type?"Email address":"Name")),o=m((()=>"email"===s.value.type?"@SunnyvaleSERV.org":null));const c=d(!1);const b=d(null);return{deleteModal:b,loading:l,list:s,nameError:function(e){return e?s.value.name?"email"!==s.value.type||s.value.name.match(/^[a-z][-a-z0-9]*$/)?i.value===s.value.name?"Another list has this name.":"":"The email address must start with a lowercase letter and consist of lowercase letters and digits.":"email"===s.value.type?"The email address is required.":"The list name is required.":""},nameHelp:o,nameLabel:n,onDelete:async function(){if(b.value){await b.value.show()&&(c.value=!0,await u.delete(`/api/lists/${e.params.lid}`),c.value=!1,a.push("/admin/lists"))}},onSubmit:async function(){const t=new FormData;t.append("type",s.value.type),t.append("name",s.value.name),c.value=!0;const l=(await u.post(`/api/lists/${e.params.lid}`,t)).data;c.value=!1,l&&l.duplicateName?i.value=s.value.name:a.push("/admin/lists")},submitLabel:t,submitting:c,typeOptions:L}}});const E={id:"list"},F=f("Delete List"),V=f("Are you sure you want to delete this list? All associated data, including role associations, manual subscribes, and unsubscribes will be permanently lost.");w.render=function(e,a,t,l,s,i){const n=c("SSpinner"),o=c("SFRadioGroup"),d=c("SFInput"),u=c("SButton"),m=c("MessageBox"),r=c("SForm");return y(),b("div",E,[e.loading?(y(),b(n,{key:0})):(y(),b(r,{key:1,submitLabel:e.submitLabel,disabled:e.submitting,onSubmit:e.onSubmit},v({default:S((()=>[h(o,{id:"list-type",label:"Type",inline:"",options:e.typeOptions,modelValue:e.list.type,"onUpdate:modelValue":a[1]||(a[1]=a=>e.list.type=a)},null,8,["options","modelValue"]),h(d,{id:"list-name",label:e.nameLabel,help:e.nameHelp,trim:"",modelValue:e.list.name,"onUpdate:modelValue":a[2]||(a[2]=a=>e.list.name=a),errorFn:e.nameError},null,8,["label","help","modelValue","errorFn"]),h(m,{ref:"deleteModal",title:"Delete List",cancelLabel:"Keep",okLabel:"Delete",variant:"danger"},{default:S((()=>[V])),_:1},512)])),_:2},[e.list.id?{name:"extraButtons",fn:S((()=>[h(u,{onClick:e.onDelete,variant:"danger",disabled:e.submitting},{default:S((()=>[F])),_:1},8,["onClick","disabled"])]))}:void 0]),1032,["submitLabel","disabled","onSubmit"]))])};export default w;