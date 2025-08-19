package utils

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func ExtractOGImage(articleURL string) (string, error) {
	resp, err := http.Get(articleURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			fmt.Printf("No image found")
			return "", nil // no image found
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "meta" {
				var property, content string
				for _, attr := range t.Attr {
					if attr.Key == "property" && attr.Val == "og:image" {
						property = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if property == "og:image" && content != "" {
					return content, nil
				}
			}
		}
	}
}
