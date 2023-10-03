package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// ---------------------------------------------------------
// CLIENT
// ---------------------------------------------------------

func getVocabulary(cfg Configuration, client *http.Client) {
	addr := cfg.IP_Address + ":" + cfg.Listen_Port
	url := "https://" + addr + "/words"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to request url")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body")
	}

	log.Printf("Response: %s", string(body))
}

func putVocabulary(cfg Configuration, client *http.Client) {
	addr := cfg.IP_Address + ":" + cfg.Listen_Port
	url := "https://" + addr + "/words"

	newVocab := Word{
		ID:          1,
		Vocabulary:  "Test",
		Translation: "Ein Test",
	}
	raw, err := json.Marshal(newVocab)
	if err != nil {
		log.Fatal("Failed to convert vocab to JSON format")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(raw))
	if err != nil {
		log.Fatal("Failed to request url")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body")
	}

	log.Printf("Response: %s", string(body))
}

func removeVocabulary(cfg Configuration, client *http.Client) {
	addr := cfg.IP_Address + ":" + cfg.Listen_Port
	url := "https://" + addr + "/words/1"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Failed to request url: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body")
	}

	log.Printf("Response: %s", string(body))
}

func startingClient(cfg Configuration) error {
	if cfg.Overwrite {
		swapExistingVocabulary()
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	getVocabulary(cfg, client)
	putVocabulary(cfg, client)
	getVocabulary(cfg, client)
	removeVocabulary(cfg, client)

	return nil
}
