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

func saveVocabulary(vocab *[]Word) {
	log.Print("Storing the vocabulary")
	// Do this every time due to wrong read or remove operation
	fixIndexing(vocab)
	rawData, err := json.MarshalIndent(*vocab, "", "\t")
	if err != nil {
		log.Print("Failed to convert data to JSON!")
		return
	}
	writeData(rawData)
}

func fixIndexing(list *[]Word) {
	log.Print("Fixing the indexing")
	for idx := range *list {
		(*list)[idx].ID = idx
	}
}

func readData() []Word {
	log.Print("Reading existing vocabulary")
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
	if len(vocabulary) > 10 {
		log.Printf("Loaded vocabulary:\n%+v", vocabulary[:10])
	} else {
		log.Printf("Loaded vocabulary:\n%+v", vocabulary)
	}
	saveVocabulary(&vocabulary)
	return vocabulary
}

func swapExistingVocabulary() {
	log.Print("Swapping the existing vocabulary file")
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
	saveVocabulary(&vocabulary)
	c.IndentedJSON(http.StatusCreated, vocabulary)
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
	if compare != updatedWord.ID {
		log.Print("incorrect word id and url id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ""})
		return
	}
	if compare >= len(vocabulary) || compare < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given index does not exist"})
		return
	}

	vocabulary[compare].Vocabulary = updatedWord.Vocabulary
	vocabulary[compare].Translation = updatedWord.Translation

	log.Printf("Updated %d to %x", compare, updatedWord)
	saveVocabulary(&vocabulary)
	c.IndentedJSON(http.StatusCreated, vocabulary)
}

func removeDataItem(c *gin.Context) {
	id := c.Param("id")
	compare, _ := strconv.Atoi(id)
	var removeWord Word
	err := c.ShouldBindJSON(&removeWord)
	if err != nil {
		log.Print("given body does not contain valid word")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ""})
		return
	}
	if compare != removeWord.ID {
		log.Print("incorrect word id and url id")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ""})
		return
	}
	if compare >= len(vocabulary) || compare < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given index does not exist"})
		return
	}

	// log.Printf("Full Vocab: %+v", vocabulary)

	wordToRemove := vocabulary[compare]
	if wordToRemove != removeWord {
		log.Print("remove vocab word does not match")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ""})
		return
	}

	vocabulary = append(vocabulary[:compare], vocabulary[compare+1:]...)

	log.Printf("Removed item at index %d", compare)
	// log.Printf("Full Vocab: %+v", vocabulary)
	saveVocabulary(&vocabulary)
	c.IndentedJSON(http.StatusOK, vocabulary)
}

func startingServer(cfg Configuration) error {
	gin.SetMode(gin.DebugMode)
	if cfg.Overwrite {
		swapExistingVocabulary()
	}
	vocabulary = readData()

	router := gin.Default()
	router.GET("/words", getData)
	router.GET("/words/:id", getDataItem)
	router.POST("words", postData)
	router.POST("/words/:id", modifyDataItem)
	router.DELETE("/words/:id", removeDataItem)

	address := cfg.IP_Address + ":" + cfg.Listen_Port
	// router.Run(address)
	router.RunTLS(address, "vocabulary.cer", "vocabulary.key")

	return nil
}
