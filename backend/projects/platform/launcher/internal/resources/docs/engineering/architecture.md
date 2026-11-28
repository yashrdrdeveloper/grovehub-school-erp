# Platform Launcher

The root launcher.go exposes lifecycle and operator bootstrap entry points; shared security contracts are imported from the standalone `grovehub/libs/platform/session` capability. All Launcher runtime logic is under internal/src, tests under internal/resources/tests. Session library internals are private to their owning capability.

Health routes are declared through @Routes and @GET in controller/health.go. The generator only accepts GET /health and GET /ready for the Launcher, not arbitrary untiered service routes. Test suites may use @Testing, @BeforeEach and @Mock. Generated files are ignored and overlay-mapped from launcher/.generated. The fallback is not a business endpoint.

Bootstrap builds the runtime, records resource cleanup in reverse order, and delegates listening and shutdown to the server interface. Workload route scope policies are separated by owning domain. No request authorization, tenant, cookie, or credential rules are intentionally relaxed.

Existing Core annotation white-box checks of unexported implementation details remain five explicitly named CI exceptions; other adjacent Core tests are rejected.


## Annotation and coverage policy

- Launcher configuration uses the existing `@Configuration`, `@Env`, `@Default`, and `@Secret` declarations with the Core typed configuration loader. Invalid tenant-required flags fail startup rather than silently disabling enforcement.
- Keep `@Routes` and HTTP method annotations on the controller interface; keep `@Service` and `@Strategy` declarations on the corresponding concrete implementations. Unsupported catalog-only annotations must not replace explicit security or lifecycle enforcement.
- All Launcher unit suites use `@Testing`; use generated `@Mock` interfaces for external dependencies and `@BeforeEach` for supported setup. Keep real cryptographic checks in security tests rather than mocking the policy being tested.
- Production composition is defined by `bootstrap/impl.DefaultDependencies`. The builder owns reverse-order rollback, database pool ownership, worker lifecycle and server draining. Unit tests supply alternative transport adapters to exercise genuine failure paths without a live broker or database.
- Run `make unit-coverage package=launcher` for exact 100% statement coverage. The existing coverage runner materializes generated source, includes all Launcher-owned production packages and rejects every uncovered statement. The standalone Platform library has its own strict coverage gate. Generated routes and facade methods are included. Run `make platform-lib-coverage` independently for the shared library.

## Operator command ownership

The executable at `cmd/main.go` delegates to `launcher.Execute`. The launcher composes service bootstraps in the established migration order, rejects duplicate command names, and dispatches a command without constructing server runtimes. Identity and Audit own their command strategies and the other service bootstraps expose their migrations. The launcher owns generic serve and migration strategies and connection lifecycle. New service commands are registered in the owning service bootstrap; no launcher dispatch change is required.
