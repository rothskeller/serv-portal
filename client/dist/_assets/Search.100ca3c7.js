let e=document.createElement("style");e.innerHTML="#search{margin:1.5rem .75rem}#search-query-row{text-align:center}#search-query{display:inline;margin-right:.25rem;width:10rem;vertical-align:bottom}#search-error{margin-top:.75rem;color:red;text-align:center}#search-results{margin-top:.75rem}.search-result-type{padding-top:.75rem;text-decoration:underline}.search-result{margin-left:2rem;text-indent:-2rem;line-height:1.2}.search-result-path{padding-left:1rem;font-style:italic;font-size:.75rem}",document.head.appendChild(e);import{d as t,S as s,P as r,f as l,_ as a,n,u,r as o,c as i,a as c,w as d,A as h,t as p,h as m,F as v,g as f,o as y,b as g}from"./index.4821d9ec.js";import{s as x}from"./page.63f8121d.js";var q=t({components:{SButton:s,SInput:r},setup(){const e=n(),t=u();x({title:"Search"});const s=l("");e.query.q&&(s.value=e.query.q,v());const r=l([]),o=l([]),i=l([]),c=l([]),d=l([]),h=l([]),p=l(""),m=l(!1);async function v(){if(s.value=s.value.trim(),!s.value)return r.value=o.value=i.value=d.value=c.value=h.value=[],void(p.value="");s.value!==e.query.q&&t.replace({path:"/search",query:{q:s.value}});const l=(await a.get("/api/search",{params:{q:s.value}})).data;p.value=l.error||"",r.value=l.results.filter((e=>"document"===e.type)),o.value=l.results.filter((e=>"event"===e.type)),i.value=l.results.filter((e=>"folder"===e.type)),c.value=l.results.filter((e=>"person"===e.type)),d.value=l.results.filter((e=>"role"===e.type)),h.value=l.results.filter((e=>"textMessage"===e.type)),m.value=!0}return{documents:r,error:p,events:o,folders:i,onSubmit:v,people:c,query:s,resultPath:function(e){return e.path.length?"in "+e.path.join(" > "):""},roles:d,submitted:m,textMessages:h}}});const b={id:"search"},S={id:"search-query-row"},k=g("Search"),_={id:"search-results"},C=c("div",{class:"search-result-type"},"Roles",-1),M={class:"search-result"},w=c("div",{class:"search-result-type"},"People",-1),P={class:"search-result"},$=c("div",{class:"search-result-type"},"Events",-1),F={class:"search-result"},j=c("div",{class:"search-result-type"},"Folders",-1),V={class:"search-result"},B=c("div",{class:"search-result-type"},"Files",-1),E={class:"search-result"},I=c("div",{class:"search-result-type"},"Text Messages",-1),N={class:"search-result"},T={key:6,class:"search-result"};q.render=function(e,t,s,r,l,a){const n=o("SInput"),u=o("SButton"),x=o("router-link");return y(),i("div",b,[c("form",{id:"search-form",onSubmit:t[2]||(t[2]=h(((...t)=>e.onSubmit&&e.onSubmit(...t)),["prevent"]))},[c("div",S,[c(n,{id:"search-query",autofocus:"",modelValue:e.query,"onUpdate:modelValue":t[1]||(t[1]=t=>e.query=t)},null,8,["modelValue"]),c(u,{type:"submit",variant:"primary"},{default:d((()=>[k])),_:1})])],32),e.error?(y(),i("div",{key:0,id:"search-error",textContent:p(e.error)},null,8,["textContent"])):m("",!0),c("div",_,[e.roles.length?(y(),i(v,{key:0},[C,(y(!0),i(v,null,f(e.roles,(e=>(y(),i("div",M,[c(x,{to:`/people/list?role=${e.id}`},{default:d((()=>[g(p(e.name),1)])),_:2},1032,["to"])])))),256))],64)):m("",!0),e.people.length?(y(),i(v,{key:1},[w,(y(!0),i(v,null,f(e.people,(e=>(y(),i("div",P,[c(x,{to:`/people/${e.id}`},{default:d((()=>[g(p(e.informalName),1)])),_:2},1032,["to"])])))),256))],64)):m("",!0),e.events.length?(y(),i(v,{key:2},[$,(y(!0),i(v,null,f(e.events,(e=>(y(),i("div",F,[c(x,{to:`/events/${e.id}`},{default:d((()=>[g(p(e.date)+" "+p(e.name),1)])),_:2},1032,["to"])])))),256))],64)):m("",!0),e.folders.length?(y(),i(v,{key:3},[j,(y(!0),i(v,null,f(e.folders,(t=>(y(),i("div",V,[c(x,{to:`/files${t.url}`},{default:d((()=>[g(p(t.name),1)])),_:2},1032,["to"]),c("span",{class:"search-result-path",textContent:p(e.resultPath(t))},null,8,["textContent"])])))),256))],64)):m("",!0),e.documents.length?(y(),i(v,{key:4},[B,(y(!0),i(v,null,f(e.documents,(t=>(y(),i("div",E,[c("a",{href:t.url,target:t.newtab?"_blank":null},p(t.name),9,["href","target"]),c("span",{class:"search-result-path",textContent:p(e.resultPath(t))},null,8,["textContent"])])))),256))],64)):m("",!0),e.textMessages.length?(y(),i(v,{key:5},[I,(y(!0),i(v,null,f(e.textMessages,(e=>(y(),i("div",N,[c(x,{to:`/texts/${e.id}`},{default:d((()=>[g("From "+p(e.sender)+" on "+p(e.timestamp.substr(0,10)),1)])),_:2},1032,["to"])])))),256))],64)):m("",!0),!e.submitted||e.error||e.roles.length||e.people.length||e.events.length||e.folders.length||e.documents.length||e.textMessages.length?m("",!0):(y(),i("div",T,"No results found."))])])};export default q;