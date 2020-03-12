package pubsub

func NewRecorder() *Recorder {
	return &Recorder{}
}

type Recorder struct {
	Messages []Message
}

func (r *Recorder) Record(topic Topic, data interface{}) {
	r.Messages = append(r.Messages, Message{Topic: topic, Data: data})
}

func (r *Recorder) Reset() {
	r.Messages = nil
}

type Message struct {
	Topic Topic
	Data  interface{}
}
