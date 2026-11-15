package impl

import (
    config "grovehub/services/integration/config"
    configimpl "grovehub/services/integration/internal/app/config/impl"
    "grovehub/services/integration/internal/app/providers"
    mailfactory "grovehub/services/integration/internal/app/providers/mailprovider/impl"
    "grovehub/services/integration/internal/app/providers/mailprovider/providers/gmail"
    "grovehub/services/integration/internal/app/providers/mailprovider/providers/outlook"
)

type providerDependencies struct {
    registry providers.Registry
    actions providers.ActionRegistry
    oauth providers.OAuthRegistry
    gmail *gmail.Client
    outlook *outlook.Client
    key []byte
    push config.GmailPush
}

func (b *builder) newProviders() (providerDependencies,error) {
    loader:=configimpl.NewEnvironmentLoader()
    gmailConfig,gmailErr:=loader.Gmail()
    key,keyErr:=loader.CredentialKey()
    var gmailClient *gmail.Client
    if keyErr!=nil {b.logger.Warn("integration encryption key not configured; OAuth providers disabled")}
    if gmailErr==nil && keyErr==nil {gmailClient=gmail.NewClient(gmailConfig)} else {b.logger.Warn("gmail provider not configured; OAuth and tests disabled")}

    outlookConfig,outlookErr:=loader.Outlook()
    var outlookClient *outlook.Client
    if outlookErr==nil && keyErr==nil {outlookClient=outlook.NewClient(outlookConfig)} else {b.logger.Warn("Outlook OAuth not configured; Outlook actions disabled")}

    runtime,err:=mailfactory.NewFactory().Build(gmail.NewStrategy(gmailClient),outlook.NewStrategy(outlookClient))
    if err!=nil{return providerDependencies{},err}
    return providerDependencies{
        registry:runtime.Registry,
        actions:runtime.Actions,
        oauth:runtime.OAuth,
        gmail:gmailClient,
        outlook:outlookClient,
        key:key,
        push:loader.GmailPush(),
    },nil
}
