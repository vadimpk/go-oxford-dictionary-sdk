package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	
	host = "https://od-api.oxforddictionaries.com/api/v2"

	endpointThesaurus = "/thesaurus/en"
	endpointWords = "/words/en-gb"
	endpointEntries = "/entries/en-gb"
	endpointSentences = "/sentences/en"
	endpointTranslations = "/translations"

	paramsEntries = "?strictMatch=false"
	paramsSentences = "?strictMatch=false"
	paramsThesaurus = "?strictMatch=false"
	paramsWords = ""
	paramsTranslations = "?strictMatch=false"

	defaultTimeout = 10 * time.Second
)


type (

	OxfordResponse struct {
		Meta 			map[string]interface{} `json:"metadata"`
		Query 			string `json:"query"`
		Results 		[]Result `json:"results"`
	}

	Result struct {
		Id 						string `json:"id"`
		Language 				string `json:"language"`
		LexicalEntries 			[]LexicalEntry `json:"lexicalEntries"`
		Type 					string `json:"type"`
		Word 					string `json:"word"`
	}

	LexicalEntry struct {
		Entries 				[]Entry `json:"entries"`
		Language 				string `json:"language"`
		LexicalCategory			LexicalCategory `json:"lexicalCategory"`
		Text 					string `json:"text"`
		Phrases 				[]Phrase `json:"phrases"`
		Sentences 				[]Sentence `json:"sentences"`
	}

	Entry struct {
		Senses 					[]Sense `json:"senses"`
		Pronunciations 			[]Pronunciation `json:"pronunciations"`
		Etymologies 			[]string `json:"etymologies"`
		Inflections				[]Inflection `json:"inflections"`
		GrammaticalFeatures 	[]GrammaticalFeature `json:"grammaticalFeatures"`
	}

	Sense struct {
		Antonyms 				[]Antonym `json:"antonyms"`
		Synonyms 				[]Synonym `json:"synonyms"`
		Definitions 			[]string `json:"definitions"`
		ShortDefinitions 		[]string `json:"shortDefinitions"`
		Examples 				[]Example `json:"examples"`
		Constructions			[]Construction `json:"constructions"`
		Inflections				[]Inflection `json:"inflections"`
		DomainClasses 			[]DomainClass `json:"domainClasses"`
		Notes					[]Note `json:"notes"`
		Translations			[]Translation `json:"translations"`
		SemanticClasses 		[]SemanticClass `json:"semanticClasses"`
		Id 						string `json:"id"`
		Subsenses 				[]Sense `json:"subsenses"`
	}

	Pronunciation struct {
		AudioFileLink 			string `json:"audioFile"`
		Dialects 				[]string `json:"dialects"`
		PhoneticNotation 		string `json:"phoneticNotation"`
		PhoneticSpelling 		string `json:"phoneticSpelling"`
	}

	Synonym struct {
		Language 	string `json:"language"`
		Text 		string `json:"text"`
	}

	Antonym struct {
		Language 	string `json:"language"`
		Text 		string `json:"text"`
	}

	Example struct {
		Text 			string `json:"text"`
		Notes 			[]Note `json:"notes"`
		Registers 		[]Register `json:"registers"`
		Translations	[]Translation `json:"translations"`
	}

	Sentence struct {
		Regions []Region `json:"regions"`
		Text 	string `json:"text"`
	}

	Inflection struct {
		InflectedForm string `json:"inflectedForm"`
	}

	Construction struct {
		Text string `json:"text"`
	}

	Translation struct {
		Language 	string `json:"language"`
		Text 		string `json:"text"`
	}

	LexicalCategory struct {
		Id 		string `json:"id"`
		Text 	string `json:"text"`
	}

	Region struct {
		Id 		string `json:"id"`
		Text 	string `json:"text"`
	}

	DomainClass struct {
		Id 		string `json:"id"`
		Text 	string `json:"text"`
	}

	Note struct {
		Text string `json:"text"`
		Type string `json:"type"`
	}

	SemanticClass struct {
		Id 		string `json:"id"`
		Text 	string `json:"text"`
	}

	Register struct {
		Id 		string `json:"id"`
		Text 	string `json:"text"`
	}

	Phrase struct {
		Id 		string `json:"id"`
		Text 	string `json:"text"`
	}

	GrammaticalFeature struct {
		Id 		string `json:"id"`
		Text 	string `json:"text"`
		Type 	string `json:"type"`
	}
)

// Client is a Oxford Dictionary API client
type Client struct {
	client      *http.Client
	appID 		string
	appKEY		string
}

// NewClient creates a new client instance with your credentials
func NewClient(appID, appKEY string) (*Client, error) {
	if appID == "" || appKEY == "" {
		return nil, errors.New("credentials are empty")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		appID: appID,
		appKEY: appKEY,
	}, nil
}


func (c *Client) Thesaurus(word string) (OxfordResponse, error) {
	endpoint := endpointThesaurus + "/" + word + paramsThesaurus
	return c.doHTTP(endpoint)
}

func (c *Client) WordInfo(word string) (OxfordResponse, error) {
	endpoint := endpointWords + "?q=" + strings.ReplaceAll(word, " ", "_") + paramsWords
	return c.doHTTP(endpoint)
}

func (c *Client) Entry(word string) (OxfordResponse, error) {
	endpoint := endpointEntries + "/" + word + paramsEntries
	return c.doHTTP(endpoint)
}

func (c *Client) Sentences(word string) (OxfordResponse, error) {
	endpoint := endpointSentences + "/" + word + paramsSentences
	return c.doHTTP(endpoint)
}

func (c *Client) Translation(word, sourceLang, targetLang string) (OxfordResponse, error) {
	endpoint := endpointTranslations + "/" + sourceLang + "/" + targetLang + "/" + word + paramsTranslations
	return c.doHTTP(endpoint)
}

// generates the request with the given endpoint and sets the headers
func (c *Client) generateRequest(endpoint, HTTPMethod string) *http.Request {

	url := host + endpoint

	req, _ := http.NewRequest(HTTPMethod, url, nil)
	req.Header.Set("app_id", c.appID)
	req.Header.Set("app_key", c.appKEY)
	req.Header.Set("Accept", "application/json")

	return req
}

// makes the request with given endpoint and parses the response
func (c *Client) doHTTP(endpoint string) (OxfordResponse, error) {
	req := c.generateRequest(endpoint, http.MethodGet)
	var resp OxfordResponse

	if err := c.fetchJSON(req, &resp); err != nil {
		return OxfordResponse{}, err
	}

	return resp, nil
}

// decodes response body into struct
func (c *Client) fetchJSON(req *http.Request, data *OxfordResponse) (error) {

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("API Error: %s", resp.Status))
	}

	return json.NewDecoder(resp.Body).Decode(data)
}
