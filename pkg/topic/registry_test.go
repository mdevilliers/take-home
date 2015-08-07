package topic

import (
	"testing"
)

func TestDeletingNonExistentTopicErrors(t *testing.T) {

	registry := NewTopicRegistry()

	err := registry.Delete("does-not-exist")

	if err != UnknownTopic {
		t.Error("A UnknownTopic error should have been returned.")
	}
}

func TestAddedTopicCanBeRemoved(t *testing.T) {

	registry := NewTopicRegistry()
	topic := registry.Get("exists")

	if topic == nil {
		t.Error("Registry should always add a new topic.")
	}

	err := registry.Delete("exists")

	if err != nil {
		t.Error("Topic should have been deleted.")
	}

	err = registry.Delete("exists")

	if err != UnknownTopic {
		t.Error("A UnknownTopic error should have been returned.")
	}
}

func TestTopicRegistryReturnSameTopic(t *testing.T) {

	registry := NewTopicRegistry()

	topic1 := registry.Get("topic-one")
	topic2 := registry.Get("topic-one")

	if topic1 != topic2 {
		t.Error("Repeated Gets should return the same object")
	}
}

func TestAnAddedTopicExistsInTheRegistry(t *testing.T) {
	registry := NewTopicRegistry()

	registry.Get("exists")

	if !registry.Contains("exists") {
		t.Error("Registry should contain the topic called 'exists'")
	}
	registry.Delete("exists")

	if registry.Contains("exists") {
		t.Error("Registry should not contain the topic called 'exists'")
	}
}
