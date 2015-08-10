package app

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zenazn/goji/web"
)

// Subscribe to a topic
// Request: POST /<topic>/<username>
// Response codes:
// ● 200: Subscription succeeded.
func TestSubscribeToTopicReturns200(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()

	//POST /<topic>/<username>
	res, _ := http.Post(instance.URL+"/topic-one/user-one", "text", nil)

	_, status := parseResponse(res)
	if status != http.StatusOK {
		t.Error("Subscribing should return 200 but returned ", status)
	}

}

// Unsubscribe from a topic
// Request: DELETE /<topic>/<username>
// Response codes:
// ● 200: Unsubscribe succeeded.
// ● 404: The subscription does not exist.
func TestUnSubscribeFromTopicReturns200(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()

	url := instance.URL + "/topic-one/user-one"

	//subscribe
	http.Post(url, "text", nil)

	//unsubscribe
	req, _ := http.NewRequest("DELETE", url, nil)
	res, _ := http.DefaultClient.Do(req)

	_, status := parseResponse(res)
	if status != http.StatusOK {
		t.Error("Subscribing then unsubscribing should return 200 but returned ", status)
	}
}

func TestUnSubscribeFromUnknownTopicReturns404(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()

	url := instance.URL + "/topic-one/user-one"

	//unsubscribe
	req, _ := http.NewRequest("DELETE", url, nil)
	res, _ := http.DefaultClient.Do(req)

	_, status := parseResponse(res)
	if status != 404 {
		t.Error("Unsubscribing without subscribing should return 404 but returned ", status)
	}
}

// Publish a message
// Request: POST /<topic>
// Request body: The message being published.
// Response codes:
// ● 200: Publish succeeded.
func TestPublishToExistingTopicReturns200(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()

	url := instance.URL + "/topic-one/user-one"

	//POST /<topic>/<username>
	http.Post(url, "text", nil)

	//POST /<topic>
	res, _ := http.Post(instance.URL+"/topic-one", "text", bytes.NewBuffer([]byte("message-one")))

	_, status := parseResponse(res)
	if status != http.StatusOK {
		t.Error("Posting to an existing topic should return 200 but returned ", status)
	}

}

func TestPublishToNonExistingTopicReturns200(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()

	//POST /<topic>
	res, _ := http.Post(instance.URL+"/topic-one", "text", bytes.NewBuffer([]byte("message-one")))

	_, status := parseResponse(res)

	if status != 200 {
		t.Error("Posting to an non-existing topic should return 200 but returned ", status)
	}

}

// Request: GET /<topic>/<username>
// Response codes:
// ● 200: Retrieval succeeded.
// ● 204: There are no messages available for this topic on this user.
// ● 404: The subscription does not exist.
// Response body: The body of the next message, if one exists.
func TestRetrieveFromExistingTopicReturns200(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()

	url := instance.URL + "/topic-one/user-one"

	//POST /<topic>/<username>
	http.Post(url, "text", nil)

	//POST /<topic>
	http.Post(instance.URL+"/topic-one", "text", bytes.NewBuffer([]byte("message-one")))

	//GET /<topic>/<username>
	res, _ := http.Get(url)

	content, status := parseResponse(res)

	if status != http.StatusOK {
		t.Error("Retreiving a posted message should return 200 but returned ", status)
	}

	if content != "message-one" {
		t.Error("Wrong message returned ", content)
	}

}

func TestRetrieveFromExistingTopicWithNoMessagesReturns204(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()
	url := instance.URL + "/topic-one/user-one"

	//POST /<topic>/<username>
	http.Post(url, "text", nil)

	//GET /<topic>/<username>
	res, _ := http.Get(url)

	_, status := parseResponse(res)

	if status != 204 {
		t.Error("Trying to retreive a non-existant message should return 204 but returned ", status)
	}

}

func TestRetrieveFromNonExistingTopicReturns404(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()

	url := instance.URL + "/topic-one/user-one"

	//GET /<topic>/<username>
	res, _ := http.Get(url)

	_, status := parseResponse(res)

	if status != 404 {
		t.Error("Trying to retreive from non-existant topic should return 404 but returned ", status)
	}

}

func TestRetrieveFromExistingTopicButNonExistantUserReturns404(t *testing.T) {

	instance := getServerInstance()
	defer instance.Close()

	//POST /<topic>/<username>
	http.Post(instance.URL+"/topic-one/user-one", "text", nil)

	//GET /<topic>/<username>
	res, _ := http.Get(instance.URL + "/topic-one/does-not-exist")

	_, status := parseResponse(res)

	if status != 404 {
		t.Error("Trying to retreive from existant topic but non-existant user should return 404 but returned ", status)
	}

}

func getServerInstance() *httptest.Server {
	api := NewApi()
	mux := web.New()
	api.Route(mux)

	return httptest.NewServer(mux)
}

func parseResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(contents), res.StatusCode
}
