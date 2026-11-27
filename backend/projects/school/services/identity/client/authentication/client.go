// Package authentication owns the transitional typed Platform auth client.
// Replace its wire request types with generated types when the restricted
// OpenAPI specification is verified.
package authentication

import (
 "bytes"
 "context"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "net/url"
)

type Login struct { Identifier string `json:"identifier"`; Password string `json:"password"` }
type ForgotPassword struct { Email string `json:"email"` }
type Result struct { Status int; Body []byte; Cookies []*http.Cookie }
type Client struct { origin *url.URL; http *http.Client }

func New(origin *url.URL, client *http.Client) *Client { return &Client{origin:origin,http:client} }
func (c *Client) Login(ctx context.Context, input Login)(Result,error) {
 return c.post(ctx,"/identity/restricted/auth/login",input)
}
func (c *Client) ForgotPassword(ctx context.Context,input ForgotPassword)(Result,error) {
 return c.post(ctx,"/identity/restricted/auth/forgot-password",input)
}
func (c *Client) post(ctx context.Context,path string,input any)(Result,error) {
 payload,err:=json.Marshal(input)
 if err!=nil {return Result{},err}
 target:=c.origin.ResolveReference(&url.URL{Path:path})
 req,err:=http.NewRequestWithContext(ctx,http.MethodPost,target.String(),bytes.NewReader(payload))
 if err!=nil {return Result{},err}
 req.Header.Set("Content-Type","application/json")
 req.Header.Set("Accept","application/json")
 resp,err:=c.http.Do(req)
 if err!=nil {return Result{},err}
 defer resp.Body.Close()
 if resp.StatusCode<200 || resp.StatusCode>=600 || (resp.StatusCode>=300 && resp.StatusCode<400) {return Result{},fmt.Errorf("unexpected platform response")}
 body,err:=io.ReadAll(io.LimitReader(resp.Body,1<<20+1))
 if err!=nil || len(body)>1<<20 || !json.Valid(body) {return Result{},fmt.Errorf("invalid platform response")}
 return Result{Status:resp.StatusCode,Body:body,Cookies:resp.Cookies()},nil
}
