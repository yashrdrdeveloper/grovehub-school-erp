import { Column } from '../Column';

describe('Column', () => {
  it('builds immutable definitions without modifying previous configurations', () => {
    const base = new Column<{ email: string }>()
      .setId('email')
      .setHeader('Email')
      .setAccessor('email');
    expect(base.build().width).toBeUndefined();
    expect(base.setFillRemainingWidth(true).build().fillRemainingWidth).toBe(true);
    expect(base.build().fillRemainingWidth).toBeUndefined();
    expect(base.setWidth(200).build().width).toBe(200);
    expect(base.build().accessor({ email: 'a@school.com' })).toBe('a@school.com');
    expect(Object.isFrozen(base.build())).toBe(true);
  });

  it('rejects incomplete definitions', () => {
    expect(() => new Column().build()).toThrow();
  });
});
