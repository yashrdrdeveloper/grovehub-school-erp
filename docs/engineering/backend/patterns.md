# Patterns, Operations and Design Decisions

Use patterns only when a **tested need** makes plain Go less maintainable. This is a decision guide, not an instruction to instantiate every GoF pattern in every service.

## 1. Operation structure

```text
services/<resource>/
├── <resource>.go                      # narrow use-case interface
├── impl/<resource>.go                 # @Service dependency composition only
├── helpers/                           # shared by 2+ operations
└── operations/<verb>/
    ├── <verb>.go                      # entry and optional chain
    ├── flows/<verb>.go                # request-scoped state of multi-step flow
    ├── helpers/<purpose>.go           # operation-private computation
    └── strategies/
        ├── strategy.go                # typed step/algorithm interface
        └── impl/<verb>.go             # actual behavior
```

A simple read needs no chain or flow. A multistep command can use **Prepare -> Load -> Authorize -> Persist -> Emit**, with documented short-circuit/error behavior. Constructors should assemble reusable dependencies; request-scoped flow state belongs to each invocation, never the singleton service.

## 2. Pattern reference

| Pattern | Use | Avoid |
|---|---|---|
| Strategy | Replaceable algorithm/provider | One-line forwarding interface |
| Factory | Nontrivial valid dependency family | Wrapping one constructor |
| Abstract Factory | Interchangeable implementation families | Universal all-purpose interface |
| Adapter | Translate vendor protocol/SDK | Merely renaming equal types |
| Facade | Narrow stable public client | Exposing private graph |
| Registry | Stable ID/capability lookup | Global mutable service locator |
| Command | Typed executable capability action | Giant string switch |
| Chain of Responsibility | Ordered short-circuiting steps | Single straightforward action |
| Decorator | Explicit cross-cutting wrapper | Hidden double-publish |
| Repository | Resource data access/queries | Generic `Save(any)` |
| UnitOfWork | Atomic multi-repo local DB writes | Pretend remote APIs are atomic |
| Outbox | DB mutation + durable event intent | Best-effort telemetry |
| Observer / Pub-sub | Async decoupling | Synchronous consistency requirement |
| State machine | Allowed lifecycle transitions | Status checks spread everywhere |
| Saga | Distributed compensation | One local SQL transaction |
| CQRS | Proven read/write scaling differences | Premature mirrored models |
| Specification | Reusable typed domain predicate | Caller-supplied SQL |
| Builder | Complex valid construction | Hide required fields |
| Proxy | Typed client delegation | Assuming gateway middleware ran |

## 3. The GoF taxonomy, with Go caution

**Creational:** Factory Method, Abstract Factory, Builder, Prototype, Singleton. Prefer explicit constructors over global state.

**Structural:** Adapter, Bridge, Composite, Decorator, Facade, Flyweight, Proxy. Prefer typed interface composition to inheritance simulation.

**Behavioral:** Chain of Responsibility, Command, Interpreter, Iterator, Mediator, Memento, Observer, State, Strategy, Template Method, Visitor. Prefer small typed functions/interfaces when sufficient.

Cataloging a design-pattern annotation does not implement it.

## 4. Operation rules

1. One operation owns one command/read behavior.
2. Each step has explicit input/output, error and cancellation behavior.
3. Construct strategies once at startup when stateless.
4. Put each request's mutable flow state in its own invocation.
5. Use typed repositories with transaction-bound context, not HTTP or SQL in domain helpers.
6. Emit events exactly once per successful meaningful transition.
7. Promote helpers one ownership level only after stable shared semantics appear.
8. Keep vendor switches out of business services.
9. Avoid duplicate service/strategy packages with no real responsibility.
10. Document owner-specific algorithms in the module engineering docs, not here.

## 5. Code smells and correction

**Too many wrappers:** Collapse packages that only delegate without maintaining policy. **Too few boundaries:** Split HTTP, business, SQL and provider SDK concerns. **Shared helper dumpster:** Place each helper with its resource/operation; promote to core only for truly generic use. **Untracked worker:** Add lifecycle cancellation and joining. **Annotation-only transaction:** Add a real UnitOfWork and rollback test. **Routes as permission:** Add real auth/tenant middleware and ownership checks.
