package impl

import (
 "grovehub/services/integration/providers"
 "grovehub/services/integration/providers/mailprovider"
 providerimpl "grovehub/services/integration/providers/impl"
)

type factory struct{}

func NewFactory()mailprovider.Factory{return factory{}}

func(factory)Build(strategies ...mailprovider.Strategy)(mailprovider.Runtime,error){
 catalog:=make([]providers.Provider,0,len(strategies))
 actions:=providerimpl.NewActionRegistry()
 oauth:=providerimpl.NewOAuthRegistry()
 for _,strategy:=range strategies{
  catalog=append(catalog,strategy)
 }
 registry,err:=providerimpl.NewRegistry(catalog...)
 if err!=nil{return mailprovider.Runtime{},err}
 for _,strategy:=range strategies{
  if err:=strategy.Register(actions,oauth);err!=nil{return mailprovider.Runtime{},err}
 }
 return mailprovider.Runtime{Registry:registry,Actions:actions,OAuth:oauth},nil
}
