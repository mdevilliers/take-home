package app

import (
	"errors"
	"log"

	"github.com/mdevilliers/take-home/pkg/topic"
)

var (
	UnknownTopic        = errors.New("Unknown topic")
	UnknownUser         = errors.New("Unknown user")
	NoMessagesAvailable = errors.New("No messages available for user")
)

// Service serializes access to topic registry, and topics
type Service struct {
	registry              topic.Registry
	subscribeChannel      chan *request
	unSubscribeChannel    chan *request
	publishMessageChannel chan *request
	getMessageChannel     chan *request
}

// Returns a new Service instance
func NewService() *Service {
	service := &Service{
		registry:              topic.NewTopicRegistry(),
		subscribeChannel:      make(chan *request),
		unSubscribeChannel:    make(chan *request),
		publishMessageChannel: make(chan *request),
		getMessageChannel:     make(chan *request),
	}
	go service.loop()
	return service
}

type request struct {
	topic           string
	user            string
	message         []byte
	responseChannel chan *response
}

type response struct {
	err     error
	message []byte
}

// subscribes a user to a topic
func (s *Service) Subscribe(topic string, username string) error {

	returnChannel := make(chan *response)
	request := &request{
		topic:           topic,
		user:            username,
		responseChannel: returnChannel,
	}

	go func() { s.subscribeChannel <- request }()
	response := <-returnChannel

	return response.err
}

// deletes a user subscription from a topic
func (s *Service) UnSubscribe(topic string, username string) error {

	returnChannel := make(chan *response)

	request := &request{
		topic:           topic,
		user:            username,
		responseChannel: returnChannel,
	}

	go func() { s.unSubscribeChannel <- request }()

	response := <-returnChannel

	return response.err
}

// allows publication of messages to an existing topic
func (s *Service) PublishMessage(topic string, message []byte) error {

	returnChannel := make(chan *response)
	request := &request{
		topic:           topic,
		message:         message,
		responseChannel: returnChannel,
	}

	go func() { s.publishMessageChannel <- request }()

	response := <-returnChannel
	return response.err
}

// retrieves messages from an existing topic for a user
func (s *Service) GetMessage(topic string, username string) ([]byte, error) {

	returnChannel := make(chan *response)
	request := &request{
		topic:           topic,
		user:            username,
		responseChannel: returnChannel,
	}

	go func() { s.getMessageChannel <- request }()

	response := <-returnChannel
	return response.message, response.err
}

func (s *Service) loop() {

	for {
		select {
		case subscribe := <-s.subscribeChannel:

			log.Print("Message recieved on subscribeChannel")

			topicToSubscribeTo := s.registry.Get(subscribe.topic)
			topicToSubscribeTo.AddChannel(subscribe.user)
			subscribe.responseChannel <- &response{err: nil}

		case unSubscribe := <-s.unSubscribeChannel:

			log.Print("Message recieved on unSubscribeChannel")

			exists := s.registry.Contains(unSubscribe.topic)

			if !exists {
				unSubscribe.responseChannel <- &response{err: UnknownTopic}
				break
			}

			existingTopic := s.registry.Get(unSubscribe.topic)
			err := existingTopic.RemoveChannel(unSubscribe.user)

			if err != nil {
				if err == topic.ChannelNotFoundError {
					unSubscribe.responseChannel <- &response{err: UnknownUser}
					break
				}

				// unexpected error
				unSubscribe.responseChannel <- &response{err: err}
				break
			}

			unSubscribe.responseChannel <- &response{err: nil}

		case getMessage := <-s.getMessageChannel:

			log.Print("Message recieved on getMessageChannel")

			exists := s.registry.Contains(getMessage.topic)

			if !exists {
				getMessage.responseChannel <- &response{err: UnknownTopic}
				break
			}

			topicToReadFrom := s.registry.Get(getMessage.topic)
			message, err := topicToReadFrom.GetNextMessage(getMessage.user)

			if err != nil {

				if err == topic.ChannelNotFoundError {
					getMessage.responseChannel <- &response{err: UnknownUser}
					break
				}

				if err == topic.NoMessagesAvailable {
					getMessage.responseChannel <- &response{err: NoMessagesAvailable}
					break
				}

				// unexpected error
				getMessage.responseChannel <- &response{err: err}
				break
			}

			getMessage.responseChannel <- &response{err: nil, message: message.Bytes()}

		case publishMessage := <-s.publishMessageChannel:

			log.Print("Message recieved on publishMessageChannel")

			topicToPostTo := s.registry.Get(publishMessage.topic)
			topicToPostTo.PublishMessage(topic.NewMessage(publishMessage.message))

			publishMessage.responseChannel <- &response{err: nil}

		}
	}
}
