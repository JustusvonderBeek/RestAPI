package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var config Configuration = Configuration{
	IP_Address:  "127.0.0.1",
	Listen_Port: "50002",
	Overwrite:   false,
	Client:      false,
}

func checkLinearIdIncrease(raw []byte) bool {
	var list []Word
	err := json.Unmarshal(raw, &list)
	if err != nil {
		log.Print("Failed to convert raw list to struct")
		return false
	}
	for idx := range list {
		if list[idx].ID != idx {
			return false
		}
	}
	return true
}

func compareWordLists(listExpected []Word, listGiven []Word) bool {
	if len(listExpected) != len(listGiven) {
		return false
	}
	for idx := range listExpected {
		if listExpected[idx] != listGiven[idx] {
			return false
		}
	}
	return true
}

func compareRawWordList(listExpected []byte, listGiven []byte) bool {
	var le []Word
	err := json.Unmarshal(listExpected, &le)
	if err != nil {
		log.Print("Failed to convert expected list to struct")
		return false
	}
	var lg []Word
	err = json.Unmarshal(listGiven, &lg)
	if err != nil {
		log.Print("Failed to convert expected list to struct")
		return false
	}
	return compareWordLists(le, lg)
}

func compareAddedCorrectly(originalList []byte, alteredList []byte, word Word) bool {
	var orgwordlist []Word
	err := json.Unmarshal(originalList, &orgwordlist)
	if err != nil {
		log.Print("Failed to convert orignal list to struct")
		return false
	}
	var wordlist []Word
	err = json.Unmarshal(alteredList, &wordlist)
	if err != nil {
		log.Print("Failed to convert modified list to struct")
		return false
	}
	length := len(wordlist)
	if length == 0 {
		return false
	}
	if len(orgwordlist) != len(wordlist)-1 {
		return false
	}
	for idx := range orgwordlist {
		if orgwordlist[idx] != wordlist[idx] {
			return false
		}
	}
	addedWord := wordlist[length-1]
	if addedWord.Vocabulary != word.Vocabulary {
		return false
	}
	if addedWord.Translation != word.Translation {
		return false
	}
	return true
}

func compareRemovedCorrectly(originalList []byte, alteredList []byte, removedIndex int) bool {
	var orgList, altList []Word
	err := json.Unmarshal(originalList, &orgList)
	if err != nil {
		log.Print("Failed to convert original list to struct")
		return false
	}
	err = json.Unmarshal(alteredList, &altList)
	if err != nil {
		log.Print("Failed to convert altered list to struct")
		return false
	}
	if len(orgList) != len(altList)+1 {
		log.Print("Lengths do not match")
		return false
	}
	for idx := range altList {
		if idx < removedIndex {
			if orgList[idx] != altList[idx] {
				return false
			}
		} else {
			if orgList[idx+1].Vocabulary != altList[idx].Vocabulary || orgList[idx+1].Translation != altList[idx].Translation {
				log.Printf("Expected: %+v Got: %+v", orgList[idx+1], altList[idx])
				return false
			}
		}
	}
	return true
}

func TestAddWord(t *testing.T) {
	addr := config.IP_Address + ":" + config.Listen_Port
	url := "https://" + addr + "/words"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	unmodified, err := os.ReadFile("vocabulary.json")
	if err != nil {
		log.Print("Failed to read old vocabulary file!")
		t.FailNow()
	}
	var originalList []Word
	err = json.Unmarshal(unmodified, &originalList)
	if err != nil {
		log.Print("Existing vocabulary is not in JSON format")
		t.FailNow()
	}
	newVocab := Word{
		ID:          1337,
		Vocabulary:  strconv.Itoa(len(originalList)),
		Translation: strconv.Itoa(len(originalList) + 1),
	}
	raw, err := json.Marshal(newVocab)
	if err != nil {
		log.Print("Failed to convert vocab to JSON format")
		t.FailNow()
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(raw))
	if err != nil {
		log.Print("Failed to post data to url")
		t.FailNow()
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Failed to post request")
		t.FailNow()
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Failed to read response body")
		t.FailNow()
	}
	expected, err := os.ReadFile("vocabulary.json")
	if err != nil {
		log.Print("Failed to read comparison file!")
		t.FailNow()
	}
	equal := compareRawWordList(expected, body)
	if !equal {
		log.Print("Lists are not equal!")
		t.FailNow()
	}
	equal = compareAddedCorrectly(unmodified, expected, newVocab)
	if !equal {
		log.Print("Word not added correctly")
		t.FailNow()
	}
	equal = checkLinearIdIncrease(body)
	if !equal {
		log.Print("IDs do not increase linear")
		t.FailNow()
	}
}

func TestRemoveWord(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	addr := config.IP_Address + ":" + config.Listen_Port
	// Removing the last word (just added)
	currentList, err := os.ReadFile("vocabulary.json")
	if err != nil {
		log.Print("Failed to read vocabulary file!")
		t.FailNow()
	}
	var wordlist []Word
	err = json.Unmarshal(currentList, &wordlist)
	if err != nil {
		log.Print("Failed to convert expected list to struct")
		return
	}
	removeIndex := 0
	if len(wordlist) == 0 {
		url := "https://" + addr + "/words/0"
		word := Word{
			ID:          0,
			Vocabulary:  "0",
			Translation: "1",
		}
		cont, _ := json.Marshal(word)
		req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(cont))
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
		if resp.StatusCode != http.StatusBadRequest {
			log.Print("word not correct")
			t.FailNow()
		}
		return
	}
	if len(wordlist) > 1 {
		removeIndex = int(math.Ceil(float64(len(wordlist)) / 2.0))
		log.Printf("Removing index: %d", removeIndex)
	}
	url := "https://" + addr + "/words/" + strconv.Itoa(removeIndex)
	removeWord := wordlist[removeIndex]
	log.Printf("Removing word: %+v", removeWord)
	raw, err := json.Marshal(removeWord)
	if err != nil {
		log.Print("Failed to convert remove word to JSON")
		t.FailNow()
	}
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(raw))
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
	if resp.StatusCode == http.StatusBadRequest {
		log.Print("word not correct")
		t.FailNow()
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body")
		t.FailNow()
	}
	alteredList, err := os.ReadFile("vocabulary.json")
	if err != nil {
		log.Print("Failed to read vocabulary file!")
		t.FailNow()
	}
	equal := compareRawWordList(alteredList, body)
	if !equal {
		log.Print("Different state at server and client")
		t.FailNow()
	}
	equal = compareRemovedCorrectly(currentList, body, removeIndex)
	if !equal {
		log.Print("Item not correctly removed")
		t.FailNow()
	}
	log.Printf("Response: %s", string(body))
}
