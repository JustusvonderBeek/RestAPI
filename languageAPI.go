package main

import (
	"bytes"
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

func swapExistingVocabulary() {
	log.Print("Swapping the existing vocabulary file...")
	vocab := "vocabulary.json"
	_, err := os.Open(vocab)
	if err != nil {
		log.Print("Vocabulary does not exist.")
		return
	}
	counter := 1
	filename := "vocabulary_" + strconv.Itoa(counter) + ".json"
	for {
		_, err = os.Open(filename)
		if err != nil {
			break
		}
		counter += 1
		filename = "vocabulary_" + strconv.Itoa(counter) + ".json"
	}
	cnt, err := os.ReadFile(vocab)
	if err != nil {
		log.Fatal("Failed to read the content of the existing vocabulary!")
	}
	os.WriteFile(filename, cnt, os.ModePerm)
	os.Create(vocab)
}

func getData(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, vocabulary)
}

func postData(c *gin.Context) {
	var newVocab Word
	if err := c.BindJSON(&newVocab); err != nil {
		return
	}

	vocabulary = append(vocabulary, newVocab)
	c.IndentedJSON(http.StatusCreated, vocabulary)
	rawData, err := json.MarshalIndent(vocabulary, "", "\t")
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

func modifyDataItem(c *gin.Context) {
	id := c.Param("id")
	compare, _ := strconv.Atoi(id)
	var updatedWord Word
	err := c.ShouldBindJSON(&updatedWord)
	if err != nil {
		log.Print("Failed to bind to Word")
		return
	}

	for idx, word := range vocabulary {
		if word.ID == compare {
			vocabulary[idx].Vocabulary = updatedWord.Vocabulary
			vocabulary[idx].Translation = updatedWord.Translation
		}
	}
	log.Printf("Updated %d to %x", compare, updatedWord)
	rawData, err := json.MarshalIndent(vocabulary, "", "\t")
	if err != nil {
		log.Print("Failed to process internal data!")
		return
	}
	writeData(rawData)
}

func startingServer(cfg Configuration) error {
	if cfg.Overwrite {
		swapExistingVocabulary()
	}
	vocabulary = readData()

	router := gin.Default()
	router.GET("/words", getData)
	router.POST("words", postData)
	router.GET("/words/:id", getDataItem)
	router.POST("/words/:id", modifyDataItem)

	address := cfg.IP_Address + ":" + cfg.Listen_Port
	// router.Run(address)
	router.RunTLS(address, "vocabulary.cer", "vocabulary.key")

	return nil
}

// ---------------------------------------------------------
// CLIENT
// ---------------------------------------------------------

func getVocabulary(cfg Configuration, client *http.Client) {
	addr := cfg.IP_Address + ":" + cfg.Listen_Port
	url := "http://" + addr + "/words"
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
	url := "http://" + addr + "/words"

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

func startingClient(cfg Configuration) error {
	if cfg.Overwrite {
		swapExistingVocabulary()
	}
	client := &http.Client{}

	getVocabulary(cfg, client)
	putVocabulary(cfg, client)
	getVocabulary(cfg, client)

	return nil
}
