package helpers

import (
    "testing"

    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
    eventflows "grovehub/services/foundation/events/producers/person/operations/changed/flows"
)

func TestChangeRules(t *testing.T) {
    old := &persondto.Person{GID:"person1",Version:1}
    current := &persondto.Person{GID:"person1",Version:2}
    cases := []struct{name string; op foundationevents.Operation; old,next *persondto.Person; want bool}{
        {"created",foundationevents.Created,nil,current,true},
        {"replay",foundationevents.Created,current,current,false},
        {"changed",foundationevents.Updated,old,current,true},
        {"unchanged",foundationevents.Updated,old,old,false},
        {"deactivated",foundationevents.Deactivated,old,current,true},
        {"already deactivated",foundationevents.Deactivated,current,current,false},
        {"unknown",foundationevents.Operation("UNKNOWN"),old,current,false},
    }
    for _,tc := range cases {
        if got := HasChanged(eventflows.Change{Operation:tc.op,Old:tc.old,New:tc.next}); got != tc.want {
            t.Errorf("%s: got %v want %v",tc.name,got,tc.want)
        }
    }
}
