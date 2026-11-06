package controller

type Controller interface{
 ConnectionController
 ProviderController
 OAuthController
 TestController
 WatchController
 WebhookController
 ActionController
}
