package search

import (
	"WebpageArchiver/common"
	"encoding/json"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"strings"
)

type Result struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
	Domain   string `json:"domain"`
}

func QueryByKeyword(keyword string, pageNum int64) (string, map[string]int) {
	pageAndHits := make(map[string]int)
	QueryResults := make([]Result, 10)
	preTag := "<span style=\"color: red;\">"
	postTag := "</span>"

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   common.MEILIHOST,
		APIKey: common.MEILIAPIKey,
	})
	meiliReqOpt := &meilisearch.SearchRequest{
		Page:                  pageNum,
		HitsPerPage:           10,
		AttributesToHighlight: []string{"content"},
		ShowMatchesPosition:   true,
		HighlightPreTag:       preTag,
		HighlightPostTag:      postTag,
		AttributesToCrop:      []string{"content"},
		CropLength:            150,
	}

	meiliResp, err := client.Index("blogs").Search(keyword, meiliReqOpt)

	TotalHits := meiliResp.TotalHits
	TotalPages := meiliResp.TotalPages
	pageAndHits["TotalHits"] = int(TotalHits)
	pageAndHits["TotalPages"] = int(TotalPages)

	hits := meiliResp.Hits

	for hitIndex, hit := range hits {
		var result Result
		singleContent, _ := hit.(map[string]interface{})

		// 获取高亮内容
		formattedContent, _ := singleContent["_formatted"].(map[string]interface{})
		formattedContentStr, _ := formattedContent["content"].(string)

		result.Id = formattedContent["id"].(string)
		result.Filename = singleContent["filename"].(string)
		result.Domain = singleContent["domain"].(string)
		result.Title = singleContent["title"].(string)
		result.Content = formattedContentStr

		QueryResults[hitIndex] = result
	}
	resultJson, _ := json.MarshalIndent(QueryResults, "", "    ")
	// fmt.Println(string(resultJson))
	resultJsonString := strings.ReplaceAll(string(resultJson), "\n", "")

	if err == nil {
		return resultJsonString, pageAndHits
	} else {
		log.Println("Error Occur: " + err.Error())
		return "Error", nil
	}
}
