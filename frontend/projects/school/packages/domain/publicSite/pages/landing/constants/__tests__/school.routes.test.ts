import { SCHOOL_ROUTE } from '../school.routes';

describe('SCHOOL_ROUTE', () => {
  it('retains direct-addressable legacy routes and admission enquiry anchors', () => {
    expect(SCHOOL_ROUTE.HOME).toBe('/');
    expect(SCHOOL_ROUTE.DISCLOSURE).toBe('/mandatory-public-disclosure');
    expect(SCHOOL_ROUTE.ENQUIRY).toBe('/admissions#enquiry');
  });
});
