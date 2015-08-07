package topic

import (
	"testing"
)

func TestUnkownChannelShouldErrorWhenRemoved(t *testing.T) {
	topic := NewTopic("topic-1")

	err := topic.RemoveChannel("unknown-subscriber")

	if err != ChannelNotFoundError {
		t.Error("Removing an unknown channel should error.")
	}

}

func TestChannelForSubscriberExists(t *testing.T) {

	topic := NewTopic("topic-1")

	topic.AddChannel("subscriber")
	exists := topic.ChannelExists("subscriber")

	if !exists {
		t.Error("A Channel should exist for this subscriber.")
	}

	err := topic.RemoveChannel("subscriber")

	if err != nil {
		t.Error("Removing an existing user should not fail.")
	}

	exists = topic.ChannelExists("subscriber")

	if exists {
		t.Error("A Channel should not exist for this subscriber as they have been removed.")
	}

}

func TestMessagesSentToTopicAreAddedToEachSubscribedChannel(t *testing.T) {

	topic := NewTopic("topic-1")

	topic.AddChannel("subscriber-1")
	topic.AddChannel("subscriber-2")

	topic.PublishMessage(NewMessage([]byte("message-1")))

	subscriber1message1, err := topic.GetNextMessage("subscriber-1")

	if err != nil {
		t.Error(err.Error())
	}

	subscriber2message1, err := topic.GetNextMessage("subscriber-2")

	if err != nil {
		t.Error(err.Error())
	}

	if subscriber1message1.String() != "message-1" {
		t.Error("Subscriber1 should have received 'message-1' actually got ", subscriber1message1.String())
	}

	if subscriber2message1.String() != "message-1" {
		t.Error("Subscriber1 should have received 'message-1' actually got ", subscriber2message1.String())
	}

	topic.RemoveChannel("subscriber-2")

	topic.PublishMessage(NewMessage([]byte("message-2")))

	_, err = topic.GetNextMessage("subscriber-1")

	if err != nil {
		t.Error("Subscriber1 should have recieved their meesage without error.")
	}

	m, err := topic.GetNextMessage("subscriber-2")

	if err == nil {

		t.Error("Subscriber2 should have had an error as they are not subscribed.", m.String())
	}
}
