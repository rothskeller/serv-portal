let t=document.createElement("style");t.innerHTML="#texts-list{padding:.75rem}#texts-list-table{display:grid;line-height:1.2;grid-auto-columns:10rem 1fr}.texts-list-heading{font-weight:700}.texts-list-heading:nth-child(3){display:none}.texts-list-timestamp{margin-top:.75rem;font-variant:tabular-nums}.texts-list-groups{margin-top:.75rem;white-space:nowrap}.texts-list-group{overflow:hidden;text-overflow:ellipsis}.texts-list-message{padding-left:4em;font-size:.75rem;grid-column:1/3}@media (min-width:700px){#texts-list-table{grid-auto-columns:10rem 10rem 1fr}.texts-list-heading:nth-child(3){display:block}.texts-list-message{padding-top:.75rem;padding-left:0;font-size:1rem;grid-column:3/4}}",document.head.appendChild(t);import{d as e,e as s,f as i,_ as n,r as l,c as a,a as d,F as o,g as r,o as m,t as g}from"./index.7c6b9706.js";import{s as p}from"./page.dded804e.js";var x=e({components:{SSpinner:s},setup(){p({title:"Text Messages"});const t=i(!0),e=i([]);return n.get("/api/sms").then((s=>{e.value=s.data.messages,t.value=!1})),{loading:t,messages:e}}});const u={id:"texts-list"},c={key:0,id:"texts-list-spinner"},h={key:1,id:"texts-list-table"},v=d("div",{class:"texts-list-timestamp texts-list-heading"},"Time Sent",-1),f=d("div",{class:"texts-list-groups texts-list-heading"},"Recipients",-1),C=d("div",{class:"texts-list-message texts-list-heading"},"Message",-1),b={class:"texts-list-timestamp"},w={class:"texts-list-groups"};x.render=function(t,e,s,i,n,p){const x=l("SSpinner"),y=l("router-link");return m(),a("div",u,[t.loading?(m(),a("div",c,[d(x)])):(m(),a("div",h,[v,f,C,(m(!0),a(o,null,r(t.messages,(t=>(m(),a(o,null,[d("div",b,[d(y,{to:"/texts/"+t.id,textContent:g(t.timestamp)},null,8,["to","textContent"])]),d("div",w,[(m(!0),a(o,null,r(t.groups,(t=>(m(),a("div",{class:"texts-list-group",textContent:g(t)},null,8,["textContent"])))),256))]),d("div",{class:"texts-list-message",textContent:g(t.message)},null,8,["textContent"])],64)))),256))]))])};export default x;