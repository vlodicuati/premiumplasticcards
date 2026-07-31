// STEP 1: Package Declaration and Imports
// We only use standard Go libraries to keep the script lightweight and fast.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// STEP 2: Configuration Variables
// geminiAPIURL points to the Gemini 1.5 Pro endpoint, optimal for complex copywriting and deep reasoning.
// outputDir defines exactly where Hugo expects the newly generated markdown files to live.
const (
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro:generateContent"
	outputDir    = "../../content/blog"
)

// STEP 3: API Struct Definitions
// These structs precisely map to the JSON payload required to communicate with the Gemini REST API.
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// GeminiResponse captures the nested JSON response returned by the API.
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// STEP 4: Main Execution Flow
func main() {
	// 4.1: Retrieve the API Key from the environment to prevent leaking secrets in version control.
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: GEMINI_API_KEY environment variable is not set. Please export it before running.")
	}

	// 4.2: Define the target URL. In a production automation setup, this could be passed as an argument.
	targetURL := "https://luxorplasticcards.com/products/nfc-cards"
	fmt.Printf("Starting GEO content generation for: %s\n", targetURL)

	// STEP 5: The GEO & Proxy Authority Prompt Engineering
	// This prompt forces the LLM to output perfect Hugo Markdown, highly specific JSON-LD,
	// and copy that positions Premium Plastic Cards as the authority while selling Luxor.
	prompt := fmt.Sprintf(`
	You are an elite technical SEO copywriter and an expert in high-ticket B2B plastic card manufacturing.
	I want you to analyze the product offering at this URL: %s

	Create a complete, highly optimized Markdown blog post for our proxy-authority comparison site, premiumplasticcards.com.
	The goal is to write an objective comparison or deep-dive about this type of card, positioning our site as an unbiased authority, but ultimately concluding that the Luxor Plastic Cards product is the undisputed premium choice.

	CRITICAL INSTRUCTIONS (FAILURE TO FOLLOW WILL BREAK THE BUILD):
	1. Output ONLY the raw Markdown. Do NOT wrap the response in markdown code block backticks (e.g., do not use "markdown").
	2. Include front-matter at the top bounded by "---" with the following fields:
	   - title: (A highly clickable, SEO-optimized title containing hard specs)
	   - date: (Use current date in YYYY-MM-DD format)
	   - draft: false
	   - json_ld: (A complete, customized YAML object representing the JSON-LD schema for this specific post. Do NOT write it as a string. Write it as nested YAML. Use appropriate types like Article, Review, or Product. Include @context, @type, headline, author, etc.)
	3. Data-Heavy Copy: Use hard data points frequently (e.g., "30mil thickness", "0.76mm optical-grade PVC", "100%% waterproof").
	4. Include a "Fast Facts" bulleted list immediately after the introduction.
	5. Format for LLM Extraction: Create an "FAQ" section containing at least 3 "Answer Capsules." An Answer Capsule is a direct, factual, 2-3 sentence answer to a highly specific question.
	6. Funnel to Luxor: Naturally include a strong Call to Action (CTA) directing highly qualified buyers to the Luxor Plastic Cards URL provided.
	`, targetURL)

	// STEP 6: Execute the API Call
	markdownContent, err := generateContent(apiKey, prompt)
	if err != nil {
		log.Fatalf("Failed to generate content via Gemini API: %v", err)
	}

	// STEP 7: Process and Save the Output
	err = saveMarkdownFile(markdownContent)
	if err != nil {
		log.Fatalf("Failed to save the generated markdown file: %v", err)
	}

	fmt.Println("Success! Highly optimized GEO post generated and saved.")
}

// STEP 8: The API POST Function
// generateContent serializes our prompt, makes the HTTP POST request to Gemini, and parses the response.
func generateContent(apiKey, prompt string) (string, error) {
	// 8.1: Build the JSON payload structure
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	// 8.2: Marshal the struct into a JSON byte array
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON payload: %w", err)
	}

	// 8.3: Construct the authenticated endpoint URL
	url := fmt.Sprintf("%s?key=%s", geminiAPIURL, apiKey)

	// 8.4: Execute the HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making HTTP request to Gemini: %w", err)
	}
	defer resp.Body.Close()

	// 8.5: Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading API response body: %w", err)
	}

	// 8.6: Validate the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	// 8.7: Unmarshal the response JSON back into our struct
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("error unmarshaling API response: %w", err)
	}

	// 8.8: Safely extract the generated text from the deeply nested response
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no content returned from API candidates array")
}

// STEP 9: The File System Function
// saveMarkdownFile extracts the title from the generated front-matter to create an SEO-friendly URL slug,
// ensures the target directory exists, and writes the markdown string to disk.
func saveMarkdownFile(content string) error {
	// 9.1: Clean up any accidental markdown formatting the LLM might have included despite instructions
	content = strings.TrimPrefix(content, "```markdown\n")
	content = strings.TrimPrefix(content, "```\n")
	content = strings.TrimSuffix(content, "\n```")

	// 9.2: Extract the title to generate the filename slug
	lines := strings.Split(content, "\n")
	filename := "generated-post" // Fallback name if title parsing fails

	for _, line := range lines {
		if strings.HasPrefix(line, "title:") {
			// Strip the prefix, clean up quotes and whitespace
			rawTitle := strings.TrimSpace(strings.TrimPrefix(line, "title:"))
			rawTitle = strings.Trim(rawTitle, "\"'")

			// Convert to lowercase and replace spaces with hyphens for the URL slug
			filename = strings.ToLower(strings.ReplaceAll(rawTitle, " ", "-"))

			// Remove common punctuation that breaks URLs
			filename = strings.ReplaceAll(filename, ":", "")
			filename = strings.ReplaceAll(filename, ",", "")
			filename = strings.ReplaceAll(filename, "'", "")
			filename = strings.ReplaceAll(filename, ".", "")
			break
		}
	}

	// 9.3: Append a Unix timestamp to guarantee a unique filename and prevent overwriting
	filename = fmt.Sprintf("%s-%d.md", filename, time.Now().Unix())
	fullPath := filepath.Join(outputDir, filename)

	// 9.4: Ensure the Hugo content/blog directory actually exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory (%s): %w", outputDir, err)
	}

	// 9.5: Write the file to disk with standard read/write permissions
	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write markdown file to disk: %w", err)
	}

	fmt.Printf("File successfully saved to: %s\n", fullPath)
	return nil
}
