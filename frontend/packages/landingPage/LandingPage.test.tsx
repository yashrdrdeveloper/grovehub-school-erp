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
});
