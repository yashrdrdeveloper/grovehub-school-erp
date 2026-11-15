// Package impl assembles Connection use-case operations, following Foundation.
package impl

import (
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
    mailcontracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
    oauth "grovehub/services/integration/internal/app/providers/mailprovider/operations/oauth"
    eventproducer "grovehub/services/integration/internal/app/events/producers/connection"
    "grovehub/services/integration/internal/app/externalService/mail/gmail"
    "grovehub/services/integration/internal/app/externalService/mail/outlook"
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    watchrepo "grovehub/services/integration/internal/app/repository/watch"
    connectionservices "grovehub/services/integration/internal/app/services/connection"
    helpers "grovehub/services/integration/internal/app/services/connection/helpers"
    connectionmapper "grovehub/services/integration/internal/app/mappers/connection"
    connectionvalidator "grovehub/services/integration/internal/app/validators/connection"

    action "grovehub/services/integration/internal/app/services/connection/operations/action"
    actionimpl "grovehub/services/integration/internal/app/services/connection/operations/action/strategies/impl"
    create "grovehub/services/integration/internal/app/services/connection/operations/create"
    createimpl "grovehub/services/integration/internal/app/services/connection/operations/create/strategies/impl"
    disconnect "grovehub/services/integration/internal/app/services/connection/operations/disconnect"
    disconnectimpl "grovehub/services/integration/internal/app/services/connection/operations/disconnect/strategies/impl"
    get "grovehub/services/integration/internal/app/services/connection/operations/get"
    getimpl "grovehub/services/integration/internal/app/services/connection/operations/get/strategies/impl"
    list "grovehub/services/integration/internal/app/services/connection/operations/list"
    listimpl "grovehub/services/integration/internal/app/services/connection/operations/list/strategies/impl"
    oauth "grovehub/services/integration/internal/app/services/connection/operations/oauth"
    oauthimpl "grovehub/services/integration/internal/app/services/connection/operations/oauth/strategies/impl"
    test "grovehub/services/integration/internal/app/services/connection/operations/test"
    testimpl "grovehub/services/integration/internal/app/services/connection/operations/test/strategies/impl"
    watch "grovehub/services/integration/internal/app/services/connection/operations/watch"
    watchimpl "grovehub/services/integration/internal/app/services/connection/operations/watch/strategies/impl"
)

var (
    ErrInvalidConnection = helpers.ErrInvalidConnection
    ErrConnectionNotReady = helpers.ErrConnectionNotReady
    ErrOAuthUnavailable = helpers.ErrOAuthUnavailable
    ErrOAuthInvalid = helpers.ErrOAuthInvalid
    ErrConnectionInactive = helpers.ErrConnectionInactive
    ErrPushUnavailable = helpers.ErrPushUnavailable
)

// @Service
type ConnectionService struct {
    *create.CreateOperation
    *get.GetOperation
    *list.ListOperation
    *disconnect.DisconnectOperation
    *oauth.OAuthOperation
    *action.ActionOperation
    *test.TestOperation
    *watch.WatchOperation
}

var _ connectionservices.ConnectionService = (*ConnectionService)(nil)

func NewConnectionService(repository connectionrepo.ConnectionRepository, watchRepository watchrepo.WatchRepository, registry providercontracts.Registry, actions providercontracts.ActionRegistry, oauthRegistry mailcontracts.OAuthRegistry, gmailClient *gmail.Client, outlookClient *outlook.Client, encryptionKey []byte, pushTopic string, producer eventproducer.ChangedProducer, validator connectionvalidator.ConnectionValidator, mapper connectionmapper.ConnectionMapper) (connectionservices.ConnectionService,error) {
    cipher,err:=helpers.NewTokenCipher(encryptionKey)
    if err!=nil && (gmailClient!=nil || outlookClient!=nil) {return nil,err}

    oauthOperation:=oauth.New(oauthimpl.New(repository,oauthRegistry,cipher))
    return &ConnectionService{
        CreateOperation:create.New(createimpl.New(repository,registry,validator,mapper)),
        GetOperation:get.New(getimpl.New(repository,validator,mapper)),
        ListOperation:list.New(listimpl.New(repository,validator,mapper)),
        DisconnectOperation:disconnect.New(disconnectimpl.New(repository,validator,mapper)),
        OAuthOperation:oauthOperation,
        ActionOperation:action.New(actionimpl.New(repository,actions,oauthOperation)),
        TestOperation:test.New(testimpl.New(repository,actions,oauthOperation,mapper)),
        WatchOperation:watch.New(watchimpl.New(repository,watchRepository,gmailClient,pushTopic,producer,oauthOperation)),
    },nil
}
