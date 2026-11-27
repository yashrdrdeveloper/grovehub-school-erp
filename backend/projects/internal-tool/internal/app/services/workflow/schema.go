package workflow

import (
 "bytes"
 "encoding/json"
 "io"
 "strings"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

// Schema is a bounded subset, not a complete JSON Schema implementation.
// Reject unknown properties and ambiguous duplicate object keys at every level.
func validFields(fields map[string]entity.Field,depth int,budget *int)bool{
 if len(fields)>32{return false}
 for key,field:=range fields{
  if !fieldPattern.MatchString(key)||!validField(field,depth,budget){return false}
 }
 return true
}
func validField(f entity.Field,depth int,budget *int)bool{
 if depth>5||*budget<=0{return false};*budget--
 switch f.Type{
 case "string","integer","boolean","number":
  return f.Items==nil&&len(f.Properties)==0
 case "object":
  return f.Items==nil&&f.Properties!=nil&&validFields(f.Properties,depth+1,budget)
 case "array":
  return f.Items!=nil&&len(f.Properties)==0&&!f.Items.Required&&validField(*f.Items,depth+1,budget)
 default:return false
 }
}
func object(raw json.RawMessage)(map[string]json.RawMessage,bool){
 d:=json.NewDecoder(bytes.NewReader(raw))
 token,err:=d.Token();if err!=nil||token!=json.Delim('{'){return nil,false}
 values:=map[string]json.RawMessage{}
 for d.More(){
  key,err:=d.Token();if err!=nil{return nil,false}
  name,ok:=key.(string);if !ok{return nil,false}
  if _,exists:=values[name];exists{return nil,false}
  var value json.RawMessage
  if d.Decode(&value)!=nil{return nil,false}
  values[name]=value
  if len(values)>32{return nil,false}
 }
 token,err=d.Token();if err!=nil||token!=json.Delim('}'){return nil,false}
 if _,err=d.Token();err!=io.EOF{return nil,false}
 return values,true
}
func validInput(raw json.RawMessage,schema map[string]entity.Field)bool{
 if len(raw)==0||len(raw)>16384{return false}
 budget:=128;if !validFields(schema,0,&budget){return false}
 return validObject(raw,schema)
}
func validObject(raw json.RawMessage,schema map[string]entity.Field)bool{
 values,ok:=object(raw);if !ok{return false}
 for key,value:=range values{field,ok:=schema[key];if !ok||!validValue(value,field){return false}}
 for key,field:=range schema{if field.Required{if _,ok:=values[key];!ok{return false}}}
 return true
}
func validValue(raw json.RawMessage,field entity.Field)bool{
 if strings.TrimSpace(string(raw))=="null"{return field.Nullable}
 switch field.Type{
 case "string":var v string;return json.Unmarshal(raw,&v)==nil&&len(v)<=4096
 case "integer":var v int64;return json.Unmarshal(raw,&v)==nil
 case "boolean":var v bool;return json.Unmarshal(raw,&v)==nil
 case "number":var v json.Number;d:=json.NewDecoder(bytes.NewReader(raw));d.UseNumber();var value any;if d.Decode(&value)!=nil{return false};v,ok:=value.(json.Number);return ok&&len(v.String())<=128
 case "object":return validObject(raw,field.Properties)
 case "array":
  var values []json.RawMessage;if json.Unmarshal(raw,&values)!=nil||values==nil||len(values)>100||field.Items==nil{return false}
  for _,value:=range values{if !validValue(value,*field.Items){return false}};return true
 default:return false
 }
}
