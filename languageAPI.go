package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Word struct {
	ID          int
	Vocabulary  string
	Translation string
}

var vocabulary = []Word{}

func getData(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, vocabulary)
}

func postData(c *gin.Context) {
	var newVocab Word
	if err := c.BindJSON(&newVocab); err != nil {
		return
	}

	totalData := append(vocabulary, newVocab)
	c.IndentedJSON(http.StatusCreated, totalData)
	rawData, err := json.Marshal(totalData)
	if err != nil {
		log.Printf("Failed to write data to disk because of: %s", err)
	}
	writeData(rawData)
}

func getDataItem(c *gin.Context) {
	id := c.Param("id")
	compare, _ := strconv.Atoi(id)

	for _, word := range vocabulary {
		if word.ID == compare {
			c.IndentedJSON(http.StatusOK, word)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "word not found"})
}

func writeData(data []byte) error {
	file := "vocabulary.json"
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Failed to create file \"%s\"", file)
		return err
	}
	// Don't forget to close the file
	defer f.Close()
	f.Write(data)

	return nil
}

func readData() []Word {
	filename := "vocabulary.json"
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("No vocabulary found. Creating new one...")
		return []Word{}
	}
	var vocabulary []Word
	err = json.Unmarshal(content, &vocabulary)
	if err != nil {
		log.Print("The given file does not contain a valid vocabulary!")
		return []Word{}
	}
	return vocabulary
}

func startingServer(cfg Configuration) error {
	vocabulary = readData()

	router := gin.Default()
	router.GET("/words", getData)
	router.POST("words", postData)
	router.GET("/words/:id", getDataItem)

	address := cfg.IP_Address + ":" + cfg.Listen_Port
	router.Run(address)

	return nil
}

func startingClient(cfg Configuration) error {
	client := &http.Client{}

	addr := cfg.IP_Address + ":" + cfg.Listen_Port
	url := "http://" + addr + "/albums"
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

	return nil
}
