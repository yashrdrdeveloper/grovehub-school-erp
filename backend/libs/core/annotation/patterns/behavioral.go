// Package patterns catalogs GoF design patterns without generating their business behavior.
package patterns

import "grovehub/libs/core/annotation/contract"

var behavioral = []contract.Pattern{
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
