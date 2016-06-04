package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	vocabularyFile = flag.String("vocabulary", "vocabularyExample.txt", "File with foreign words")
	numOfQuestions = flag.Int("questions", 10, "Number of checking words")
	langValue      = flag.String("lang", "en-ru", "Translating-translated languages in format \"en-ru\"")
	debug          = flag.Bool("debug", false, "Debug mode")
)

// Yandex.Dictionary API stuff
const (
	apiUrl      string = "https://translate.yandex.net/api/v1.5/tr/translate"
	keyAPIvalue string = "trnsl.1.1.20160603T190621Z.22dbd77262d17f5a.151d8f03b42dbf98cb5a3741eaa291f4d4a7b30f"
)

type queryWord struct {
	Text string `xml:"text"`
}

var vocabulary = make(map[string]string)

// Words' format: wordInForeignLang - wordInNativeLang
func checkTranslateStatus(rawWord string) bool {

	//two words: foreign - native
	if strings.Contains(rawWord, " - ") {
		return true
	}
	//one word: no translation of the word
	return false
}

func checkFormat(rawWord *string) {
	// Deletes extra spaces and brackets
}

// Translating word with Yandex.Dictionary API
func translateWord(rawWord string) string {

	client := &http.Client{}

	data := `key=` + keyAPIvalue + `&lang=` + *langValue + `&text=` + rawWord

	body := strings.NewReader(data)
	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if *debug == true {
		log.Println(resp.Status)
	}

	bodyEncoding, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error;", err)
	}
	var q = queryWord{}
	err = xml.Unmarshal(bodyEncoding, &q)
	if err != nil {
		log.Println("Error;", err)
	}

	if *debug == true {
		log.Printf("response Body: %v", q)
	}

	return q.Text
}

// Translate untranslated words from vocabularyFile
// and initialize program vocabilary map
func checkVocabulary() {
	input, err := ioutil.ReadFile(*vocabularyFile)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(input), "\n")
	var numOfUntraslatedWords int = 0

	for i, rawWord := range lines {

		if *debug == true {
			log.Println(i+1, rawWord)
		}

		wordTranslated := checkTranslateStatus(rawWord)

		if wordTranslated == true {
			checkFormat(&rawWord)
			parsedWordRaw := strings.Split(rawWord, " - ")
			vocabulary[parsedWordRaw[0]] = parsedWordRaw[1]
		} else {
			numOfUntraslatedWords++
			vocabulary[rawWord] = translateWord(rawWord)
			lines[i] = rawWord + " - " + vocabulary[rawWord]

		}
		if *debug == true {
			log.Printf("Splitted rawWord. Key:-%s-, value:-%s-", rawWord, vocabulary[rawWord])
		}
	}
	//TODO: make two loggers!
	fmt.Printf("Untranslated words in %s: %d.\n", *vocabularyFile, numOfUntraslatedWords)
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(*vocabularyFile, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	if *debug == true {
		for k, v := range vocabulary {
			log.Printf("%s - %s", k, v)
		}
	}

}

func startTest() {

	fmt.Printf("Number of words in database: %d\n", len(vocabulary))
	fmt.Printf("Inter number of words for test (default - %d): ", *numOfQuestions)

	_, err := fmt.Scanf("%d", numOfQuestions)
	if err != nil {
		fmt.Printf("Wrong value; will use default\n")
	}

	var showList string
	fmt.Printf("Do you want to see all words before? (y/n) ")
	fmt.Scanf("%s", &showList)
	if showList == "y" {
		for k, v := range vocabulary {
			if k == "" {
				continue
			}
			fmt.Println(k, " - ", v, "\t")
		}
	}

	fmt.Printf("Number of words in database: %d\n", len(vocabulary))
	fmt.Println("-----------------------------------------")

	i := 0
	right := 0
	var translation string
	for k, v := range vocabulary {
		if i == *numOfQuestions-1 {
			break
		}
		fmt.Println("Word: ", k)
		fmt.Printf("A translation of a word: ")
		fmt.Scanf("%s", &translation)
		if translation == v {
			if k == "" {
				fmt.Println("***(heh): take a gift from me")
			}
			fmt.Println("***Good job! Right")
			right++
		} else {
			fmt.Println("***Wrong! :(")
		}
		fmt.Println(k, " - ", v)
		fmt.Printf("***Scope: %d/%d\n", right, *numOfQuestions)
		fmt.Println("-----------------------------------------")
		i++
	}
	log.Printf("Final scope: %d/%d", right, *numOfQuestions)
}

func main() {

	flag.Parse()
	flag.Parsed()

	if *debug {
		log.Println("Parsed args: ", *vocabularyFile, *numOfQuestions)
	}

	checkVocabulary()
	startTest()
}
