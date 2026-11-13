// Package strategy catalogs typed service strategy declarations.
package strategy

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Strategy",Syntax:"@Strategy",Category:"Service pipeline",Purpose:"Mark typed strategy implementation explicitly composed by the service",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P4",Target:contract.TypeTarget},
    {Name:"Strategy",Syntax:"@Strategy",Category:"Service pipeline",Purpose:"Mark the typed strategy interface implemented in strategies/impl",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P4",Target:contract.InterfaceTarget},
}
