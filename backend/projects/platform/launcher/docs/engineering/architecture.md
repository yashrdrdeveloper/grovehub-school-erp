# Platform Launcher

The root launcher.go exposes lifecycle and operator bootstrap entry points; session/ is the only public session contract boundary. All runtime logic is under internal/src, tests under internal/resources/tests. Session types remain public so code generation cannot leak a Go internal type into another service's mock.

Health routes are declared through @Routes and @GET in controller/health.go. The generator only accepts GET /health and GET /ready for the Launcher, not arbitrary untiered service routes. Test suites may use @Testing, @BeforeEach and @Mock. Generated files are ignored and overlay-mapped from launcher/.generated. The fallback is not a business endpoint.

Bootstrap builds the runtime, records resource cleanup in reverse order, and delegates listening and shutdown to the server interface. Workload route scope policies are separated by owning domain. No request authorization, tenant, cookie, or credential rules are intentionally relaxed.

Existing Core annotation white-box checks of unexported implementation details remain five explicitly named CI exceptions; other adjacent Core tests are rejected.
