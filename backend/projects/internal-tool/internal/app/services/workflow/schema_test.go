package workflow

import (
 "encoding/json"
 "testing"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
func TestNestedSchema(t *testing.T){
 schema:=map[string]entity.Field{"profile":{Type:"object",Required:true,Properties:map[string]entity.Field{"name":{Type:"string",Required:true},"roles":{Type:"array",Items:&entity.Field{Type:"string"}},"score":{Type:"number",Nullable:true}}}}
 for _,raw:=range []string{`{"profile":{"name":"A","roles":["L1"],"score":1.25}}`,`{"profile":{"name":"A","score":null}}`}{if !validInput(json.RawMessage(raw),schema){t.Errorf("rejected %s",raw)}}
 for _,raw:=range []string{`{"profile":{"name":"A","name":"B"}}`,`{"profile":{"name":null}}`,`{"profile":{"name":"A","roles":[1]}}`,`{"profile":{"name":"A","score":"1"}}`,`{"profile":{"name":"A","extra":true}}`,`{"profile":{},"profile":{"name":"A"}}`}{if validInput(json.RawMessage(raw),schema){t.Errorf("accepted %s",raw)}}
}
