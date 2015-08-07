package topic

import (
	"testing"
)

func TestNoMessagesAvailableForSubscriber(t *testing.T) {

	channel := NewChannel()

	_, err := channel.Pop()

	if err != NoMessagesAvailable {
		t.Error("A NoMessagesAvailable error should have been returned.")
	}
}

func TestMessagesAddedCanBeRetrieved(t *testing.T) {

	channel := NewChannel()

	assertChannelLength(t, channel, 0)
	channel.Push(NewMessage([]byte("message-1")))

	assertChannelLength(t, channel, 1)
	channel.Push(NewMessage([]byte("message-2")))

	assertChannelLength(t, channel, 2)
	channel.Push(NewMessage([]byte("message-3")))

	assertChannelLength(t, channel, 3)
	assertMessageRetreivedWithExpectedContent(t, channel, "message-1")

	assertChannelLength(t, channel, 2)
	assertMessageRetreivedWithExpectedContent(t, channel, "message-2")

	assertChannelLength(t, channel, 1)
	assertMessageRetreivedWithExpectedContent(t, channel, "message-3")

	assertChannelLength(t, channel, 0)
}

func assertChannelLength(t *testing.T, channel Channel, expectedCount int) {

	actualCount := channel.Count()

	if actualCount != expectedCount {
		t.Error("Incorrect Channel Message Count. Expected : ", expectedCount, "Actual : ", actualCount)
	}
}

func assertMessageRetreivedWithExpectedContent(t *testing.T, channel Channel, expectedContent string) {

	message, err := channel.Pop()

	if err != nil {
		t.Error("No error should have been thrown.")
	}

	if message.String() != expectedContent {
		t.Error("Incorrect content for message. Expected :", expectedContent, " Actual :", message.String())
	}
}
