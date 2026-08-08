import{d as Z,h as p,u as Rt,c as zt,r as _,a as Pe,i as Lt,b as oe,e as $t,f as Bt,g as _t,j as Ie,t as Wt,m as At,F as kt,N as Et,k as Vt,l as jt,n as Q,o as It,p as r,q as d,s as T,v as B,w as Mt,x as de,y as Re,V as be,z as Ht,A as Me,B as Ot,C as ce,D as Ft,E as Gt,G as Nt,H as Dt,I as ze,J as Ut,K as Xt,L as Kt,T as qt,M as Yt,O as fe,P as M,Q as re,R as Jt,S as H,U as ne,W as Qt,X as Zt,Y as ea,Z as L,_ as E,$ as z,a0 as ta,a1 as aa,a2 as ra,a3 as na,a4 as oa,a5 as Le,a6 as $e,a7 as Be,a8 as ia,a9 as la,aa as sa}from"./index-nWgbspkO.js";import{u as da,N as _e,a as Y}from"./FormItem-DP2e30S_.js";import{A as ba}from"./Add-B2IlRCGM.js";import{N as J}from"./Input-Ck5InYTz.js";const ca=Pe(".v-x-scroll",{overflow:"auto",scrollbarWidth:"none"},[Pe("&::-webkit-scrollbar",{width:0,height:0})]),fa=Z({name:"XScroll",props:{disabled:Boolean,onScroll:Function},setup(){const e=_(null);function n(s){!(s.currentTarget.offsetWidth<s.currentTarget.scrollWidth)||s.deltaY===0||(s.currentTarget.scrollLeft+=s.deltaY+s.deltaX,s.preventDefault())}const i=Rt();return ca.mount({id:"vueuc/x-scroll",head:!0,anchorMetaName:zt,ssr:i}),Object.assign({selfRef:e,handleWheel:n},{scrollTo(...s){var y;(y=e.value)===null||y===void 0||y.scrollTo(...s)}})},render(){return p("div",{ref:"selfRef",onScroll:this.onScroll,onWheel:this.disabled?void 0:this.handleWheel,class:"v-x-scroll"},this.$slots)}});var pa=/\s/;function ua(e){for(var n=e.length;n--&&pa.test(e.charAt(n)););return n}var va=/^\s+/;function ga(e){return e&&e.slice(0,ua(e)+1).replace(va,"")}var We=NaN,ha=/^[-+]0x[0-9a-f]+$/i,xa=/^0b[01]+$/i,ma=/^0o[0-7]+$/i,ya=parseInt;function Ae(e){if(typeof e=="number")return e;if(Lt(e))return We;if(oe(e)){var n=typeof e.valueOf=="function"?e.valueOf():e;e=oe(n)?n+"":n}if(typeof e!="string")return e===0?e:+e;e=ga(e);var i=xa.test(e);return i||ma.test(e)?ya(e.slice(2),i?2:8):ha.test(e)?We:+e}var pe=function(){return $t.Date.now()},Ca="Expected a function",Sa=Math.max,wa=Math.min;function Ta(e,n,i){var g,s,y,u,l,h,C=0,v=!1,f=!1,x=!0;if(typeof e!="function")throw new TypeError(Ca);n=Ae(n)||0,oe(i)&&(v=!!i.leading,f="maxWait"in i,y=f?Sa(Ae(i.maxWait)||0,n):y,x="trailing"in i?!!i.trailing:x);function S(b){var k=g,G=s;return g=s=void 0,C=b,u=e.apply(G,k),u}function w(b){return C=b,l=setTimeout(A,n),v?S(b):u}function P(b){var k=b-h,G=b-C,N=n-k;return f?wa(N,y-G):N}function $(b){var k=b-h,G=b-C;return h===void 0||k>=n||k<0||f&&G>=y}function A(){var b=pe();if($(b))return W(b);l=setTimeout(A,P(b))}function W(b){return l=void 0,x&&g?S(b):(g=s=void 0,u)}function O(){l!==void 0&&clearTimeout(l),C=0,g=h=s=l=void 0}function I(){return l===void 0?u:W(pe())}function m(){var b=pe(),k=$(b);if(g=arguments,s=this,h=b,k){if(l===void 0)return w(h);if(f)return clearTimeout(l),l=setTimeout(A,n),S(h)}return l===void 0&&(l=setTimeout(A,n)),u}return m.cancel=O,m.flush=I,m}var Pa="Expected a function";function Ra(e,n,i){var g=!0,s=!0;if(typeof e!="function")throw new TypeError(Pa);return oe(i)&&(g="leading"in i?!!i.leading:g,s="trailing"in i?!!i.trailing:s),Ta(e,n,{leading:g,maxWait:n,trailing:s})}const za={tabFontSizeSmall:"14px",tabFontSizeMedium:"14px",tabFontSizeLarge:"16px",tabGapSmallLine:"36px",tabGapMediumLine:"36px",tabGapLargeLine:"36px",tabGapSmallLineVertical:"8px",tabGapMediumLineVertical:"8px",tabGapLargeLineVertical:"8px",tabPaddingSmallLine:"6px 0",tabPaddingMediumLine:"10px 0",tabPaddingLargeLine:"14px 0",tabPaddingVerticalSmallLine:"6px 12px",tabPaddingVerticalMediumLine:"8px 16px",tabPaddingVerticalLargeLine:"10px 20px",tabGapSmallBar:"36px",tabGapMediumBar:"36px",tabGapLargeBar:"36px",tabGapSmallBarVertical:"8px",tabGapMediumBarVertical:"8px",tabGapLargeBarVertical:"8px",tabPaddingSmallBar:"4px 0",tabPaddingMediumBar:"6px 0",tabPaddingLargeBar:"10px 0",tabPaddingVerticalSmallBar:"6px 12px",tabPaddingVerticalMediumBar:"8px 16px",tabPaddingVerticalLargeBar:"10px 20px",tabGapSmallCard:"4px",tabGapMediumCard:"4px",tabGapLargeCard:"4px",tabGapSmallCardVertical:"4px",tabGapMediumCardVertical:"4px",tabGapLargeCardVertical:"4px",tabPaddingSmallCard:"8px 16px",tabPaddingMediumCard:"10px 20px",tabPaddingLargeCard:"12px 24px",tabPaddingSmallSegment:"4px 0",tabPaddingMediumSegment:"6px 0",tabPaddingLargeSegment:"8px 0",tabPaddingVerticalLargeSegment:"0 8px",tabPaddingVerticalSmallCard:"8px 12px",tabPaddingVerticalMediumCard:"10px 16px",tabPaddingVerticalLargeCard:"12px 20px",tabPaddingVerticalSmallSegment:"0 4px",tabPaddingVerticalMediumSegment:"0 6px",tabGapSmallSegment:"0",tabGapMediumSegment:"0",tabGapLargeSegment:"0",tabGapSmallSegmentVertical:"0",tabGapMediumSegmentVertical:"0",tabGapLargeSegmentVertical:"0",panePaddingSmall:"8px 0 0 0",panePaddingMedium:"12px 0 0 0",panePaddingLarge:"16px 0 0 0",closeSize:"18px",closeIconSize:"14px"};function La(e){const{textColor2:n,primaryColor:i,textColorDisabled:g,closeIconColor:s,closeIconColorHover:y,closeIconColorPressed:u,closeColorHover:l,closeColorPressed:h,tabColor:C,baseColor:v,dividerColor:f,fontWeight:x,textColor1:S,borderRadius:w,fontSize:P,fontWeightStrong:$}=e;return Object.assign(Object.assign({},za),{colorSegment:C,tabFontSizeCard:P,tabTextColorLine:S,tabTextColorActiveLine:i,tabTextColorHoverLine:i,tabTextColorDisabledLine:g,tabTextColorSegment:S,tabTextColorActiveSegment:n,tabTextColorHoverSegment:n,tabTextColorDisabledSegment:g,tabTextColorBar:S,tabTextColorActiveBar:i,tabTextColorHoverBar:i,tabTextColorDisabledBar:g,tabTextColorCard:S,tabTextColorHoverCard:S,tabTextColorActiveCard:i,tabTextColorDisabledCard:g,barColor:i,closeIconColor:s,closeIconColorHover:y,closeIconColorPressed:u,closeColorHover:l,closeColorPressed:h,closeBorderRadius:w,tabColor:C,tabColorSegment:v,tabBorderColor:f,tabFontWeightActive:x,tabFontWeight:x,tabBorderRadius:w,paneTextColor:n,fontWeightStrong:$})}const $a={common:Bt,self:La},he=_t("n-tabs"),He={tab:[String,Number,Object,Function],name:{type:[String,Number],required:!0},disabled:Boolean,displayDirective:{type:String,default:"if"},closable:{type:Boolean,default:void 0},tabProps:Object,label:[String,Number,Object,Function]},ke=Z({__TAB_PANE__:!0,name:"TabPane",alias:["TabPanel"],props:He,slots:Object,setup(e){const n=Ie(he,null);return n||Wt("tab-pane","`n-tab-pane` must be placed inside `n-tabs`."),{style:n.paneStyleRef,class:n.paneClassRef,mergedClsPrefix:n.mergedClsPrefixRef}},render(){return p("div",{class:[`${this.mergedClsPrefix}-tab-pane`,this.class],style:this.style},this.$slots)}}),Ba=Object.assign({internalLeftPadded:Boolean,internalAddable:Boolean,internalCreatedByPane:Boolean},It(He,["displayDirective"])),ge=Z({__TAB__:!0,inheritAttrs:!1,name:"Tab",props:Ba,setup(e){const{mergedClsPrefixRef:n,valueRef:i,typeRef:g,closableRef:s,tabStyleRef:y,addTabStyleRef:u,tabClassRef:l,addTabClassRef:h,tabChangeIdRef:C,onBeforeLeaveRef:v,triggerRef:f,handleAdd:x,activateTab:S,handleClose:w}=Ie(he);return{trigger:f,mergedClosable:Q(()=>{if(e.internalAddable)return!1;const{closable:P}=e;return P===void 0?s.value:P}),style:y,addStyle:u,tabClass:l,addTabClass:h,clsPrefix:n,value:i,type:g,handleClose(P){P.stopPropagation(),!e.disabled&&w(e.name)},activateTab(){if(e.disabled)return;if(e.internalAddable){x();return}const{name:P}=e,$=++C.id;if(P!==i.value){const{value:A}=v;A?Promise.resolve(A(e.name,i.value)).then(W=>{W&&C.id===$&&S(P)}):S(P)}}}},render(){const{internalAddable:e,clsPrefix:n,name:i,disabled:g,label:s,tab:y,value:u,mergedClosable:l,trigger:h,$slots:{default:C}}=this,v=s??y;return p("div",{class:`${n}-tabs-tab-wrapper`},this.internalLeftPadded?p("div",{class:`${n}-tabs-tab-pad`}):null,p("div",Object.assign({key:i,"data-name":i,"data-disabled":g?!0:void 0},At({class:[`${n}-tabs-tab`,u===i&&`${n}-tabs-tab--active`,g&&`${n}-tabs-tab--disabled`,l&&`${n}-tabs-tab--closable`,e&&`${n}-tabs-tab--addable`,e?this.addTabClass:this.tabClass],onClick:h==="click"?this.activateTab:void 0,onMouseenter:h==="hover"?this.activateTab:void 0,style:e?this.addStyle:this.style},this.internalCreatedByPane?this.tabProps||{}:this.$attrs)),p("span",{class:`${n}-tabs-tab__label`},e?p(kt,null,p("div",{class:`${n}-tabs-tab__height-placeholder`}," "),p(Et,{clsPrefix:n},{default:()=>p(ba,null)})):C?C():typeof v=="object"?v:Vt(v??i)),l&&this.type==="card"?p(jt,{clsPrefix:n,class:`${n}-tabs-tab__close`,onClick:this.handleClose,disabled:g}):null))}}),_a=r("tabs",`
 box-sizing: border-box;
 width: 100%;
 display: flex;
 flex-direction: column;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
`,[d("segment-type",[r("tabs-rail",[T("&.transition-disabled",[r("tabs-capsule",`
 transition: none;
 `)])])]),d("top",[r("tab-pane",`
 padding: var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left);
 `)]),d("left",[r("tab-pane",`
 padding: var(--n-pane-padding-right) var(--n-pane-padding-bottom) var(--n-pane-padding-left) var(--n-pane-padding-top);
 `)]),d("left, right",`
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
 `)]),d("right",`
 flex-direction: row-reverse;
 `,[r("tab-pane",`
 padding: var(--n-pane-padding-left) var(--n-pane-padding-top) var(--n-pane-padding-right) var(--n-pane-padding-bottom);
 `),r("tabs-bar",`
 left: 0;
 `)]),d("bottom",`
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
 `,[d("active",`
 font-weight: var(--n-font-weight-strong);
 color: var(--n-tab-text-color-active);
 `),T("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])])]),d("flex",[r("tabs-nav",`
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
 `),B("prefix","padding-right: 16px;"),B("suffix","padding-left: 16px;")]),d("top, bottom",[T(">",[r("tabs-nav",[r("tabs-nav-scroll-wrapper",[T("&::before",`
 top: 0;
 bottom: 0;
 left: 0;
 width: 20px;
 `),T("&::after",`
 top: 0;
 bottom: 0;
 right: 0;
 width: 20px;
 `),d("shadow-start",[T("&::before",`
 box-shadow: inset 10px 0 8px -8px rgba(0, 0, 0, .12);
 `)]),d("shadow-end",[T("&::after",`
 box-shadow: inset -10px 0 8px -8px rgba(0, 0, 0, .12);
 `)])])])])]),d("left, right",[r("tabs-nav-scroll-content",`
 flex-direction: column;
 `),T(">",[r("tabs-nav",[r("tabs-nav-scroll-wrapper",[T("&::before",`
 top: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),T("&::after",`
 bottom: 0;
 left: 0;
 right: 0;
 height: 20px;
 `),d("shadow-start",[T("&::before",`
 box-shadow: inset 0 10px 8px -8px rgba(0, 0, 0, .12);
 `)]),d("shadow-end",[T("&::after",`
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
 `,[T("&::-webkit-scrollbar, &::-webkit-scrollbar-track-piece, &::-webkit-scrollbar-thumb",`
 width: 0;
 height: 0;
 display: none;
 `)]),T("&::before, &::after",`
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
 `,[d("disabled",{cursor:"not-allowed"}),B("close",`
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
 `,[T("&.transition-disabled",`
 transition: none;
 `),d("disabled",`
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
 `,[T("&.next-transition-leave-active, &.prev-transition-leave-active, &.next-transition-enter-active, &.prev-transition-enter-active",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 transform .2s var(--n-bezier),
 opacity .2s var(--n-bezier);
 `),T("&.next-transition-leave-active, &.prev-transition-leave-active",`
 position: absolute;
 `),T("&.next-transition-enter-from, &.prev-transition-leave-to",`
 transform: translateX(32px);
 opacity: 0;
 `),T("&.next-transition-leave-to, &.prev-transition-enter-from",`
 transform: translateX(-32px);
 opacity: 0;
 `),T("&.next-transition-leave-from, &.next-transition-enter-to, &.prev-transition-leave-from, &.prev-transition-enter-to",`
 transform: translateX(0);
 opacity: 1;
 `)]),r("tabs-tab-pad",`
 box-sizing: border-box;
 width: var(--n-tab-gap);
 flex-grow: 0;
 flex-shrink: 0;
 `),d("line-type, bar-type",[r("tabs-tab",`
 font-weight: var(--n-tab-font-weight);
 box-sizing: border-box;
 vertical-align: bottom;
 `,[T("&:hover",{color:"var(--n-tab-text-color-hover)"}),d("active",`
 color: var(--n-tab-text-color-active);
 font-weight: var(--n-tab-font-weight-active);
 `),d("disabled",{color:"var(--n-tab-text-color-disabled)"})])]),r("tabs-nav",[d("line-type",[d("top",[B("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 bottom: -1px;
 `)]),d("left",[B("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 right: -1px;
 `)]),d("right",[B("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-nav-scroll-content",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-bar",`
 left: -1px;
 `)]),d("bottom",[B("prefix, suffix",`
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
 `)]),d("card-type",[B("prefix, suffix",`
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
 `,[d("addable",`
 padding-left: 8px;
 padding-right: 8px;
 font-size: 16px;
 justify-content: center;
 `,[B("height-placeholder",`
 width: 0;
 font-size: var(--n-tab-font-size);
 `),Mt("disabled",[T("&:hover",`
 color: var(--n-tab-text-color-hover);
 `)])]),d("closable","padding-right: 8px;"),d("active",`
 background-color: #0000;
 font-weight: var(--n-tab-font-weight-active);
 color: var(--n-tab-text-color-active);
 `),d("disabled","color: var(--n-tab-text-color-disabled);")])]),d("left, right",`
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
 `)])]),d("top",[d("card-type",[r("tabs-scroll-padding","border-bottom: 1px solid var(--n-tab-border-color);"),B("prefix, suffix",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-top-right-radius: var(--n-tab-border-radius);
 `,[d("active",`
 border-bottom: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-bottom: 1px solid var(--n-tab-border-color);
 `)])]),d("left",[d("card-type",[r("tabs-scroll-padding","border-right: 1px solid var(--n-tab-border-color);"),B("prefix, suffix",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-left-radius: var(--n-tab-border-radius);
 border-bottom-left-radius: var(--n-tab-border-radius);
 `,[d("active",`
 border-right: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-right: 1px solid var(--n-tab-border-color);
 `)])]),d("right",[d("card-type",[r("tabs-scroll-padding","border-left: 1px solid var(--n-tab-border-color);"),B("prefix, suffix",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-top-right-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[d("active",`
 border-left: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-left: 1px solid var(--n-tab-border-color);
 `)])]),d("bottom",[d("card-type",[r("tabs-scroll-padding","border-top: 1px solid var(--n-tab-border-color);"),B("prefix, suffix",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-tab",`
 border-bottom-left-radius: var(--n-tab-border-radius);
 border-bottom-right-radius: var(--n-tab-border-radius);
 `,[d("active",`
 border-top: 1px solid #0000;
 `)]),r("tabs-tab-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `),r("tabs-pad",`
 border-top: 1px solid var(--n-tab-border-color);
 `)])])])]),ue=Ra,Wa=Object.assign(Object.assign({},Me.props),{value:[String,Number],defaultValue:[String,Number],trigger:{type:String,default:"click"},type:{type:String,default:"bar"},closable:Boolean,justifyContent:String,size:String,placement:{type:String,default:"top"},tabStyle:[String,Object],tabClass:String,addTabStyle:[String,Object],addTabClass:String,barWidth:Number,paneClass:String,paneStyle:[String,Object],paneWrapperClass:String,paneWrapperStyle:[String,Object],addable:[Boolean,Object],tabsPadding:{type:Number,default:0},animated:Boolean,onBeforeLeave:Function,onAdd:Function,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onClose:[Function,Array],labelSize:String,activeName:[String,Number],onActiveNameChange:[Function,Array]}),Aa=Z({name:"Tabs",props:Wa,slots:Object,setup(e,{slots:n}){var i,g,s,y;const{mergedClsPrefixRef:u,inlineThemeDisabled:l,mergedComponentPropsRef:h}=Ht(e),C=Me("Tabs","-tabs",_a,$a,e,u),v=_(null),f=_(null),x=_(null),S=_(null),w=_(null),P=_(null),$=_(!0),A=_(!0),W=ze(e,["labelSize","size"]),O=Q(()=>{var t,a;if(W.value)return W.value;const o=(a=(t=h==null?void 0:h.value)===null||t===void 0?void 0:t.Tabs)===null||a===void 0?void 0:a.size;return o||"medium"}),I=ze(e,["activeName","value"]),m=_((g=(i=I.value)!==null&&i!==void 0?i:e.defaultValue)!==null&&g!==void 0?g:n.default?(y=(s=de(n.default())[0])===null||s===void 0?void 0:s.props)===null||y===void 0?void 0:y.name:null),b=Ot(I,m),k={id:0},G=Q(()=>{if(!(!e.justifyContent||e.type==="card"))return{display:"flex",justifyContent:e.justifyContent}});ce(b,()=>{k.id=0,ee(),me()});function N(){var t;const{value:a}=b;return a===null?null:(t=v.value)===null||t===void 0?void 0:t.querySelector(`[data-name="${a}"]`)}function Oe(t){if(e.type==="card")return;const{value:a}=f;if(!a)return;const o=a.style.opacity==="0";if(t){const c=`${u.value}-tabs-bar--disabled`,{barWidth:R,placement:V}=e;if(t.dataset.disabled==="true"?a.classList.add(c):a.classList.remove(c),["top","bottom"].includes(V)){if(xe(["top","maxHeight","height"]),typeof R=="number"&&t.offsetWidth>=R){const j=Math.floor((t.offsetWidth-R)/2)+t.offsetLeft;a.style.left=`${j}px`,a.style.maxWidth=`${R}px`}else a.style.left=`${t.offsetLeft}px`,a.style.maxWidth=`${t.offsetWidth}px`;a.style.width="8192px",o&&(a.style.transition="none"),a.offsetWidth,o&&(a.style.transition="",a.style.opacity="1")}else{if(xe(["left","maxWidth","width"]),typeof R=="number"&&t.offsetHeight>=R){const j=Math.floor((t.offsetHeight-R)/2)+t.offsetTop;a.style.top=`${j}px`,a.style.maxHeight=`${R}px`}else a.style.top=`${t.offsetTop}px`,a.style.maxHeight=`${t.offsetHeight}px`;a.style.height="8192px",o&&(a.style.transition="none"),a.offsetHeight,o&&(a.style.transition="",a.style.opacity="1")}}}function Fe(){if(e.type==="card")return;const{value:t}=f;t&&(t.style.opacity="0")}function xe(t){const{value:a}=f;if(a)for(const o of t)a.style[o]=""}function ee(){if(e.type==="card")return;const t=N();t?Oe(t):Fe()}function me(){var t;const a=(t=w.value)===null||t===void 0?void 0:t.$el;if(!a)return;const o=N();if(!o)return;const{scrollLeft:c,offsetWidth:R}=a,{offsetLeft:V,offsetWidth:j}=o;c>V?a.scrollTo({top:0,left:V,behavior:"smooth"}):V+j>c+R&&a.scrollTo({top:0,left:V+j-R,behavior:"smooth"})}const te=_(null);let ie=0,F=null;function Ge(t){const a=te.value;if(a){ie=t.getBoundingClientRect().height;const o=`${ie}px`,c=()=>{a.style.height=o,a.style.maxHeight=o};F?(c(),F(),F=null):F=c}}function Ne(t){const a=te.value;if(a){const o=t.getBoundingClientRect().height,c=()=>{document.body.offsetHeight,a.style.maxHeight=`${o}px`,a.style.height=`${Math.max(ie,o)}px`};F?(F(),F=null,c()):F=c}}function De(){const t=te.value;if(t){t.style.maxHeight="",t.style.height="";const{paneWrapperStyle:a}=e;if(typeof a=="string")t.style.cssText=a;else if(a){const{maxHeight:o,height:c}=a;o!==void 0&&(t.style.maxHeight=o),c!==void 0&&(t.style.height=c)}}}const ye={value:[]},Ce=_("next");function Ue(t){const a=b.value;let o="next";for(const c of ye.value){if(c===a)break;if(c===t){o="prev";break}}Ce.value=o,Xe(t)}function Xe(t){const{onActiveNameChange:a,onUpdateValue:o,"onUpdate:value":c}=e;a&&ne(a,t),o&&ne(o,t),c&&ne(c,t),m.value=t}function Ke(t){const{onClose:a}=e;a&&ne(a,t)}function Se(){const{value:t}=f;if(!t)return;const a="transition-disabled";t.classList.add(a),ee(),t.classList.remove(a)}const D=_(null);function le({transitionDisabled:t}){const a=v.value;if(!a)return;t&&a.classList.add("transition-disabled");const o=N();o&&D.value&&(D.value.style.width=`${o.offsetWidth}px`,D.value.style.height=`${o.offsetHeight}px`,D.value.style.transform=`translateX(${o.offsetLeft-Ut(getComputedStyle(a).paddingLeft)}px)`,t&&D.value.offsetWidth),t&&a.classList.remove("transition-disabled")}ce([b],()=>{e.type==="segment"&&fe(()=>{le({transitionDisabled:!1})})}),Ft(()=>{e.type==="segment"&&le({transitionDisabled:!0})});let we=0;function qe(t){var a;if(t.contentRect.width===0&&t.contentRect.height===0||we===t.contentRect.width)return;we=t.contentRect.width;const{type:o}=e;if((o==="line"||o==="bar")&&Se(),o!=="segment"){const{placement:c}=e;se((c==="top"||c==="bottom"?(a=w.value)===null||a===void 0?void 0:a.$el:P.value)||null)}}const Ye=ue(qe,64);ce([()=>e.justifyContent,()=>e.size],()=>{fe(()=>{const{type:t}=e;(t==="line"||t==="bar")&&Se()})});const U=_(!1);function Je(t){var a;const{target:o,contentRect:{width:c,height:R}}=t,V=o.parentElement.parentElement.offsetWidth,j=o.parentElement.parentElement.offsetHeight,{placement:K}=e;if(!U.value)K==="top"||K==="bottom"?V<c&&(U.value=!0):j<R&&(U.value=!0);else{const{value:q}=S;if(!q)return;K==="top"||K==="bottom"?V-c>q.$el.offsetWidth&&(U.value=!1):j-R>q.$el.offsetHeight&&(U.value=!1)}se(((a=w.value)===null||a===void 0?void 0:a.$el)||null)}const Qe=ue(Je,64);function Ze(){const{onAdd:t}=e;t&&t(),fe(()=>{const a=N(),{value:o}=w;!a||!o||o.scrollTo({left:a.offsetLeft,top:0,behavior:"smooth"})})}function se(t){if(!t)return;const{placement:a}=e;if(a==="top"||a==="bottom"){const{scrollLeft:o,scrollWidth:c,offsetWidth:R}=t;$.value=o<=0,A.value=o+R>=c}else{const{scrollTop:o,scrollHeight:c,offsetHeight:R}=t;$.value=o<=0,A.value=o+R>=c}}const et=ue(t=>{se(t.target)},64);Jt(he,{triggerRef:H(e,"trigger"),tabStyleRef:H(e,"tabStyle"),tabClassRef:H(e,"tabClass"),addTabStyleRef:H(e,"addTabStyle"),addTabClassRef:H(e,"addTabClass"),paneClassRef:H(e,"paneClass"),paneStyleRef:H(e,"paneStyle"),mergedClsPrefixRef:u,typeRef:H(e,"type"),closableRef:H(e,"closable"),valueRef:b,tabChangeIdRef:k,onBeforeLeaveRef:H(e,"onBeforeLeave"),activateTab:Ue,handleClose:Ke,handleAdd:Ze}),Gt(()=>{ee(),me()}),Nt(()=>{const{value:t}=x;if(!t)return;const{value:a}=u,o=`${a}-tabs-nav-scroll-wrapper--shadow-start`,c=`${a}-tabs-nav-scroll-wrapper--shadow-end`;$.value?t.classList.remove(o):t.classList.add(o),A.value?t.classList.remove(c):t.classList.add(c)});const tt={syncBarPosition:()=>{ee()}},at=()=>{le({transitionDisabled:!0})},Te=Q(()=>{const{value:t}=O,{type:a}=e,o={card:"Card",bar:"Bar",line:"Line",segment:"Segment"}[a],c=`${t}${o}`,{self:{barColor:R,closeIconColor:V,closeIconColorHover:j,closeIconColorPressed:K,tabColor:q,tabBorderColor:rt,paneTextColor:nt,tabFontWeight:ot,tabBorderRadius:it,tabFontWeightActive:lt,colorSegment:st,fontWeightStrong:dt,tabColorSegment:bt,closeSize:ct,closeIconSize:ft,closeColorHover:pt,closeColorPressed:ut,closeBorderRadius:vt,[M("panePadding",t)]:ae,[M("tabPadding",c)]:gt,[M("tabPaddingVertical",c)]:ht,[M("tabGap",c)]:xt,[M("tabGap",`${c}Vertical`)]:mt,[M("tabTextColor",a)]:yt,[M("tabTextColorActive",a)]:Ct,[M("tabTextColorHover",a)]:St,[M("tabTextColorDisabled",a)]:wt,[M("tabFontSize",t)]:Tt},common:{cubicBezierEaseInOut:Pt}}=C.value;return{"--n-bezier":Pt,"--n-color-segment":st,"--n-bar-color":R,"--n-tab-font-size":Tt,"--n-tab-text-color":yt,"--n-tab-text-color-active":Ct,"--n-tab-text-color-disabled":wt,"--n-tab-text-color-hover":St,"--n-pane-text-color":nt,"--n-tab-border-color":rt,"--n-tab-border-radius":it,"--n-close-size":ct,"--n-close-icon-size":ft,"--n-close-color-hover":pt,"--n-close-color-pressed":ut,"--n-close-border-radius":vt,"--n-close-icon-color":V,"--n-close-icon-color-hover":j,"--n-close-icon-color-pressed":K,"--n-tab-color":q,"--n-tab-font-weight":ot,"--n-tab-font-weight-active":lt,"--n-tab-padding":gt,"--n-tab-padding-vertical":ht,"--n-tab-gap":xt,"--n-tab-gap-vertical":mt,"--n-pane-padding-left":re(ae,"left"),"--n-pane-padding-right":re(ae,"right"),"--n-pane-padding-top":re(ae,"top"),"--n-pane-padding-bottom":re(ae,"bottom"),"--n-font-weight-strong":dt,"--n-tab-color-segment":bt}}),X=l?Dt("tabs",Q(()=>`${O.value[0]}${e.type[0]}`),Te,e):void 0;return Object.assign({mergedClsPrefix:u,mergedValue:b,renderedNames:new Set,segmentCapsuleElRef:D,tabsPaneWrapperRef:te,tabsElRef:v,barElRef:f,addTabInstRef:S,xScrollInstRef:w,scrollWrapperElRef:x,addTabFixed:U,tabWrapperStyle:G,handleNavResize:Ye,mergedSize:O,handleScroll:et,handleTabsResize:Qe,cssVars:l?void 0:Te,themeClass:X==null?void 0:X.themeClass,animationDirection:Ce,renderNameListRef:ye,yScrollElRef:P,handleSegmentResize:at,onAnimationBeforeLeave:Ge,onAnimationEnter:Ne,onAnimationAfterEnter:De,onRender:X==null?void 0:X.onRender},tt)},render(){const{mergedClsPrefix:e,type:n,placement:i,addTabFixed:g,addable:s,mergedSize:y,renderNameListRef:u,onRender:l,paneWrapperClass:h,paneWrapperStyle:C,$slots:{default:v,prefix:f,suffix:x}}=this;l==null||l();const S=v?de(v()).filter(m=>m.type.__TAB_PANE__===!0):[],w=v?de(v()).filter(m=>m.type.__TAB__===!0):[],P=!w.length,$=n==="card",A=n==="segment",W=!$&&!A&&this.justifyContent;u.value=[];const O=()=>{const m=p("div",{style:this.tabWrapperStyle,class:`${e}-tabs-wrapper`},W?null:p("div",{class:`${e}-tabs-scroll-padding`,style:i==="top"||i==="bottom"?{width:`${this.tabsPadding}px`}:{height:`${this.tabsPadding}px`}}),P?S.map((b,k)=>(u.value.push(b.props.name),ve(p(ge,Object.assign({},b.props,{internalCreatedByPane:!0,internalLeftPadded:k!==0&&(!W||W==="center"||W==="start"||W==="end")}),b.children?{default:b.children.tab}:void 0)))):w.map((b,k)=>(u.value.push(b.props.name),ve(k!==0&&!W?je(b):b))),!g&&s&&$?Ve(s,(P?S.length:w.length)!==0):null,W?null:p("div",{class:`${e}-tabs-scroll-padding`,style:{width:`${this.tabsPadding}px`}}));return p("div",{ref:"tabsElRef",class:`${e}-tabs-nav-scroll-content`},$&&s?p(be,{onResize:this.handleTabsResize},{default:()=>m}):m,$?p("div",{class:`${e}-tabs-pad`}):null,$?null:p("div",{ref:"barElRef",class:`${e}-tabs-bar`}))},I=A?"top":i;return p("div",{class:[`${e}-tabs`,this.themeClass,`${e}-tabs--${n}-type`,`${e}-tabs--${y}-size`,W&&`${e}-tabs--flex`,`${e}-tabs--${I}`],style:this.cssVars},p("div",{class:[`${e}-tabs-nav--${n}-type`,`${e}-tabs-nav--${I}`,`${e}-tabs-nav`]},Re(f,m=>m&&p("div",{class:`${e}-tabs-nav__prefix`},m)),A?p(be,{onResize:this.handleSegmentResize},{default:()=>p("div",{class:`${e}-tabs-rail`,ref:"tabsElRef"},p("div",{class:`${e}-tabs-capsule`,ref:"segmentCapsuleElRef"},p("div",{class:`${e}-tabs-wrapper`},p("div",{class:`${e}-tabs-tab`}))),P?S.map((m,b)=>(u.value.push(m.props.name),p(ge,Object.assign({},m.props,{internalCreatedByPane:!0,internalLeftPadded:b!==0}),m.children?{default:m.children.tab}:void 0))):w.map((m,b)=>(u.value.push(m.props.name),b===0?m:je(m))))}):p(be,{onResize:this.handleNavResize},{default:()=>p("div",{class:`${e}-tabs-nav-scroll-wrapper`,ref:"scrollWrapperElRef"},["top","bottom"].includes(I)?p(fa,{ref:"xScrollInstRef",onScroll:this.handleScroll},{default:O}):p("div",{class:`${e}-tabs-nav-y-scroll`,onScroll:this.handleScroll,ref:"yScrollElRef"},O()))}),g&&s&&$?Ve(s,!0):null,Re(x,m=>m&&p("div",{class:`${e}-tabs-nav__suffix`},m))),P&&(this.animated&&(I==="top"||I==="bottom")?p("div",{ref:"tabsPaneWrapperRef",style:C,class:[`${e}-tabs-pane-wrapper`,h]},Ee(S,this.mergedValue,this.renderedNames,this.onAnimationBeforeLeave,this.onAnimationEnter,this.onAnimationAfterEnter,this.animationDirection)):Ee(S,this.mergedValue,this.renderedNames)))}});function Ee(e,n,i,g,s,y,u){const l=[];return e.forEach(h=>{const{name:C,displayDirective:v,"display-directive":f}=h.props,x=w=>v===w||f===w,S=n===C;if(h.key!==void 0&&(h.key=C),S||x("show")||x("show:lazy")&&i.has(C)){i.has(C)||i.add(C);const w=!x("if");l.push(w?Xt(h,[[Kt,S]]):h)}}),u?p(qt,{name:`${u}-transition`,onBeforeLeave:g,onEnter:s,onAfterEnter:y},{default:()=>l}):l}function Ve(e,n){return p(ge,{ref:"addTabInstRef",key:"__addable",name:"__addable",internalCreatedByPane:!0,internalAddable:!0,internalLeftPadded:n,disabled:typeof e=="object"&&e.disabled})}function je(e){const n=Yt(e);return n.props?n.props.internalLeftPadded=!0:n.props={internalLeftPadded:!0},n}function ve(e){return Array.isArray(e.dynamicProps)?e.dynamicProps.includes("internalLeftPadded")||e.dynamicProps.push("internalLeftPadded"):e.dynamicProps=["internalLeftPadded"],e}const ka={class:"login-bg"},Ea={class:"brand"},Va=Z({__name:"LoginView",setup(e){const n=aa(),i=da(),g=Qt(),s=Zt();s.load();const y=_("login"),u=_(!1),l=_({username:"admin_user",password:"admin_user",user:"",pwd:"",confirm:""});async function h(){u.value=!0;try{const v=await ia(l.value.username,l.value.password);g.setAuth(v.token,v.name,v.role_id,v.requirePasswordChange),i.success("登录成功"),n.push(v.requirePasswordChange?"/profile":"/dashboard")}catch(v){i.error(v.message)}finally{u.value=!1}}async function C(){if(l.value.pwd!==l.value.confirm){i.error("两次密码不一致");return}u.value=!0;try{await la(l.value.user,l.value.pwd),i.success("注册成功,请登录"),y.value="login"}catch(v){i.error(v.message)}finally{u.value=!1}}return(v,f)=>(ra(),ea("div",ka,[L(z(ta),{class:"login-card",bordered:!1},{default:E(()=>[na("h1",Ea,oa(z(s).appName),1),L(z(Aa),{value:y.value,"onUpdate:value":f[5]||(f[5]=x=>y.value=x),type:"line",animated:""},{default:E(()=>[L(z(ke),{name:"login",tab:"登录"},{default:E(()=>[L(z(_e),{onKeyup:Le(h,["enter"])},{default:E(()=>[L(z(Y),{label:"用户名"},{default:E(()=>[L(z(J),{value:l.value.username,"onUpdate:value":f[0]||(f[0]=x=>l.value.username=x),placeholder:"用户名"},null,8,["value"])]),_:1}),L(z(Y),{label:"密码"},{default:E(()=>[L(z(J),{value:l.value.password,"onUpdate:value":f[1]||(f[1]=x=>l.value.password=x),type:"password","show-password-on":"click",placeholder:"密码"},null,8,["value"])]),_:1}),L(z($e),{type:"primary",block:"",loading:u.value,onClick:h},{default:E(()=>[...f[6]||(f[6]=[Be("登 录",-1)])]),_:1},8,["loading"])]),_:1})]),_:1}),L(z(ke),{name:"register",tab:"注册"},{default:E(()=>[L(z(_e),{onKeyup:Le(C,["enter"])},{default:E(()=>[L(z(Y),{label:"用户名"},{default:E(()=>[L(z(J),{value:l.value.user,"onUpdate:value":f[2]||(f[2]=x=>l.value.user=x),placeholder:"3-20 位"},null,8,["value"])]),_:1}),L(z(Y),{label:"密码"},{default:E(()=>[L(z(J),{value:l.value.pwd,"onUpdate:value":f[3]||(f[3]=x=>l.value.pwd=x),type:"password","show-password-on":"click",placeholder:"6-32 位"},null,8,["value"])]),_:1}),L(z(Y),{label:"确认密码"},{default:E(()=>[L(z(J),{value:l.value.confirm,"onUpdate:value":f[4]||(f[4]=x=>l.value.confirm=x),type:"password","show-password-on":"click"},null,8,["value"])]),_:1}),L(z($e),{type:"primary",block:"",loading:u.value,onClick:C},{default:E(()=>[...f[7]||(f[7]=[Be("注 册",-1)])]),_:1},8,["loading"])]),_:1})]),_:1})]),_:1},8,["value"])]),_:1})]))}}),Oa=sa(Va,[["__scopeId","data-v-0171ae2e"]]);export{Oa as default};
