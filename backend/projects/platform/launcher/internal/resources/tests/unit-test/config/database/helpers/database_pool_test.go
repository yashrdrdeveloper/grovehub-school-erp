package helpers_test
import (ginkgo "github.com/onsi/ginkgo/v2";"grovehub/launcher/internal/src/config/database/helpers")
func verifyPoolLimits() {
	t := ginkgo.GinkgoT()
 for _,tc:=range []struct{open,idle string;wantOpen,wantIdle int;valid bool}{
 {"","",20,5,true},{"3","",3,3,true},{"8","0",8,0,true},{"0","",0,0,false},{"101","",0,0,false},{"5","6",0,0,false},{"bad","",0,0,false},
 }{
 t.Setenv("TEST_DB_MAX_OPEN",tc.open);t.Setenv("TEST_DB_MAX_IDLE",tc.idle)
 open,idle,err:=helpers.PoolLimits("TEST_DB",20,5)
 if (err==nil)!=tc.valid||(tc.valid&&(open!=tc.wantOpen||idle!=tc.wantIdle)){t.Fatal(tc,open,idle,err)}
 }
}

// @Testing suite=LauncherDatabasePool
func PoolLimitsSpecs() {
	ginkgo.Describe("PoolLimits", func() {
		ginkgo.It("Pool Limits", verifyPoolLimits)
	})
}
