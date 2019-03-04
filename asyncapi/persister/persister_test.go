package persister

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/repos/messages_repo"
	"github.com/uesteibar/asyncapi-watcher/asyncapi/spec"
	"github.com/uesteibar/asyncapi-watcher/storage/db"
	"testing"
	"time"
)

func TestPersist(t *testing.T) {
	database := db.GetUniqueDB()
	repo := messages_repo.New(database)
	repo.Migrate()
	chIn := make(chan spec.MessageSpec)
	chOut := make(chan spec.MessageSpec)
	p := New(chIn, chOut, database)
	go p.Watch()

	topic := uuid.New().String()

	msg := spec.MessageSpec{
		Topic: topic,
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "name", Type: "string"},
				spec.FieldSpec{Name: "age", Type: "float"},
			},
		},
	}
	chIn <- msg
	select {
	case _, _ = <-chOut:
		m, err := repo.Find(topic)
		assert.Nil(t, err)
		assert.Equal(t, msg, m)
	case <-time.After(3 * time.Second):
		t.Error("Expected to receive message, didn't receive any.")
	}

	// Update the message
	uMsg := spec.MessageSpec{
		Topic: topic,
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []spec.FieldSpec{
				spec.FieldSpec{Name: "name", Type: "string"},
				spec.FieldSpec{Name: "age", Type: "number"},
			},
		},
	}
	chIn <- uMsg
	select {
	case _, _ = <-chOut:
		um, uerr := repo.Find(topic)
		assert.Nil(t, uerr)
		assert.Equal(t, uMsg, um)
	case <-time.After(3 * time.Second):
		t.Error("Expected to receive message, didn't receive any.")
	}
}