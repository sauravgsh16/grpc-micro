package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	address := flag.String("server", "http://localhost:8080", "HTTP gateway url")
	flag.Parse()

	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	var body string

	// Create Call
	resp, err := http.Post(*address+"/v1/todo", "application/json", strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
			"toDo": {
				"title":"title (%s)",
				"description":"description (%s)",
				"reminder":"%s"
			}
		}
	`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("failed to call Create method: %v", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to Create response body: %v", err)
	} else {
		body = string(b)
	}
	log.Printf("Create resp: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Parse fields returned from create resp body
	var created struct {
		API string `json:"api"`
		ID  string `json:"id"`
	}
	err = json.Unmarshal(b, &created)
	if err != nil {
		log.Fatalf("failed to unmarshal json of create resp: %v", err)
	}

	// Call Read
	resp, err = http.Get(fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID))
	if err != nil {
		log.Fatalf("failed to call Read method: %v", err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("Failed to create resp body for READ: %v", err)
	} else {
		body = string(b)
	}
	log.Printf("Read resp: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Update Call
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID),
		strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
			"todo": {
				"title":"title (%s) + updated",
				"decription":"description (%s) + updated",
				"reminder":"%s"
			}
		}
		
	`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call UPDATE method %v", err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read Update resp: %v", err)
	} else {
		body = string(b)
	}
	fmt.Printf("Update resp: Code=%d, Body=%v\n\n", resp.StatusCode, body)

	// Read All Call

	url := fmt.Sprintf("%s%s", *address, "/v1/todo/all")
	log.Printf("%s\n", url)
	resp, err = http.Get(url) //fmt.Sprintf("%s%s", *address, "/v1/todo/al"))
	if err != nil {
		log.Fatalf("failed to call ReadAll method: %v", err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read ReadAll response body: %v", err)
	} else {
		body = string(b)
	}
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Delete Call

	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Delete method: %v", err)
	}
	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Delete response body: %v", err)
	} else {
		body = string(b)
	}
	log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
