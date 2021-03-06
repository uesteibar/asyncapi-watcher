package messagesrepo

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uesteibar/scribano/asyncapi/spec"
	"github.com/uesteibar/scribano/storage/db"
)

func TestRepo(t *testing.T) {
	repo := New(db.GetUniqueDB())
	repo.Migrate()

	topic := uuid.New().String()
	msg := spec.MessageSpec{
		Topic: topic,
		Payload: spec.PayloadSpec{
			Type: "object",
			Fields: []*spec.FieldSpec{
				&spec.FieldSpec{Name: "name", Type: "string"},
				&spec.FieldSpec{Name: "age", Type: "float"},
				&spec.FieldSpec{Name: "emptyHash", Type: "object"},
				&spec.FieldSpec{Name: "fines", Type: "array", Item: &spec.FieldSpec{Type: "string"}},
				&spec.FieldSpec{Name: "emptyHashes", Type: "array", Item: &spec.FieldSpec{Type: "object"}},
				&spec.FieldSpec{
					Name: "friends",
					Type: "array",
					Item: &spec.FieldSpec{
						Type:   "object",
						Fields: []*spec.FieldSpec{&spec.FieldSpec{Name: "name", Type: "string"}},
					},
				},
				&spec.FieldSpec{
					Name:   "car",
					Type:   "object",
					Fields: []*spec.FieldSpec{&spec.FieldSpec{Name: "brand", Type: "string"}},
				},
			},
		},
	}

	err := repo.Create(msg)
	assert.Nil(t, err)

	m, err := repo.Find(topic)

	assert.Nil(t, err)
	assert.Equal(t, topic, m.Topic)
	assert.Equal(t, msg.Payload, m.Payload)

	notFoundMsg, err := repo.Find("wrong.topic")
	assert.Equal(t, "", notFoundMsg.Topic)
	switch err.(type) {
	case *ErrNotFound:
		assert.Equal(t, "NOT_FOUND", err.Error())
	default:
		t.Error("Expected error, got nothing")
	}

	newPayload := spec.PayloadSpec{
		Type: "object",
		Fields: []*spec.FieldSpec{
			&spec.FieldSpec{Name: "name", Type: "string"},
			&spec.FieldSpec{Name: "age", Type: "number"},
		},
	}
	m.Payload = newPayload

	err = repo.Update(m)
	assert.Nil(t, err)
	um, _ := repo.Find(topic)
	assert.Equal(t, m, um)
}

func TestFindAll(t *testing.T) {
	repo := New(db.GetUniqueDB())
	repo.Migrate()

	msg1 := spec.MessageSpec{
		Topic: uuid.New().String(),
		Payload: spec.PayloadSpec{
			Type:   "object",
			Fields: []*spec.FieldSpec{},
		},
	}
	repo.Create(msg1)

	msg2 := spec.MessageSpec{
		Topic: uuid.New().String(),
		Payload: spec.PayloadSpec{
			Type:   "object",
			Fields: []*spec.FieldSpec{},
		},
	}
	repo.Create(msg2)

	msgs, err := repo.FindAll()
	var expectedMsgs []spec.MessageSpec = []spec.MessageSpec{msg1, msg2}
	assert.Nil(t, err)
	assert.Equal(t, expectedMsgs, msgs)
}

func TestFindByExchange(t *testing.T) {
	repo := New(db.GetUniqueDB())
	repo.Migrate()

	msg1 := spec.MessageSpec{
		Topic:    uuid.New().String(),
		Exchange: "correct",
		Payload: spec.PayloadSpec{
			Type:   "object",
			Fields: []*spec.FieldSpec{},
		},
	}
	repo.Create(msg1)

	msg2 := spec.MessageSpec{
		Topic:    uuid.New().String(),
		Exchange: "other",
		Payload: spec.PayloadSpec{
			Type:   "object",
			Fields: []*spec.FieldSpec{},
		},
	}
	repo.Create(msg2)

	msgs, err := repo.FindByExchange("correct")
	assert.Nil(t, err)
	assert.Equal(t, []spec.MessageSpec{msg1}, msgs)

	msgs, err = repo.FindByExchange("incorrect")
	assert.Nil(t, err)
	assert.Equal(t, []spec.MessageSpec{}, msgs)
}
