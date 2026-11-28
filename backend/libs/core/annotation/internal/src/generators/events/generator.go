// Package events generates stable, typed event identity constants from
// domain-owned @Event contracts. Publication, transactions, payloads and
// subscription semantics remain explicit service responsibilities.
package events

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/events/helpers"
)

func GenerateAll(backend string)error{
    services,err:=os.ReadDir(filepath.Join(backend,"services"));if err!=nil{return err}
    for _,service:=range services {
        if !service.IsDir(){continue}
        owner:=filepath.Join(backend,"services",service.Name())
        relative:="events"
        if info,err:=os.Stat(filepath.Join(owner,"internal","app","events"));err==nil&&info.IsDir(){
            relative=filepath.Join("internal","app","events")
        }
        if err:=layout.WalkOwnedDirectories(owner,relative,nil,func(dir,rel string)error{
            out:=filepath.Join(owner,".generated",rel,"zz_events.gen.go")
            defs,pkg,err:=helpers.Discover(dir);if err!=nil{return fmt.Errorf("%s: %w",dir,err)}
            if len(defs)==0 {
                if err:=os.Remove(out);err!=nil&&!os.IsNotExist(err){return err}
                return nil
            }
            return helpers.Write(out,pkg,defs)
        });err!=nil{return err}
    }
    return nil
}
