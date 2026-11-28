package coverage_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/testing/coverage"
)

type closeFailureReader struct{io.Reader;failure error}
func(r closeFailureReader)Close()error{return r.failure}

var _=Describe("Profile and manifest persistence",func(){
	It("returns a read-close failure instead of claiming coverage",func(){
		failure:=errors.New("reader close failed")
		data,total,covered,err:=coverage.ReadProfile(closeFailureReader{strings.NewReader("mode: atomic\na.go:1.1,1.2 1 1\n"),failure},"fixture")
		Expect(err).To(MatchError(failure));Expect(data).To(BeNil());Expect(total).To(BeZero());Expect(covered).To(BeZero())
	})
	It("removes only the materialized output when durable manifest persistence fails",func(){
		root:=GinkgoT().TempDir()
		overlay,manifest:=filepath.Join(root,"overlay.json"),filepath.Join(root,"manifest")
		source,destination:=filepath.Join(root,"source.go"),filepath.Join(root,"zz_generated.go")
		writeFixture(source,"package example")
		writeOverlay(overlay,map[string]string{destination:source})
		failure:=errors.New("manifest unavailable")
		clean,err:=coverage.MaterializeWithManifest(overlay,manifest,func(path,logical string)error{
			Expect(path).To(Equal(manifest));Expect(logical).To(Equal(destination))
			Expect(destination).To(BeAnExistingFile())
			return failure
		})
		Expect(err).To(MatchError(failure));Expect(clean).To(BeNil())
		Expect(destination).NotTo(BeAnExistingFile());Expect(manifest).NotTo(BeAnExistingFile())
		contents,err:=os.ReadFile(source);Expect(err).NotTo(HaveOccurred());Expect(contents).To(Equal([]byte("package example")))
	})
	It("reports a blocked materialization parent and a read-only normalized profile",func(){
		root:=GinkgoT().TempDir()
		overlay,manifest:=filepath.Join(root,"overlay.json"),filepath.Join(root,"manifest")
		source:=filepath.Join(root,"source.go")
		writeFixture(source,"package example")
		link:=filepath.Join(root,"link")
		Expect(os.Symlink(filepath.Join(root,"missing"),link)).To(Succeed())
		writeOverlay(overlay,map[string]string{filepath.Join(link,"child","zz_generated.go"):source})
		_,err:=coverage.Materialize(overlay,manifest);Expect(err).To(HaveOccurred());Expect(manifest).NotTo(BeAnExistingFile())
		if os.Geteuid()!=0{
			dir:=filepath.Join(root,"readonly")
			profile:=filepath.Join(dir,"coverage.out")
			writeFixture(profile,"mode: atomic\na.go:1.1,1.2 1 1\n")
			Expect(os.Chmod(dir,0555)).To(Succeed())
			defer os.Chmod(dir,0755)
			_,_,err=coverage.Score(profile);Expect(err).To(HaveOccurred())
		}
	})
})
