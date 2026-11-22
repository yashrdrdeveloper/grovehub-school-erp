import { _map } from '@Base/utilities/map';
import { TextAreaControl } from '@Components/components/TextAreaControl';
import { SelectControl } from '@Components/components/SelectControl';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { FieldLabel } from '@Components/components/FieldLabel';
import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
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
  let rootTestId: string = ENQUIRY_TEST_IDS.PAGE;
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
    <Box as="div" className={rootClass}
      id={rootId} data-testid={rootTestId}>
      <Box as="div" className="rdr-enquiry__head">
        <Box as="span" className="rdr-eyebrow">LET’S BEGIN</Box>
        <Box as="h3">{heading}</Box>
        <Box as="p">{description}</Box>
      </Box>
      {ready && (
        <Box as="div" className="rdr-enquiry__ready" role="status">
          <Box as="span" className="rdr-enquiry__success"><Check size={22}/></Box>
          <Box as="p">We’ve prepared your enquiry, but it has <Box as="strong">not been submitted</Box> to the school yet. Choose Send enquiry to open your email application.</Box>
          <ExternalAnchor className="rdr-btn rdr-btn--solid" data-testid={ENQUIRY_TEST_IDS.SEND_EMAIL} href={mailto}>Send enquiry by email <Mail size={17}/></ExternalAnchor>
          <Button unstyled className="rdr-link-button" data-testid={ENQUIRY_TEST_IDS.EDIT} type="button" onClick={edit}>Edit details</Button>
        </Box>
      )}
      {!ready && (
        <Form unstyled dataTestId={ENQUIRY_TEST_IDS.FORM} onSubmit={handleSubmit} noValidate>
          <Box as="div" className="rdr-enquiry__step"><Box as="span">STEP {step + 1} OF 2</Box><Box as="span">{stepHeading}</Box></Box>
          {step === 0 && <>
            <FieldLabel>Parent / guardian name <Box as="span">*</Box></FieldLabel>
            <Box as="div" className="rdr-field-grid rdr-field-grid--three">
              <InputControl unstyled aria-label="Guardian first name" dataTestId={ENQUIRY_TEST_IDS.FIRST_NAME} placeholder="First name *" autoComplete="given-name" value={data.firstName} onChange={event => patch('firstName', event.target.value)}/>
              <InputControl unstyled aria-label="Guardian middle name" dataTestId={ENQUIRY_TEST_IDS.MIDDLE_NAME} placeholder="Middle" value={data.middleName} onChange={event => patch('middleName', event.target.value)}/>
              <InputControl unstyled aria-label="Guardian last name" dataTestId={ENQUIRY_TEST_IDS.LAST_NAME} placeholder="Last name" autoComplete="family-name" value={data.lastName} onChange={event => patch('lastName', event.target.value)}/>
            </Box>
            <FieldLabel htmlFor={`${fieldId}-phone`}>Phone number</FieldLabel>
            <InputControl unstyled id={`${fieldId}-phone`} dataTestId={ENQUIRY_TEST_IDS.PHONE} type="tel" placeholder="+91 98765 43210" autoComplete="tel" value={data.phone} onChange={event => patch('phone', event.target.value)}/>
            <FieldLabel htmlFor={`${fieldId}-email`}>Email address</FieldLabel>
            <InputControl unstyled id={`${fieldId}-email`} dataTestId={ENQUIRY_TEST_IDS.EMAIL} type="email" placeholder="parent@example.com" autoComplete="email" value={data.email} onChange={event => patch('email', event.target.value)}/>
            <Box as="small" className="rdr-enquiry__hint">Please provide at least one: phone number or email.</Box>
            <FieldLabel htmlFor={`${fieldId}-type`}>I'm enquiring about <Box as="span">*</Box></FieldLabel>
            <SelectControl id={`${fieldId}-type`} data-testid={ENQUIRY_TEST_IDS.TYPE} value={data.type} onChange={event => patch('type', event.target.value)}>
              {_map(enquiryTypes, type => <option key={type.value} value={type.value}>{type.label}</option>)}
            </SelectControl>
          </>}
          {step !== 0 && <>
            {data.type === 'ADMISSION' && <><FieldLabel htmlFor={`${fieldId}-class`}>Class interested in</FieldLabel>
              <SelectControl id={`${fieldId}-class`} data-testid={ENQUIRY_TEST_IDS.CLASS} value={data.interestedClass} onChange={event => patch('interestedClass', event.target.value)}>
                <option value="">Select a class (optional)</option>
                {_map(['Nursery','KG','I','II','III','IV','V','VI','VII','VIII','IX','X','XI','XII'], (value) => <option key={value} value={value}>{value}</option>)}
              </SelectControl></>}
            <FieldLabel htmlFor={`${fieldId}-message`}>Anything you'd like us to know?</FieldLabel>
            <TextAreaControl id={`${fieldId}-message`} data-testid={ENQUIRY_TEST_IDS.MESSAGE} rows={4} placeholder="Tell us a little about your enquiry..." value={data.message} onChange={event => patch('message', event.target.value)}/>
            <FieldLabel className="rdr-enquiry__consent"><InputControl unstyled dataTestId={ENQUIRY_TEST_IDS.CONSENT} type="checkbox" checked={data.privacy} onChange={event => patch('privacy', event.target.checked)}/><Box as="span">I agree to share these details with RDR School to receive a response to my enquiry.</Box></FieldLabel>
          </>}
          {error && <Box as="p" className="rdr-enquiry__error" role="alert">{error}</Box>}
          <Box as="div" className="rdr-enquiry__actions">
            {step === 1 && <Button unstyled className="rdr-btn rdr-btn--outline" data-testid={ENQUIRY_TEST_IDS.BACK} type="button" onClick={back}>Back</Button>}
            <Button unstyled type="submit" className="rdr-btn rdr-btn--solid" data-testid={ENQUIRY_TEST_IDS.CONTINUE}>{submitLabel} <ArrowRight size={17}/></Button>
          </Box>
        </Form>
      )}
    </Box>
  );
};
