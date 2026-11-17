package helpers
import ("crypto/rand";"crypto/sha256";"encoding/base64";"encoding/hex";"errors")
func NewSecret(n int)(string,error){b:=make([]byte,n);if _,err:=rand.Read(b);err!=nil{return "",err};return base64.RawURLEncoding.EncodeToString(b),nil}
func Digest(s string,n int)(string,error){
 b,err:=base64.RawURLEncoding.DecodeString(s)
 if err!=nil||len(b)!=n{return "",errors.New("invalid credential")}
 h:=sha256.Sum256(b);return hex.EncodeToString(h[:]),nil
}
