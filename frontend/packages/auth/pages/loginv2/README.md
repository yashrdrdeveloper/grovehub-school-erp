# GrooveHub LoginV2

A responsive, photo-free sign-in screen with the Soffit WebGL2 gradient restricted to the left-side hero on desktop and a compact banner on mobile.

## Included

- `LoginV2.tsx`: login form, light/dark toggle, password reveal, validation and submit state
- `LoginV2.css`: responsive design with no forced viewport clipping
- `SoffitGradient.tsx`: native WebGL2 fragment shader with one fullscreen triangle, reduced-motion support and CSS fallback

## Integration

Register this page at `/loginv2` in the frontend router. For example:

```tsx
import LoginV2 from '@Auth/pages/loginv2/LoginV2';

<Route
  path="/loginv2"
  element={
    <LoginV2
      onSignIn={({ email, password, remember }) => authService.login({ email, password, remember })}
      onForgotPassword={() => navigate('/forgot-password')}
    />
  }
/>
```

Replace `authService` and navigation callbacks with the real integration. The page deliberately does **not** invent an authentication API.

The connected `develop` branch did not contain a frontend application when the page was added, so this is a standalone page ready for integration rather than an already-routed screen. Do not regard it as runtime-tested until it is installed in the application.

The fragment uses the user-supplied Soffit palette and scene settings and native WebGL2 rendering.
