package diff_test

import (
    "encoding/json"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    corediff "grovehub/libs/core/diff"
)

type testContact struct {
    Value string `json:"value"`
}
type testPerson struct {
    Name string `json:"name"`
    Contacts []testContact `json:"contacts"`
    Version int `json:"version"`
}

var _ = Describe("Compare", func() {
    It("reports deterministic flattened JSON names, nested arrays and numbers", func() {
        before := testPerson{Name: "Old", Contacts: []testContact{{Value:"a@example.com"}}, Version:1}
        after := testPerson{Name: "New", Contacts: []testContact{{Value:"b@example.com"}}, Version:2}
        changes, err := corediff.Compare(before, after)
        Expect(err).NotTo(HaveOccurred())
        Expect(changes).To(HaveLen(3))
        Expect(changes[0].Field).To(Equal("contacts[0].value"))
        Expect(changes[1].Field).To(Equal("name"))
        Expect(changes[2].Field).To(Equal("version"))
        Expect(changes[2].OldValue).To(Equal(json.Number("1")))
        Expect(changes[2].NewValue).To(Equal(json.Number("2")))
    })

    It("returns no changes for equivalent objects", func() {
        value := testPerson{Name: "Same", Version: 3}
        changes, err := corediff.Compare(value, value)
        Expect(err).NotTo(HaveOccurred())
        Expect(changes).To(BeEmpty())
    })

    It("represents added and removed values with null", func() {
        old := map[string]any{"gone":"value"}
        next := map[string]any{"added":"value"}
        changes, err := corediff.Compare(old, next)
        Expect(err).NotTo(HaveOccurred())
        Expect(changes).To(HaveLen(2))
        Expect(changes[0]).To(Equal(corediff.Change{Field:"added", OldValue:nil, NewValue:"value"}))
        Expect(changes[1]).To(Equal(corediff.Change{Field:"gone", OldValue:"value", NewValue:nil}))
    })
})
