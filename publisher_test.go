package pubsub

import "testing"

func TestSubject(t *testing.T) {
	rec := NewRecorder()
	pub := NewPublisher()
	sub := pub.Subscribe("#", rec.Record)

	pub.Publish("bäm", nil)

	if len(rec.Messages) != 1 {
		t.Errorf("want: %d, got: %d", 1, len(rec.Messages))
	}
	if rec.Messages[0].Topic != "bäm" {
		t.Errorf("want: %s, got: %s", "bäm", rec.Messages[0].Topic)
	}
	if rec.Messages[0].Data != nil {
		t.Errorf("want: %#v, got: %#v", nil, rec.Messages[0].Data)
	}

	sub.Cancel()
	rec.Reset()

	pub.Publish("bäm", nil)

	if len(rec.Messages) != 0 {
		t.Fail()
	}
}
