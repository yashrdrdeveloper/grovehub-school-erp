// Package patterns catalogs GoF design patterns without generating their business behavior.
package patterns

import "grovehub/libs/core/annotation/contract"

var structural = []contract.Pattern{
    {Name:"Adapter", Family:"Structural", Pattern:"Adapter", Purpose:"Bind vendor to stable interface", Decision:"Useful registry"},
    {Name:"Bridge", Family:"Structural", Pattern:"Bridge", Purpose:"Map abstraction to independent implementor", Decision:"Optional metadata"},
    {Name:"Composite", Family:"Structural", Pattern:"Composite", Purpose:"Declare child composition, no hidden recursion", Decision:"Optional metadata"},
    {Name:"Decorator", Family:"Structural", Pattern:"Decorator", Purpose:"Explicit typed wrapper order", Decision:"Useful generated wiring"},
    {Name:"Facade", Family:"Structural", Pattern:"Facade", Purpose:"Mark stable public capability API", Decision:"Check only"},
    {Name:"Flyweight", Family:"Structural", Pattern:"Flyweight", Purpose:"Object sharing semantics require manual proof", Decision:"Avoid generation"},
    {Name:"Proxy", Family:"Structural", Pattern:"Proxy", Purpose:"Typed forwarding for client/resilience boundary", Decision:"Useful generation"},
}
