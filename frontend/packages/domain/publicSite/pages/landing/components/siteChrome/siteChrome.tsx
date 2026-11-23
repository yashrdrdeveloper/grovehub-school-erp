import { _map } from '@Base/utilities/map';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Button } from '@Components/components/Button';
import { Box , BOX_TAG} from '@Components/components/Box';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { useEffect, useState } from 'react';
import { AppLink } from '@Components/components/AppLink';
import { useLocation } from '@Components/router';
import { ArrowRight, ArrowUpRight, MapPin, X } from '@Components/icons';
import { SchoolMonogram } from '@Components';
import { APP_ROUTE } from '@Business/app/routes';
import { AUTH_ROUTE } from '@Business/auth/routes';
import { getSiteNavClassName } from './helpers/siteChrome.getNavClassName';
import { getSiteChromeMenuView } from './helpers/siteChrome.getMenuView';
import { school } from '../../constants/school.content';
import { SCHOOL_NAV_ITEMS as navItems } from '../../constants/school.navigation';
import { EnquiryLink } from '../siteActions';


const SiteHeader = () => {
  const [open,setOpen] = useState(false);
  const { pathname } = useLocation();
  const { menuLabel, MenuIcon } = getSiteChromeMenuView(open);
  useEffect(() => setOpen(false),[pathname]);
  useEffect(() => {
    if (!open) return;
    const previous=document.body.style.overflow;
    document.body.style.overflow='hidden';
    return () => { document.body.style.overflow=previous; };
  },[open]);
  return (
    <>
      <Box as={BOX_TAG.DIV} className="rdr-topline">
        <Box as={BOX_TAG.DIV} className="rdr-container rdr-topline__inner">
          <Box as={BOX_TAG.SPAN}>{__("HALUHERA, REWARI · HARYANA")}</Box>
          <Box as={BOX_TAG.SPAN}>{__("ESTABLISHED 2006")}</Box>
        </Box>
      </Box>
      <Box as={BOX_TAG.HEADER} className="rdr-header">
        <Box as={BOX_TAG.DIV} className="rdr-container rdr-header__inner">
          <AppLink to={APP_ROUTE.HOME} className="rdr-logo" aria-label={__("RDR School home")}>
            <SchoolMonogram/>
            <Box as={BOX_TAG.SPAN}><Box as={BOX_TAG.STRONG}>{__("RDR School")}</Box><Box as={BOX_TAG.SMALL}>{__("Rao Deena Ram Vidya Vihar")}</Box></Box>
          </AppLink>
          <Box as={BOX_TAG.NAV} className="rdr-nav" aria-label={__("Main navigation")}>
            {_map(navItems, item => <AppLink className={getSiteNavClassName(pathname, item.to)} to={item.to} key={item.to}>{item.label}</AppLink>)}
          </Box>
          <Box as={BOX_TAG.DIV} className="rdr-header__actions">
            <EnquiryLink className="rdr-btn rdr-btn--solid rdr-header__enquire">{__("Enquire now ")}<ArrowUpRight size={16}/></EnquiryLink>
            <Button unstyled type="button" className="rdr-menu-button" aria-label={menuLabel}
              aria-expanded={open} aria-controls="rdr-mobile-navigation" onClick={()=>setOpen(v=>!v)}>
              <MenuIcon size={22}/>
            </Button>
          </Box>
        </Box>
      </Box>
      {open && <>
        <Button unstyled className="rdr-mobile-nav__scrim" type="button" aria-label={__("Close menu")} onClick={()=>setOpen(false)}/>
        <Box as={BOX_TAG.NAV} className="rdr-mobile-nav" id="rdr-mobile-navigation" aria-label={__("Mobile navigation")}>
          <Box as={BOX_TAG.DIV} className="rdr-mobile-nav__top"><Box as={BOX_TAG.STRONG}>{__("Explore RDR School")}</Box><Button unstyled type="button" onClick={()=>setOpen(false)} aria-label={__("Close menu")}><X size={21}/></Button></Box>
          {_map(navItems, item=><AppLink key={item.to} to={item.to} onClick={()=>setOpen(false)}>{item.label}<ArrowUpRight size={16}/></AppLink>)}
          <AppLink to={SCHOOL_ROUTE.DIRECTOR}>{__("Director's message ")}<ArrowUpRight size={16}/></AppLink>
          <AppLink to={SCHOOL_ROUTE.PRINCIPAL}>{__("Principal's message ")}<ArrowUpRight size={16}/></AppLink>
          <AppLink to={SCHOOL_ROUTE.CONTACT}>{__("Contact & directions ")}<ArrowUpRight size={16}/></AppLink>
          <EnquiryLink className="rdr-mobile-nav__action" onClick={()=>setOpen(false)}>{__("Start an enquiry ")}<ArrowRight size={16}/></EnquiryLink>
        </Box>
      </>}
    </>
  );
};

