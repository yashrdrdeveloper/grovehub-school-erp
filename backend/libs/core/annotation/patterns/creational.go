// Package patterns catalogs GoF design patterns without generating their business behavior.
package patterns

import "grovehub/libs/core/annotation/contract"

var creational = []contract.Pattern{
    {Name:"AbstractFactory", Family:"Creational", Pattern:"Abstract Factory", Purpose:"Registry of related provider families; explicit factory contract", Decision:"Optional registry"},
    {Name:"Builder", Family:"Creational", Pattern:"Builder", Purpose:"Typed fluent construction", Decision:"Supported"},
    {Name:"FactoryMethod", Family:"Creational", Pattern:"Factory Method", Purpose:"Register constructors for a capability", Decision:"Useful registry"},
    {Name:"Prototype", Family:"Creational", Pattern:"Prototype", Purpose:"Explicit clone/copy contract only", Decision:"Avoid default auto-clone"},
    {Name:"Singleton", Family:"Creational", Pattern:"Singleton", Purpose:"Lifetime constraint in composition root", Decision:"Avoid global instance"},
}
