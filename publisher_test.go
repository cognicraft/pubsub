package pubsub

import "testing"

func TestPublisher(t *testing.T) {
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

func TestSubscriptions(t *testing.T) {

	pub := NewPublisher()

	rec1 := NewRecorder()
	sub1 := pub.Subscribe("#", rec1.Record)
	rec2 := NewRecorder()
	sub2 := pub.Subscribe("#", rec2.Record)

	sub := NewSubscriptions(sub1, sub2)

	pub.Publish("bäm", nil)

	if len(rec1.Messages) != 1 {
		t.Errorf("want: %d, got: %d", 1, len(rec1.Messages))
	}
	if len(rec2.Messages) != 1 {
		t.Errorf("want: %d, got: %d", 1, len(rec2.Messages))
	}
	rec1.Reset()
	rec2.Reset()

	sub.Cancel()

	pub.Publish("bäm", nil)

	if len(rec1.Messages) != 0 {
		t.Errorf("want: %d, got: %d", 0, len(rec1.Messages))
	}
	if len(rec2.Messages) != 0 {
		t.Errorf("want: %d, got: %d", 0, len(rec2.Messages))
	}

}
