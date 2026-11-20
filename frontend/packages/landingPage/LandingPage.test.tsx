import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { LandingPage } from './LandingPage';

// Isolate public-site navigation from the unrelated component barrel's
// ESM-only Lodash implementation under Jest's CommonJS test runner.
jest.mock('@Components', () => ({ SchoolMonogram: () => 'RDR' }));

describe('RDR public-site quick actions', () => {
  beforeEach(() => {
    window.scrollTo = jest.fn();
  });

  it('opens the admission enquiry in a full-page dialog without losing the homepage', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    fireEvent.click(screen.getByRole('link',{name:/enquire now/i}));
    await waitFor(() => expect(screen.getByRole('dialog')).toBeTruthy());
    expect(screen.getByRole('heading',{name:'Enquire at RDR School'})).toBeTruthy();
    expect(screen.getByRole('textbox',{name:'Guardian first name'})).toBeTruthy();
    fireEvent.click(screen.getByRole('button',{name:/close and return to the page/i}));
    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
    expect(screen.getByRole('heading',{name:/where curious minds/i})).toBeTruthy();
  });

  it('opens historical school leadership information without redirecting to the legacy site', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    const director = screen.getAllByRole('link',{name:/director.s message/i})[0];
    expect(director).toBeDefined();
    fireEvent.click(director!);
    await waitFor(() => expect(screen.getByRole('dialog')).toBeTruthy());
    expect(screen.getByRole('heading',{name:"Director's message"})).toBeTruthy();
    expect(screen.getByText(/education for the whole child/i)).toBeTruthy();
  });

  it('keeps CBSE public disclosure as a proper route, not a modal', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    fireEvent.click(screen.getAllByRole('link',{name:/cbse disclosure/i})[0]!);
    await waitFor(() => expect(screen.getByRole('heading',{name:/mandatory public disclosure/i})).toBeTruthy());
    expect(screen.queryByRole('dialog')).toBeNull();
  });

  it('validates guardian details within the enquiry sheet and never claims a CRM submission', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    fireEvent.click(screen.getByRole('link',{name:/enquire now/i}));
    await waitFor(() => expect(screen.getByRole('dialog')).toBeTruthy());
    fireEvent.click(screen.getByRole('button',{name:/continue/i}));
    expect(screen.getByRole('alert').textContent).toMatch(/guardian.*first name/i);
    fireEvent.change(screen.getByRole('textbox',{name:'Guardian first name'}),{target:{value:'Parent'}});
    fireEvent.click(screen.getByRole('button',{name:/continue/i}));
    expect(screen.getByRole('alert').textContent).toMatch(/phone number or email/i);
    fireEvent.change(screen.getByRole('textbox',{name:'Phone number'}),{target:{value:'9876543210'}});
    fireEvent.click(screen.getByRole('button',{name:/continue/i}));
    expect(screen.getByText(/a little more detail/i)).toBeTruthy();
    fireEvent.click(screen.getByRole('checkbox'));
    fireEvent.click(screen.getByRole('button',{name:/prepare enquiry/i}));
    expect(screen.getByText(/not been submitted/i)).toBeTruthy();
    const link=screen.getByRole('link',{name:/send enquiry by email/i});
    expect(link.getAttribute('href')).toMatch(/^mailto:/);
  });

  it('offers an accessible close control in the mobile navigation drawer', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    fireEvent.click(screen.getByRole('button',{name:'Open menu'}));
    expect(screen.getByRole('navigation',{name:'Mobile navigation'})).toBeTruthy();
    fireEvent.click(screen.getAllByRole('button',{name:'Close menu'})[0]!);
    await waitFor(() => expect(screen.queryByRole('navigation',{name:'Mobile navigation'})).toBeNull());
  });

  it.each([
    ['header', 'Enquire now'],
    ['hero', 'Enquire for admissions'],
    ['mobile action', 'Admission enquiry'],
    ['homepage prompt', 'Open quick enquiry'],
  ])('opens the same working enquiry sheet from the %s CTA', async (_source, label) => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    fireEvent.click(screen.getByRole('link',{name:label}));
    await waitFor(()=>expect(screen.getByRole('dialog')).toBeTruthy());
    expect(screen.getByRole('heading',{name:'Enquire at RDR School'})).toBeTruthy();
    expect(screen.getByRole('textbox',{name:'Guardian first name'})).toBeTruthy();
    fireEvent.click(screen.getByRole('button',{name:/close and return to the page/i}));
    await waitFor(()=>expect(screen.queryByRole('dialog')).toBeNull());
    expect(screen.getByRole('heading',{name:/where curious minds/i})).toBeTruthy();
  });

  it('opens a live enquiry when Admissions is selected in the footer', async () => {
    const {container}=render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    const footerAdmissions=container.querySelector<HTMLAnchorElement>('.rdr-footer a[href="/admissions"]');
    expect(footerAdmissions).not.toBeNull();
    fireEvent.click(footerAdmissions!);
    await waitFor(()=>expect(screen.getByRole('dialog')).toBeTruthy());
    expect(screen.getByRole('textbox',{name:'Guardian first name'})).toBeTruthy();
  });

  it('opens an enquiry and closes the mobile navigation drawer', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    fireEvent.click(screen.getByRole('button',{name:'Open menu'}));
    fireEvent.click(screen.getByRole('link',{name:/start an enquiry/i}));
    await waitFor(()=>expect(screen.getByRole('dialog')).toBeTruthy());
    expect(screen.queryByRole('navigation',{name:'Mobile navigation'})).toBeNull();
  });

  it('opens an enquiry from the director page without losing that page', async () => {
    render(<MemoryRouter initialEntries={['/about/director']}><LandingPage/></MemoryRouter>);
    expect(screen.getByRole('heading',{name:/director.s message/i})).toBeTruthy();
    fireEvent.click(screen.getByRole('link',{name:/ask the school/i}));
    await waitFor(()=>expect(screen.getByRole('dialog')).toBeTruthy());
    fireEvent.click(screen.getByRole('button',{name:/close and return to the page/i}));
    await waitFor(()=>expect(screen.queryByRole('dialog')).toBeNull());
    expect(screen.getByRole('heading',{name:/director.s message/i})).toBeTruthy();
  });

  it('keeps the standalone admissions page usable from a direct URL', () => {
    render(<MemoryRouter initialEntries={['/admissions']}><LandingPage/></MemoryRouter>);
    expect(screen.getByRole('textbox',{name:'Guardian first name'})).toBeTruthy();
    fireEvent.click(screen.getByRole('button',{name:/continue/i}));
    expect(screen.getByRole('alert').textContent).toMatch(/guardian.*first name/i);
  });


});
