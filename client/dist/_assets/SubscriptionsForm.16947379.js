import{d as e,m as s,l as a,e as n,f as i,I as t,_ as o,u as r,r as l,c as u,o as d,w as c,F as b,a as m,t as p,z as g}from"./index.3e114086.js";var S=e({components:{SFCheckGroup:s,SForm:a,SSpinner:n},props:{inModal:{type:Boolean,default:!1},pid:{type:[Number,String],required:!0},email:{type:String}},emits:["done"],setup(e,{emit:s}){const a=i([]),n=i(new Set);t((async()=>{a.value=[],n.value.clear(),a.value=(await o.get(`/api/people/${e.pid}/lists`)).data,n.value.clear(),a.value.filter((s=>s.subscribed&&(!e.email||s.name!==e.email+"@SunnyvaleSERV.org"))).forEach((e=>{n.value.add(e.id)}))}));const l=r((()=>a.value.filter((e=>!n.value.has(e.id)&&e.subWarn.length)).map((e=>{switch(e.subWarn.length){case 1:return`Messages sent to ${e.name} are considered required for the ${e.subWarn[0]} role.  Unsubscribing from it may cause you to lose that role.`;case 2:return`Messages sent to ${e.name} are considered required for the ${e.subWarn[0]} and ${e.subWarn[1]} roles.  Unsubscribing from it may cause you to lose those roles.`;default:return`Messages sent to ${e.name} are considered required for the ${e.subWarn.slice(0,-1).join(", ")}, and ${e.subWarn[e.subWarn.length-1]} roles.  Unsubscribing from it may cause you to lose those roles.`}})).join("\n\n"))),u=i(!1);return{lists:a,onCancel:function(){s("done",!1)},onSubmit:async function(){var a=new FormData;n.value.forEach((e=>{a.append("list",e.toString())})),u.value=!0,await o.post(`/api/people/${e.pid}/lists`,a),u.value=!1,s("done",!0)},submitting:u,subscribed:n,unsubscribeWarnings:l}}});S.render=function(e,s,a,n,i,t){const o=l("SSpinner"),r=l("SFCheckGroup"),S=l("SForm");return d(),u(S,{dialog:e.inModal,title:e.inModal?"Edit List Subscriptions":null,submitLabel:"Save",cancelLabel:e.inModal?"Cancel":"",disabled:e.submitting,onSubmit:e.onSubmit,onCancel:e.onCancel},{default:c((()=>[e.lists.length?(d(),u(b,{key:1},[m(r,{id:"person-lists",label:"Subscriptions",options:e.lists,valueKey:"id",labelKey:"name",modelValue:e.subscribed,"onUpdate:modelValue":s[1]||(s[1]=s=>e.subscribed=s)},null,8,["options","modelValue"]),e.unsubscribeWarnings?(d(),u("div",{key:0,class:"form-item",id:"subscriptions-warning",textContent:p(e.unsubscribeWarnings)},null,8,["textContent"])):g("",!0)],64)):(d(),u(o,{key:0}))])),_:1},8,["dialog","title","cancelLabel","disabled","onSubmit","onCancel"])};export{S as s};