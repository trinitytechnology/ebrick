package messaging

import (
	"fmt"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

func GenerateConsumerName(name string) string {
	return fmt.Sprintf("%s-%s", name, uuid.NewString())
}

func CreateEvent(eventSource string, eventType EventType, data any) event.Event {
	ev := event.New()
	ev.SetSource(eventSource)
	ev.SetType(string(eventType))
	ev.SetTime(time.Now())
	ev.SetData(*event.StringOfApplicationJSON(), data)
	ev.SetID(uuid.NewString())
	return ev
}
