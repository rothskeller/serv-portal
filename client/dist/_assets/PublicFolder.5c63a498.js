let e=document.createElement("style");e.innerHTML="#public-sf{position:relative;height:100%}.mouse #public-sf{padding:1.5rem .75rem}#public-sf-name{font-weight:700;font-size:1.25rem}.public-sf-line{display:flex;align-items:center}.touch .public-sf-line{padding:.25rem .5rem;min-height:40px;border-bottom:1px solid #ccc}.public-sf-icon{display:flex;flex:none;justify-content:center;align-items:center;margin-right:.5rem;padding:0;width:1rem;height:1rem;border:none;background-color:#fff;color:#000}.public-sf-icon:active,.public-sf-icon:focus,.public-sf-icon:hover{background-color:#fff!important;color:#000!important}.touch .public-sf-icon{width:1.5rem;height:1.5rem}.public-sf-icon svg{width:100%;height:100%}.public-sf-name{flex:1 1 auto;overflow:hidden;min-width:0;text-overflow:ellipsis;white-space:nowrap}.touch .public-sf-name{white-space:normal;line-height:1.2}",document.head.appendChild(e);import{d as l,e as c,f as n,q as t,u as a,v as s,_ as i,r as o,c as r,a as h,t as p,F as d,g as u,o as m,w as v,b as f}from"./index.67d5d48d.js";import{s as w}from"./page.fbc41cd1.js";var g=l({components:{SSpinner:c},setup(){const e=s();w({title:""});const l=n(null);t((async()=>{const c=e.params.rest?`${e.params.path}/${e.params.rest}`:e.params.path;l.value=(await i.get("/api/folders/",{params:{path:c}})).data,l.value.docDownload&&(location.href=l.value.docDownload),l.value.children||(l.value.children=[]),l.value.documents||(l.value.documents=[])}));const c=a((()=>l.value&&l.value.parent&&l.value.parent.url?l.value.parent.name:"Home Page")),o=a((()=>l.value&&l.value.parent&&l.value.parent.url?l.value.parent.url:"/"));return{docPath:function(e){return`${l.value.url}/${encodeURIComponent(e.name)}`},folder:l,parentName:c,parentURL:o}}});const x={key:0},z={key:1,id:"public-sf"},b={class:"public-sf-line"},C=h("div",{class:"public-sf-icon"},[h("svg",{xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 448 512"},[h("path",{fill:"currentColor",d:"M34.9 289.5l-22.2-22.2c-9.4-9.4-9.4-24.6 0-33.9L207 39c9.4-9.4 24.6-9.4 33.9 0l194.3 194.3c9.4 9.4 9.4 24.6 0 33.9L413 289.4c-9.5 9.5-25 9.3-34.3-.4L264 168.6V456c0 13.3-10.7 24-24 24h-32c-13.3 0-24-10.7-24-24V168.6L69.2 289.1c-9.3 9.8-24.8 10-34.3.4z"})])],-1),H=h("div",{class:"public-sf-icon"},[h("svg",{xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 512 512"},[h("path",{fill:"currentColor",d:"M464 128H272l-54.63-54.63c-6-6-14.14-9.37-22.63-9.37H48C21.49 64 0 85.49 0 112v288c0 26.51 21.49 48 48 48h416c26.51 0 48-21.49 48-48V176c0-26.51-21.49-48-48-48zm0 272H48V112h140.12l54.63 54.63c6 6 14.14 9.37 22.63 9.37H464v224z"})])],-1),V={class:"public-sf-name"},M={class:"public-sf-icon"},L={key:0,xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 384 512"},y=h("path",{fill:"currentColor",d:"M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm220.1-208c-5.7 0-10.6 4-11.7 9.5-20.6 97.7-20.4 95.4-21 103.5-.2-1.2-.4-2.6-.7-4.3-.8-5.1.3.2-23.6-99.5-1.3-5.4-6.1-9.2-11.7-9.2h-13.3c-5.5 0-10.3 3.8-11.7 9.1-24.4 99-24 96.2-24.8 103.7-.1-1.1-.2-2.5-.5-4.2-.7-5.2-14.1-73.3-19.1-99-1.1-5.6-6-9.7-11.8-9.7h-16.8c-7.8 0-13.5 7.3-11.7 14.8 8 32.6 26.7 109.5 33.2 136 1.3 5.4 6.1 9.1 11.7 9.1h25.2c5.5 0 10.3-3.7 11.6-9.1l17.9-71.4c1.5-6.2 2.5-12 3-17.3l2.9 17.3c.1.4 12.6 50.5 17.9 71.4 1.3 5.3 6.1 9.1 11.6 9.1h24.7c5.5 0 10.3-3.7 11.6-9.1 20.8-81.9 30.2-119 34.5-136 1.9-7.6-3.8-14.9-11.6-14.9h-15.8z"},null,-1),k={key:1,xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 384 512"},W=h("path",{fill:"currentColor",d:"M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm250.2-143.7c-12.2-12-47-8.7-64.4-6.5-17.2-10.5-28.7-25-36.8-46.3 3.9-16.1 10.1-40.6 5.4-56-4.2-26.2-37.8-23.6-42.6-5.9-4.4 16.1-.4 38.5 7 67.1-10 23.9-24.9 56-35.4 74.4-20 10.3-47 26.2-51 46.2-3.3 15.8 26 55.2 76.1-31.2 22.4-7.4 46.8-16.5 68.4-20.1 18.9 10.2 41 17 55.8 17 25.5 0 28-28.2 17.5-38.7zm-198.1 77.8c5.1-13.7 24.5-29.5 30.4-35-19 30.3-30.4 35.7-30.4 35zm81.6-190.6c7.4 0 6.7 32.1 1.8 40.8-4.4-13.9-4.3-40.8-1.8-40.8zm-24.4 136.6c9.7-16.9 18-37 24.7-54.7 8.3 15.1 18.9 27.2 30.1 35.5-20.8 4.3-38.9 13.1-54.8 19.2zm131.6-5s-5 6-37.3-7.8c35.1-2.6 40.9 5.4 37.3 7.8z"},null,-1),B={key:2,xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 384 512"},j=h("path",{fill:"currentColor",d:"M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm72-60V236c0-6.6 5.4-12 12-12h69.2c36.7 0 62.8 27 62.8 66.3 0 74.3-68.7 66.5-95.5 66.5V404c0 6.6-5.4 12-12 12H132c-6.6 0-12-5.4-12-12zm48.5-87.4h23c7.9 0 13.9-2.4 18.1-7.2 8.5-9.8 8.4-28.5.1-37.8-4.1-4.6-9.9-7-17.4-7h-23.9v52z"},null,-1),S={key:3,xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 384 512"},$=h("path",{fill:"currentColor",d:"M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm212-240h-28.8c-4.4 0-8.4 2.4-10.5 6.3-18 33.1-22.2 42.4-28.6 57.7-13.9-29.1-6.9-17.3-28.6-57.7-2.1-3.9-6.2-6.3-10.6-6.3H124c-9.3 0-15 10-10.4 18l46.3 78-46.3 78c-4.7 8 1.1 18 10.4 18h28.9c4.4 0 8.4-2.4 10.5-6.3 21.7-40 23-45 28.6-57.7 14.9 30.2 5.9 15.9 28.6 57.7 2.1 3.9 6.2 6.3 10.6 6.3H260c9.3 0 15-10 10.4-18L224 320c.7-1.1 30.3-50.5 46.3-78 4.7-8-1.1-18-10.3-18z"},null,-1),P={key:4,xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 384 512"},R=h("path",{fill:"currentColor",d:"M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48zm32-48h224V288l-23.5-23.5c-4.7-4.7-12.3-4.7-17 0L176 352l-39.5-39.5c-4.7-4.7-12.3-4.7-17 0L80 352v64zm48-240c-26.5 0-48 21.5-48 48s21.5 48 48 48 48-21.5 48-48-21.5-48-48-48z"},null,-1),U={key:5,xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 384 512"},_=h("path",{fill:"currentColor",d:"M369.9 97.9L286 14C277 5 264.8-.1 252.1-.1H48C21.5 0 0 21.5 0 48v416c0 26.5 21.5 48 48 48h288c26.5 0 48-21.5 48-48V131.9c0-12.7-5.1-25-14.1-34zM332.1 128H256V51.9l76.1 76.1zM48 464V48h160v104c0 13.3 10.7 24 24 24h104v288H48z"},null,-1),D={class:"public-sf-name"};g.render=function(e,l,c,n,t,a){const s=o("SSpinner"),i=o("router-link");return e.folder?(m(),r("div",z,[h("div",{id:"public-sf-name",textContent:p(e.folder.name)},null,8,["textContent"]),h("div",b,[C,h(i,{class:"public-sf-name",to:e.parentURL,textContent:p(e.parentName)},null,8,["to","textContent"])]),(m(!0),r(d,null,u(e.folder.children,(e=>(m(),r("div",{class:"public-sf-line",key:"f"+e.id},[H,h("div",V,[h(i,{to:e.url},{default:v((()=>[f(p(e.name),1)])),_:2},1032,["to"])])])))),128)),(m(!0),r(d,null,u(e.folder.documents,(l=>(m(),r("div",{class:"public-sf-line",key:"d"+l.id},[h("div",M,[l.name.endsWith(".docx")||l.name.endsWith(".doc")?(m(),r("svg",L,[y])):l.name.endsWith(".pdf")?(m(),r("svg",k,[W])):l.name.endsWith(".ppt")||l.name.endsWith(".pptx")?(m(),r("svg",B,[j])):l.name.endsWith(".xls")||l.name.endsWith(".xlsx")?(m(),r("svg",S,[$])):l.name.endsWith(".jpg")||l.name.endsWith(".jpeg")?(m(),r("svg",P,[R])):(m(),r("svg",U,[_]))]),h("span",D,[h(i,{to:e.docPath(l)},{default:v((()=>[f(p(l.name),1)])),_:2},1032,["to"])])])))),128))])):(m(),r("div",x,[h(s)]))};export default g;