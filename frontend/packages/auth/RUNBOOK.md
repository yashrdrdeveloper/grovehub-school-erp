# Auth feature
Owns sign-in, forgot password, invitation setup and session runtime. `@Business/auth` owns transport, `packages/auth/pages` owns UI, and the shell only imports the public `@Domain/auth` facade. Session secrets are HttpOnly cookies; do not store them in localStorage, sessionStorage or Redux.
