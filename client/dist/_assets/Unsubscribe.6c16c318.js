let e=document.createElement("style");e.innerHTML="#unsubscribe{margin:1.5rem .75rem}#unsub-head{font-weight:700;font-size:1.25rem}#unsub-intro{margin:.5rem 0}#unsub-warn{margin:.5rem 0 .75rem;max-width:40rem;line-height:1.2}",document.head.appendChild(e);import{d as n,s as a,Z as i,e as u,f as s,_ as o,v as t,r as l,c as r,a as d,F as m,g as b,w as c,b as p,o as v}from"./index.d1c103e1.js";import{s as g}from"./page.b52cf230.js";var h=n({components:{SButton:a,SCheck:i,SSpinner:u},setup(){const e=t();g({title:"SunnyvaleSERV Unsubscribe"});const n=s(!0),a=s(!1),i=s([]);o.get("/api/unsubscribe/"+e.params.token).then((u=>{a.value=u.data.noEmail,i.value=u.data.groups,n.value=!1,e.params.email&&i.value.forEach((n=>{n.email===e.params.email&&(n.unsub=!0)}))}));const u=s(!1);return{groups:i,loading:n,noEmail:a,onSubmit:async function(){const n=new FormData;n.append("noEmail",a.value.toString()),i.value.forEach((e=>{n.append("unsub:"+e.id,e.unsub.toString())})),await o.post("/api/unsubscribe/"+e.params.token,n),u.value=!0},submitted:u}}});const f={key:0,id:"unsubscribe"},y={key:1,id:"unsubscribe"},S=d("div",{id:"unsub-head"},"Unsubscribe",-1),k={key:0,id:"unsub-intro"},w={key:1,id:"unsub-intro"},E=d("div",{id:"unsub-warn"},[p("We would appreciate it if you’d drop us a note at "),d("a",{href:"mailto:admin@sunnyvaleserv.org"},"admin@sunnyvaleserv.org"),p(" and let us know why you unsubscribed.")],-1),V=d("div",{id:"unsub-warn"},"If you ever want to get back on the email lists, come back to this page and let us know.",-1),U={key:2,id:"unsubscribe"},C=d("div",{id:"unsub-head"},"Unsubscribe",-1),x=d("div",{id:"unsub-intro"},"Which email list(s) do you want to unsubscribe from?",-1),R={style:{"margin-top":"0.5rem"}},j=d("div",{id:"unsub-warn"},"Please note that, if you unsubscribe from a critical mailing list for one of our volunteer groups, you may no longer be able to participate in that group.",-1),B={id:"unsub-buttons"},F=p("Unsubscribe");h.render=function(e,n,a,i,u,s){const o=l("SSpinner"),t=l("SCheck"),p=l("SButton");return e.loading?(v(),r("div",f,[d(o)])):e.submitted?(v(),r("div",y,[S,e.noEmail?(v(),r("div",k,"You have been removed from all of our email lists.")):(v(),r("div",w,"You have been removed from the selected email lists.")),E,V])):(v(),r("div",U,[C,x,(v(!0),r(m,null,b(e.groups,(e=>(v(),r("div",null,[d(t,{id:"unsub-"+e.id,label:e.email+"@SunnyvaleSERV.org",modelValue:e.unsub,"onUpdate:modelValue":n=>e.unsub=n},null,8,["id","label","modelValue","onUpdate:modelValue"])])))),256)),d("div",R,[d(t,{id:"unsub-all",label:"All SunnyvaleSERV email lists",modelValue:e.noEmail,"onUpdate:modelValue":n[1]||(n[1]=n=>e.noEmail=n)},null,8,["modelValue"])]),j,d("div",B,[d(p,{variant:"primary",onClick:e.onSubmit},{default:c((()=>[F])),_:1},8,["onClick"])])]))};export default h;