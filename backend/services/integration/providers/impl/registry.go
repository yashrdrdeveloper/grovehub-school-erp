package impl

import (
 "sort"
 "strings"
 "grovehub/services/integration/providers"
)

type registry struct{items map[string]providers.Provider}

func NewRegistry(items ...providers.Provider)(providers.Registry,error){
 r:=&registry{items:make(map[string]providers.Provider,len(items))}
 for _,item:=range items{
  descriptor:=item.Descriptor()
  id:=strings.ToUpper(strings.TrimSpace(descriptor.ID))
  if id==""{return nil,providers.ErrUnknownProvider}
  if _,exists:=r.items[id];exists{return nil,providers.ErrDuplicateProvider}
  r.items[id]=item
 }
 return r,nil
}
func(r *registry)Find(id string)(providers.Provider,error){
 p,ok:=r.items[strings.ToUpper(strings.TrimSpace(id))]
 if !ok{return nil,providers.ErrUnknownProvider}
 return p,nil
}
func(r *registry)List()[]providers.Descriptor{
 result:=make([]providers.Descriptor,0,len(r.items))
 for _,p:=range r.items{
  d:=p.Descriptor()
  d.Capabilities=append([]providers.Capability(nil),d.Capabilities...)
  result=append(result,d)
 }
 sort.Slice(result,func(i,j int)bool{return result[i].ID<result[j].ID})
 return result
}
