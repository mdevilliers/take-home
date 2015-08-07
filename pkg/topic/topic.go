package topic

import (
	"errors"
	"sync"
)

var (
	ChannelNotFoundError = errors.New("Channel not found")
)

// A Topic is the 'broker' type object distrubuting messages to connected Channels
// Safe for use via goroutines
type Topic struct {
	sync.RWMutex
	channels map[string]Channel
	name     string
}

func NewTopic(name string) *Topic {
	return &Topic{
		name:     name,
		channels: make(map[string]Channel),
	}
}

// Adds a channel to a topic. If it doesn't exist a channel is created for the topic
func (t *Topic) AddChannel(channelName string) {
	t.Lock()
	defer t.Unlock()

	_, exists := t.channels[channelName]

	if !exists {
		t.channels[channelName] = NewChannel()
	}
}

// Test for whether a specific channel exists
func (t *Topic) ChannelExists(channelName string) bool {

	t.Lock()
	defer t.Unlock()

	_, exists := t.channels[channelName]
	return exists
}

// Removes a channel for a channel. If the channel does not exist returns a ChannelNotFoundError
func (t *Topic) RemoveChannel(channelName string) error {

	t.Lock()
	defer t.Unlock()

	_, exists := t.channels[channelName]

	if !exists {
		return ChannelNotFoundError
	}

	delete(t.channels, channelName)

	return nil
}

// Appends message to all known channels
func (t *Topic) PublishMessage(message *Message) {

	t.Lock()
	defer t.Unlock()

	for _, channel := range t.channels {
		channel.Push(message)
	}
}

// Returns the next message for the channel. If the channel does not exist returns a ChannelNotFoundError
func (t *Topic) GetNextMessage(channelName string) (*Message, error) {

	t.Lock()
	defer t.Unlock()

	channel, exists := t.channels[channelName]

	if !exists {
		return nil, ChannelNotFoundError
	}

	return channel.Pop()
}
