import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { ArrowRight, ArrowUpRight, MapPin, Menu, X } from '@Components/icons';
import { SchoolMonogram } from '@Components';
import { APP_ROUTE } from '@Business/app/routes';
import { AUTH_ROUTE } from '@Business/auth/routes';
import { getSiteNavClassName } from './helpers/siteChrome.getNavClassName';
import { school } from '../../constants/school.content';
import { SCHOOL_NAV_ITEMS as navItems } from '../../constants/school.navigation';
import { EnquiryLink } from '../siteActions';


const SiteHeader = () => {
  const [open,setOpen] = useState(false);
  const { pathname } = useLocation();
  let menuLabel = 'Open menu';
  let MenuIcon = Menu;
  if (open) { menuLabel = 'Close menu'; MenuIcon = X; }
  useEffect(() => setOpen(false),[pathname]);
  useEffect(() => {
    if (!open) return;
    const previous=document.body.style.overflow;
    document.body.style.overflow='hidden';
    return () => { document.body.style.overflow=previous; };
  },[open]);
  return (
    <>
      <Box as="div" className="rdr-topline">
        <Box as="div" className="rdr-container rdr-topline__inner">
          <Box as="span">HALUHERA, REWARI · HARYANA</Box>
          <Box as="span">ESTABLISHED 2006</Box>
        </Box>
      </Box>
      <Box as="header" className="rdr-header">
        <Box as="div" className="rdr-container rdr-header__inner">
          <Link to={APP_ROUTE.HOME} className="rdr-logo" aria-label="RDR School home">
            <SchoolMonogram/>
            <Box as="span"><Box as="strong">RDR School</Box><Box as="small">Rao Deena Ram Vidya Vihar</Box></Box>
          </Link>
          <Box as="nav" className="rdr-nav" aria-label="Main navigation">
            {navItems.map(item => <Link className={getSiteNavClassName(pathname, item.to)} to={item.to} key={item.to}>{item.label}</Link>)}
          </Box>
          <Box as="div" className="rdr-header__actions">
            <EnquiryLink className="rdr-btn rdr-btn--solid rdr-header__enquire">Enquire now <ArrowUpRight size={16}/></EnquiryLink>
            <Button unstyled type="button" className="rdr-menu-button" aria-label={menuLabel}
              aria-expanded={open} aria-controls="rdr-mobile-navigation" onClick={()=>setOpen(v=>!v)}>
              <MenuIcon size={22}/>
            </Button>
          </Box>
        </Box>
      </Box>
      {open && <>
        <Button unstyled className="rdr-mobile-nav__scrim" type="button" aria-label="Close menu" onClick={()=>setOpen(false)}/>
        <Box as="nav" className="rdr-mobile-nav" id="rdr-mobile-navigation" aria-label="Mobile navigation">
          <Box as="div" className="rdr-mobile-nav__top"><Box as="strong">Explore RDR School</Box><Button unstyled type="button" onClick={()=>setOpen(false)} aria-label="Close menu"><X size={21}/></Button></Box>
          {navItems.map(item=><Link key={item.to} to={item.to} onClick={()=>setOpen(false)}>{item.label}<ArrowUpRight size={16}/></Link>)}
          <Link to={SCHOOL_ROUTE.DIRECTOR}>Director's message <ArrowUpRight size={16}/></Link>
          <Link to={SCHOOL_ROUTE.PRINCIPAL}>Principal's message <ArrowUpRight size={16}/></Link>
          <Link to={SCHOOL_ROUTE.CONTACT}>Contact & directions <ArrowUpRight size={16}/></Link>
          <EnquiryLink className="rdr-mobile-nav__action" onClick={()=>setOpen(false)}>Start an enquiry <ArrowRight size={16}/></EnquiryLink>
        </Box>
      </>}
    </>
  );
};

const SiteFooter = () => (
  <Box as="footer" className="rdr-footer">
    <Box as="div" className="rdr-container">
      <Box as="div" className="rdr-footer__top">
        <Box as="div" className="rdr-footer__intro">
          <Box as="div" className="rdr-logo rdr-logo--footer"><SchoolMonogram/><Box as="span"><Box as="strong">RDR School</Box><Box as="small">Rao Deena Ram Vidya Vihar</Box></Box></Box>
          <Box as="p">Learning with purpose. Growing with confidence. Serving the community since 2006.</Box>
          <Box as="span" className="rdr-footer__location"><MapPin size={16}/>{school.address}</Box>
        </Box>
        <Box as="div"><Box as="h3">Discover</Box><Link to={SCHOOL_ROUTE.ABOUT}>Our story</Link><Link to={SCHOOL_ROUTE.ACADEMICS}>Academics</Link><Link to={SCHOOL_ROUTE.CAMPUS}>Facilities</Link><Link to={SCHOOL_ROUTE.GALLERY}>School life</Link></Box>
        <Box as="div"><Box as="h3">For families</Box><Link to={SCHOOL_ROUTE.ADMISSIONS}>Admissions</Link><Link to={SCHOOL_ROUTE.RESOURCES}>Resources</Link><Link to={SCHOOL_ROUTE.CONTACT}>Contact</Link><Link to={AUTH_ROUTE.LOGIN}>School workspace</Link></Box>
        <Box as="div"><Box as="h3">Transparency</Box><Link to={SCHOOL_ROUTE.DISCLOSURE}>Mandatory Public Disclosure</Link><ExternalAnchor href={school.cbseSaras} target="_blank" rel="noreferrer">CBSE affiliation</ExternalAnchor><ExternalAnchor href={school.disclosurePdf} target="_blank" rel="noreferrer">2023 disclosure archive</ExternalAnchor></Box>
      </Box>
      <Box as="div" className="rdr-footer__bottom"><Box as="span">© {new Date().getFullYear()} RDR School. All rights reserved.</Box><Box as="span">Haluhera · Rewari · Haryana</Box></Box>
    </Box>
  </Box>
);

export { SiteHeader };
export { SiteFooter };
