package translateapi

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type TranslatedText struct {
	Text           string `json:"translatedText"`
	DetectLanguage DetectLanguage
}

type DetectLanguage struct {
	Confidence float32 `json:"confidence"`
	Language   string  `json:"language"`
}

func TranslateAPI(query string, baseurl string) string {
	queryurl := baseurl + "/translate"
	httpClient := &http.Client{}
	ctx := context.Background()
	data := `{ "q": "` + query + `","source": "auto","target": "en"}`
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(3))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, queryurl, strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	text, _ := io.ReadAll(resp.Body)
	translated := TranslatedText{}
	_ = json.Unmarshal([]byte(string(text)), &translated)
	log.Printf("%s", queryurl)
	log.Printf("%s", "Translate "+translated.Text+" from language "+translated.DetectLanguage.Language)
	return translated.Text
}
