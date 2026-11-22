import { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { ArrowRight, ArrowUpRight, MapPin, Menu, X } from 'lucide-react';
import { SchoolMonogram } from '@Components';
import { school } from '../../constants/school.content';
import { SCHOOL_NAV_ITEMS as navItems } from '../../constants/school.navigation';
import { EnquiryLink } from '../siteActions';


const SiteHeader = () => {
  const [open,setOpen] = useState(false);
  const { pathname } = useLocation();
  useEffect(() => setOpen(false),[pathname]);
  useEffect(() => {
    if (!open) return;
    const previous=document.body.style.overflow;
    document.body.style.overflow='hidden';
    return () => { document.body.style.overflow=previous; };
  },[open]);
  return (
    <>
      <div className="rdr-topline">
        <div className="rdr-container rdr-topline__inner">
          <span>HALUHERA, REWARI · HARYANA</span>
          <span>ESTABLISHED 2006</span>
        </div>
      </div>
      <header className="rdr-header">
        <div className="rdr-container rdr-header__inner">
          <Link to="/" className="rdr-logo" aria-label="RDR School home">
            <SchoolMonogram/>
            <span><strong>RDR School</strong><small>Rao Deena Ram Vidya Vihar</small></span>
          </Link>
          <nav className="rdr-nav" aria-label="Main navigation">
            {navItems.map(item => <Link className={pathname===item.to?'active':''} to={item.to} key={item.to}>{item.label}</Link>)}
          </nav>
          <div className="rdr-header__actions">
            <EnquiryLink className="rdr-btn rdr-btn--solid rdr-header__enquire">Enquire now <ArrowUpRight size={16}/></EnquiryLink>
            <button type="button" className="rdr-menu-button" aria-label={open?'Close menu':'Open menu'}
              aria-expanded={open} aria-controls="rdr-mobile-navigation" onClick={()=>setOpen(v=>!v)}>
              {open?<X size={22}/>:<Menu size={22}/>}
            </button>
          </div>
        </div>
      </header>
      {open && <>
        <button className="rdr-mobile-nav__scrim" type="button" aria-label="Close menu" onClick={()=>setOpen(false)}/>
        <nav className="rdr-mobile-nav" id="rdr-mobile-navigation" aria-label="Mobile navigation">
          <div className="rdr-mobile-nav__top"><strong>Explore RDR School</strong><button type="button" onClick={()=>setOpen(false)} aria-label="Close menu"><X size={21}/></button></div>
          {navItems.map(item=><Link key={item.to} to={item.to} onClick={()=>setOpen(false)}>{item.label}<ArrowUpRight size={16}/></Link>)}
          <Link to="/about/director">Director's message <ArrowUpRight size={16}/></Link>
          <Link to="/about/principal">Principal's message <ArrowUpRight size={16}/></Link>
          <Link to="/contact">Contact & directions <ArrowUpRight size={16}/></Link>
          <EnquiryLink className="rdr-mobile-nav__action" onClick={()=>setOpen(false)}>Start an enquiry <ArrowRight size={16}/></EnquiryLink>
        </nav>
      </>}
    </>
  );
};

const SiteFooter = () => (
  <footer className="rdr-footer">
    <div className="rdr-container">
      <div className="rdr-footer__top">
        <div className="rdr-footer__intro">
          <div className="rdr-logo rdr-logo--footer"><SchoolMonogram/><span><strong>RDR School</strong><small>Rao Deena Ram Vidya Vihar</small></span></div>
          <p>Learning with purpose. Growing with confidence. Serving the community since 2006.</p>
          <span className="rdr-footer__location"><MapPin size={16}/>{school.address}</span>
        </div>
        <div><h3>Discover</h3><Link to="/about">Our story</Link><Link to="/academics">Academics</Link><Link to="/campus">Facilities</Link><Link to="/gallery">School life</Link></div>
        <div><h3>For families</h3><Link to="/admissions">Admissions</Link><Link to="/resources">Resources</Link><Link to="/contact">Contact</Link><Link to="/login">School workspace</Link></div>
        <div><h3>Transparency</h3><Link to="/mandatory-public-disclosure">Mandatory Public Disclosure</Link><a href={school.cbseSaras} target="_blank" rel="noreferrer">CBSE affiliation</a><a href={school.disclosurePdf} target="_blank" rel="noreferrer">2023 disclosure archive</a></div>
      </div>
      <div className="rdr-footer__bottom"><span>© {new Date().getFullYear()} RDR School. All rights reserved.</span><span>Haluhera · Rewari · Haryana</span></div>
    </div>
  </footer>
);

export { SiteHeader };
export { SiteFooter };
