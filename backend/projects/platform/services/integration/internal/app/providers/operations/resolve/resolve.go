package resolve

import (
    "strings"
    providers "grovehub/services/integration/internal/app/providers"
    helpers "grovehub/services/integration/internal/app/providers/helpers"
)

func Find(items map[string]providers.Provider, id string) (providers.Provider,error) {
    provider,ok:=items[strings.ToUpper(strings.TrimSpace(id))]
    if !ok {return nil,helpers.ErrUnknownProvider}
    return provider,nil
}
