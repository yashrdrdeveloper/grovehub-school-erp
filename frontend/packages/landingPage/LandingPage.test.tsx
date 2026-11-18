import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { LandingPage } from './LandingPage';

describe('RDR public-site quick actions', () => {
  beforeEach(() => {
    window.scrollTo = jest.fn();
  });

  it('opens the admission enquiry in a full-page dialog without losing the homepage', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    fireEvent.click(screen.getByRole('link',{name:/enquire now/i}));
    await waitFor(() => expect(screen.getByRole('dialog')).toBeInTheDocument());
    expect(screen.getByRole('heading',{name:'Enquire at RDR School'})).toBeInTheDocument();
    expect(screen.getByRole('textbox',{name:'Guardian first name'})).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button',{name:/close and return to the page/i}));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    expect(screen.getByRole('heading',{name:/where curious minds/i})).toBeInTheDocument();
  });

  it('opens historical school leadership information without redirecting to the legacy site', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    const director = screen.getAllByRole('link',{name:/director.s message/i})[0];
    expect(director).toBeDefined();
    fireEvent.click(director!);
    await waitFor(() => expect(screen.getByRole('dialog')).toBeInTheDocument());
    expect(screen.getByRole('heading',{name:"Director's message"})).toBeInTheDocument();
    expect(screen.getByText(/education for the whole child/i)).toBeInTheDocument();
  });

  it('keeps CBSE public disclosure as a proper route, not a modal', async () => {
    render(<MemoryRouter initialEntries={['/']}><LandingPage/></MemoryRouter>);
    fireEvent.click(screen.getAllByRole('link',{name:/cbse disclosure/i})[0]!);
    await waitFor(() => expect(screen.getByRole('heading',{name:/mandatory public disclosure/i})).toBeInTheDocument());
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });
});
