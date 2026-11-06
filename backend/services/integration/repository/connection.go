package repository

import (
 "context"
 "errors"
 "time"
 "grovehub/services/integration/entity"
)
var (
 ErrConnectionNotFound = errors.New("integration connection not found")
 ErrOAuthStateInvalid = errors.New("oauth state invalid or expired")
)
type ConnectionRepository interface {
 Create(ctx context.Context, connection entity.Connection) error
 Get(ctx context.Context, tenantID, id string) (entity.Connection, error)
 GetByID(ctx context.Context,id string)(entity.Connection,error)
 List(ctx context.Context, tenantID, ownerGID string) ([]entity.Connection, error)
 Disconnect(ctx context.Context, tenantID, id string) (entity.Connection, error)
 SaveOAuthState(ctx context.Context, state entity.OAuthState) error
 ConsumeOAuthState(ctx context.Context,hash string,now time.Time)(entity.OAuthState,error)
 Activate(ctx context.Context,connectionID,provider,externalID string,ciphertext []byte) error
 Credential(ctx context.Context,connectionID string)([]byte,error)
 UpdateCredential(ctx context.Context,connectionID string,ciphertext []byte) error
 SaveTest(ctx context.Context,test entity.ConnectionTest) error
 ListTests(ctx context.Context,connectionID string)([]entity.ConnectionTest,error)
}
