package google

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

func GetGoogleDocContent(docID string) (string, error) {
	// Path to the service account key file (api credentials)
	serviceAccountKeyFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	// Create a context
	ctx := context.Background()

	// Create a Docs API client using the service account
	service, err := docs.NewService(ctx, option.WithCredentialsFile(serviceAccountKeyFile))
	if err != nil {
		log.Fatalf("Failed to create Docs service: %v", err)
	}

	// Access a specific Google Doc
	doc, err := service.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve document: %v", err)
	}

	var content string
	lineCount := 0

	// Iterate over the document content
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			// Only process the first 50 lines (paragraphs)
			if lineCount >= 10 {
				break
			}

			// Process each paragraph element
			for _, paraElement := range elem.Paragraph.Elements {
				if paraElement.TextRun != nil {
					content += paraElement.TextRun.Content
				}
			}

			// Increment the line counter after processing each paragraph
			lineCount++
		}
	}

	// for _, elem := range doc.Body.Content {
	// 	if elem.Paragraph != nil {
	// 		for _, paraElement := range elem.Paragraph.Elements {
	// 			if paraElement.TextRun != nil {
	// 				content += paraElement.TextRun.Content
	// 				// fmt.Print(paraElement.TextRun.Content)
	// 			}
	// 		}
	// 	}
	// }

	return content, nil
}

// Example usage
// Initialize API key and Google Doc ID
// docID := os.Getenv("DOC_ID")

// // Fetch content from Google Doc
// text, err := google.GetGoogleDocContent(docID)
// if err != nil {
// 	fmt.Printf("Error fetching Google Doc content: %v\n", err)
// 	return
// }
