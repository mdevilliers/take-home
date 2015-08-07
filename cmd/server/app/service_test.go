package app

import (
	"testing"
)

// Path 1
// 'user-1' subscribes to 'topic-one'
// 'message-one' published to 'topic-one'
// 'user-1' gets next message on 'topic-one'
// 'user-1' expects this to be 'message-one'
func TestPath1(t *testing.T) {

	service := NewService()

	err := service.Subscribe("topic-one", "user-1")

	if err != nil {
		t.Error("Subscribe unsuccesful", err.Error())
	}

	err = service.PublishMessage("topic-one", []byte("message-one"))

	if err != nil {
		t.Error("PublishMessage unsuccesful", err.Error())
	}

	message, err := service.GetMessage("topic-one", "user-1")

	if string(message) != "message-one" {
		t.Error("Expected 'message-one'")
	}

}

// Path 2
// 'message-one' published to 'topic-one'
// expects this to error with unknown subscription
func TestPath2(t *testing.T) {

	service := NewService()

	err := service.PublishMessage("topic-one", []byte("message-one"))

	if err != UnknownTopic {
		t.Error("PublishMessage should have errored with UnknownTopic :", err.Error())
	}

}

// Path3
// 'user-1' subscribes to 'topic-one'
// 'user-2' subscribes to 'topic-one'
// 'message-one' published to 'topic-one'
// 'user-1' gets next message on 'topic-one'
// 'user-2' gets next message on 'topic-one'
//  both messages are 'message-one'
func TestPath3(t *testing.T) {

	service := NewService()
	service.Subscribe("topic-one", "user-1")
	service.Subscribe("topic-one", "user-2")

	service.PublishMessage("topic-one", []byte("message-one"))

	message1, _ := service.GetMessage("topic-one", "user-1")
	message2, _ := service.GetMessage("topic-one", "user-2")

	if string(message1) != string(message2) {
		t.Error("user-1 and user-2 should have recieved the same message")
	}

	if string(message1) != "message-one" {
		t.Error("user-1 should have recieved the same message they sent")
	}
}

// Path4
// 'user-1' subscribes to 'topic-one'
// 'user-2' unsubscribes to 'topic-one'
func TestPath4(t *testing.T) {

	service := NewService()
	service.Subscribe("topic-one", "user-1")

	err := service.UnSubscribe("topic-one", "user-1")

	if err != nil {
		t.Error("UnSubscribe unsuccessful : ", err.Error())
	}
}
