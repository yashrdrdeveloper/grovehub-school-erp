package catalog

import (
    "sort"
    providers "grovehub/services/integration/internal/app/providers"
    dto "grovehub/services/integration/internal/app/providers/dto"
)

func List(items map[string]providers.Provider) []dto.Descriptor {
    result:=make([]dto.Descriptor,0,len(items))
    for _,provider:=range items {
        descriptor:=provider.Descriptor()
        descriptor.Capabilities=append([]dto.Capability(nil),descriptor.Capabilities...)
        result=append(result,descriptor)
    }
    sort.Slice(result,func(i,j int)bool{return result[i].ID<result[j].ID})
    return result
}
