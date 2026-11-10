// Package configuration owns the annotation metadata for this backend responsibility.
package configuration

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Configuration",Syntax:"@Configuration",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Typed config struct registered for startup loading",Mode:"REG/GEN",Status:contract.Pilot,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"ConfigPrefix",Syntax:"@ConfigPrefix",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Namespace environment properties by service",Mode:"GEN",Status:contract.Proposed,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"Env",Syntax:"@Env",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Bind explicitly named environment variable",Mode:"GEN",Status:contract.Pilot,Roadmap:"P5",Target:contract.FieldTarget},
    {Name:"Default",Syntax:"@Default",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Typed default with explicit precedence",Mode:"GEN",Status:contract.Pilot,Roadmap:"P5",Target:contract.FieldTarget},
    {Name:"Required",Syntax:"@Required",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Configuration missing error at startup",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P5",Target:contract.FieldTarget},
    {Name:"ValidateConfig",Syntax:"@ValidateConfig",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Run technical config validation before module startup",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"Secret",Syntax:"@Secret",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Mark sensitive config redaction and source restrictions",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P5",Target:contract.FieldTarget},
    {Name:"ConfigSource",Syntax:"@ConfigSource",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Describe allowed source (env, file, secret store)",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"Profile",Syntax:"@Profile",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Explicit development/test/production selection",Mode:"REG/CHECK",Status:contract.Optional,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"Reloadable",Syntax:"@Reloadable",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Controlled reload with atomic switch and validation",Mode:"WRAP",Status:contract.Optional,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"FeatureFlag",Syntax:"@FeatureFlag",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Typed feature selection; not a security boundary",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.FieldTarget},
}
