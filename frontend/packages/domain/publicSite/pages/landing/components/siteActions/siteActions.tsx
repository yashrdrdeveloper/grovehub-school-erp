import { SCHOOL_ROUTE } from '../../constants/school.routes';
import { createContext, useContext, type ReactNode } from 'react';
import { Link } from 'react-router-dom';

/** Explicit enquiry actions do not depend on catching React Router links
 * at the page root. Their href remains useful for opening a new tab or
 * loading the fully accessible admissions page without JavaScript. */
const EnquiryActionContext = createContext<(() => void) | null>(null);

const EnquiryLink = ({className,children,onClick}:{
  className?:string;children:ReactNode;onClick?:()=>void;
}) => {
  const openEnquiry=useContext(EnquiryActionContext);
  return <Link to={SCHOOL_ROUTE.ENQUIRY} className={className} data-rdr-enquiry-action="true"
    onClick={event=>{
      onClick?.();
      if(event.defaultPrevented || event.button!==0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey)return;
      if(!openEnquiry)return;
      event.preventDefault();
      openEnquiry();
    }}>{children}</Link>;
};

export { EnquiryActionContext };
export { EnquiryLink };
