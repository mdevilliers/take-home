package topic

import (
	"errors"
	"sync"
)

var (
	UnknownTopic = errors.New("Unknown Topic")
)

type Registry interface {
	Delete(topicName string) error
	Contains(topicName string) bool
	Get(topicName string) *Topic
}

// Registry implementatoin which maintains an in memory index
type InMemoryRegistry struct {
	sync.RWMutex
	topics map[string]*Topic
}

// Returns an instance of the InMemoryRegistry
func NewTopicRegistry() Registry {
	return &InMemoryRegistry{
		topics: make(map[string]*Topic),
	}
}

// removes an existing topic
func (r *InMemoryRegistry) Delete(topicName string) error {
	r.Lock()
	defer r.Unlock()

	if !r.exists(topicName) {
		return UnknownTopic
	}

	delete(r.topics, topicName)

	return nil

}

// Return true if a topic exists in the registry
func (r *InMemoryRegistry) Contains(topicName string) bool {

	r.Lock()
	defer r.Unlock()

	return r.exists(topicName)
}

// Finds a Topic of a given name. If the topic does not exists it create it.
func (r *InMemoryRegistry) Get(topicName string) *Topic {
	r.Lock()
	defer r.Unlock()

	if !r.exists(topicName) {
		r.topics[topicName] = NewTopic(topicName)
	}
	return r.topics[topicName]

}

// only to be called when locked
func (r *InMemoryRegistry) exists(topicName string) bool {
	_, exists := r.topics[topicName]
	return exists
}
