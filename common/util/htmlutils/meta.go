package htmlutils

import (
	"io"

	"golang.org/x/net/html"
)

func ExtractMetaTags(resp io.Reader) []map[string]interface{} {
	tokenizer := html.NewTokenizer(resp)

	var metaList []map[string]interface{}

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return metaList
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "meta" {
				attributeMap := map[string]interface{}{}
				for _, attr := range token.Attr {
					attributeMap[attr.Key] = attr.Val
				}

				metaList = append(metaList, attributeMap)
			}
		}
	}
	return metaList
}

func extractMetaProperty(token html.Token, prop string) (content string, ok bool) {
	for _, attr := range token.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}
