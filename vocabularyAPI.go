package main

import (
	"encoding/json"
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

// -------------------------------------------------------------------------------
// Auxiliary Functions
// -------------------------------------------------------------------------------

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

func saveVocabulary() {
	rawData, err := json.MarshalIndent(vocabulary, "", "\t")
	if err != nil {
		log.Print("Failed to convert data to JSON!")
		return
	}
	writeData(rawData)
}

func fixIndexing(list *[]Word) {
	// log.Print("Fixing the indexing")
	for idx := range *list {
		(*list)[idx].ID = idx
	}
}

func readData() []Word {
	filename := "vocabulary.json"
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("No vocabulary found. Creating new one...")
		return []Word{}
	}
	if string(content) == "" {
		return []Word{}
	}
	var vocabulary []Word
	err = json.Unmarshal(content, &vocabulary)
	if err != nil {
		log.Print("The given file does not contain a valid vocabulary!")
		return []Word{}
	}
	// Fixing the indexing
	fixIndexing(&vocabulary)
	saveVocabulary()
	return vocabulary
}

func swapExistingVocabulary() {
	log.Print("Swapping the existing vocabulary file...")
	vocab := "vocabulary.json"
	_, err := os.Open(vocab)
	if err != nil {
		log.Print("Vocabulary file does not exist")
		return
	}
	content, err := os.ReadFile("vocabulary.json")
	if err != nil {
		log.Printf("Vocabulary seems to be corrupted: %s", err)
		return
	}
	if string(content) == "" {
		log.Print("Vocabulary is empty")
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

// -------------------------------------------------------------------------------
// API Implementation
// -------------------------------------------------------------------------------

func getData(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, vocabulary)
}

func postData(c *gin.Context) {
	var newVocab Word
	if err := c.BindJSON(&newVocab); err != nil {
		return
	}

	// Fixing the ID in the received Word
	newVocab.ID = len(vocabulary)
	vocabulary = append(vocabulary, newVocab)
	c.IndentedJSON(http.StatusCreated, vocabulary)
	saveVocabulary()
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
	saveVocabulary()
}

func removeDataItem(c *gin.Context) {
	id := c.Param("id")
	compare, _ := strconv.Atoi(id)
	for idx, word := range vocabulary {
		if word.ID == compare {
			if idx == 0 {
				vocabulary = vocabulary[idx+1:]
				break
			} else {
				vocabulary = append(vocabulary[:idx-1], vocabulary[idx+1:]...)
				break
			}
		}
	}
	log.Printf("Removed item at index %d", compare)
	saveVocabulary()
	c.IndentedJSON(http.StatusOK, vocabulary)
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
	router.DELETE("/words/:id", removeDataItem)

	address := cfg.IP_Address + ":" + cfg.Listen_Port
	// router.Run(address)
	router.RunTLS(address, "vocabulary.cer", "vocabulary.key")

	return nil
}
