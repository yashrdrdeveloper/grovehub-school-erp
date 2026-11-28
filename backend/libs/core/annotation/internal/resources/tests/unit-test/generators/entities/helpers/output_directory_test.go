package helpers_test

import (
"grovehub/libs/core/annotation/internal/src/generators/entities/helpers"
fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"
)

var _=Describe("Entity output directory",func(){It("reports an existing file in place of the output directory",func(){Expect(helpers.Generate(GinkgoT().TempDir(),fixture.Write(GinkgoT().TempDir(),"blocked","file"))).NotTo(Succeed())})})
