package search

import (
	"WebpageArchiver/common"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"os"
)

func AddDocFile(fileName string, originDomain string) (err error) {
	htmlFilePath := common.ARCHIVEFILELOACTION + "Temporary/" + fileName
	_, err = os.Stat(htmlFilePath)
	if err != nil {
		return err
	}
	HTMLContent, err := common.GetHTMLFileContent(htmlFilePath)
	if err != nil {
		return err
	}
	title, err := common.GetHTMLTitle(HTMLContent)
	if err != nil {
		return err
	}
	HTMLPureText, err := common.ExtractHTMLText(HTMLContent)
	if err != nil {
		return err
	}

	documents := []map[string]interface{}{
		{
			"id":       uuid.New().String(),
			"title":    title,
			"filename": fileName,
			"domain":   originDomain,
			"content":  HTMLPureText,
		},
	}
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   common.MEILIHOST,
		APIKey: common.MEILIAPIKey,
	})
	_, err = client.Index(common.MEILIBlogsIndex).AddDocuments(documents)
	if err != nil {
		return err
	}

	// 成功添加索引内容后，移动文件到域名目录
	_, err = os.Stat(common.ARCHIVEFILELOACTION + originDomain)
	if os.IsNotExist(err) {
		err = os.Mkdir(common.ARCHIVEFILELOACTION+originDomain, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = os.Rename(htmlFilePath, common.ARCHIVEFILELOACTION+originDomain+"/"+fileName)
	if err != nil {
		return err
	}
	return nil
}
