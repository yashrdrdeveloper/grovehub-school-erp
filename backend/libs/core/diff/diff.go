// Package diff computes JSON-field-oriented changes without depending on domain DTOs.
package diff

import (
    "bytes"
    "encoding/json"
    "fmt"
    "reflect"
    "sort"
    "strconv"
)

type Change struct {
    Field string `json:"field"`
    OldValue any `json:"oldValue"`
    NewValue any `json:"newValue"`
}

// Compare returns deterministic flat changes using JSON field names and array indices.
// Array elements are compared by position, not by domain identity.
func Compare[T any](old, next T) ([]Change, error) {
    before, err := toJSON(old)
    if err != nil { return nil, fmt.Errorf("marshal old value: %w", err) }
    after, err := toJSON(next)
    if err != nil { return nil, fmt.Errorf("marshal new value: %w", err) }
    changes := make([]Change, 0)
    walk("", before, after, &changes)
    return changes, nil
}

func toJSON(value any) (any, error) {
    payload, err := json.Marshal(value)
    if err != nil { return nil, err }
    decoder := json.NewDecoder(bytes.NewReader(payload))
    decoder.UseNumber()
    var result any
    if err := decoder.Decode(&result); err != nil { return nil, err }
    return result, nil
}

func walk(path string, old, next any, changes *[]Change) {
    if reflect.DeepEqual(old, next) { return }
    oldMap, oldIsMap := old.(map[string]any)
    newMap, newIsMap := next.(map[string]any)
    if oldIsMap && newIsMap {
        names := make(map[string]struct{}, len(oldMap)+len(newMap))
        for name := range oldMap { names[name] = struct{}{} }
        for name := range newMap { names[name] = struct{}{} }
        keys := make([]string, 0, len(names))
        for name := range names { keys = append(keys, name) }
        sort.Strings(keys)
        for _, name := range keys {
            child := name
            if path != "" { child = path + "." + name }
            walk(child, oldMap[name], newMap[name], changes)
        }
        return
    }
    oldList, oldIsList := old.([]any)
    newList, newIsList := next.([]any)
    if oldIsList && newIsList {
        length := len(oldList)
        if len(newList) > length { length = len(newList) }
        for i := 0; i < length; i++ {
            var left, right any
            if i < len(oldList) { left = oldList[i] }
            if i < len(newList) { right = newList[i] }
            walk(path+"["+strconv.Itoa(i)+"]", left, right, changes)
        }
        return
    }
    if path == "" { path = "$" }
    *changes = append(*changes, Change{Field: path, OldValue: old, NewValue: next})
}
