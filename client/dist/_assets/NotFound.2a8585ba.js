let t=document.createElement("style");t.innerHTML="#notfound-top{margin:0 auto;padding:0 .75rem;max-width:calc(.75rem + 10rem + 20rem + .75rem + 1.5rem)}#notfound-banner{margin-top:1rem;text-align:center;font-weight:700;font-size:1.5rem}#notfound-intro{margin-top:.5rem;text-align:center;font-size:.9rem;line-height:1.2}#notfound-button{margin-top:1rem;text-align:center}",document.head.appendChild(t);import{d as e,S as n,m as o,i as a,r,c as i,a as d,t as m,b as u,o as s}from"./index.4821d9ec.js";import{s as p}from"./page.63f8121d.js";var c=e({components:{SButton:n},setup(){p({title:""});const t=a("me");return{buttonLabel:o((()=>t.value?"Go to Home Page":"Go to Login Page"))}}});const f={id:"notfound-top"},g=d("div",{id:"notfound-banner"},"Page Not Found",-1),l={id:"notfound-intro"},b=u("The page you attempted to reach does not exist."),x={id:"notfound-button"};c.render=function(t,e,n,o,a,u){const p=r("SButton");return s(),i("div",f,[g,d("div",l,[b,d("div",x,[d(p,{to:"/",variant:"primary",textContent:m(t.buttonLabel)},null,8,["textContent"])])])])};export default c;