const SiteFooter = () => (
  <Box as={BOX_TAG.FOOTER} className="rdr-footer">
    <Box as={BOX_TAG.DIV} className="rdr-container">
      <Box as={BOX_TAG.DIV} className="rdr-footer__top">
        <Box as={BOX_TAG.DIV} className="rdr-footer__intro">
          <Box as={BOX_TAG.DIV} className="rdr-logo rdr-logo--footer"><SchoolMonogram/><Box as={BOX_TAG.SPAN}><Box as={BOX_TAG.STRONG}>{__("RDR School")}</Box><Box as={BOX_TAG.SMALL}>{__("Rao Deena Ram Vidya Vihar")}</Box></Box></Box>
          <Box as={BOX_TAG.P}>{__("Learning with purpose. Growing with confidence. Serving the community since 2006.")}</Box>
          <Box as={BOX_TAG.SPAN} className="rdr-footer__location"><MapPin size={16}/>{school.address}</Box>
        </Box>
        <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.H3}>{__("Discover")}</Box><AppLink to={SCHOOL_ROUTE.ABOUT}>{__("Our story")}</AppLink><AppLink to={SCHOOL_ROUTE.ACADEMICS}>{__("Academics")}</AppLink><AppLink to={SCHOOL_ROUTE.CAMPUS}>{__("Facilities")}</AppLink><AppLink to={SCHOOL_ROUTE.GALLERY}>{__("School life")}</AppLink></Box>
        <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.H3}>{__("For families")}</Box><AppLink to={SCHOOL_ROUTE.ADMISSIONS}>{__("Admissions")}</AppLink><AppLink to={SCHOOL_ROUTE.RESOURCES}>{__("Resources")}</AppLink><AppLink to={SCHOOL_ROUTE.CONTACT}>{__("Contact")}</AppLink><AppLink to={AUTH_ROUTE.LOGIN}>{__("School workspace")}</AppLink></Box>
        <Box as={BOX_TAG.DIV}><Box as={BOX_TAG.H3}>{__("Transparency")}</Box><AppLink to={SCHOOL_ROUTE.DISCLOSURE}>{__("Mandatory Public Disclosure")}</AppLink><ExternalAnchor href={school.cbseSaras} target="_blank" rel="noreferrer">{__("CBSE affiliation")}</ExternalAnchor><ExternalAnchor href={school.disclosurePdf} target="_blank" rel="noreferrer">{__("2023 disclosure archive")}</ExternalAnchor></Box>
      </Box>
      <Box as={BOX_TAG.DIV} className="rdr-footer__bottom"><Box as={BOX_TAG.SPAN}>© {new Date().getFullYear()}{__("RDR School. All rights reserved.")}</Box><Box as={BOX_TAG.SPAN}>{__("Haluhera · Rewari · Haryana")}</Box></Box>
    </Box>
  </Box>
);

export { SiteHeader };
export { SiteFooter };
