# Platform Launcher

The root launcher.go exposes lifecycle and operator bootstrap entry points; session/ is the only public session contract boundary. All runtime logic is under internal/src, tests under internal/resources/tests. Session types remain public so code generation cannot leak a Go internal type into another service's mock.

Health routes are declared through @Routes and @GET in controller/health.go. The generator only accepts GET /health and GET /ready for the Launcher, not arbitrary untiered service routes. Test suites may use @Testing, @BeforeEach and @Mock. Generated files are ignored and overlay-mapped from launcher/.generated. The fallback is not a business endpoint.

Bootstrap builds the runtime, records resource cleanup in reverse order, and delegates listening and shutdown to the server interface. Workload route scope policies are separated by owning domain. No request authorization, tenant, cookie, or credential rules are intentionally relaxed.

Existing Core annotation white-box checks of unexported implementation details remain five explicitly named CI exceptions; other adjacent Core tests are rejected.


## Annotation and coverage policy

- Launcher configuration uses the existing `@Configuration`, `@Env`, `@Default`, and `@Secret` declarations with the Core typed configuration loader. Invalid tenant-required flags fail startup rather than silently disabling enforcement.
- Keep `@Routes` and HTTP method annotations on the controller interface; keep `@Service` and `@Strategy` declarations on the corresponding concrete implementations. Unsupported catalog-only annotations must not replace explicit security or lifecycle enforcement.
- All Launcher unit suites use `@Testing`; use generated `@Mock` interfaces for external dependencies and `@BeforeEach` for supported setup. Keep real cryptographic checks in security tests rather than mocking the policy being tested.
- Production composition is defined by `bootstrap/impl.DefaultDependencies`. The builder owns reverse-order rollback, database pool ownership, worker lifecycle and server draining. Unit tests supply alternative transport adapters to exercise genuine failure paths without a live broker or database.
- Run `make unit-coverage package=launcher` for exact 100% statement coverage. The existing coverage runner materializes generated source, includes all Launcher production packages and the shared Platform-private session capability, and rejects every uncovered statement. Generated routes and facade methods are included.
