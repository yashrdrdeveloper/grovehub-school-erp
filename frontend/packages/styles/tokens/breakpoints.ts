/** This remains paired with Tailwind's --breakpoint-phone: 45rem in index.css. */
export const STYLE_BREAKPOINT = {
  PHONE_PX: 720,
} as const;

export const STYLE_MEDIA_QUERY = {
  PHONE: '(max-width: 720px)',
} as const;
