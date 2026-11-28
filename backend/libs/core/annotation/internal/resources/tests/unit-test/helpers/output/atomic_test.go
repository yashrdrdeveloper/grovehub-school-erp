package output_test

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/output"
)

var persistenceFailure=errors.New("persistence failed")

type faultFile struct{
	*os.File
	failWrite,failClose bool
}
func(f faultFile) Write(data []byte)(int,error){
	if f.failWrite{return 0,persistenceFailure}
	return f.File.Write(data)
}
func(f faultFile) Close()error{
	err:=f.File.Close()
	if f.failClose{return persistenceFailure}
	return err
}

var _=Describe("AtomicWriter",func(){
	It("replaces only after a complete successful write and close",func(){
		root:=GinkgoT().TempDir()
		path:=filepath.Join(root,"profile.out")
		Expect(os.WriteFile(path,[]byte("old"),0644)).To(Succeed())
		Expect(output.WriteAtomic(path,"profile-*.out",[]byte("new"))).To(Succeed())
		data,err:=os.ReadFile(path);Expect(err).NotTo(HaveOccurred());Expect(string(data)).To(Equal("new"))
		entries,err:=os.ReadDir(root);Expect(err).NotTo(HaveOccurred());Expect(entries).To(HaveLen(1))
	})
	It("propagates create, write, close and rename failures without replacing the old file",func(){
		root:=GinkgoT().TempDir()
		path:=filepath.Join(root,"profile.out")
		Expect(os.WriteFile(path,[]byte("old"),0644)).To(Succeed())
		for _,failure:=range []string{"create","write","close","rename"}{
			writer:=output.AtomicWriter{
				CreateTemp:func(dir,pattern string)(output.TempFile,error){
					if failure=="create"{return nil,persistenceFailure}
					file,err:=os.CreateTemp(dir,pattern)
					if err!=nil{return nil,err}
					return faultFile{File:file,failWrite:failure=="write",failClose:failure=="close"},nil
				},
				Rename:func(old,next string)error{
					Expect(failure).To(Equal("rename"))
					return persistenceFailure
				},
				Remove:os.Remove,
			}
			Expect(writer.Write(path,"profile-*.out",[]byte("new"))).To(MatchError(persistenceFailure),failure)
			data,err:=os.ReadFile(path);Expect(err).NotTo(HaveOccurred());Expect(string(data)).To(Equal("old"))
			entries,err:=os.ReadDir(root);Expect(err).NotTo(HaveOccurred());Expect(entries).To(HaveLen(1))
		}
	})
})

var _=Describe("WriteJSON",func(){
	It("serializes deterministically, skips identical files and propagates encoding errors",func(){
		root:=GinkgoT().TempDir()
		path:=filepath.Join(root,"generated","overlay.json")
		value:=map[string]string{"z":"last","a":"first"}
		Expect(output.WriteJSON(path,"overlay-*.json",value)).To(Succeed())
		data,err:=os.ReadFile(path);Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("{\n  \"a\": \"first\",\n  \"z\": \"last\"\n}\n"))
		stamp:=time.Unix(1000,0);Expect(os.Chtimes(path,stamp,stamp)).To(Succeed())
		Expect(output.WriteJSON(path,"overlay-*.json",value)).To(Succeed())
		info,err:=os.Stat(path);Expect(err).NotTo(HaveOccurred());Expect(info.ModTime()).To(Equal(stamp))
		Expect(output.WriteJSON(path,"overlay-*.json",make(chan int))).NotTo(Succeed())
		current,err:=os.ReadFile(path);Expect(err).NotTo(HaveOccurred());Expect(current).To(Equal(data))
		Expect(output.WriteJSON(filepath.Join(path,"child.json"),"overlay-*.json",value)).NotTo(Succeed())
		Expect(output.WriteJSON(filepath.Dir(path),"overlay-*.json",value)).NotTo(Succeed())
	})
})
