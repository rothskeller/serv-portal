let e=document.createElement("style");e.innerHTML="#event-view{padding:1.5rem .75rem}#event-view-name{font-weight:700;font-size:1.25rem;line-height:1.2}#event-view-orgtype{color:#888}#event-view-date-time{margin-top:.75rem;white-space:pre-line;line-height:1.2}#event-view-venue{margin-top:.75rem;line-height:1.2}#event-view-venue-address{font-size:.875rem}#event-view-details{margin-top:.75rem;max-width:40rem;white-space:pre-line;line-height:1.2}",document.head.appendChild(e);import{d as t,e as n,f as a,_ as i,u as v,v as r,r as d,c as o,a as m,t as s,z as l,b as u,o as p}from"./index.d1c103e1.js";import{s as w}from"./page.b52cf230.js";import{m as c}from"./moment.min.6b5db032.js";var f=t({components:{SSpinner:n},props:{onLoadEvent:{type:Function,required:!0}},setup(e){const t=r();w({title:"Events"});const n=a(null);i.get("/api/events/"+t.params.id).then((t=>{n.value=t.data.event,w({title:`${n.value.date} ${n.value.name}`,browserTitle:n.value.date}),e.onLoadEvent(n.value)}));const d=v((()=>{if(!n.value)return"";const e=c(n.value.date,"YYYY-MM-DD"),t=c(n.value.start,"HH:mm"),a=c(n.value.end,"HH:mm");return t.format("a")!==a.format("a")?`${e.format("dddd, MMMM D, YYYY")}\n${t.format("h:mma")} to ${a.format("h:mma")}`:`${e.format("dddd, MMMM D, YYYY")}\n${t.format("h:mm")} to ${a.format("h:mma")}`}));return{event:n,dateTimeFmt:d}}});const h={key:0,id:"event-view"},g={key:1,id:"event-view"},x={id:"event-view-venue"},y={key:0,id:"event-view-venue-address"},C={key:1,id:"event-view-venue-map"},M=u(" ("),Y=u(")");f.render=function(e,t,n,a,i,v){const r=d("SSpinner");return e.event?(p(),o("div",g,[m("div",{id:"event-view-name",textContent:s(e.event.name)},null,8,["textContent"]),m("div",{id:"event-view-orgtype",textContent:s(`${e.event.organization} ${e.event.type}`)},null,8,["textContent"]),m("div",{id:"event-view-date-time",textContent:s(e.dateTimeFmt)},null,8,["textContent"]),m("div",x,[m("div",{id:"event-view-venue-name",textContent:s(e.event.venue?e.event.venue.name:"Location TBD")},null,8,["textContent"]),e.event.venue?(p(),o("div",y,[m("span",{textContent:s(e.event.venue.address)},null,8,["textContent"]),e.event.venue.city?(p(),o("span",{key:0,textContent:s(", "+e.event.venue.city)},null,8,["textContent"])):l("",!0),e.event.venue.url?(p(),o("span",C,[M,m("a",{target:"_blank",href:e.event.venue.url},"map",8,["href"]),Y])):l("",!0)])):l("",!0)]),e.event.details?(p(),o("div",{key:0,id:"event-view-details",innerHTML:e.event.details},null,8,["innerHTML"])):l("",!0)])):(p(),o("div",h,[m(r)]))};export default f;