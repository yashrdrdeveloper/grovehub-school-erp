package helpers
import(
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "html/template"
 "grovehub/services/identity/internal/resources"
 integrationclient "grovehub/services/integration/client/connection"
)
type Mail struct{Title,Name,Description,Button,URL string;Minutes int}
func RenderMail(mail Mail)(string,error){
 t,err:=template.ParseFS(resources.Templates,"templates/registration/password_setup.html")
 if err!=nil{return "",err}
 var b bytes.Buffer
 if err:=t.Execute(&b,mail);err!=nil{return "",err}
 return b.String(),nil
}
// Send uses Integration's public typed client, not a raw vendor HTTP request.
func Send(ctx context.Context,client integrationclient.Client,tenantID,connectionID,to string,mail Mail)error{
 if client==nil||tenantID==""||connectionID==""{return errors.New("system mail is not configured")}
 html,err:=RenderMail(mail);if err!=nil{return err}
 input:=map[string]any{"to":[]string{to},"subject":mail.Title,"body":map[string]string{"html":html,"text":mail.Description+"\n"+mail.URL}}
 b,err:=json.Marshal(input);if err!=nil{return err}
 _,err=client.ExecuteAction(ctx,tenantID,connectionID,integrationclient.Action("mail.send"),json.RawMessage(b))
 return err
}
