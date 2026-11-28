package helpers_test
import ("testing";"grovehub/launcher/internal/src/config/database/helpers")
func TestPoolLimits(t *testing.T){
 for _,tc:=range []struct{open,idle string;wantOpen,wantIdle int;valid bool}{
 {"","",20,5,true},{"3","",3,3,true},{"8","0",8,0,true},{"0","",0,0,false},{"101","",0,0,false},{"5","6",0,0,false},{"bad","",0,0,false},
 }{
 t.Setenv("TEST_DB_MAX_OPEN",tc.open);t.Setenv("TEST_DB_MAX_IDLE",tc.idle)
 open,idle,err:=helpers.PoolLimits("TEST_DB",20,5)
 if (err==nil)!=tc.valid||(tc.valid&&(open!=tc.wantOpen||idle!=tc.wantIdle)){t.Fatal(tc,open,idle,err)}
 }
}
