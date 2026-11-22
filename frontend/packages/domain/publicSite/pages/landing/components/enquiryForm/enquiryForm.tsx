import { ArrowRight, Check, Mail } from '@Components/icons';
import { ENQUIRY_TEST_IDS } from '../../constants/enquiry.testIds';
import { InputControl } from '@Components/components/InputControl';
import { Form } from '@Components/components/Form';
import { enquiryTypes } from '../../constants/school.content';
import type { EnquiryFormComponent } from './types/enquiryForm';

/** Presentational only; validation and transitions belong to useEnquiryForm. */
export const EnquiryForm: EnquiryFormComponent = ({ compact = false, model }) => {
  const { fieldId, data, step, error, ready, mailto, patch, next, submit, edit, back } = model;
  let rootClass = 'rdr-enquiry';
  let rootId: string | undefined = 'enquiry';
  let rootTestId = ENQUIRY_TEST_IDS.PAGE;
  if (compact) {
    rootClass += ' rdr-enquiry--compact';
    rootId = undefined;
    rootTestId = ENQUIRY_TEST_IDS.SHEET;
  }
  let heading = 'Have a question? We’re here.';
  let description = 'A quick conversation is the first step. We’ll help you find the answers.';
  if (ready) {
    heading = 'Your enquiry is ready';
    description = 'Open your email application to send the enquiry to the school.';
  }
  let handleSubmit = submit;
  let stepHeading = 'A little more detail';
  let submitLabel = 'Prepare enquiry';
  if (step === 0) { handleSubmit = next; stepHeading = 'Your details'; submitLabel = 'Continue'; }
  return (
    <div className={rootClass}
      id={rootId} data-testid={rootTestId}>
      <div className="rdr-enquiry__head">
        <span className="rdr-eyebrow">LET’S BEGIN</span>
        <h3>{heading}</h3>
        <p>{description}</p>
      </div>
      {ready && (
        <div className="rdr-enquiry__ready" role="status">
          <span className="rdr-enquiry__success"><Check size={22}/></span>
          <p>We’ve prepared your enquiry, but it has <strong>not been submitted</strong> to the school yet. Choose Send enquiry to open your email application.</p>
          <a className="rdr-btn rdr-btn--solid" data-testid={ENQUIRY_TEST_IDS.SEND_EMAIL} href={mailto}>Send enquiry by email <Mail size={17}/></a>
          <button className="rdr-link-button" data-testid={ENQUIRY_TEST_IDS.EDIT} type="button" onClick={edit}>Edit details</button>
        </div>
      )}
      {!ready && (
        <Form unstyled dataTestId={ENQUIRY_TEST_IDS.FORM} onSubmit={handleSubmit} noValidate>
          <div className="rdr-enquiry__step"><span>STEP {step + 1} OF 2</span><span>{stepHeading}</span></div>
          {step === 0 && <>
            <label>Parent / guardian name <span>*</span></label>
            <div className="rdr-field-grid rdr-field-grid--three">
              <InputControl unstyled aria-label="Guardian first name" dataTestId={ENQUIRY_TEST_IDS.FIRST_NAME} placeholder="First name *" autoComplete="given-name" value={data.firstName} onChange={event => patch('firstName', event.target.value)}/>
              <InputControl unstyled aria-label="Guardian middle name" dataTestId={ENQUIRY_TEST_IDS.MIDDLE_NAME} placeholder="Middle" value={data.middleName} onChange={event => patch('middleName', event.target.value)}/>
              <InputControl unstyled aria-label="Guardian last name" dataTestId={ENQUIRY_TEST_IDS.LAST_NAME} placeholder="Last name" autoComplete="family-name" value={data.lastName} onChange={event => patch('lastName', event.target.value)}/>
            </div>
            <label htmlFor={`${fieldId}-phone`}>Phone number</label>
            <InputControl unstyled id={`${fieldId}-phone`} dataTestId={ENQUIRY_TEST_IDS.PHONE} type="tel" placeholder="+91 98765 43210" autoComplete="tel" value={data.phone} onChange={event => patch('phone', event.target.value)}/>
            <label htmlFor={`${fieldId}-email`}>Email address</label>
            <InputControl unstyled id={`${fieldId}-email`} dataTestId={ENQUIRY_TEST_IDS.EMAIL} type="email" placeholder="parent@example.com" autoComplete="email" value={data.email} onChange={event => patch('email', event.target.value)}/>
            <small className="rdr-enquiry__hint">Please provide at least one: phone number or email.</small>
            <label htmlFor={`${fieldId}-type`}>I'm enquiring about <span>*</span></label>
            <select id={`${fieldId}-type`} data-testid={ENQUIRY_TEST_IDS.TYPE} value={data.type} onChange={event => patch('type', event.target.value)}>
              {enquiryTypes.map(type => <option key={type.value} value={type.value}>{type.label}</option>)}
            </select>
          </>}
          {step !== 0 && <>
            {data.type === 'ADMISSION' && <><label htmlFor={`${fieldId}-class`}>Class interested in</label>
              <select id={`${fieldId}-class`} data-testid={ENQUIRY_TEST_IDS.CLASS} value={data.interestedClass} onChange={event => patch('interestedClass', event.target.value)}>
                <option value="">Select a class (optional)</option>
                {['Nursery','KG','I','II','III','IV','V','VI','VII','VIII','IX','X','XI','XII'].map(value => <option key={value} value={value}>{value}</option>)}
              </select></>}
            <label htmlFor={`${fieldId}-message`}>Anything you'd like us to know?</label>
            <textarea id={`${fieldId}-message`} data-testid={ENQUIRY_TEST_IDS.MESSAGE} rows={4} placeholder="Tell us a little about your enquiry..." value={data.message} onChange={event => patch('message', event.target.value)}/>
            <label className="rdr-enquiry__consent"><InputControl unstyled dataTestId={ENQUIRY_TEST_IDS.CONSENT} type="checkbox" checked={data.privacy} onChange={event => patch('privacy', event.target.checked)}/><span>I agree to share these details with RDR School to receive a response to my enquiry.</span></label>
          </>}
          {error && <p className="rdr-enquiry__error" role="alert">{error}</p>}
          <div className="rdr-enquiry__actions">
            {step === 1 && <button className="rdr-btn rdr-btn--outline" data-testid={ENQUIRY_TEST_IDS.BACK} type="button" onClick={back}>Back</button>}
            <button type="submit" className="rdr-btn rdr-btn--solid" data-testid={ENQUIRY_TEST_IDS.CONTINUE}>{submitLabel} <ArrowRight size={17}/></button>
          </div>
        </Form>
      )}
    </div>
  );
};
