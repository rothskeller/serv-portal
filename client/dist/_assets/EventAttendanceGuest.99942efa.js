import{d as e,q as t,S as a,P as n,f as u,z as s,_ as l,r as o,c as i,o as d,w as r,a as c,G as v,H as m,F as p,g as f,t as g,A as S,b}from"./index.be4dc8e3.js";var h=e({components:{Modal:t,SButton:a,SInput:n},setup(e,{emit:t}){const a=u(null);const n=u(""),o=u(0),i=u([]);return s(n,(async()=>{if(!n.value.trim())return i.value=[],void(o.value=0);i.value=(await l.get("/api/people",{params:{search:n.value.trim()}})).data,o.value=1===i.value.length?i.value[0].id:0})),{guest:o,modal:a,onSubmit:function(){o.value&&a.value.close(i.value.find((e=>e.id===o.value)))},options:i,search:n,show:function(){return a.value.show()}}}});const C=c("div",{id:"event-attendance-guest-title"},"Add Guest",-1),w={id:"event-attendance-guest-buttons"},x=b("Cancel"),V=b("OK");h.render=function(e,t,a,n,u,s){const l=o("SInput"),b=o("SButton"),h=o("Modal");return d(),i(h,{ref:"modal"},{default:r((({close:a})=>[C,c("form",{id:"event-attendance-guest",onSubmit:t[3]||(t[3]=S(((...t)=>e.onSubmit&&e.onSubmit(...t)),["prevent"]))},[c(l,{id:"event-attendance-guest-input",placeholder:"Guest Name",autofocus:"",modelValue:e.search,"onUpdate:modelValue":t[1]||(t[1]=t=>e.search=t)},null,8,["modelValue"]),v(c("select",{size:"10","onUpdate:modelValue":t[2]||(t[2]=t=>e.guest=t)},[(d(!0),i(p,null,f(e.options,(e=>(d(),i("option",{value:e.id,textContent:g(e.sortName)},null,8,["value","textContent"])))),256))],512),[[m,e.guest]]),c("div",w,[c(b,{onClick:e=>a(null)},{default:r((()=>[x])),_:2},1032,["onClick"]),c(b,{type:"submit",variant:"primary"},{default:r((()=>[V])),_:1})])],32)])),_:1},512)};export{h as s};