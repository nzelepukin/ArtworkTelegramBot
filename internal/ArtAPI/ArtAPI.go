package artapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ArtworkList struct {
	Data []Artwork `json:"data"`
}
type Artwork struct {
	Api_model string `json:"api_model"`
	ID        int64  `json:"id"`
	Title     string `json:"title"`
}

type ArtworkInfo struct {
	Data ArtworkImageID `json:"data"`
}

type ArtworkImageID struct {
	ImageID string `json:"image_id"`
}

func GetArtAPI(query string) string {
	BaseURL := "https://api.artic.edu/api/v1/artworks/"
	query = strings.ToLower(strings.ReplaceAll(query, " ", "+"))
	fmt.Println(query)
	artList := getid(query, BaseURL)
	if len(artList.Data) > 0 {
		for _, art := range artList.Data[:1] {
			return "https://www.artic.edu/iiif/2/" + getImageId(art.ID, BaseURL) + "/full/843,/0/default.jpg"
		}
	}
	return "Can't find paintings"
	//if len(artList.Data) > 0 {
	//	for _, art := range artList.Data[:3] {
	//		fmt.Println(art.Title)
	//		saveFile("https://www.artic.edu/iiif/2/"+getImageId(art.ID, BaseURL)+"/full/843,/0/default.jpg", art.Title)
	//	}
	//} else {
	//	fmt.Println("No artworks with query " + query)
	//}
}

func getid(query string, baseurl string) ArtworkList {
	url := baseurl + "search?q=" + query
	httpClient := &http.Client{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(3))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
	output := ArtworkList{}
	_ = json.Unmarshal([]byte(string(text)), &output)
	return output
}

func getImageId(artid int64, baseurl string) string {
	url := baseurl + strconv.FormatInt(int64(artid), 10) + "?fields=image_id"
	httpClient := &http.Client{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(3))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
	output := ArtworkInfo{}
	_ = json.Unmarshal([]byte(string(text)), &output)
	return output.Data.ImageID
}

func saveFile(url string, title string) {
	httpClient := &http.Client{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(3))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	out, err := os.Create("imgs/" + title + ".jpg")
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

}
