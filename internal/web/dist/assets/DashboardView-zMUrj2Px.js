import{u as Q}from"./index-DhMb_Ktv.js";import{N as U,a as I}from"./Grid-x4-dfNlT.js";import{f as E,d as N,h as t,N as V,n as $,a7 as X,a8 as H,a9 as Y,aa as Z,ab as ee,ac as B,s as W,p as u,q as P,z as J,A as j,H as K,P as G,v as D,y as T,ad as re,D as te,X as ie,Y as b,Z as C,_ as x,$ as R,ae as oe,r as L,a1 as le}from"./index-DDP37JBl.js";import{N as M}from"./DataTable-CeYApLa5.js";function ne(e){const{infoColor:i,successColor:c,warningColor:g,errorColor:a,textColor2:l,progressRailColor:r,fontSize:o,fontWeight:f}=e;return{fontSize:o,fontSizeCircle:"28px",fontWeightCircle:f,railColor:r,railHeight:"8px",iconSizeCircle:"36px",iconSizeLine:"18px",iconColor:i,iconColorInfo:i,iconColorSuccess:c,iconColorWarning:g,iconColorError:a,textColorCircle:l,textColorLineInner:"rgb(255, 255, 255)",textColorLineOuter:l,fillColor:i,fillColorInfo:i,fillColorSuccess:c,fillColorWarning:g,fillColorError:a,lineBgProcessing:"linear-gradient(90deg, rgba(255, 255, 255, .3) 0%, rgba(255, 255, 255, .5) 100%)"}}const ae={common:E,self:ne};function se(e){const{textColor2:i,textColor3:c,fontSize:g,fontWeight:a}=e;return{labelFontSize:g,labelFontWeight:a,valueFontWeight:a,valueFontSize:"24px",labelTextColor:c,valuePrefixTextColor:i,valueSuffixTextColor:i,valueTextColor:i}}const ce={common:E,self:se},ue={success:t(Z,null),error:t(Y,null),warning:t(H,null),info:t(X,null)},de=N({name:"ProgressCircle",props:{clsPrefix:{type:String,required:!0},status:{type:String,required:!0},strokeWidth:{type:Number,required:!0},fillColor:[String,Object],railColor:String,railStyle:[String,Object],percentage:{type:Number,default:0},offsetDegree:{type:Number,default:0},showIndicator:{type:Boolean,required:!0},indicatorTextColor:String,unit:String,viewBoxWidth:{type:Number,required:!0},gapDegree:{type:Number,required:!0},gapOffsetDegree:{type:Number,default:0}},setup(e,{slots:i}){const c=$(()=>{const l="gradient",{fillColor:r}=e;return typeof r=="object"?`${l}-${ee(JSON.stringify(r))}`:l});function g(l,r,o,f){const{gapDegree:p,viewBoxWidth:h,strokeWidth:v}=e,s=50,m=0,d=s,n=0,S=2*s,w=50+v/2,y=`M ${w},${w} m ${m},${d}
      a ${s},${s} 0 1 1 ${n},${-S}
      a ${s},${s} 0 1 1 ${-n},${S}`,z=Math.PI*2*s,k={stroke:f==="rail"?o:typeof e.fillColor=="object"?`url(#${c.value})`:o,strokeDasharray:`${Math.min(l,100)/100*(z-p)}px ${h*8}px`,strokeDashoffset:`-${p/2}px`,transformOrigin:r?"center":void 0,transform:r?`rotate(${r}deg)`:void 0};return{pathString:y,pathStyle:k}}const a=()=>{const l=typeof e.fillColor=="object",r=l?e.fillColor.stops[0]:"",o=l?e.fillColor.stops[1]:"";return l&&t("defs",null,t("linearGradient",{id:c.value,x1:"0%",y1:"100%",x2:"100%",y2:"0%"},t("stop",{offset:"0%","stop-color":r}),t("stop",{offset:"100%","stop-color":o})))};return()=>{const{fillColor:l,railColor:r,strokeWidth:o,offsetDegree:f,status:p,percentage:h,showIndicator:v,indicatorTextColor:s,unit:m,gapOffsetDegree:d,clsPrefix:n}=e,{pathString:S,pathStyle:w}=g(100,0,r,"rail"),{pathString:y,pathStyle:z}=g(h,f,l,"fill"),k=100+o;return t("div",{class:`${n}-progress-content`,role:"none"},t("div",{class:`${n}-progress-graph`,"aria-hidden":!0},t("div",{class:`${n}-progress-graph-circle`,style:{transform:d?`rotate(${d}deg)`:void 0}},t("svg",{viewBox:`0 0 ${k} ${k}`},a(),t("g",null,t("path",{class:`${n}-progress-graph-circle-rail`,d:S,"stroke-width":o,"stroke-linecap":"round",fill:"none",style:w})),t("g",null,t("path",{class:[`${n}-progress-graph-circle-fill`,h===0&&`${n}-progress-graph-circle-fill--empty`],d:y,"stroke-width":o,"stroke-linecap":"round",fill:"none",style:z}))))),v?t("div",null,i.default?t("div",{class:`${n}-progress-custom-content`,role:"none"},i.default()):p!=="default"?t("div",{class:`${n}-progress-icon`,"aria-hidden":!0},t(V,{clsPrefix:n},{default:()=>ue[p]})):t("div",{class:`${n}-progress-text`,style:{color:s},role:"none"},t("span",{class:`${n}-progress-text__percentage`},h),t("span",{class:`${n}-progress-text__unit`},m))):null)}}}),fe={success:t(Z,null),error:t(Y,null),warning:t(H,null),info:t(X,null)},ge=N({name:"ProgressLine",props:{clsPrefix:{type:String,required:!0},percentage:{type:Number,default:0},railColor:String,railStyle:[String,Object],fillColor:[String,Object],status:{type:String,required:!0},indicatorPlacement:{type:String,required:!0},indicatorTextColor:String,unit:{type:String,default:"%"},processing:{type:Boolean,required:!0},showIndicator:{type:Boolean,required:!0},height:[String,Number],railBorderRadius:[String,Number],fillBorderRadius:[String,Number]},setup(e,{slots:i}){const c=$(()=>B(e.height)),g=$(()=>{var r,o;return typeof e.fillColor=="object"?`linear-gradient(to right, ${(r=e.fillColor)===null||r===void 0?void 0:r.stops[0]} , ${(o=e.fillColor)===null||o===void 0?void 0:o.stops[1]})`:e.fillColor}),a=$(()=>e.railBorderRadius!==void 0?B(e.railBorderRadius):e.height!==void 0?B(e.height,{c:.5}):""),l=$(()=>e.fillBorderRadius!==void 0?B(e.fillBorderRadius):e.railBorderRadius!==void 0?B(e.railBorderRadius):e.height!==void 0?B(e.height,{c:.5}):"");return()=>{const{indicatorPlacement:r,railColor:o,railStyle:f,percentage:p,unit:h,indicatorTextColor:v,status:s,showIndicator:m,processing:d,clsPrefix:n}=e;return t("div",{class:`${n}-progress-content`,role:"none"},t("div",{class:`${n}-progress-graph`,"aria-hidden":!0},t("div",{class:[`${n}-progress-graph-line`,{[`${n}-progress-graph-line--indicator-${r}`]:!0}]},t("div",{class:`${n}-progress-graph-line-rail`,style:[{backgroundColor:o,height:c.value,borderRadius:a.value},f]},t("div",{class:[`${n}-progress-graph-line-fill`,d&&`${n}-progress-graph-line-fill--processing`],style:{maxWidth:`${e.percentage}%`,background:g.value,height:c.value,lineHeight:c.value,borderRadius:l.value}},r==="inside"?t("div",{class:`${n}-progress-graph-line-indicator`,style:{color:v}},i.default?i.default():`${p}${h}`):null)))),m&&r==="outside"?t("div",null,i.default?t("div",{class:`${n}-progress-custom-content`,style:{color:v},role:"none"},i.default()):s==="default"?t("div",{role:"none",class:`${n}-progress-icon ${n}-progress-icon--as-text`,style:{color:v}},p,h):t("div",{class:`${n}-progress-icon`,"aria-hidden":!0},t(V,{clsPrefix:n},{default:()=>fe[s]}))):null)}}});function A(e,i,c=100){return`m ${c/2} ${c/2-e} a ${e} ${e} 0 1 1 0 ${2*e} a ${e} ${e} 0 1 1 0 -${2*e}`}const pe=N({name:"ProgressMultipleCircle",props:{clsPrefix:{type:String,required:!0},viewBoxWidth:{type:Number,required:!0},percentage:{type:Array,default:[0]},strokeWidth:{type:Number,required:!0},circleGap:{type:Number,required:!0},showIndicator:{type:Boolean,required:!0},fillColor:{type:Array,default:()=>[]},railColor:{type:Array,default:()=>[]},railStyle:{type:Array,default:()=>[]}},setup(e,{slots:i}){const c=$(()=>e.percentage.map((l,r)=>`${Math.PI*l/100*(e.viewBoxWidth/2-e.strokeWidth/2*(1+2*r)-e.circleGap*r)*2}, ${e.viewBoxWidth*8}`)),g=(a,l)=>{const r=e.fillColor[l],o=typeof r=="object"?r.stops[0]:"",f=typeof r=="object"?r.stops[1]:"";return typeof e.fillColor[l]=="object"&&t("linearGradient",{id:`gradient-${l}`,x1:"100%",y1:"0%",x2:"0%",y2:"100%"},t("stop",{offset:"0%","stop-color":o}),t("stop",{offset:"100%","stop-color":f}))};return()=>{const{viewBoxWidth:a,strokeWidth:l,circleGap:r,showIndicator:o,fillColor:f,railColor:p,railStyle:h,percentage:v,clsPrefix:s}=e;return t("div",{class:`${s}-progress-content`,role:"none"},t("div",{class:`${s}-progress-graph`,"aria-hidden":!0},t("div",{class:`${s}-progress-graph-circle`},t("svg",{viewBox:`0 0 ${a} ${a}`},t("defs",null,v.map((m,d)=>g(m,d))),v.map((m,d)=>t("g",{key:d},t("path",{class:`${s}-progress-graph-circle-rail`,d:A(a/2-l/2*(1+2*d)-r*d,l,a),"stroke-width":l,"stroke-linecap":"round",fill:"none",style:[{strokeDashoffset:0,stroke:p[d]},h[d]]}),t("path",{class:[`${s}-progress-graph-circle-fill`,m===0&&`${s}-progress-graph-circle-fill--empty`],d:A(a/2-l/2*(1+2*d)-r*d,l,a),"stroke-width":l,"stroke-linecap":"round",fill:"none",style:{strokeDasharray:c.value[d],strokeDashoffset:0,stroke:typeof f[d]=="object"?`url(#gradient-${d})`:f[d]}})))))),o&&i.default?t("div",null,t("div",{class:`${s}-progress-text`},i.default())):null)}}}),he=W([u("progress",{display:"inline-block"},[u("progress-icon",`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 `),P("line",`
 width: 100%;
 display: block;
 `,[u("progress-content",`
 display: flex;
 align-items: center;
 `,[u("progress-graph",{flex:1})]),u("progress-custom-content",{marginLeft:"14px"}),u("progress-icon",`
 width: 30px;
 padding-left: 14px;
 height: var(--n-icon-size-line);
 line-height: var(--n-icon-size-line);
 font-size: var(--n-icon-size-line);
 `,[P("as-text",`
 color: var(--n-text-color-line-outer);
 text-align: center;
 width: 40px;
 font-size: var(--n-font-size);
 padding-left: 4px;
 transition: color .3s var(--n-bezier);
 `)])]),P("circle, dashboard",{width:"120px"},[u("progress-custom-content",`
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 display: flex;
 align-items: center;
 justify-content: center;
 `),u("progress-text",`
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 display: flex;
 align-items: center;
 color: inherit;
 font-size: var(--n-font-size-circle);
 color: var(--n-text-color-circle);
 font-weight: var(--n-font-weight-circle);
 transition: color .3s var(--n-bezier);
 white-space: nowrap;
 `),u("progress-icon",`
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 display: flex;
 align-items: center;
 color: var(--n-icon-color);
 font-size: var(--n-icon-size-circle);
 `)]),P("multiple-circle",`
 width: 200px;
 color: inherit;
 `,[u("progress-text",`
 font-weight: var(--n-font-weight-circle);
 color: var(--n-text-color-circle);
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 display: flex;
 align-items: center;
 justify-content: center;
 transition: color .3s var(--n-bezier);
 `)]),u("progress-content",{position:"relative"}),u("progress-graph",{position:"relative"},[u("progress-graph-circle",[W("svg",{verticalAlign:"bottom"}),u("progress-graph-circle-fill",`
 stroke: var(--n-fill-color);
 transition:
 opacity .3s var(--n-bezier),
 stroke .3s var(--n-bezier),
 stroke-dasharray .3s var(--n-bezier);
 `,[P("empty",{opacity:0})]),u("progress-graph-circle-rail",`
 transition: stroke .3s var(--n-bezier);
 overflow: hidden;
 stroke: var(--n-rail-color);
 `)]),u("progress-graph-line",[P("indicator-inside",[u("progress-graph-line-rail",`
 height: 16px;
 line-height: 16px;
 border-radius: 10px;
 `,[u("progress-graph-line-fill",`
 height: inherit;
 border-radius: 10px;
 `),u("progress-graph-line-indicator",`
 background: #0000;
 white-space: nowrap;
 text-align: right;
 margin-left: 14px;
 margin-right: 14px;
 height: inherit;
 font-size: 12px;
 color: var(--n-text-color-line-inner);
 transition: color .3s var(--n-bezier);
 `)])]),P("indicator-inside-label",`
 height: 16px;
 display: flex;
 align-items: center;
 `,[u("progress-graph-line-rail",`
 flex: 1;
 transition: background-color .3s var(--n-bezier);
 `),u("progress-graph-line-indicator",`
 background: var(--n-fill-color);
 font-size: 12px;
 transform: translateZ(0);
 display: flex;
 vertical-align: middle;
 height: 16px;
 line-height: 16px;
 padding: 0 10px;
 border-radius: 10px;
 position: absolute;
 white-space: nowrap;
 color: var(--n-text-color-line-inner);
 transition:
 right .2s var(--n-bezier),
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `)]),u("progress-graph-line-rail",`
 position: relative;
 overflow: hidden;
 height: var(--n-rail-height);
 border-radius: 5px;
 background-color: var(--n-rail-color);
 transition: background-color .3s var(--n-bezier);
 `,[u("progress-graph-line-fill",`
 background: var(--n-fill-color);
 position: relative;
 border-radius: 5px;
 height: inherit;
 width: 100%;
 max-width: 0%;
 transition:
 background-color .3s var(--n-bezier),
 max-width .2s var(--n-bezier);
 `,[P("processing",[W("&::after",`
 content: "";
 background-image: var(--n-line-bg-processing);
 animation: progress-processing-animation 2s var(--n-bezier) infinite;
 `)])])])])])]),W("@keyframes progress-processing-animation",`
 0% {
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 right: 100%;
 opacity: 1;
 }
 66% {
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 right: 0;
 opacity: 0;
 }
 100% {
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 right: 0;
 opacity: 0;
 }
 `)]),ve=Object.assign(Object.assign({},j.props),{processing:Boolean,type:{type:String,default:"line"},gapDegree:Number,gapOffsetDegree:Number,status:{type:String,default:"default"},railColor:[String,Array],railStyle:[String,Array],color:[String,Array,Object],viewBoxWidth:{type:Number,default:100},strokeWidth:{type:Number,default:7},percentage:[Number,Array],unit:{type:String,default:"%"},showIndicator:{type:Boolean,default:!0},indicatorPosition:{type:String,default:"outside"},indicatorPlacement:{type:String,default:"outside"},indicatorTextColor:String,circleGap:{type:Number,default:1},height:Number,borderRadius:[String,Number],fillBorderRadius:[String,Number],offsetDegree:Number}),me=N({name:"Progress",props:ve,setup(e){const i=$(()=>e.indicatorPlacement||e.indicatorPosition),c=$(()=>{if(e.gapDegree||e.gapDegree===0)return e.gapDegree;if(e.type==="dashboard")return 75}),{mergedClsPrefixRef:g,inlineThemeDisabled:a}=J(e),l=j("Progress","-progress",he,ae,e,g),r=$(()=>{const{status:f}=e,{common:{cubicBezierEaseInOut:p},self:{fontSize:h,fontSizeCircle:v,railColor:s,railHeight:m,iconSizeCircle:d,iconSizeLine:n,textColorCircle:S,textColorLineInner:w,textColorLineOuter:y,lineBgProcessing:z,fontWeightCircle:k,[G("iconColor",f)]:q,[G("fillColor",f)]:_}}=l.value;return{"--n-bezier":p,"--n-fill-color":_,"--n-font-size":h,"--n-font-size-circle":v,"--n-font-weight-circle":k,"--n-icon-color":q,"--n-icon-size-circle":d,"--n-icon-size-line":n,"--n-line-bg-processing":z,"--n-rail-color":s,"--n-rail-height":m,"--n-text-color-circle":S,"--n-text-color-line-inner":w,"--n-text-color-line-outer":y}}),o=a?K("progress",$(()=>e.status[0]),r,e):void 0;return{mergedClsPrefix:g,mergedIndicatorPlacement:i,gapDeg:c,cssVars:a?void 0:r,themeClass:o==null?void 0:o.themeClass,onRender:o==null?void 0:o.onRender}},render(){const{type:e,cssVars:i,indicatorTextColor:c,showIndicator:g,status:a,railColor:l,railStyle:r,color:o,percentage:f,viewBoxWidth:p,strokeWidth:h,mergedIndicatorPlacement:v,unit:s,borderRadius:m,fillBorderRadius:d,height:n,processing:S,circleGap:w,mergedClsPrefix:y,gapDeg:z,gapOffsetDegree:k,themeClass:q,$slots:_,onRender:F}=this;return F==null||F(),t("div",{class:[q,`${y}-progress`,`${y}-progress--${e}`,`${y}-progress--${a}`],style:i,"aria-valuemax":100,"aria-valuemin":0,"aria-valuenow":f,role:e==="circle"||e==="line"||e==="dashboard"?"progressbar":"none"},e==="circle"||e==="dashboard"?t(de,{clsPrefix:y,status:a,showIndicator:g,indicatorTextColor:c,railColor:l,fillColor:o,railStyle:r,offsetDegree:this.offsetDegree,percentage:f,viewBoxWidth:p,strokeWidth:h,gapDegree:z===void 0?e==="dashboard"?75:0:z,gapOffsetDegree:k,unit:s},_):e==="line"?t(ge,{clsPrefix:y,status:a,showIndicator:g,indicatorTextColor:c,railColor:l,fillColor:o,railStyle:r,percentage:f,processing:S,indicatorPlacement:v,unit:s,fillBorderRadius:d,railBorderRadius:m,height:n},_):e==="multiple-circle"?t(pe,{clsPrefix:y,strokeWidth:h,railColor:l,fillColor:o,railStyle:r,viewBoxWidth:p,percentage:f,showIndicator:g,circleGap:w},_):null)}}),be=u("statistic",[D("label",`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),u("statistic-value",`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[D("prefix",`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[u("icon",{verticalAlign:"-0.125em"})]),D("content",`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),D("suffix",`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[u("icon",{verticalAlign:"-0.125em"})])])]),xe=Object.assign(Object.assign({},j.props),{tabularNums:Boolean,label:String,value:[String,Number]}),O=N({name:"Statistic",props:xe,slots:Object,setup(e){const{mergedClsPrefixRef:i,inlineThemeDisabled:c,mergedRtlRef:g}=J(e),a=j("Statistic","-statistic",be,ce,e,i),l=re("Statistic",g,i),r=$(()=>{const{self:{labelFontWeight:f,valueFontSize:p,valueFontWeight:h,valuePrefixTextColor:v,labelTextColor:s,valueSuffixTextColor:m,valueTextColor:d,labelFontSize:n},common:{cubicBezierEaseInOut:S}}=a.value;return{"--n-bezier":S,"--n-label-font-size":n,"--n-label-font-weight":f,"--n-label-text-color":s,"--n-value-font-weight":h,"--n-value-font-size":p,"--n-value-prefix-text-color":v,"--n-value-suffix-text-color":m,"--n-value-text-color":d}}),o=c?K("statistic",void 0,r,e):void 0;return{rtlEnabled:l,mergedClsPrefix:i,cssVars:c?void 0:r,themeClass:o==null?void 0:o.themeClass,onRender:o==null?void 0:o.onRender}},render(){var e;const{mergedClsPrefix:i,$slots:{default:c,label:g,prefix:a,suffix:l}}=this;return(e=this.onRender)===null||e===void 0||e.call(this),t("div",{class:[`${i}-statistic`,this.themeClass,this.rtlEnabled&&`${i}-statistic--rtl`],style:this.cssVars},T(g,r=>t("div",{class:`${i}-statistic__label`},this.label||r)),t("div",{class:`${i}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?"tabular-nums":""}},T(a,r=>r&&t("span",{class:`${i}-statistic-value__prefix`},r)),this.value!==void 0?t("span",{class:`${i}-statistic-value__content`},this.value):T(c,r=>r&&t("span",{class:`${i}-statistic-value__content`},r)),T(l,r=>r&&t("span",{class:`${i}-statistic-value__suffix`},r))))}}),ye={key:0},ze=N({__name:"DashboardView",setup(e){const i=L(null),c=L([]),g=[{title:"时间",key:"time"},{title:"增量",key:"flow",render:r=>`${(r.flow/1024/1024).toFixed(2)} MB`},{title:"累计",key:"totalFlow",render:r=>`${(r.totalFlow/1024/1024/1024).toFixed(2)} GB`}],a=[{title:"名称",key:"name"},{title:"入口",key:"entry",render:r=>`${r.inIp}:${r.inPort}`},{title:"状态",key:"status",render:r=>r.status===1?"运行中":"已暂停"}];te(async()=>{try{i.value=await Q(),c.value=i.value.statisticsFlows||[]}catch{}});function l(r){var f,p,h,v,s,m;const o=((p=(f=i.value)==null?void 0:f.userInfo)==null?void 0:p.flow)||1;return Math.min(100,Math.round((((v=(h=i.value)==null?void 0:h.userInfo)==null?void 0:v.inFlow)+((m=(s=i.value)==null?void 0:s.userInfo)==null?void 0:m.outFlow))/(o*1024**3)*100))}return(r,o)=>i.value?(le(),ie("div",ye,[b(x(U),{cols:4,"x-gap":12,responsive:"screen"},{default:C(()=>[b(x(I),null,{default:C(()=>[b(x(R),null,{default:C(()=>[b(x(O),{label:"账户",value:i.value.userInfo.user},null,8,["value"])]),_:1})]),_:1}),b(x(I),null,{default:C(()=>[b(x(R),null,{default:C(()=>[b(x(O),{label:"流量使用",value:l(0),suffix:"%"},{prefix:C(()=>[...o[0]||(o[0]=[])]),_:1},8,["value"]),b(x(me),{type:"line",percentage:l(0),"show-indicator":!1},null,8,["percentage"])]),_:1})]),_:1}),b(x(I),null,{default:C(()=>[b(x(R),null,{default:C(()=>[b(x(O),{label:"到期时间",value:new Date(i.value.userInfo.expTime).toLocaleDateString()},null,8,["value"])]),_:1})]),_:1}),b(x(I),null,{default:C(()=>[b(x(R),null,{default:C(()=>[b(x(O),{label:"转发数量",value:i.value.forwards.length,suffix:`/ ${i.value.userInfo.num}`},null,8,["value","suffix"])]),_:1})]),_:1})]),_:1}),b(x(R),{title:"我的转发",style:{"margin-top":"12px"}},{default:C(()=>[b(x(M),{columns:a,data:i.value.forwards,bordered:!1,size:"small"},null,8,["data"])]),_:1}),b(x(R),{title:"近 24 小时流量",style:{"margin-top":"12px"}},{default:C(()=>[b(x(M),{columns:g,data:c.value.slice(-24),bordered:!1,size:"small"},null,8,["data"])]),_:1})])):oe("",!0)}});export{ze as default};
