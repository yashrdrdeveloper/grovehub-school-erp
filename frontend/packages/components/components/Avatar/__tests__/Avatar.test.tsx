import { fireEvent, render, screen } from '@testing-library/react';

import { Avatar } from '../Avatar';
import { getAvatarHue } from '../helpers/avatar.getHue';

describe('Avatar', () => {
  it('applies shared color tokens with a stable per-name hue', () => {
    render(<Avatar name="Alex Example" dataTestId="avatar" compact />);

    const avatar = screen.getByTestId('avatar');
    expect(avatar.classList.contains('gh-avatar-colors')).toBe(true);
    expect(avatar.style.getPropertyValue('--avatar-hue')).toBe(
      String(getAvatarHue('Alex Example')),
    );
    expect(avatar.getAttribute('data-compact')).toBe('true');
    expect(avatar.textContent).toBe('AE');
  });

  it('retains the fallback when the image fails', () => {
    render(<Avatar name="Alex Example" src="/avatar.png" dataTestId="avatar" />);
    const avatar = screen.getByTestId('avatar');
    const image = avatar.querySelector('img');
    if (!image) throw new Error('Expected avatar image');

    fireEvent.error(image);

    expect(avatar.querySelector('img')).toBeNull();
    expect(avatar.textContent).toBe('AE');
  });
});
