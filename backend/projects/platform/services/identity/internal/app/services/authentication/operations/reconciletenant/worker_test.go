package reconciletenant

import (
 "context"
 "errors"
 "testing"
 schoolgateway "grovehub/internal/session"
)

type testStore struct {
 calls []string
 cutoffs map[string]int64
 failWrite bool
}

func (s *testStore) ListTenants(_ context.Context, after string, limit int) ([]string,error) {
 s.calls = append(s.calls, after)
 if limit != batchSize { return nil,errors.New("unbounded scan") }
 return []string{"a","b","c"},nil
}
func (s *testStore) RevokeBefore(_ context.Context, tenant string, generation int64)(int64,error) {
 if tenant=="b" && s.failWrite { return 0,errors.New("database unavailable") }
 s.cutoffs[tenant]=generation
 return 1,nil
}

func TestReconciliationRetriesWithoutBlockingOtherTenants(t *testing.T) {
 store:=&testStore{cutoffs:map[string]int64{},failWrite:true}
 failLookup:=true
 worker:=New(store,func(_ context.Context, tenant string)(schoolgateway.TenantState,error){
  if tenant=="a"&&failLookup{return schoolgateway.TenantState{},errors.New("registry unavailable")}
  return schoolgateway.TenantState{Generation:2,Active:tenant!="c"},nil
 },nil)
 if err:=worker.Step(context.Background());err==nil{t.Fatal("failures not reported")}
 if len(store.cutoffs)!=1||store.cutoffs["c"]!=2{t.Fatal("unrelated disabled tenant not reconciled",store.cutoffs)}
 failLookup=false;store.failWrite=false
 if err:=worker.Step(context.Background());err!=nil{t.Fatal(err)}
 if len(store.cutoffs)!=3||store.cutoffs["a"]!=2||store.cutoffs["b"]!=2{t.Fatal("retry missing",store.cutoffs)}
 if len(store.calls)!=2||store.calls[1]!=""{t.Fatal("sweep did not wrap",store.calls)}
}
