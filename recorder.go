package pubsub

func NewRecorder() *Recorder {
	return &Recorder{}
}

type Recorder struct {
	Messages []Message
}

func (r *Recorder) Record(topic Topic, args ...interface{}) {
	r.Messages = append(r.Messages, Message{Topic: topic, Args: args})
}

func (r *Recorder) Reset() {
	r.Messages = nil
}

type Message struct {
	Topic Topic
	Args  []interface{}
}
