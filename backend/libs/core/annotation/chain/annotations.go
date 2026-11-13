// Package chain catalogs the typed service chain marker.
// The marker is checked by the annotation scanner; the service runner executes
// the behavior explicitly, without reflection or invisible interception.
package chain

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"ChainOfResponsibility",Syntax:"@ChainOfResponsibility",Category:"Service pipeline",Purpose:"Mark a method implemented through the typed core service chain runner",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P4",Target:contract.MethodTarget},
}
