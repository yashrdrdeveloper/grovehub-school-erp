package helpers

import (
    "context"
    "encoding/json"
    "testing"

    coreevents "grovehub/libs/core/events"
    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
    changed "grovehub/services/foundation/events/producers/person/operations/changed"
    eventflows "grovehub/services/foundation/events/producers/person/operations/changed/flows"
)

// Moving the contract must not change the public event topic, type, or snapshots.
func TestChangedMessagePreservesWireContract(t *testing.T) {
    previous := &persondto.Person{GID: "person-one", FirstName: "Before", Version: 3}
    current := &persondto.Person{GID: "person-one", FirstName: "After", Version: 4}
    message, err := BuildMessage(context.Background(), eventflows.Change{
        Operation: foundationevents.Updated,
        Old: previous,
        New: current,
    })
    if err != nil { t.Fatal(err) }
    if message.Topic != "foundation_person_changed" || message.Topic != changed.ChangedTopic {
        t.Fatalf("unexpected topic: %q", message.Topic)
    }
    if message.Metadata["event-id"] != message.EventID ||
        message.Metadata["event-type"] != changed.ChangedType ||
        message.Metadata["source"] != changed.Source {
        t.Fatalf("unexpected broker metadata: %+v", message.Metadata)
    }

    var envelope coreevents.Envelope[changed.Changed]
    if err := json.Unmarshal(message.Payload, &envelope); err != nil { t.Fatal(err) }
    if envelope.EventID != message.EventID ||
        envelope.EventType != "foundation.person.changed.v1" ||
        envelope.Version != 1 || envelope.AggregateID != current.GID {
        t.Fatalf("unexpected envelope identity: %+v", envelope)
    }
    if envelope.Data.Operation != foundationevents.Updated ||
        envelope.Data.Old == nil || envelope.Data.New == nil ||
        envelope.Data.Old.FirstName != "Before" || envelope.Data.New.FirstName != "After" {
        t.Fatalf("unexpected change snapshots: %+v", envelope.Data)
    }
}
