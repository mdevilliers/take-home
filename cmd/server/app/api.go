package app

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
)

type Api struct {
	service *Service
}

func NewApi() *Api {
	return &Api{
		service: NewService(),
	}
}

// Sets up the routes
func (api *Api) Route(m *web.Mux) {

	m.Get("/:topic/:username", api.NextMessage)
	m.Post("/:topic/:username", api.SubscribeToTopic)
	m.Delete("/:topic/:username", api.UnsubscribeFromTopic)

	m.Post("/:topic", api.PublishMessage)

}

//  POST /<topic>/<username>
func (api *Api) SubscribeToTopic(c web.C, w http.ResponseWriter, r *http.Request) {

	topicFromRequest := c.URLParams["topic"]
	usernameFromRequest := c.URLParams["username"]

	if isEmptyString(topicFromRequest) || isEmptyString(usernameFromRequest) {
		w.WriteHeader(500)
		return
	}

	log.Println("SubscribeToTopic : topic", topicFromRequest, "username", usernameFromRequest)

	err := api.service.Subscribe(topicFromRequest, usernameFromRequest)

	if err == nil {
		w.WriteHeader(200)

	} else {
		
		w.WriteHeader(500)
		log.Print("SubscribeToTopic : unexpected error : ", err.Error())
	}

	log.Println("SubscribeToTopic : finsihed topic", topicFromRequest, "username", usernameFromRequest)

}

// DELETE /<topic>/<username>
func (api *Api) UnsubscribeFromTopic(c web.C, w http.ResponseWriter, r *http.Request) {

	topicFromRequest := c.URLParams["topic"]
	usernameFromRequest := c.URLParams["username"]

	if isEmptyString(topicFromRequest) || isEmptyString(usernameFromRequest) {
		w.WriteHeader(500)
		return
	}

	log.Println("UnsubscribeFromTopic : topic", topicFromRequest, "username", usernameFromRequest)

	err := api.service.UnSubscribe(topicFromRequest, usernameFromRequest)

	if err != nil {

		if err == UnknownTopic || err == UnknownUser {

			w.WriteHeader(404)

		} else {

			// unexpected error
			log.Print("UnsubscribeFromTopic : unexpected error : ", err.Error())
			w.WriteHeader(500)

		}

	} else {
		w.WriteHeader(200)
	}
}

// POST /<topic>
func (api *Api) PublishMessage(c web.C, w http.ResponseWriter, r *http.Request) {

	topicFromRequest := c.URLParams["topic"]
	messageFromRequest, err := ioutil.ReadAll(r.Body)

	if err != nil {

		log.Print("PublishMessage : error parsing body : ", err.Error())
		w.WriteHeader(500)
		return
	}

	if isEmptyString(topicFromRequest) || len(messageFromRequest) == 0 {
		w.WriteHeader(500)
		return
	}

	log.Println("PublishMessage : topic", topicFromRequest, "message", messageFromRequest)
	err = api.service.PublishMessage(topicFromRequest, messageFromRequest)

	if err != nil {

		if err == UnknownTopic {
			w.WriteHeader(404)

		} else {

			// unexpected error
			log.Print("PublishMessage : unexpected error : ", err.Error())
			w.WriteHeader(500)

		}

		return

	} else {
		w.WriteHeader(200)
	}

}

// GET /<topic>/<username>
func (api *Api) NextMessage(c web.C, w http.ResponseWriter, r *http.Request) {

	topicFromRequest := c.URLParams["topic"]
	usernameFromRequest := c.URLParams["username"]

	if isEmptyString(topicFromRequest) || isEmptyString(usernameFromRequest) {
		w.WriteHeader(500)
		return
	}

	log.Println("NextMessage : topic", topicFromRequest, "username", usernameFromRequest)

	message, err := api.service.GetMessage(topicFromRequest, usernameFromRequest)

	if err != nil {

		if err == UnknownUser || err == UnknownTopic {
			w.WriteHeader(404)
			return
		}

		if err == NoMessagesAvailable {
			w.WriteHeader(204)
			return
		}

		// unexpected error
		w.WriteHeader(500)
		return
	}

	// REVIEW : html encode response? my feeling is no as this is a "service" rather than a front end web application
	io.WriteString(w, string(message))
	w.WriteHeader(200)
}

func isEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
