package pack_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/tools/modulepack/internal/src/operations/pack"
)

type closeFailure struct { io.WriteCloser }
func (writer closeFailure) Close() error {
	_ = writer.WriteCloser.Close()
	return errors.New("archive close failed")
}

var _=Describe("Module archive edge cases",func(){
	It("rejects invalid proxy names and path-major version mismatches",func(){
		source,_,args:=fixture()
		mod:=filepath.Join(source,"go.mod")
		Expect(os.WriteFile(mod,[]byte("module example.com/Invalid!\n\ngo 1.27.0\n"),0644)).To(Succeed())
		Expect(pack.Run(args)).NotTo(Succeed())
		Expect(os.WriteFile(mod,[]byte(moduleFile),0644)).To(Succeed())
		args[2]="v0.1.0!invalid"
		Expect(pack.Run(args)).NotTo(Succeed())
		args[2]="v2.0.0"
		Expect(pack.Run(args)).NotTo(Succeed())
	})
	It("removes a partial archive when the module contains conflicting paths",func(){
		source,proxy,args:=fixture()
		Expect(os.WriteFile(filepath.Join(source,"A.go"),[]byte("package core\n"),0644)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(source,"a.go"),[]byte("package core\n"),0644)).To(Succeed())
		Expect(pack.Run(args)).NotTo(Succeed())
		Expect(filepath.Join(versionDirectory(proxy),version+".zip")).NotTo(BeAnExistingFile())
	})
	It("propagates archive close errors and discards the incomplete archive",func(){
		_,proxy,args:=fixture()
		err:=pack.RunWithArchive(args,func(path string)(io.WriteCloser,error){
			file,err:=os.OpenFile(path,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0644)
			if err!=nil{return nil,err}
			return closeFailure{file},nil
		})
		Expect(err).To(MatchError("archive close failed"))
		Expect(filepath.Join(versionDirectory(proxy),version+".zip")).NotTo(BeAnExistingFile())
	})
	It("reports UTC timestamps outside JSON's representable year range",func(){
		_,_,args:=fixture()
		args[3]="0000-01-01T00:00:00+01:00"
		Expect(pack.Run(args)).To(MatchError(ContainSubstring("year outside")))
	})
})
