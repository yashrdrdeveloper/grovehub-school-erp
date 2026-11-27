import { _map } from '@Base/utilities/map';
import { TextAreaControl } from '@Components/components/TextAreaControl';
import { SelectControl } from '@Components/components/SelectControl';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { FieldLabel } from '@Components/components/FieldLabel';
import { Button } from '@Components/components/Button';
import { Box , BOX_TAG} from '@Components/components/Box';
import { ArrowRight, Check, Mail } from '@Components/icons';
import { ENQUIRY_TEST_IDS } from '../../constants/enquiry.testIds';
import { InputControl } from '@Components/components/InputControl';
import { Form } from '@Components/components/Form';
import { enquiryTypes } from '../../constants/school.content';
import { getEnquiryFormView } from './helpers/enquiryForm.getView';
import type { EnquiryFormComponent } from './types/enquiryForm';

/** Presentational only; validation and transitions belong to useEnquiryForm. */
export const EnquiryForm: EnquiryFormComponent = ({ compact = false, model }) => {
  const { fieldId, data, step, error, ready, mailto, patch, edit, back } = model;
  const { rootClass, rootId, rootTestId, heading, description, handleSubmit,
    stepHeading, submitLabel } = getEnquiryFormView(compact, model);
  return (
    <Box as={BOX_TAG.DIV} className={rootClass}
      id={rootId} data-testid={rootTestId}>
      <Box as={BOX_TAG.DIV} className="rdr-enquiry__head">
        <Box as={BOX_TAG.SPAN} className="rdr-eyebrow">{__("LET’S BEGIN")}</Box>
        <Box as={BOX_TAG.H3}>{heading}</Box>
        <Box as={BOX_TAG.P}>{description}</Box>
      </Box>
      {ready && (
        <Box as={BOX_TAG.DIV} className="rdr-enquiry__ready" role="status">
          <Box as={BOX_TAG.SPAN} className="rdr-enquiry__success"><Check size={22}/></Box>
          <Box as={BOX_TAG.P}>{__("We’ve prepared your enquiry, but it has ")}<Box as={BOX_TAG.STRONG}>{__("not been submitted")}</Box>{__("to the school yet. Choose Send enquiry to open your email application.")}</Box>
          <ExternalAnchor className="rdr-btn rdr-btn--solid" data-testid={ENQUIRY_TEST_IDS.SEND_EMAIL} href={mailto}>{__("Send enquiry by email ")}<Mail size={17}/></ExternalAnchor>
          <Button unstyled className="rdr-link-button" data-testid={ENQUIRY_TEST_IDS.EDIT} type="button" onClick={edit}>{__("Edit details")}</Button>
        </Box>
      )}
      {!ready && (
        <Form unstyled dataTestId={ENQUIRY_TEST_IDS.FORM} onSubmit={handleSubmit} noValidate>
          <Box as={BOX_TAG.DIV} className="rdr-enquiry__step"><Box as={BOX_TAG.SPAN}>{__("STEP ")}{step + 1}{__("OF 2")}</Box><Box as={BOX_TAG.SPAN}>{stepHeading}</Box></Box>
          {step === 0 && <>
            <FieldLabel>{__("Parent / guardian name ")}<Box as={BOX_TAG.SPAN}>*</Box></FieldLabel>
            <Box as={BOX_TAG.DIV} className="rdr-field-grid rdr-field-grid--three">
              <InputControl unstyled aria-label={__("Guardian first name")} dataTestId={ENQUIRY_TEST_IDS.FIRST_NAME} placeholder={__("First name *")} autoComplete="given-name" value={data.firstName} onChange={event => patch('firstName', event.target.value)}/>
              <InputControl unstyled aria-label={__("Guardian middle name")} dataTestId={ENQUIRY_TEST_IDS.MIDDLE_NAME} placeholder={__("Middle")} value={data.middleName} onChange={event => patch('middleName', event.target.value)}/>
              <InputControl unstyled aria-label={__("Guardian last name")} dataTestId={ENQUIRY_TEST_IDS.LAST_NAME} placeholder={__("Last name")} autoComplete="family-name" value={data.lastName} onChange={event => patch('lastName', event.target.value)}/>
            </Box>
            <FieldLabel htmlFor={`${fieldId}-phone`}>{__("Phone number")}</FieldLabel>
            <InputControl unstyled id={`${fieldId}-phone`} dataTestId={ENQUIRY_TEST_IDS.PHONE} type="tel" placeholder="+91 98765 43210" autoComplete="tel" value={data.phone} onChange={event => patch('phone', event.target.value)}/>
            <FieldLabel htmlFor={`${fieldId}-email`}>{__("Email address")}</FieldLabel>
            <InputControl unstyled id={`${fieldId}-email`} dataTestId={ENQUIRY_TEST_IDS.EMAIL} type="email" placeholder={__("parent@example.com")} autoComplete="email" value={data.email} onChange={event => patch('email', event.target.value)}/>
            <Box as={BOX_TAG.SMALL} className="rdr-enquiry__hint">{__("Please provide at least one: phone number or email.")}</Box>
            <FieldLabel htmlFor={`${fieldId}-type`}>{__("I'm enquiring about ")}<Box as={BOX_TAG.SPAN}>*</Box></FieldLabel>
            <SelectControl id={`${fieldId}-type`} data-testid={ENQUIRY_TEST_IDS.TYPE} value={data.type} onChange={event => patch('type', event.target.value)}>
              {_map(enquiryTypes, type => <option key={type.value} value={type.value}>{type.label}</option>)}
            </SelectControl>
          </>}
          {step !== 0 && <>
            {data.type === 'ADMISSION' && <><FieldLabel htmlFor={`${fieldId}-class`}>{__("Class interested in")}</FieldLabel>
              <SelectControl id={`${fieldId}-class`} data-testid={ENQUIRY_TEST_IDS.CLASS} value={data.interestedClass} onChange={event => patch('interestedClass', event.target.value)}>
                <option value="">{__("Select a class (optional)")}</option>
                {_map(['Nursery','KG','I','II','III','IV','V','VI','VII','VIII','IX','X','XI','XII'], (value) => <option key={value} value={value}>{value}</option>)}
              </SelectControl></>}
            <FieldLabel htmlFor={`${fieldId}-message`}>{__("Anything you'd like us to know?")}</FieldLabel>
            <TextAreaControl id={`${fieldId}-message`} data-testid={ENQUIRY_TEST_IDS.MESSAGE} rows={4} placeholder={__("Tell us a little about your enquiry...")} value={data.message} onChange={event => patch('message', event.target.value)}/>
            <FieldLabel className="rdr-enquiry__consent"><InputControl unstyled dataTestId={ENQUIRY_TEST_IDS.CONSENT} type="checkbox" checked={data.privacy} onChange={event => patch('privacy', event.target.checked)}/><Box as={BOX_TAG.SPAN}>{__("I agree to share these details with RDR School to receive a response to my enquiry.")}</Box></FieldLabel>
          </>}
          {error && <Box as={BOX_TAG.P} className="rdr-enquiry__error" role="alert">{error}</Box>}
          <Box as={BOX_TAG.DIV} className="rdr-enquiry__actions">
            {step === 1 && <Button unstyled className="rdr-btn rdr-btn--outline" data-testid={ENQUIRY_TEST_IDS.BACK} type="button" onClick={back}>{__("Back")}</Button>}
            <Button unstyled type="submit" className="rdr-btn rdr-btn--solid" data-testid={ENQUIRY_TEST_IDS.CONTINUE}>{submitLabel} <ArrowRight size={17}/></Button>
          </Box>
        </Form>
      )}
    </Box>
  );
};
