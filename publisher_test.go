package pubsub

import "testing"

func TestSubject(t *testing.T) {
	rec := NewRecorder()
	pub := NewPublisher()
	sub := pub.Subscribe("#", rec.Record)

	pub.Publish("bäm", nil)

	if len(rec.Messages) != 1 {
		t.Fail()
	}
	if rec.Messages[0].Topic != "bäm" {
		t.Fail()
	}
	if rec.Messages[0].Data == nil {
		t.Fail()
	}

	sub.Cancel()
	rec.Reset()

	pub.Publish("bäm", nil)

	if len(rec.Messages) != 0 {
		t.Fail()
	}
}
