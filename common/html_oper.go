package common

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"regexp"
	"strings"
)

func ExtractHTMLText(htmlContent string) (string, error) {
	r := strings.NewReader(htmlContent)
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}
	removeScriptAndStyle(doc)
	pureHTML := extractText(doc)
	stringReplacer := strings.NewReplacer(" ", "", "\n", "", "\r", "", "\\n", "")
	pureHTML = stringReplacer.Replace(pureHTML)
	// 二次解析保证去除所有标签
	r = strings.NewReader(pureHTML)
	doc, err = html.Parse(r)
	if err != nil {
		return "", err
	}
	removeScriptAndStyle(doc)
	pureHTML = extractText(doc)
	pureHTML = stringReplacer.Replace(pureHTML)
	return pureHTML, nil
}

func removeScriptAndStyle(doc *html.Node) {
	// BFS
	queue := []*html.Node{doc}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == nil {
			continue
		}

		for child := current.FirstChild; child != nil; {
			next := child.NextSibling
			if child.Type == html.ElementNode && (child.Data == "script" || child.Data == "style") {
				// 移除后，不应再让 child = child.NextSibling
				current.RemoveChild(child)
			} else {
				// 如果 child 未被移除，才放入队列
				queue = append(queue, child)
			}
			// 切换到下一个兄弟节点
			child = next
		}
	}
}

func extractText(doc *html.Node) string {
	var builder strings.Builder
	queue := []*html.Node{doc}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Type == html.TextNode {
			builder.WriteString(current.Data)
		}

		for c := current.FirstChild; c != nil; c = c.NextSibling {
			queue = append(queue, c)
		}
	}
	return builder.String()
}

func GetHTMLTitle(htmlContent string) (title string, err error) {
	re := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	matches := re.FindStringSubmatch(htmlContent)
	if len(matches) < 2 {
		return "", fmt.Errorf("no title found")
	}
	return matches[1], nil
}

func GetHTMLFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
