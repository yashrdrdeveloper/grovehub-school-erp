import { _find } from '@Base/utilities/find';
import { _keys } from '@Base/utilities/keys';
import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { useEffect, useState, type MouseEvent } from 'react';
import { useLocation, useNavigate } from '@Components/router';
import type { SchoolPanel } from '../../components/publicSheet';
import { HomePage } from '../../components/homeContent';
import { pageContent } from '../../registry/schoolPages.registry';
import { pathPanel } from '../../constants/school.panels';
import { pageMetadata } from '../../constants/school.seo';
import type { UseLandingPage } from './types/useLandingPage';

/** Owns route, enquiry and SEO effects. The public page has no local state. */
export const useLandingPage: UseLandingPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [enquiryOpen, setEnquiryOpen] = useState(false);
  const state = location.state as { background?: string; panel?: SchoolPanel } | null;
  const panel: SchoolPanel | null = (enquiryOpen && 'enquiry') || null;
  const backgroundPath = location.pathname;
  const CurrentPage = pageContent[backgroundPath] ?? HomePage;

  const closePanel = () => {
    if (enquiryOpen) {
      setEnquiryOpen(false);
      return;
    }
    navigate(state?.background ?? SCHOOL_ROUTE.HOME, { replace: true });
  };
  const selectPanel = (next: SchoolPanel) => {
    if (next === 'enquiry') {
      setEnquiryOpen(true);
      return;
    }
    setEnquiryOpen(false);
    const path = _find(_keys(pathPanel), (key) => pathPanel[key] === next) ?? SCHOOL_ROUTE.ABOUT;
    navigate(path, { replace: true, state: { background: backgroundPath, panel: next } });
  };

  // Public navigation opens full, addressable pages. The enquiry CTA
  // continues to use the existing accessible sheet until Lead CRM is ready.
  const handleInternalLink = (event: MouseEvent<HTMLDivElement>) => {
    if (
      event.defaultPrevented ||
      event.button !== 0 ||
      event.metaKey ||
      event.ctrlKey ||
      event.shiftKey ||
      event.altKey
    )
      return;
    if (!(event.target instanceof Element)) return;
    const anchor = event.target.closest<HTMLAnchorElement>('a[href]');
    if (!anchor || anchor.hasAttribute('download') || (anchor.target && anchor.target !== '_self'))
      return;
    if (anchor.dataset.rdrEnquiryAction === 'true') return;
    const url = new URL(anchor.href, window.location.href);
    if (url.origin !== window.location.origin) return;
    if (url.pathname === SCHOOL_ROUTE.ADMISSIONS && url.hash === '#enquiry') {
      event.preventDefault();
      setEnquiryOpen(true);
    }
  };
  useEffect(() => {
    if (panel) return;
    if (location.hash) {
      window.requestAnimationFrame(() =>
        document
          .getElementById(location.hash.slice(1))
          ?.scrollIntoView({ behavior: 'smooth', block: 'start' }),
      );
    } else window.scrollTo({ top: 0, behavior: 'instant' });
  }, [backgroundPath, location.hash, panel]);
  useEffect(() => {
    const fallback =
      (location.pathname.startsWith('/about/') && pageMetadata['/about']) || pageMetadata['/'];
    const details = pageMetadata[location.pathname] ?? fallback;
    if (!details) return;
    document.title = details.title;
    const description = document.querySelector<HTMLMetaElement>('meta[name="description"]');
    if (description) description.content = details.description;
    const existingCanonical = document.querySelector<HTMLLinkElement>('link[rel="canonical"]');
    const canonical = existingCanonical || document.createElement('link');
    if (!existingCanonical) {
      canonical.rel = 'canonical';
      document.head.appendChild(canonical);
    }
    canonical.href = `https://rdrschool.com${location.pathname}`;
  }, [location.pathname]);

  const openEnquiry = () => setEnquiryOpen(true);
  return { CurrentPage, panel, closePanel, selectPanel, handleInternalLink, openEnquiry };
};
