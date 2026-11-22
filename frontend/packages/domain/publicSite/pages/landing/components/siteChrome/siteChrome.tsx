import { _map } from '@Base/utilities/map';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { useEffect, useState } from 'react';
import { AppLink } from '@Components/components/AppLink';
import { useLocation } from '@Components/router';
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
          <AppLink to={APP_ROUTE.HOME} className="rdr-logo" aria-label="RDR School home">
            <SchoolMonogram/>
            <Box as="span"><Box as="strong">RDR School</Box><Box as="small">Rao Deena Ram Vidya Vihar</Box></Box>
          </AppLink>
          <Box as="nav" className="rdr-nav" aria-label="Main navigation">
            {_map(navItems, item => <AppLink className={getSiteNavClassName(pathname, item.to)} to={item.to} key={item.to}>{item.label}</AppLink>)}
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
          {_map(navItems, item=><AppLink key={item.to} to={item.to} onClick={()=>setOpen(false)}>{item.label}<ArrowUpRight size={16}/></AppLink>)}
          <AppLink to={SCHOOL_ROUTE.DIRECTOR}>Director's message <ArrowUpRight size={16}/></AppLink>
          <AppLink to={SCHOOL_ROUTE.PRINCIPAL}>Principal's message <ArrowUpRight size={16}/></AppLink>
          <AppLink to={SCHOOL_ROUTE.CONTACT}>Contact & directions <ArrowUpRight size={16}/></AppLink>
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
        <Box as="div"><Box as="h3">Discover</Box><AppLink to={SCHOOL_ROUTE.ABOUT}>Our story</AppLink><AppLink to={SCHOOL_ROUTE.ACADEMICS}>Academics</AppLink><AppLink to={SCHOOL_ROUTE.CAMPUS}>Facilities</AppLink><AppLink to={SCHOOL_ROUTE.GALLERY}>School life</AppLink></Box>
        <Box as="div"><Box as="h3">For families</Box><AppLink to={SCHOOL_ROUTE.ADMISSIONS}>Admissions</AppLink><AppLink to={SCHOOL_ROUTE.RESOURCES}>Resources</AppLink><AppLink to={SCHOOL_ROUTE.CONTACT}>Contact</AppLink><AppLink to={AUTH_ROUTE.LOGIN}>School workspace</AppLink></Box>
        <Box as="div"><Box as="h3">Transparency</Box><AppLink to={SCHOOL_ROUTE.DISCLOSURE}>Mandatory Public Disclosure</AppLink><ExternalAnchor href={school.cbseSaras} target="_blank" rel="noreferrer">CBSE affiliation</ExternalAnchor><ExternalAnchor href={school.disclosurePdf} target="_blank" rel="noreferrer">2023 disclosure archive</ExternalAnchor></Box>
      </Box>
      <Box as="div" className="rdr-footer__bottom"><Box as="span">© {new Date().getFullYear()} RDR School. All rights reserved.</Box><Box as="span">Haluhera · Rewari · Haryana</Box></Box>
    </Box>
  </Box>
);

export { SiteHeader };
export { SiteFooter };
