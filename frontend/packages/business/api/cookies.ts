// The backend issues opaque HttpOnly access and refresh cookies. Only the
// non-HttpOnly CSRF double-submit cookie is intentionally readable by the UI.
export const readCSRF = (): string => {
  const cookie = document.cookie
    .split('; ')
    .find((entry) => entry.startsWith('gh_csrf='));

  if (!cookie) return '';
  return decodeURIComponent(cookie.slice('gh_csrf='.length));
};
