let e=document.createElement("style");e.innerHTML="#roles2-list{margin:1.5rem .75rem}#roles2-list-table td{padding-left:1rem;vertical-align:middle}#roles2-list-table td:first-child{padding-left:0}.touch #roles2-list-table td{height:40px}.roles2-list-heading{font-weight:700}.roles2-list-people{color:#888}#roles2-list-buttons{margin-top:1.5rem}#roles2-list-saveOrder{margin-right:.5rem}",document.head.appendChild(e);import{d as t,s as r,e as n,f as a,_ as l,x as o,r as s,c as i,F as d,a as u,g,w as v,z as c,o as p,t as f,b as D}from"./index.69270c2b.js";import{s as m}from"./page.0f9cb5de.js";const b={admin:"Admin","cert-d":"CERT-D","cert-t":"CERT-T",listos:"Listos",sares:"SARES",snap:"SNAP"},x={student:"Student",member:"Member",leader:"Leader"};var S=t({components:{SButton:r,SSpinner:n},setup(){const e=o();m({title:"Roles"});const t=a(!0),r=a([]);l.get("/api/roles2").then((e=>{r.value=e.data,t.value=!1}));const n=a(null);let s=0;const i=a(!1),d=a(!1);return{dragOverStyle:function(e){return e===n.value?{borderBottom:"2px solid #888"}:null},loading:t,onAdd:function(){e.push("/admin/roles2/NEW")},onDragEnd:function(){n.value=null},onDragEnter:function(e,t){e.dataTransfer.types.includes("x-serv-role")&&(e.preventDefault(),n.value===t?s++:(n.value=t,s=1))},onDragLeave:function(e,t){e.dataTransfer.types.includes("x-serv-role")&&(n.value===t&&s>1?s--:n.value===t&&(n.value=null))},onDragOver:function(e,t){e.dataTransfer.types.includes("x-serv-role")&&e.preventDefault()},onDragStart:function(e,t){e.dataTransfer.setData("x-serv-role",t.id.toString()),e.dataTransfer.effectAllowed="move"},onDrop:function(e,t){const n=parseInt(e.dataTransfer.getData("x-serv-role")),a=r.value.findIndex((e=>e.id===n));let l="TOP"===t?-1:r.value.findIndex((e=>e===t));if(a===l||a===l+1)return;const o=r.value.splice(a,1);l<a&&l++,r.value.splice(l,0,o[0]),i.value=!0},onSaveOrder:async function(){const e=new FormData;r.value.forEach((t=>{e.append("role",t.id.toString())})),d.value=!0,await l.post("/api/roles2",e),d.value=!1,i.value=!1},orderChanged:i,orgNames:b,privLevelNames:x,roles:r,submitting:d}}});const h={id:"roles2-list"},O={key:0,id:"roles2-list-table"},y=u("td",{class:"roles2-list-heading"},"Org",-1),C=u("td",{class:"roles2-list-heading"},"Priv",-1),T=u("td",{class:"roles2-list-heading"},"Role",-1),E={key:1},k={id:"roles2-list-buttons"},L=D("Save Order"),P=D("Add Role");S.render=function(e,t,r,n,a,l){const o=s("SSpinner"),D=s("router-link"),m=s("SButton");return p(),i("div",h,[e.loading?(p(),i(o,{key:0})):(p(),i(d,{key:1},[e.roles.length?(p(),i("table",O,[u("tr",{style:e.dragOverStyle("TOP"),onDragenter:t[1]||(t[1]=t=>e.onDragEnter(t,"TOP")),onDragover:t[2]||(t[2]=t=>e.onDragOver(t,"TOP")),onDragleave:t[3]||(t[3]=t=>e.onDragLeave(t,"TOP")),onDrop:t[4]||(t[4]=t=>e.onDrop(t,"TOP"))},[y,C,T],36),(p(!0),i(d,null,g(e.roles,(r=>(p(),i("tr",{draggable:"true",style:e.dragOverStyle(r),onDragstart:t=>e.onDragStart(t,r),onDragenter:t=>e.onDragEnter(t,r),onDragover:t=>e.onDragOver(t,r),onDragleave:t=>e.onDragLeave(t,r),onDrop:t=>e.onDrop(t,r),onDragend:t[5]||(t[5]=t=>e.onDragEnd())},[u("td",{textContent:f(e.orgNames[r.org]||"—")},null,8,["textContent"]),u("td",{textContent:f(e.privLevelNames[r.privLevel]||"—")},null,8,["textContent"]),u("td",null,[u(D,{to:"/admin/roles2/"+r.id,draggable:"false",textContent:f(r.name)},null,8,["to","textContent"]),u("span",{class:"roles2-list-people",textContent:f(` [${r.people}]`)},null,8,["textContent"])])],44,["onDragstart","onDragenter","onDragover","onDragleave","onDrop"])))),256))])):(p(),i("div",E,"No roles currently defined.")),u("div",k,[e.orderChanged?(p(),i(m,{key:0,id:"roles2-list-saveOrder",variant:"warning",disabled:e.submitting,onClick:e.onSaveOrder},{default:v((()=>[L])),_:1},8,["disabled","onClick"])):c("",!0),u(m,{variant:"primary",disabled:e.submitting,onClick:e.onAdd},{default:v((()=>[P])),_:1},8,["disabled","onClick"])])],64))])};export default S;