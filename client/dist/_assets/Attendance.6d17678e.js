let e=document.createElement("style");e.innerHTML="#attrep-params{margin-bottom:1.5rem}#attrep-params td:first-child{padding-right:1rem}@media print{#attrep-params{display:none}}#attrep{margin:1.5rem .75rem;overflow-x:auto}.attrep-col-h{border-left:2px solid #888}.attrep-col-h2{border-left:1px solid #ccc}.attrep-col-s{border-left:2px solid #888}.attrep-col-c{border-left:1px solid #ccc}.attrep-col-t{border-left:1px solid #ccc;border-right:2px solid #888}.attrep-col-1{border-left:2px solid #888;border-right:2px solid #888}.attrep-row-h,.attrep-row-s,.attrep-row-t{border-top:2px solid #888}.attrep-row-h2{border-top:1px solid #ccc}.attrep-row-t{border-bottom:2px solid #888}.attrep-row-c2,.attrep-row-tc2{background-color:#eee}.attrep-row-h2>.attrep-col-h,.attrep-row-h2>.attrep-col-h2,.attrep-row-h>.attrep-col-h,.attrep-row-h>.attrep-col-h2{border-left:hidden;border-top:hidden}.attrep-row-h>.attrep-col-1,.attrep-row-h>.attrep-col-c,.attrep-row-h>.attrep-col-s,.attrep-row-h>.attrep-col-t{text-align:right;vertical-align:bottom;padding-right:.5rem}.attrep-vertical{-ms-writing-mode:tb-rl;writing-mode:vertical-rl;font-variant-numeric:tabular-nums;line-height:1;width:100%;padding-top:.5rem;padding-bottom:.5rem}.attrep-col-1,.attrep-col-c,.attrep-col-s,.attrep-col-t{padding-right:.5rem;min-width:4rem;text-align:right;font-variant-numeric:tabular-nums}.attrep-col-h,.attrep-col-h2{padding-left:.5rem;padding-right:.5rem;white-space:nowrap}.attrep-col-t,.attrep-row-t{font-weight:700}#attrep-buttons,#attrep-count{margin-top:1.5rem}",document.head.appendChild(e);import{d as t,S as a,J as l,a9 as o,aa as r,Q as n,f as s,T as p,z as d,n as u,u as i,_ as c,r as m,c as v,a as g,b as w,t as h,F as y,g as b,h as S,o as V,w as f}from"./index.be4dc8e3.js";import{s as x}from"./page.46d1c173.js";const C=[{value:"p",label:"Person"},{value:"o",label:"Org"},{value:"po",label:"Person, Org"},{value:"op",label:"Org, Person"}],T=[{value:"e",label:"Events"},{value:"m",label:"Months"}],k=[{value:"h",label:"Cumulative Hours"},{value:"c",label:"Attendance Counts"}];var Z=t({components:{SButton:a,SCheck:l,SCheckGroup:o,SRadioGroup:r,SSelect:n},setup(){const e=u(),t=i();x({title:"Attendance Report",browserTitle:"Attendance"});const a=s({}),l=s({}),o=s([]),r=s([]),n=s([]),m=s(new Set),v=s(new Set),g=s(new Set),w=s(0);return p((async()=>{const t=(await c.get("/api/reports/attendance",{params:e.query})).data;a.value=t.parameters,l.value=t.options,o.value=t.columns,r.value=t.rows,n.value=t.cells,m.value=new Set(t.parameters.orgs),v.value=new Set(t.parameters.eventTypes),g.value=new Set(t.parameters.attendanceTypes),w.value=t.personCount||0})),d([a,g,v,m],(()=>{const e={};e.dateRange=a.value.dateRange,e.rows=a.value.rows,e.columns=a.value.columns,e.cells=a.value.cells,e.orgs=Array.from(m.value.keys(),(e=>e.toString())).sort().join(","),e.eventTypes=Array.from(v.value.keys(),(e=>e.toString())).sort().join(","),e.attendanceTypes=Array.from(g.value.keys(),(e=>e.toString())).sort().join(","),a.value.includeZerosX&&(e.includeZerosX="true"),a.value.includeZerosY&&(e.includeZerosY="true"),t.replace({path:"/reports/attendance",query:e})}),{deep:!0}),{attendanceTypes:g,cells:n,cellOptions:k,columns:o,columnOptions:T,exportCSV:function(){const e=new URLSearchParams;e.set("dateRange",a.value.dateRange),e.set("rows",a.value.rows),e.set("columns",a.value.columns),e.set("cells",a.value.cells),e.set("orgs",Array.from(m.value.keys(),(e=>e.toString())).sort().join(",")),e.set("eventTypes",Array.from(v.value.keys(),(e=>e.toString())).sort().join(",")),e.set("attendanceTypes",Array.from(g.value.keys(),(e=>e.toString())).sort().join(",")),a.value.includeZerosX&&e.set("includeZerosX","true"),a.value.includeZerosY&&e.set("includeZerosY","true"),e.set("format","csv"),window.location.href="/api/reports/attendance?"+e.toString()},eventTypes:v,options:l,orgs:m,params:a,personCount:w,rows:r,rowOptions:C}}});const R={key:0,id:"attrep"},A={id:"attrep-params"},O=g("td",null,"Date range",-1),U=g("td",null,"Rows",-1),j=g("td",null,"Columns",-1),X=g("td",null,"Cells",-1),Y=g("td",null,"Organizations",-1),K=g("td",null,"Events",-1),E=g("td",null,"Attendees",-1),G={id:"attrep-table"},P={key:1,id:"attrep-buttons"},$=w("Export");Z.render=function(e,t,a,l,o,r){const n=m("SSelect"),s=m("SRadioGroup"),p=m("SCheck"),d=m("SCheckGroup"),u=m("SButton");return e.params.cells?(V(),v("div",R,[g("table",A,[g("tr",null,[O,g("td",null,[g(n,{options:e.options.dateRanges,valueKey:"tag",labelKey:"label",modelValue:e.params.dateRange,"onUpdate:modelValue":t[1]||(t[1]=t=>e.params.dateRange=t)},null,8,["options","modelValue"]),w(" "+h(e.params.dateFrom)+" to "+h(e.params.dateTo),1)])]),g("tr",null,[U,g("td",null,[g(s,{id:"attrep-rows",inline:"",options:e.rowOptions,modelValue:e.params.rows,"onUpdate:modelValue":t[2]||(t[2]=t=>e.params.rows=t)},null,8,["options","modelValue"]),g(p,{id:"attrep-includeZerosY",label:"Include Zeros",inline:"",modelValue:e.params.includeZerosY,"onUpdate:modelValue":t[3]||(t[3]=t=>e.params.includeZerosY=t)},null,8,["modelValue"])])]),g("tr",null,[j,g("td",null,[g(s,{id:"attrep-columns",inline:"",options:e.columnOptions,modelValue:e.params.columns,"onUpdate:modelValue":t[4]||(t[4]=t=>e.params.columns=t)},null,8,["options","modelValue"]),g(p,{id:"attrep-includeZerosX",label:"Include Zeros",inline:"",modelValue:e.params.includeZerosX,"onUpdate:modelValue":t[5]||(t[5]=t=>e.params.includeZerosX=t)},null,8,["modelValue"])])]),g("tr",null,[X,g("td",null,[g(s,{id:"attrep-cells",inline:"",options:e.cellOptions,modelValue:e.params.cells,"onUpdate:modelValue":t[6]||(t[6]=t=>e.params.cells=t)},null,8,["options","modelValue"])])]),g("tr",null,[Y,g("td",null,[g(d,{id:"attrep-orgs",inline:"",options:e.options.orgs,valueKey:"id",modelValue:e.orgs,"onUpdate:modelValue":t[7]||(t[7]=t=>e.orgs=t)},null,8,["options","modelValue"])])]),g("tr",null,[K,g("td",null,[g(d,{id:"attrep-eventTypes",inline:"",options:e.options.eventTypes,valueKey:"id",modelValue:e.eventTypes,"onUpdate:modelValue":t[8]||(t[8]=t=>e.eventTypes=t)},null,8,["options","modelValue"])])]),g("tr",null,[E,g("td",null,[g(d,{id:"attrep-attendanceTypes",inline:"",options:e.options.attendanceTypes,valueKey:"id",modelValue:e.attendanceTypes,"onUpdate:modelValue":t[9]||(t[9]=t=>e.attendanceTypes=t)},null,8,["options","modelValue"])])])]),g("table",G,[g("tbody",null,[(V(!0),v(y,null,b(e.rows,((t,a)=>(V(),v("tr",{class:`attrep-row-${t}`},[0===a?(V(!0),v(y,{key:0},b(e.cells[a],((t,a)=>(V(),v("td",{class:`attrep-col-${e.columns[a]}`},[t?(V(),v("div",{key:0,class:"attrep-vertical",textContent:h(t)},null,8,["textContent"])):S("",!0)],2)))),256)):(V(!0),v(y,{key:1},b(e.cells[a],((t,a)=>(V(),v("td",{class:`attrep-col-${e.columns[a]}`,textContent:h(t)},null,10,["textContent"])))),256))],2)))),256))])]),e.personCount?(V(),v("div",{key:0,id:"attrep-count",textContent:h(e.personCount>1?`${e.personCount} people listed`:"1 person listed")},null,8,["textContent"])):S("",!0),e.rows.length?(V(),v("div",P,[g(u,{variant:"primary",onClick:e.exportCSV},{default:f((()=>[$])),_:1},8,["onClick"])])):S("",!0)])):S("",!0)};export default Z;