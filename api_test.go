package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
)

var config Configuration = Configuration{
	IP_Address:  "127.0.0.1",
	Listen_Port: "50002",
	Overwrite:   false,
	Client:      false,
}

func TestAddWord(t *testing.T) {
	addr := config.IP_Address + ":" + config.Listen_Port
	url := "https://" + addr + "/words"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	newVocab := Word{
		ID:          1,
		Vocabulary:  "Test",
		Translation: "Ein Test",
	}
	raw, err := json.Marshal(newVocab)
	if err != nil {
		log.Print("Failed to convert vocab to JSON format")
		t.FailNow()
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(raw))
	if err != nil {
		log.Print("Failed to request url")
		t.FailNow()
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Failed to request")
		t.FailNow()
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Failed to read response body")
		t.FailNow()
	}
	log.Printf("Response: %s", string(body))
}

func TestRemoveWord(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	addr := config.IP_Address + ":" + config.Listen_Port
	url := "https://" + addr + "/words/1"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("Failed to request url: %s", err)
		t.FailNow()
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Failed to request")
		t.FailNow()
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body")
		t.FailNow()
	}
	log.Printf("Response: %s", string(body))
}
