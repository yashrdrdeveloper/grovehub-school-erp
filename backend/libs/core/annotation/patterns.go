// The 23 GoF design patterns are registered as decisions, not promises
// of automatic behavior. Unsafe pattern magic must never be enabled implicitly.
package annotation

var patternInventory = []Pattern{
    {Name:"AbstractFactory", Family:"Creational", Pattern:"Abstract Factory", Purpose:"Registry of related provider families; explicit factory contract", Decision:"Optional registry"},
    {Name:"Builder", Family:"Creational", Pattern:"Builder", Purpose:"Typed fluent construction", Decision:"Supported"},
    {Name:"FactoryMethod", Family:"Creational", Pattern:"Factory Method", Purpose:"Register constructors for a capability", Decision:"Useful registry"},
    {Name:"Prototype", Family:"Creational", Pattern:"Prototype", Purpose:"Explicit clone/copy contract only", Decision:"Avoid default auto-clone"},
    {Name:"Singleton", Family:"Creational", Pattern:"Singleton", Purpose:"Lifetime constraint in composition root", Decision:"Avoid global instance"},
    {Name:"Adapter", Family:"Structural", Pattern:"Adapter", Purpose:"Bind vendor to stable interface", Decision:"Useful registry"},
    {Name:"Bridge", Family:"Structural", Pattern:"Bridge", Purpose:"Map abstraction to independent implementor", Decision:"Optional metadata"},
    {Name:"Composite", Family:"Structural", Pattern:"Composite", Purpose:"Declare child composition, no hidden recursion", Decision:"Optional metadata"},
    {Name:"Decorator", Family:"Structural", Pattern:"Decorator", Purpose:"Explicit typed wrapper order", Decision:"Useful generated wiring"},
    {Name:"Facade", Family:"Structural", Pattern:"Facade", Purpose:"Mark stable public capability API", Decision:"Check only"},
    {Name:"Flyweight", Family:"Structural", Pattern:"Flyweight", Purpose:"Object sharing semantics require manual proof", Decision:"Avoid generation"},
    {Name:"Proxy", Family:"Structural", Pattern:"Proxy", Purpose:"Typed forwarding for client/resilience boundary", Decision:"Useful generation"},
    {Name:"Chain", Family:"Behavioral", Pattern:"Chain of Responsibility", Purpose:"Ordered middleware/handler chain", Decision:"Useful registry"},
    {Name:"Command", Family:"Behavioral", Pattern:"Command", Purpose:"Typed command registry and dispatch", Decision:"Optional registry"},
    {Name:"Interpreter", Family:"Behavioral", Pattern:"Interpreter", Purpose:"Registry of parsing grammar components", Decision:"Avoid auto-parser"},
    {Name:"Iterator", Family:"Behavioral", Pattern:"Iterator", Purpose:"Typed iteration interface only when needed", Decision:"Avoid generation"},
    {Name:"Mediator", Family:"Behavioral", Pattern:"Mediator", Purpose:"Explicit event/notification coordination", Decision:"Optional registry"},
    {Name:"Memento", Family:"Behavioral", Pattern:"Memento", Purpose:"Snapshot semantics with sensitive-state controls", Decision:"Avoid automation"},
    {Name:"Observer", Family:"Behavioral", Pattern:"Observer", Purpose:"Listener registration around event contract", Decision:"Useful registry"},
    {Name:"State", Family:"Behavioral", Pattern:"State", Purpose:"Register state transitions with guards", Decision:"Check only initially"},
    {Name:"Strategy", Family:"Behavioral", Pattern:"Strategy", Purpose:"Provider capability registration + selection", Decision:"Useful registry"},
    {Name:"TemplateMethod", Family:"Behavioral", Pattern:"Template Method", Purpose:"Common skeleton remains an explicit Go method", Decision:"Avoid inheritance emulation"},
    {Name:"Visitor", Family:"Behavioral", Pattern:"Visitor", Purpose:"Typed dispatch with explicit exhaustive cases", Decision:"Optional, rare"},
}
