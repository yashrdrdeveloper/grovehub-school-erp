package credential_test

import (
	"crypto/ed25519"
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/gateway/internal/src/contracts"
	"grovehub/libs/core/gateway/internal/src/operations/credential"
)

var _=Describe("Credential encoding",func(){
	It("never signs an output when serialization fails",func(){
		_,key,err:=ed25519.GenerateKey(nil)
		Expect(err).NotTo(HaveOccurred())
		token,err:=credential.SignWorkloadWithEncoder("key","issuer",key,contracts.WorkloadClaims{Subject:"school",Audience:"platform",Scopes:[]string{"read"},ExpiresAt:time.Now().Add(time.Minute)},func(any)([]byte,error){
			return json.Marshal(make(chan int))
		})
		Expect(err).To(HaveOccurred())
		Expect(token).To(BeEmpty())
	})
})
