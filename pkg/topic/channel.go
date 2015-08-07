package topic

import (
	"errors"
	"sync"
)

var (
	NoMessagesAvailable = errors.New("No messages available in channel")
)

// A Channel is a store of messages for a user. 
type Channel interface {
	Push(message *Message)
	Pop() (*Message, error)
	Count() int
}

type Channels []*Channel

// Create a Channel with an InMemory implementation
// Safe for use via a goroutine
// WARNING : message store length is unbounded
func NewChannel() Channel {
	return &InMemoryChannel{
		messageStore: make([]*Message, 0),
		messageCount: 0,
	}
}

type InMemoryChannel struct {
	sync.RWMutex
	// REVIEW : look at using a RingBuffer with a fixed length rather than an unbounded array
	messageStore []*Message
	messageCount int
}

// Pushes a message to the store
func (c *InMemoryChannel) Push(message *Message) {

	c.Lock()
	defer c.Unlock()

	c.messageStore = append(c.messageStore, message)
	c.messageCount++
}

// Pops a message from the store
func (c *InMemoryChannel) Pop() (*Message, error) {

	c.Lock()
	defer c.Unlock()

	length := c.messageCount

	if length > 0 {
		message := c.messageStore[0]
		c.messageStore = c.messageStore[1:length]
		c.messageCount--
		return message, nil
	}
	return nil, NoMessagesAvailable
}

// Current count of messages waiting to be delivered
func (c *InMemoryChannel) Count() int {

	c.Lock()
	defer c.Unlock()

	return c.messageCount
}
