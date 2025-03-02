package utils

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func ConvertHTML(input string) string {
	// Parse the input HTML
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return ""
	}

	// Traverse and process the HTML nodes
	var processNode func(*html.Node)
	processNode = func(n *html.Node) {
		// If the node is a <figure> tag, check for image link and replace it with <img>
		if n.Type == html.ElementNode && n.Data == "figure" {
			for _, attr := range n.Attr {
				if attr.Key == "data-trix-attachment" {
					// Extract image URL from data-trix-attachment JSON string
					// Simple way, this could be improved with a proper JSON parser if needed
					if strings.Contains(attr.Val, `"href"`) {
						startIdx := strings.Index(attr.Val, `"href":"`) + 8
						endIdx := strings.Index(attr.Val[startIdx:], `"`) + startIdx
						if startIdx > 0 && endIdx > startIdx {
							imageURL := attr.Val[startIdx:endIdx]
							//imageURL = strings.Replace(imageURL, "localhost", "172.27.67.208", 1) //for local development only
							//imageURL = strings.Replace(imageURL, "localhost", "10.0.2.2", 1) //for local development only
							imageURL = strings.Replace(imageURL, "localhost", "172.26.43.0", 1) //for local development only
							newNode := &html.Node{
								Type: html.ElementNode,
								Data: "img",
								Attr: []html.Attribute{
									{Key: "src", Val: imageURL},
									{Key: "width", Val: "514"},
									{Key: "height", Val: "514"},
								},
							}
							// Replace the <figure> tag with the new <img> node
							n.Parent.InsertBefore(newNode, n)
							n.Parent.RemoveChild(n)
						}
					}
				}
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}

	// Start processing from the root node
	processNode(doc)

	// Render the modified HTML back to a string
	var builder strings.Builder
	err = html.Render(&builder, doc)
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return ""
	}

	return builder.String()
}
