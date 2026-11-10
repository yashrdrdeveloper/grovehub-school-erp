// Phase 1 configuration, profiles, typed environment and secrets annotation metadata. No runtime behavior.
package annotation

var configurationUseCases = []UseCase{
    {Name:"Configuration",Syntax:"@Configuration",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Typed config struct registered for startup loading",Mode:"REG/GEN",Status:Pilot,Roadmap:"P5",Target:TypeTarget},
    {Name:"ConfigPrefix",Syntax:"@ConfigPrefix",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Namespace environment properties by service",Mode:"GEN",Status:Proposed,Roadmap:"P5",Target:TypeTarget},
    {Name:"Env",Syntax:"@Env",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Bind explicitly named environment variable",Mode:"GEN",Status:Pilot,Roadmap:"P5",Target:FieldTarget},
    {Name:"Default",Syntax:"@Default",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Typed default with explicit precedence",Mode:"GEN",Status:Pilot,Roadmap:"P5",Target:FieldTarget},
    {Name:"Required",Syntax:"@Required",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Configuration missing error at startup",Mode:"CHECK",Status:Pilot,Roadmap:"P5",Target:FieldTarget},
    {Name:"ValidateConfig",Syntax:"@ValidateConfig",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Run technical config validation before module startup",Mode:"WRAP",Status:Proposed,Roadmap:"P5",Target:TypeTarget},
    {Name:"Secret",Syntax:"@Secret",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Mark sensitive config redaction and source restrictions",Mode:"CHECK",Status:Pilot,Roadmap:"P5",Target:FieldTarget},
    {Name:"ConfigSource",Syntax:"@ConfigSource",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Describe allowed source (env, file, secret store)",Mode:"REG",Status:Optional,Roadmap:"P5",Target:TypeTarget},
    {Name:"Profile",Syntax:"@Profile",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Explicit development/test/production selection",Mode:"REG/CHECK",Status:Optional,Roadmap:"P5",Target:TypeTarget},
    {Name:"Reloadable",Syntax:"@Reloadable",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Controlled reload with atomic switch and validation",Mode:"WRAP",Status:Optional,Roadmap:"P5",Target:TypeTarget},
    {Name:"FeatureFlag",Syntax:"@FeatureFlag",Category:"Configuration, profiles, typed environment and secrets",Purpose:"Typed feature selection; not a security boundary",Mode:"REG",Status:Optional,Roadmap:"P5",Target:FieldTarget},
}
