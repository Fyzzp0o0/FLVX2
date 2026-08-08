import{d as Z,h as p,u as Rt,c as zt,r as _,a as Pe,i as Lt,b as oe,e as $t,f as Bt,g as _t,j as Ie,t as Wt,m as At,F as kt,N as Et,k as Vt,l as jt,n as Q,o as It,p as r,q as b,s as w,v as B,w as Mt,x as de,y as Re,V as be,z as Ht,A as Me,B as Ft,C as ce,D as Ot,E as Gt,G as Nt,H as Dt,I as ze,J as Ut,K as Xt,L as Kt,T as qt,M as Yt,O as fe,P as M,Q as re,R as Jt,S as H,U as ne,W as Qt,X as Zt,Y as R,Z as E,_ as z,$ as ea,a0 as ta,a1 as aa,a2 as ra,a3 as Le,a4 as $e,a5 as Be,a6 as na}from"./index-MlikZPC5.js";import{N as Y,l as oa,r as ia}from"./index-lSbKxEuE.js";import{u as la,N as _e,a as J}from"./FormItem-DsYUrVL0.js";import{A as sa}from"./Add-BnLQNTMQ.js";const da=Pe(".v-x-scroll",{overflow:"auto",scrollbarWidth:"none"},[Pe("&::-webkit-scrollbar",{width:0,height:0})]),ba=Z({name:"XScroll",props:{disabled:Boolean,onScroll:Function},setup(){const e=_(null);function n(s){!(s.currentTarget.offsetWidth<s.currentTarget.scrollWidth)||s.deltaY===0||(s.currentTarget.scrollLeft+=s.deltaY+s.deltaX,s.preventDefault())}const i=Rt();return da.mount({id:"vueuc/x-scroll",head:!0,anchorMetaName:zt,ssr:i}),Object.assign({selfRef:e,handleWheel:n},{scrollTo(...s){var x;(x=e.value)===null||x===void 0||x.scrollTo(...s)}})},render(){return p("div",{ref:"selfRef",onScroll:this.onScroll,onWheel:this.disabled?void 0:this.handleWheel,class:"v-x-scroll"},this.$slots)}});var ca=/\s/;function fa(e){for(var n=e.length;n--&&ca.test(e.charAt(n)););return n}var pa=/^\s+/;function ua(e){return e&&e.slice(0,fa(e)+1).replace(pa,"")}var We=NaN,va=/^[-+]0x[0-9a-f]+$/i,ga=/^0b[01]+$/i,ha=/^0o[0-7]+$/i,xa=parseInt;function Ae(e){if(typeof e=="number")return e;if(Lt(e))return We;if(oe(e)){var n=typeof e.valueOf=="function"?e.valueOf():e;e=oe(n)?n+"":n}if(typeof e!="string")return e===0?e:+e;e=ua(e);var i=ga.test(e);return i||ha.test(e)?xa(e.slice(2),i?2:8):va.test(e)?We:+e}var pe=function(){return $t.Date.now()},ma="Expected a function",ya=Math.max,Ca=Math.min;function Sa(e,n,i){var v,s,x,l,g,m,u=0,d=!1,h=!1,L=!0;if(typeof e!="function")throw new TypeError(ma);n=Ae(n)||0,oe(i)&&(d=!!i.leading,h="maxWait"in i,x=h?ya(Ae(i.maxWait)||0,n):x,L="trailing"in i?!!i.trailing:L);function C(c){var k=v,G=s;return v=s=void 0,u=c,l=e.apply(G,k),l}function S(c){return u=c,g=setTimeout(A,n),d?C(c):l}function T(c){var k=c-m,G=c-u,N=n-k;return h?Ca(N,x-G):N}function $(c){var k=c-m,G=c-u;return m===void 0||k>=n||k<0||h&&G>=x}function A(){var c=pe();if($(c))return W(c);g=setTimeout(A,T(c))}function W(c){return g=void 0,L&&v?C(c):(v=s=void 0,l)}function F(){g!==void 0&&clearTimeout(g),u=0,v=m=s=g=void 0}function I(){return g===void 0?l:W(pe())}function y(){var c=pe(),k=$(c);if(v=arguments,s=this,m=c,k){if(g===void 0)return S(m);if(h)return clearTimeout(g),g=setTimeout(A,n),C(m)}return g===void 0&&(g=setTimeout(A,n)),l}return y.cancel=F,y.flush=I,y}var wa="Expected a function";function Ta(e,n,i){var v=!0,s=!0;if(typeof e!="function")throw new TypeError(wa);return oe(i)&&(v="leading"in i?!!i.leading:v,s="trailing"in i?!!i.trailing:s),Sa(e,n,{leading:v,maxWait:n,trailing:s})}const Pa={tabFontSizeSmall:"14px",tabFontSizeMedium:"14px",tabFontSizeLarge:"16px",tabGapSmallLine:"36px",tabGapMediumLine:"36px",tabGapLargeLine:"36px",tabGapSmallLineVertical:"8px",tabGapMediumLineVertical:"8px",tabGapLargeLineVertical:"8px",tabPaddingSmallLine:"6px 0",tabPaddingMediumLine:"10px 0",tabPaddingLargeLine:"14px 0",tabPaddingVerticalSmallLine:"6px 12px",tabPaddingVerticalMediumLine:"8px 16px",tabPaddingVerticalLargeLine:"10px 20px",tabGapSmallBar:"36px",tabGapMediumBar:"36px",tabGapLargeBar:"36px",tabGapSmallBarVertical:"8px",tabGapMediumBarVertical:"8px",tabGapLargeBarVertical:"8px",tabPaddingSmallBar:"4px 0",tabPaddingMediumBar:"6px 0",tabPaddingLargeBar:"10px 0",tabPaddingVerticalSmallBar:"6px 12px",tabPaddingVerticalMediumBar:"8px 16px",tabPaddingVerticalLargeBar:"10px 20px",tabGapSmallCard:"4px",tabGapMediumCard:"4px",tabGapLargeCard:"4px",tabGapSmallCardVertical:"4px",tabGapMediumCardVertical:"4px",tabGapLargeCardVertical:"4px",tabPaddingSmallCard:"8px 16px",tabPaddingMediumCard:"10px 20px",tabPaddingLargeCard:"12px 24px",tabPaddingSmallSegment:"4px 0",tabPaddingMediumSegment:"6px 0",tabPaddingLargeSegment:"8px 0",tabPaddingVerticalLargeSegment:"0 8px",tabPaddingVerticalSmallCard:"8px 12px",tabPaddingVerticalMediumCard:"10px 16px",tabPaddingVerticalLargeCard:"12px 20px",tabPaddingVerticalSmallSegment:"0 4px",tabPaddingVerticalMediumSegment:"0 6px",tabGapSmallSegment:"0",tabGapMediumSegment:"0",tabGapLargeSegment:"0",tabGapSmallSegmentVertical:"0",tabGapMediumSegmentVertical:"0",tabGapLargeSegmentVertical:"0",panePaddingSmall:"8px 0 0 0",panePaddingMedium:"12px 0 0 0",panePaddingLarge:"16px 0 0 0",closeSize:"18px",closeIconSize:"14px"};function Ra(e){const{textColor2:n,primaryColor:i,textColorDisabled:v,closeIconColor:s,closeIconColorHover:x,closeIconColorPressed:l,closeColorHover:g,closeColorPressed:m,tabColor:u,baseColor:d,dividerColor:h,fontWeight:L,textColor1:C,borderRadius:S,fontSize:T,fontWeightStrong:$}=e;return Object.assign(Object.assign({},Pa),{colorSegment:u,tabFontSizeCard:T,tabTextColorLine:C,tabTextColorActiveLine:i,tabTextColorHoverLine:i,tabTextColorDisabledLine:v,tabTextColorSegment:C,tabTextColorActiveSegment:n,tabTextColorHoverSegment:n,tabTextColorDisabledSegment:v,tabTextColorBar:C,tabTextColorActiveBar:i,tabTextColorHoverBar:i,tabTextColorDisabledBar:v,tabTextColorCard:C,tabTextColorHoverCard:C,tabTextColorActiveCard:i,tabTextColorDisabledCard:v,barColor:i,closeIconColor:s,closeIconColorHover:x,closeIconColorPressed:l,closeColorHover:g,closeColorPressed:m,closeBorderRadius:S,tabColor:u,tabColorSegment:d,tabBorderColor:h,tabFontWeightActive:L,tabFontWeight:L,tabBorderRadius:S,paneTextColor:n,fontWeightStrong:$})}const za={common:Bt,self:Ra},he=_t("n-tabs"),He={tab:[String,Number,Object,Function],name:{type:[String,Number],required:!0},disabled:Boolean,displayDirective:{type:String,default:"if"},closable:{type:Boolean,default:void 0},tabProps:Object,label:[String,Number,Object,Function]},ke=Z({__TAB_PANE__:!0,name:"TabPane",alias:["TabPanel"],props:He,slots:Object,setup(e){const n=Ie(he,null);return n||Wt("tab-pane","`n-tab-pane` must be placed inside `n-tabs`."),{style:n.paneStyleRef,class:n.paneClassRef,mergedClsPrefix:n.mergedClsPrefixRef}},render(){return p("div",{class:[`${this.mergedClsPrefix}-tab-pane`,this.class],style:this.style},this.$slots)}}),La=Object.assign({internalLeftPadded:Boolean,internalAddable:Boolean,internalCreatedByPane:Boolean},It(He,["displayDirective"])),ge=Z({__TAB__:!0,inheritAttrs:!1,name:"Tab",props:La,setup(e){const{mergedClsPrefixRef:n,valueRef:i,typeRef:v,closableRef:s,tabStyleRef:x,addTabStyleRef:l,tabClassRef:g,addTabClassRef:m,tabChangeIdRef:u,onBeforeLeaveRef:d,triggerRef:h,handleAdd:L,activateTab:C,handleClose:S}=Ie(he);return{trigger:h,mergedClosable:Q(()=>{if(e.internalAddable)return!1;const{closable:T}=e;return T===void 0?s.value:T}),style:x,addStyle:l,tabClass:g,addTabClass:m,clsPrefix:n,value:i,type:v,handleClose(T){T.stopPropagation(),!e.disabled&&S(e.name)},activateTab(){if(e.disabled)return;if(e.internalAddable){L();return}const{name:T}=e,$=++u.id;if(T!==i.value){const{value:A}=d;A?Promise.resolve(A(e.name,i.value)).then(W=>{W&&u.id===$&&C(T)}):C(T)}}}},render(){const{internalAddable:e,clsPrefix:n,name:i,disabled:v,label:s,tab:x,value:l,mergedClosable:g,trigger:m,$slots:{default:u}}=this,d=s??x;return p("div",{class:`${n}-tabs-tab-wrapper`},this.internalLeftPadded?p("div",{class:`${n}-tabs-tab-pad`}):null,p("div",Object.assign({key:i,"data-name":i,"data-disabled":v?!0:void 0},At({class:[`${n}-tabs-tab`,l===i&&`${n}-tabs-tab--active`,v&&`${n}-tabs-tab--disabled`,g&&`${n}-tabs-tab--closable`,e&&`${n}-tabs-tab--addable`,e?this.addTabClass:this.tabClass],onClick:m==="click"?this.activateTab:void 0,onMouseenter:m==="hover"?this.activateTab:void 0,style:e?this.addStyle:this.style},this.internalCreatedByPane?this.tabProps||{}:this.$attrs)),p("span",{class:`${n}-tabs-tab__label`},e?p(kt,null,p("div",{class:`${n}-tabs-tab__height-placeholder`}," "),p(Et,{clsPrefix:n},{default:()=>p(sa,null)})):u?u():typeof d=="object"?d:Vt(d??i)),g&&this.type==="card"?p(jt,{clsPrefix:n,class:`${n}-tabs-tab__close`,onClick:this.handleClose,disabled:v}):null))}}),$a=r("tabs",`
 box-sizing: border-box;
 width: 100%;
 display: flex;
 flex-direction: column;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
`,[b("segment-type",[r("tabs-rail",[w("&.transition-disabled",[r("tabs-capsule",`
 transition: none;
 `)])])]),b("top",[r("tab-pane",`
 padding: var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left);
 `)]),b("left",[r("tab-pane",`
 padding: var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left) var(--n-pane-padding-top);
 `)]),b("left, right",`
 flex-direction: row;
 `,[r("tabs-bar",`
 width: 2px;
 right: 0;
 transition:
 top .2s var(--n-bezier),
 max-height .2s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),r("tabs-tab",`
 padding: var(--n-tab-padding-vertical); 
 `)]),b("right",`
 flex-direction: row-reverse;
 `,[r("tab-pane",`
 padding: var(--n-pane-padding-left) var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom);
 `),r("tabs-bar",`
 left: 0;
 `)]),b("bottom",`
 flex-direction: column-reverse;
 justify-content: flex-end;
 `,[r("tab-pane",`
 padding: var(--n-pane-padding-bottom) var(--n-pane-padding-right) var(--n-pane-padding-top) var(--n-pane-padding-left);
 `),r("tabs-bar",`
 top: 0;
 `)]),r("tabs-rail",`
 position: relative;
 padding: 3px;
 border-radius: var(--n-tab-border-radius);
 width: 100%;
 background-color: var(--n-color-segment);
 transition: background-color .3s var(--n-bezier);
 display: flex;
 align-items: center;
 `,[r("tabs-capsule",`
 border-radius: var(--n-tab-border-radius);
 position: absolute;
 pointer-events: none;
 background-color: var(--n-tab-color-segment);
 box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .08);
 transition: transform 0.3s var(--n-bezier);
 `),r("tabs-tab-wrapper",`
 flex-basis: 0;
 flex-grow: 1;
 display: flex;
 align-items: center;
 justify-content: center;
 `,[r("tabs-tab",`
 overflow: hidden;
 border-radius: var(--n-tab-border-radius);
 width: 100%;
 display: flex;
 align-items: center;
 justify-content: center;
 `,[b("active",`
 font-weight: var(--n-font-weight-strong);
 color: var(--n-tab-text-color-active);
 `),w("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])])]),b("flex",[r("tabs-nav",`
 width: 100%;
 position: relative;
 `,[r("tabs-wrapper",`
 width: 100%;
 `,[r("tabs-tab",`
 margin-right: 0;
 `)])])]),r("tabs-nav",`
 box-sizing: border-box;
 line-height: 1.5;
 display: flex;
 transition: border-color .3s var(--n-bezier);
 `,[B("prefix, suffix",`
 display: flex;
 align-items: center;
 `),B("prefix","padding-right: 16px;"),B("suffix","padding-left: 16px;")]),b("top, bottom",[w(">",[r("tabs-nav",[r("tabs-nav-scroll-wrapper",[w("&::before",`
 top: 0;
 bottom: 0;
 left: 0;
 width: 20px;
 `),w("&::after",`
 top: 0;
 bottom: 0;
 right: 0;
 width: 20px;
 `),b("shadow-start",[w("&::before",`
 box-shadow: inset 10px 0 8px -8px rgba(0, 0, 0, .12);
 `)]),b("shadow-end",[w("&::after",`
 box-shadow: inset -10px 0 8px -8px rgba(0, 0, 0, .12);
 `)])])])])]),b("left, right",[r("tabs-nav-scroll-content",`
 flex-direction: column;
 `),w(">",[r("tabs-nav",[r("tabs-nav-scroll-wrapper",[w("&::before",`
 top: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),w("&::after",`
 bottom: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),b("shadow-start",[w("&::before",`
 box-shadow: inset 0 10px 8px -8px rgba(0, 0, 0, .12);
 `)]),b("shadow-end",[w("&::after",`
 box-shadow: inset 0 -10px 8px -8px rgba(0, 0, 0, .12);
 `)])])])])]),r("tabs-nav-scroll-wrapper",`
 flex: 1;
 position: relative;
 overflow: hidden;
 `,[r("tabs-nav-y-scroll",`
 height: 100%;
 width: 100%;
 overflow-y: auto; 
 scrollbar-width: none;
 `,[w("&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb",`
 width: 0;
 height: 0;
 display: none;
 `)]),w("&::before, &::after",`
 transition: box-shadow .3s var(--n-bezier);
 pointer-events: none;
 content: "";
 position: absolute;
 z-index: 1;
 `)]),r("tabs-nav-scroll-content",`
 display: flex;
 position: relative;
 min-width: 100%;
 min-height: 100%;
 width: fit-content;
 box-sizing: border-box;
 `),r("tabs-wrapper",`
 display: inline-flex;
 flex-wrap: nowrap;
 position: relative;
 `),r("tabs-tab-wrapper",`
 display: flex;
 flex-wrap: nowrap;
 flex-shrink: 0;
 flex-grow: 0;
 `),r("tabs-tab",`
 cursor: pointer;
 white-space: nowrap;
 flex-wrap: nowrap;
 display: inline-flex;
 align-items: center;
 color: var(--n-tab-text-color);
 font-size: var(--n-tab-font-size);
 background-clip: padding-box;
 padding: var(--n-tab-padding);
 transition:
 box-shadow .3s var(--n-bezier),
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[b("disabled",{cursor:"not-allowed"}),B("close",`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),B("label",`
 display: flex;
 align-items: center;
 z-index: 1;
 `)]),r("tabs-bar",`
 position: absolute;
 bottom: 0;
 height: 2px;
 border-radius: 1px;
 background-color: var(--n-bar-color);
 transition:
 left .2s var(--n-bezier),
 max-width .2s var(--n-bezier),
 opacity .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `,[w("&.transition-disabled",`
 transition: none;
 `),b("disabled",`
 background-color: var(--n-tab-text-color-disabled)
 `)]),r("tabs-pane-wrapper",`
 position: relative;
 overflow: hidden;
 transition: max-height .2s var(--n-bezier);
 `),r("tab-pane",`
 color: var(--n-pane-text-color);
 width: 100%;
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 opacity .2s var(--n-bezier);
 left: 0;
 right: 0;
 top: 0;
 `,[w("&.next-transition-leave-active, &.prev-transition-leave-active, &.next-transition-enter-active, &.prev-transition-enter-active",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 transform .2s var(--n-bezier),
 opacity .2s var(--n-bezier);
 `),w("&.next-transition-leave-active, &.prev-transition-leave-active",`
 position: absolute;
 `),w("&.next-transition-enter-from, &.prev-transition-leave-to",`
 transform: translateX(32px);
 opacity: 0;
 `),w("&.next-transition-leave-to, &.prev-transition-enter-from",`
 transform: translateX(-32px);
 opacity: 0;
 `),w("&.next-transition-leave-from, &.next-transition-enter-to, &.prev-transition-leave-from, &.prev-transition-enter-to",`
 transform: translateX(0);
 opacity: 1;
 `)]),r("tabs-tab-pad",`
 box-sizing: border-box;
 width: var(--n-tab-gap);
 flex-grow: 0;
 flex-shrink: 0;
 `),b("line-type, bar-type",[r("tabs-tab",`
 font-weight: var(--n-tab-font-weight);
 box-sizing: border-box;
 vertical-align: bottom;
 `,[w("&:hover",{color:"var(--n-tab-text-color-hover)"}),b("active",`
 color: var(--n-tab-text-color-active);
 font-weight: var(--n-tab-font-weight-active);
 `),b("disabled",{color:"var(--n-tab-text-color-disabled)"})])]),r("tabs-nav",[b("line-type",[b("top",[B("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 bottom: -1px;
 `)]),b("left",[B("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 right: -1px;
 `)]),b("right",[B("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 left: -1px;
 `)]),b("bottom",[B("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 top: -1px;
 `)]),B("prefix, suffix",`
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-nav-scroll-content",`
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-bar",`
 border-radius: 0;
 `)]),b("card-type",[B("prefix, suffix",`
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-pad",`
 flex-grow: 1;
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-tab-pad",`
 transition: border-color .3s var(--n-bezier);
 `),r("tabs-tab",`
 font-weight: var(--n-tab-font-weight);
 border: 1px solid var(--n-tab-border-color);
 background-color: var(--n-tab-color);
 box-sizing: border-box;
 position: relative;
 vertical-align: bottom;
 display: flex;
 justify-content: space-between;
 font-size: var(--n-tab-font-size);
 color: var(--n-tab-text-color);
 `,[b("addable",`
 padding-left: 8px;
 padding-right: 8px;
 font-size: 16px;
 justify-content: center;
 `,[B("height-placeholder",`
 width: 0;
 font-size: var(--n-tab-font-size);
 `),Mt("disabled",[w("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])]),b("closable","padding-right: 8px;"),b("active",`
 background-color: #0000;
 font-weight: var(--n-tab-font-weight-active);
 color: var(--n-tab-text-color-active);
 `),b("disabled","color: var(--n-tab-text-color-disabled);")])]),b("left, right",`
 flex-direction: column; 
 `,[B("prefix, suffix",`
 padding: var(--n-tab-padding-vertical);
 `),r("tabs-wrapper",`
 flex-direction: column;
 `),r("tabs-tab-wrapper",`
 flex-direction: column;
 `,[r("tabs-tab-pad",`
 height: var(--n-tab-gap-vertical);
 width: 100%;
 `)])]),b("top",[b("card-type",[r("tabs-scroll-padding","border-bottom: 1px solid var(--n-tab-border-color);"),B("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-top-right-radius: var(--n-tab-border-radius);
 `,[b("active",`
 border-bottom: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `)])]),b("left",[b("card-type",[r("tabs-scroll-padding","border-right: 1px solid var(--n-tab-border-color);"),B("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-bottom-left-radius: var(--n-tab-border-radius);
 `,[b("active",`
 border-right: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `)])]),b("right",[b("card-type",[r("tabs-scroll-padding","border-left: 1px solid var(--n-tab-border-color);"),B("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-right-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[b("active",`
 border-left: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `)])]),b("bottom",[b("card-type",[r("tabs-scroll-padding","border-top: 1px solid var(--n-tab-border-color);"),B("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-bottom-left-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[b("active",`
 border-top: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `)])])])]),ue=Ta,Ba=Object.assign(Object.assign({},Me.props),{value:[String,Number],defaultValue:[String,Number],trigger:{type:String,default:"click"},type:{type:String,default:"bar"},closable:Boolean,justifyContent:String,size:String,placement:{type:String,default:"top"},tabStyle:[String,Object],tabClass:String,addTabStyle:[String,Object],addTabClass:String,barWidth:Number,paneClass:String,paneStyle:[String,Object],paneWrapperClass:String,paneWrapperStyle:[String,Object],addable:[Boolean,Object],tabsPadding:{type:Number,default:0},animated:Boolean,onBeforeLeave:Function,onAdd:Function,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onClose:[Function,Array],labelSize:String,activeName:[String,Number],onActiveNameChange:[Function,Array]}),_a=Z({name:"Tabs",props:Ba,slots:Object,setup(e,{slots:n}){var i,v,s,x;const{mergedClsPrefixRef:l,inlineThemeDisabled:g,mergedComponentPropsRef:m}=Ht(e),u=Me("Tabs","-tabs",$a,za,e,l),d=_(null),h=_(null),L=_(null),C=_(null),S=_(null),T=_(null),$=_(!0),A=_(!0),W=ze(e,["labelSize","size"]),F=Q(()=>{var t,a;if(W.value)return W.value;const o=(a=(t=m==null?void 0:m.value)===null||t===void 0?void 0:t.Tabs)===null||a===void 0?void 0:a.size;return o||"medium"}),I=ze(e,["activeName","value"]),y=_((v=(i=I.value)!==null&&i!==void 0?i:e.defaultValue)!==null&&v!==void 0?v:n.default?(x=(s=de(n.default())[0])===null||s===void 0?void 0:s.props)===null||x===void 0?void 0:x.name:null),c=Ft(I,y),k={id:0},G=Q(()=>{if(!(!e.justifyContent||e.type==="card"))return{display:"flex",justifyContent:e.justifyContent}});ce(c,()=>{k.id=0,ee(),me()});function N(){var t;const{value:a}=c;return a===null?null:(t=d.value)===null||t===void 0?void 0:t.querySelector(`[data-name="${a}"]`)}function Fe(t){if(e.type==="card")return;const{value:a}=h;if(!a)return;const o=a.style.opacity==="0";if(t){const f=`${l.value}-tabs-bar--disabled`,{barWidth:P,placement:V}=e;if(t.dataset.disabled==="true"?a.classList.add(f):a.classList.remove(f),["top","bottom"].includes(V)){if(xe(["top","maxHeight","height"]),typeof P=="number"&&t.offsetWidth>=P){const j=Math.floor((t.offsetWidth-P)/2)+t.offsetLeft;a.style.left=`${j}px`,a.style.maxWidth=`${P}px`}else a.style.left=`${t.offsetLeft}px`,a.style.maxWidth=`${t.offsetWidth}px`;a.style.width="8192px",o&&(a.style.transition="none"),a.offsetWidth,o&&(a.style.transition="",a.style.opacity="1")}else{if(xe(["left","maxWidth","width"]),typeof P=="number"&&t.offsetHeight>=P){const j=Math.floor((t.offsetHeight-P)/2)+t.offsetTop;a.style.top=`${j}px`,a.style.maxHeight=`${P}px`}else a.style.top=`${t.offsetTop}px`,a.style.maxHeight=`${t.offsetHeight}px`;a.style.height="8192px",o&&(a.style.transition="none"),a.offsetHeight,o&&(a.style.transition="",a.style.opacity="1")}}}function Oe(){if(e.type==="card")return;const{value:t}=h;t&&(t.style.opacity="0")}function xe(t){const{value:a}=h;if(a)for(const o of t)a.style[o]=""}function ee(){if(e.type==="card")return;const t=N();t?Fe(t):Oe()}function me(){var t;const a=(t=S.value)===null||t===void 0?void 0:t.$el;if(!a)return;const o=N();if(!o)return;const{scrollLeft:f,offsetWidth:P}=a,{offsetLeft:V,offsetWidth:j}=o;f>V?a.scrollTo({top:0,left:V,behavior:"smooth"}):V+j>f+P&&a.scrollTo({top:0,left:V+j-P,behavior:"smooth"})}const te=_(null);let ie=0,O=null;function Ge(t){const a=te.value;if(a){ie=t.getBoundingClientRect().height;const o=`${ie}px`,f=()=>{a.style.height=o,a.style.maxHeight=o};O?(f(),O(),O=null):O=f}}function Ne(t){const a=te.value;if(a){const o=t.getBoundingClientRect().height,f=()=>{document.body.offsetHeight,a.style.maxHeight=`${o}px`,a.style.height=`${Math.max(ie,o)}px`};O?(O(),O=null,f()):O=f}}function De(){const t=te.value;if(t){t.style.maxHeight="",t.style.height="";const{paneWrapperStyle:a}=e;if(typeof a=="string")t.style.cssText=a;else if(a){const{maxHeight:o,height:f}=a;o!==void 0&&(t.style.maxHeight=o),f!==void 0&&(t.style.height=f)}}}const ye={value:[]},Ce=_("next");function Ue(t){const a=c.value;let o="next";for(const f of ye.value){if(f===a)break;if(f===t){o="prev";break}}Ce.value=o,Xe(t)}function Xe(t){const{onActiveNameChange:a,onUpdateValue:o,"onUpdate:value":f}=e;a&&ne(a,t),o&&ne(o,t),f&&ne(f,t),y.value=t}function Ke(t){const{onClose:a}=e;a&&ne(a,t)}function Se(){const{value:t}=h;if(!t)return;const a="transition-disabled";t.classList.add(a),ee(),t.classList.remove(a)}const D=_(null);function le({transitionDisabled:t}){const a=d.value;if(!a)return;t&&a.classList.add("transition-disabled");const o=N();o&&D.value&&(D.value.style.width=`${o.offsetWidth}px`,D.value.style.height=`${o.offsetHeight}px`,D.value.style.transform=`translateX(${o.offsetLeft-Ut(getComputedStyle(a).paddingLeft)}px)`,t&&D.value.offsetWidth),t&&a.classList.remove("transition-disabled")}ce([c],()=>{e.type==="segment"&&fe(()=>{le({transitionDisabled:!1})})}),Ot(()=>{e.type==="segment"&&le({transitionDisabled:!0})});let we=0;function qe(t){var a;if(t.contentRect.width===0&&t.contentRect.height===0||we===t.contentRect.width)return;we=t.contentRect.width;const{type:o}=e;if((o==="line"||o==="bar")&&Se(),o!=="segment"){const{placement:f}=e;se((f==="top"||f==="bottom"?(a=S.value)===null||a===void 0?void 0:a.$el:T.value)||null)}}const Ye=ue(qe,64);ce([()=>e.justifyContent,()=>e.size],()=>{fe(()=>{const{type:t}=e;(t==="line"||t==="bar")&&Se()})});const U=_(!1);function Je(t){var a;const{target:o,contentRect:{width:f,height:P}}=t,V=o.parentElement.parentElement.offsetWidth,j=o.parentElement.parentElement.offsetHeight,{placement:K}=e;if(!U.value)K==="top"||K==="bottom"?V<f&&(U.value=!0):j<P&&(U.value=!0);else{const{value:q}=C;if(!q)return;K==="top"||K==="bottom"?V-f>q.$el.offsetWidth&&(U.value=!1):j-P>q.$el.offsetHeight&&(U.value=!1)}se(((a=S.value)===null||a===void 0?void 0:a.$el)||null)}const Qe=ue(Je,64);function Ze(){const{onAdd:t}=e;t&&t(),fe(()=>{const a=N(),{value:o}=S;!a||!o||o.scrollTo({left:a.offsetLeft,top:0,behavior:"smooth"})})}function se(t){if(!t)return;const{placement:a}=e;if(a==="top"||a==="bottom"){const{scrollLeft:o,scrollWidth:f,offsetWidth:P}=t;$.value=o<=0,A.value=o+P>=f}else{const{scrollTop:o,scrollHeight:f,offsetHeight:P}=t;$.value=o<=0,A.value=o+P>=f}}const et=ue(t=>{se(t.target)},64);Jt(he,{triggerRef:H(e,"trigger"),tabStyleRef:H(e,"tabStyle"),tabClassRef:H(e,"tabClass"),addTabStyleRef:H(e,"addTabStyle"),addTabClassRef:H(e,"addTabClass"),paneClassRef:H(e,"paneClass"),paneStyleRef:H(e,"paneStyle"),mergedClsPrefixRef:l,typeRef:H(e,"type"),closableRef:H(e,"closable"),valueRef:c,tabChangeIdRef:k,onBeforeLeaveRef:H(e,"onBeforeLeave"),activateTab:Ue,handleClose:Ke,handleAdd:Ze}),Gt(()=>{ee(),me()}),Nt(()=>{const{value:t}=L;if(!t)return;const{value:a}=l,o=`${a}-tabs-nav-scroll-wrapper--shadow-start`,f=`${a}-tabs-nav-scroll-wrapper--shadow-end`;$.value?t.classList.remove(o):t.classList.add(o),A.value?t.classList.remove(f):t.classList.add(f)});const tt={syncBarPosition:()=>{ee()}},at=()=>{le({transitionDisabled:!0})},Te=Q(()=>{const{value:t}=F,{type:a}=e,o={card:"Card",bar:"Bar",line:"Line",segment:"Segment"}[a],f=`${t}${o}`,{self:{barColor:P,closeIconColor:V,closeIconColorHover:j,closeIconColorPressed:K,tabColor:q,tabBorderColor:rt,paneTextColor:nt,tabFontWeight:ot,tabBorderRadius:it,tabFontWeightActive:lt,colorSegment:st,fontWeightStrong:dt,tabColorSegment:bt,closeSize:ct,closeIconSize:ft,closeColorHover:pt,closeColorPressed:ut,closeBorderRadius:vt,[M("panePadding",t)]:ae,[M("tabPadding",f)]:gt,[M("tabPaddingVertical",f)]:ht,[M("tabGap",f)]:xt,[M("tabGap",`${f}Vertical`)]:mt,[M("tabTextColor",a)]:yt,[M("tabTextColorActive",a)]:Ct,[M("tabTextColorHover",a)]:St,[M("tabTextColorDisabled",a)]:wt,[M("tabFontSize",t)]:Tt},common:{cubicBezierEaseInOut:Pt}}=u.value;return{"--n-bezier":Pt,"--n-color-segment":st,"--n-bar-color":P,"--n-tab-font-size":Tt,"--n-tab-text-color":yt,"--n-tab-text-color-active":Ct,"--n-tab-text-color-disabled":wt,"--n-tab-text-color-hover":St,"--n-pane-text-color":nt,"--n-tab-border-color":rt,"--n-tab-border-radius":it,"--n-close-size":ct,"--n-close-icon-size":ft,"--n-close-color-hover":pt,"--n-close-color-pressed":ut,"--n-close-border-radius":vt,"--n-close-icon-color":V,"--n-close-icon-color-hover":j,"--n-close-icon-color-pressed":K,"--n-tab-color":q,"--n-tab-font-weight":ot,"--n-tab-font-weight-active":lt,"--n-tab-padding":gt,"--n-tab-padding-vertical":ht,"--n-tab-gap":xt,"--n-tab-gap-vertical":mt,"--n-pane-padding-left":re(ae,"left"),"--n-pane-padding-right":re(ae,"right"),"--n-pane-padding-top":re(ae,"top"),"--n-pane-padding-bottom":re(ae,"bottom"),"--n-font-weight-strong":dt,"--n-tab-color-segment":bt}}),X=g?Dt("tabs",Q(()=>`${F.value[0]}${e.type[0]}`),Te,e):void 0;return Object.assign({mergedClsPrefix:l,mergedValue:c,renderedNames:new Set,segmentCapsuleElRef:D,tabsPaneWrapperRef:te,tabsElRef:d,barElRef:h,addTabInstRef:C,xScrollInstRef:S,scrollWrapperElRef:L,addTabFixed:U,tabWrapperStyle:G,handleNavResize:Ye,mergedSize:F,handleScroll:et,handleTabsResize:Qe,cssVars:g?void 0:Te,themeClass:X==null?void 0:X.themeClass,animationDirection:Ce,renderNameListRef:ye,yScrollElRef:T,handleSegmentResize:at,onAnimationBeforeLeave:Ge,onAnimationEnter:Ne,onAnimationAfterEnter:De,onRender:X==null?void 0:X.onRender},tt)},render(){const{mergedClsPrefix:e,type:n,placement:i,addTabFixed:v,addable:s,mergedSize:x,renderNameListRef:l,onRender:g,paneWrapperClass:m,paneWrapperStyle:u,$slots:{default:d,prefix:h,suffix:L}}=this;g==null||g();const C=d?de(d()).filter(y=>y.type.__TAB_PANE__===!0):[],S=d?de(d()).filter(y=>y.type.__TAB__===!0):[],T=!S.length,$=n==="card",A=n==="segment",W=!$&&!A&&this.justifyContent;l.value=[];const F=()=>{const y=p("div",{style:this.tabWrapperStyle,class:`${e}-tabs-wrapper`},W?null:p("div",{class:`${e}-tabs-scroll-padding`,style:i==="top"||i==="bottom"?{width:`${this.tabsPadding}px`}:{height:`${this.tabsPadding}px`}}),T?C.map((c,k)=>(l.value.push(c.props.name),ve(p(ge,Object.assign({},c.props,{internalCreatedByPane:!0,internalLeftPadded:k!==0&&(!W||W==="center"||W==="start"||W==="end")}),c.children?{default:c.children.tab}:void 0)))):S.map((c,k)=>(l.value.push(c.props.name),ve(k!==0&&!W?je(c):c))),!v&&s&&$?Ve(s,(T?C.length:S.length)!==0):null,W?null:p("div",{class:`${e}-tabs-scroll-padding`,style:{width:`${this.tabsPadding}px`}}));return p("div",{ref:"tabsElRef",class:`${e}-tabs-nav-scroll-content`},$&&s?p(be,{onResize:this.handleTabsResize},{default:()=>y}):y,$?p("div",{class:`${e}-tabs-pad`}):null,$?null:p("div",{ref:"barElRef",class:`${e}-tabs-bar`}))},I=A?"top":i;return p("div",{class:[`${e}-tabs`,this.themeClass,`${e}-tabs--${n}-type`,`${e}-tabs--${x}-size`,W&&`${e}-tabs--flex`,`${e}-tabs--${I}`],style:this.cssVars},p("div",{class:[`${e}-tabs-nav--${n}-type`,`${e}-tabs-nav--${I}`,`${e}-tabs-nav`]},Re(h,y=>y&&p("div",{class:`${e}-tabs-nav__prefix`},y)),A?p(be,{onResize:this.handleSegmentResize},{default:()=>p("div",{class:`${e}-tabs-rail`,ref:"tabsElRef"},p("div",{class:`${e}-tabs-capsule`,ref:"segmentCapsuleElRef"},p("div",{class:`${e}-tabs-wrapper`},p("div",{class:`${e}-tabs-tab`}))),T?C.map((y,c)=>(l.value.push(y.props.name),p(ge,Object.assign({},y.props,{internalCreatedByPane:!0,internalLeftPadded:c!==0}),y.children?{default:y.children.tab}:void 0))):S.map((y,c)=>(l.value.push(y.props.name),c===0?y:je(y))))}):p(be,{onResize:this.handleNavResize},{default:()=>p("div",{class:`${e}-tabs-nav-scroll-wrapper`,ref:"scrollWrapperElRef"},["top","bottom"].includes(I)?p(ba,{ref:"xScrollInstRef",onScroll:this.handleScroll},{default:F}):p("div",{class:`${e}-tabs-nav-y-scroll`,onScroll:this.handleScroll,ref:"yScrollElRef"},F()))}),v&&s&&$?Ve(s,!0):null,Re(L,y=>y&&p("div",{class:`${e}-tabs-nav__suffix`},y))),T&&(this.animated&&(I==="top"||I==="bottom")?p("div",{ref:"tabsPaneWrapperRef",style:u,class:[`${e}-tabs-pane-wrapper`,m]},Ee(C,this.mergedValue,this.renderedNames,this.onAnimationBeforeLeave,this.onAnimationEnter,this.onAnimationAfterEnter,this.animationDirection)):Ee(C,this.mergedValue,this.renderedNames)))}});function Ee(e,n,i,v,s,x,l){const g=[];return e.forEach(m=>{const{name:u,displayDirective:d,"display-directive":h}=m.props,L=S=>d===S||h===S,C=n===u;if(m.key!==void 0&&(m.key=u),C||L("show")||L("show:lazy")&&i.has(u)){i.has(u)||i.add(u);const S=!L("if");g.push(S?Xt(m,[[Kt,C]]):m)}}),l?p(qt,{name:`${l}-transition`,onBeforeLeave:v,onEnter:s,onAfterEnter:x},{default:()=>g}):g}function Ve(e,n){return p(ge,{ref:"addTabInstRef",key:"__addable",name:"__addable",internalCreatedByPane:!0,internalAddable:!0,internalLeftPadded:n,disabled:typeof e=="object"&&e.disabled})}function je(e){const n=Yt(e);return n.props?n.props.internalLeftPadded=!0:n.props={internalLeftPadded:!0},n}function ve(e){return Array.isArray(e.dynamicProps)?e.dynamicProps.includes("internalLeftPadded")||e.dynamicProps.push("internalLeftPadded"):e.dynamicProps=["internalLeftPadded"],e}const Wa={class:"login-bg"},Aa=Z({__name:"LoginView",setup(e){const n=ta(),i=la(),v=Qt(),s=_("login"),x=_(!1),l=_({username:"admin_user",password:"admin_user",user:"",pwd:"",confirm:""});async function g(){x.value=!0;try{const u=await oa(l.value.username,l.value.password);v.setAuth(u.token,u.name,u.role_id,u.requirePasswordChange),i.success("登录成功"),n.push(u.requirePasswordChange?"/profile":"/dashboard")}catch(u){i.error(u.message)}finally{x.value=!1}}async function m(){if(l.value.pwd!==l.value.confirm){i.error("两次密码不一致");return}x.value=!0;try{await ia(l.value.user,l.value.pwd),i.success("注册成功,请登录"),s.value="login"}catch(u){i.error(u.message)}finally{x.value=!1}}return(u,d)=>(aa(),Zt("div",Wa,[R(z(ea),{class:"login-card",bordered:!1},{default:E(()=>[d[8]||(d[8]=ra("h1",{class:"brand"},"FLVX2",-1)),R(z(_a),{value:s.value,"onUpdate:value":d[5]||(d[5]=h=>s.value=h),type:"line",animated:""},{default:E(()=>[R(z(ke),{name:"login",tab:"登录"},{default:E(()=>[R(z(_e),{onKeyup:Le(g,["enter"])},{default:E(()=>[R(z(J),{label:"用户名"},{default:E(()=>[R(z(Y),{value:l.value.username,"onUpdate:value":d[0]||(d[0]=h=>l.value.username=h),placeholder:"用户名"},null,8,["value"])]),_:1}),R(z(J),{label:"密码"},{default:E(()=>[R(z(Y),{value:l.value.password,"onUpdate:value":d[1]||(d[1]=h=>l.value.password=h),type:"password","show-password-on":"click",placeholder:"密码"},null,8,["value"])]),_:1}),R(z($e),{type:"primary",block:"",loading:x.value,onClick:g},{default:E(()=>[...d[6]||(d[6]=[Be("登 录",-1)])]),_:1},8,["loading"])]),_:1})]),_:1}),R(z(ke),{name:"register",tab:"注册"},{default:E(()=>[R(z(_e),{onKeyup:Le(m,["enter"])},{default:E(()=>[R(z(J),{label:"用户名"},{default:E(()=>[R(z(Y),{value:l.value.user,"onUpdate:value":d[2]||(d[2]=h=>l.value.user=h),placeholder:"3-20 位"},null,8,["value"])]),_:1}),R(z(J),{label:"密码"},{default:E(()=>[R(z(Y),{value:l.value.pwd,"onUpdate:value":d[3]||(d[3]=h=>l.value.pwd=h),type:"password","show-password-on":"click",placeholder:"6-32 位"},null,8,["value"])]),_:1}),R(z(J),{label:"确认密码"},{default:E(()=>[R(z(Y),{value:l.value.confirm,"onUpdate:value":d[4]||(d[4]=h=>l.value.confirm=h),type:"password","show-password-on":"click"},null,8,["value"])]),_:1}),R(z($e),{type:"primary",block:"",loading:x.value,onClick:m},{default:E(()=>[...d[7]||(d[7]=[Be("注 册",-1)])]),_:1},8,["loading"])]),_:1})]),_:1})]),_:1},8,["value"])]),_:1})]))}}),Ia=na(Aa,[["__scopeId","data-v-e45ba2a8"]]);export{Ia as default};
