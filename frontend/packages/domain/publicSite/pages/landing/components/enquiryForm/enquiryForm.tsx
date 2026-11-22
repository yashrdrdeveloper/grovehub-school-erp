import { useId, useState, type FormEvent } from 'react';
import { ArrowRight, Check, Mail } from 'lucide-react';
import { enquiryTypes, school } from '../../config/content';

type EnquiryFormData = {
  firstName:string;
  middleName:string;
  lastName:string;
  phone:string;
  email:string;
  type:string;
  interestedClass:string;
  message:string;
  privacy:boolean;
};

const emptyEnquiry: EnquiryFormData = {
  firstName:'', middleName:'', lastName:'', phone:'', email:'',
  type:'ADMISSION', interestedClass:'', message:'', privacy:false,
};
const EnquiryForm = ({compact=false}:{compact?:boolean}) => {
  const fieldId=useId();
  const [data,setData] = useState<EnquiryFormData>(emptyEnquiry);
  const [step,setStep] = useState(0);
  const [error,setError] = useState('');
  const [ready,setReady] = useState(false);
  const patch = (key:keyof EnquiryFormData,value:string|boolean) => setData(prev=>({...prev,[key]:value}));
  const next = (event:FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    if(!data.firstName.trim()){setError('Please enter the guardian’s first name.');return;}
    if(!data.phone.trim() && !data.email.trim()){setError('Enter a phone number or email address.');return;}
    if(data.email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email.trim())){setError('Enter a valid email address.');return;}
    if(data.phone.trim() && data.phone.replace(/\D/g,'').length<10){setError('Enter a valid phone number.');return;}
    setStep(1);
  };
  const submit = (event:FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    if(!data.privacy){setError('Please confirm you agree to be contacted about your enquiry.');return;}
    // No public Lead CRM endpoint exists in this branch. Avoid pretending that
    // a lead was created: provide a real email draft for the school instead.
    setReady(true);
  };
  const subject = encodeURIComponent(`RDR School — ${enquiryTypes.find(type=>type.value===data.type)?.label ?? 'Enquiry'}`);
  const body = encodeURIComponent([
    `Guardian: ${[data.firstName,data.middleName,data.lastName].filter(Boolean).join(' ')}`,
    `Phone: ${data.phone || 'Not provided'}`,
    `Email: ${data.email || 'Not provided'}`,
    `Enquiry: ${enquiryTypes.find(type=>type.value===data.type)?.label ?? data.type}`,
    `Interested class: ${data.interestedClass || 'Not provided'}`,
    `Message: ${data.message || '—'}`,
  ].join('\n'));
  return (
    <div className={`rdr-enquiry ${compact?'rdr-enquiry--compact':''}`} id={compact?undefined:'enquiry'}>
      <div className="rdr-enquiry__head"><span className="rdr-eyebrow">LET’S BEGIN</span><h3>{ready?'Your enquiry is ready': 'Have a question? We’re here.'}</h3><p>{ready?'Open your email application to send the enquiry to the school.': 'A quick conversation is the first step. We’ll help you find the answers.'}</p></div>
      {ready ? (
        <div className="rdr-enquiry__ready" role="status">
          <span className="rdr-enquiry__success"><Check size={22}/></span>
          <p>We’ve prepared your enquiry, but it has <strong>not been submitted</strong> to the school yet. Choose Send enquiry to open your email application.</p>
          <a className="rdr-btn rdr-btn--solid" href={`mailto:${school.email}?subject=${subject}&body=${body}`}>Send enquiry by email <Mail size={17}/></a>
          <button className="rdr-link-button" type="button" onClick={()=>{setReady(false);setStep(0);}}>Edit details</button>
        </div>
      ) : (
        <form onSubmit={step===0?next:submit} noValidate>
          <div className="rdr-enquiry__step"><span>STEP {step+1} OF 2</span><span>{step===0?'Your details':'A little more detail'}</span></div>
          {step===0 ? <>
            <label>Parent / guardian name <span>*</span></label>
            <div className="rdr-field-grid rdr-field-grid--three">
              <input aria-label="Guardian first name" placeholder="First name *" autoComplete="given-name" value={data.firstName} onChange={e=>patch('firstName',e.target.value)}/>
              <input aria-label="Guardian middle name" placeholder="Middle" value={data.middleName} onChange={e=>patch('middleName',e.target.value)}/>
              <input aria-label="Guardian last name" placeholder="Last name" autoComplete="family-name" value={data.lastName} onChange={e=>patch('lastName',e.target.value)}/>
            </div>
            <label htmlFor={`${fieldId}-phone`}>Phone number</label>
            <input id={`${fieldId}-phone`} type="tel" placeholder="+91 98765 43210" autoComplete="tel" value={data.phone} onChange={e=>patch('phone',e.target.value)}/>
            <label htmlFor={`${fieldId}-email`}>Email address</label>
            <input id={`${fieldId}-email`} type="email" placeholder="parent@example.com" autoComplete="email" value={data.email} onChange={e=>patch('email',e.target.value)}/>
            <small className="rdr-enquiry__hint">Please provide at least one: phone number or email.</small>
            <label htmlFor={`${fieldId}-type`}>I'm enquiring about <span>*</span></label>
            <select id={`${fieldId}-type`} value={data.type} onChange={e=>patch('type',e.target.value)}>
              {enquiryTypes.map(type=><option key={type.value} value={type.value}>{type.label}</option>)}
            </select>
          </>:<>
            {data.type==='ADMISSION' && <><label htmlFor={`${fieldId}-class`}>Class interested in</label><select id={`${fieldId}-class`} value={data.interestedClass} onChange={e=>patch('interestedClass',e.target.value)}>
              <option value="">Select a class (optional)</option>
              {['Nursery','KG','I','II','III','IV','V','VI','VII','VIII','IX','X','XI','XII'].map(value=><option key={value} value={value}>{value}</option>)}
            </select></>}
            <label htmlFor={`${fieldId}-message`}>Anything you'd like us to know?</label>
            <textarea id={`${fieldId}-message`} rows={4} placeholder="Tell us a little about your enquiry..." value={data.message} onChange={e=>patch('message',e.target.value)}/>
            <label className="rdr-enquiry__consent"><input type="checkbox" checked={data.privacy} onChange={e=>patch('privacy',e.target.checked)}/><span>I agree to share these details with RDR School to receive a response to my enquiry.</span></label>
          </>}
          {error && <p className="rdr-enquiry__error" role="alert">{error}</p>}
          <div className="rdr-enquiry__actions">
            {step===1 && <button className="rdr-btn rdr-btn--outline" type="button" onClick={()=>{setStep(0);setError('');}}>Back</button>}
            <button type="submit" className="rdr-btn rdr-btn--solid">{step===0?'Continue':'Prepare enquiry'} <ArrowRight size={17}/></button>
          </div>
        </form>
      )}
    </div>
  );
};

export { EnquiryForm };
