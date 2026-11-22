import { useEffect, useMemo, useState, type MouseEvent } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { pageContent, pathPanel, pageMetadata } from './config/pageRegistry';
import { SiteHeader, SiteFooter } from './components/siteChrome';
import { HomePage } from './components/homeContent';
import { EnquiryActionContext } from './components/siteActions';
import { EnquiryForm } from './components/enquiryForm';
import { PublicSheet, type SchoolPanel } from './components/publicSheet';

import './styles/landingPage.css';
import './styles/optionOne.css';
import './styles/optionOnePages.css';

export const LandingPage = () => {
  const location=useLocation();
  const navigate=useNavigate();
  const [enquiryOpen,setEnquiryOpen]=useState(false);
  const state=location.state as {background?:string;panel?:SchoolPanel}|null;
  const routedPanel:SchoolPanel|null=null;
  const panel:SchoolPanel|null=enquiryOpen?'enquiry':routedPanel;
  const backgroundPath=location.pathname;
  const CurrentPage=useMemo(()=>pageContent[backgroundPath] ?? HomePage,[backgroundPath]);

  const closePanel=()=>{
    if(enquiryOpen){setEnquiryOpen(false);return;}
    navigate(state?.background??'/',{replace:true});
  };
  const selectPanel=(next:SchoolPanel)=>{
    if(next==='enquiry'){setEnquiryOpen(true);return;}
    setEnquiryOpen(false);
    const path=Object.entries(pathPanel).find(([,value])=>value===next)?.[0]??'/about';
    navigate(path,{replace:true,state:{background:backgroundPath,panel:next}});
  };

  // Public navigation opens full, addressable pages. The enquiry CTA
  // continues to use the existing accessible sheet until Lead CRM is ready.
  const handleInternalLink=(event:MouseEvent<HTMLDivElement>)=>{
    if(event.defaultPrevented || event.button!==0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey)return;
    if(!(event.target instanceof Element))return;
    const anchor=event.target.closest<HTMLAnchorElement>('a[href]');
    if(!anchor || anchor.hasAttribute('download') || (anchor.target&&anchor.target!=='_self'))return;
    if(anchor.dataset.rdrEnquiryAction==='true')return;
    const url=new URL(anchor.href,window.location.href);
    if(url.origin!==window.location.origin)return;
    if(url.pathname==='/admissions' && url.hash==='#enquiry'){
      event.preventDefault();
      setEnquiryOpen(true);
    }
  };
  useEffect(()=>{
    if(panel)return;
    if(location.hash){
      window.requestAnimationFrame(()=>document.getElementById(location.hash.slice(1))?.scrollIntoView({behavior:'smooth',block:'start'}));
    } else window.scrollTo({top:0,behavior:'instant'});
  },[backgroundPath,location.hash,panel]);
  useEffect(()=>{
    const details=pageMetadata[location.pathname] ?? (
      location.pathname.startsWith('/about/')?pageMetadata['/about']:pageMetadata['/']
    );
    if(!details)return;
    document.title=details.title;
    const description=document.querySelector<HTMLMetaElement>('meta[name="description"]');
    if(description)description.content=details.description;
    let canonical=document.querySelector<HTMLLinkElement>('link[rel="canonical"]');
    if(!canonical){canonical=document.createElement('link');canonical.rel='canonical';document.head.appendChild(canonical);}
    canonical.href=`https://rdrschool.com${location.pathname}`;
  },[location.pathname]);

  return <EnquiryActionContext.Provider value={()=>setEnquiryOpen(true)}>
    <div className="rdr-site" onClickCapture={handleInternalLink}>
      <SiteHeader/>
      <main id="main-content"><CurrentPage/></main>
      <SiteFooter/>
      <PublicSheet panel={panel} onClose={closePanel} onNavigate={selectPanel} enquiry={<EnquiryForm compact/>}/>
    </div>
  </EnquiryActionContext.Provider>;
};
