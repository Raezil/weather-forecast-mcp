package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/genai"
)

func main() {
	s := server.NewMCPServer(
		"Weather Forecast 🚀",
		"1.0.0",
	)

	// Define weather tool with city, country, and optional date range

	tool := mcp.NewTool("weather",
		mcp.WithDescription("Get the weather forecast for a given city and country over a date range"),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("Name of the city"),
		),
		mcp.WithString("country",
			mcp.Required(),
			mcp.Description("Name of the country"),
		),
		mcp.WithString("fromDate",
			mcp.Description("Start date in YYYY-MM-DD format (defaults to today)"),
		),
		mcp.WithString("toDate",
			mcp.Description("End date in YYYY-MM-DD format (defaults to same as fromDate)"),
		),
	)

	s.AddTool(tool, weatherHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func weatherHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract and validate parameters
	city, ok := request.Params.Arguments["city"].(string)
	if !ok || city == "" {
		return nil, errors.New("city is required and must be a string")
	}
	country, ok := request.Params.Arguments["country"].(string)
	if !ok || country == "" {
		return nil, errors.New("country is required and must be a string")
	}

	// Parse dates
	today := time.Now().UTC()
	from := today
	if v, ok := request.Params.Arguments["fromDate"].(string); ok && v != "" {
		parsed, err := time.Parse("2006-01-02", v)
		if err != nil {
			return nil, fmt.Errorf("invalid fromDate: %v", err)
		}
		from = parsed
	}
	to := from
	if v, ok := request.Params.Arguments["toDate"].(string); ok && v != "" {
		parsed, err := time.Parse("2006-01-02", v)
		if err != nil {
			return nil, fmt.Errorf("invalid toDate: %v", err)
		}
		to = parsed
	}

	// Generate forecast via GenAI
	forecastText, err := generateForecast(ctx, city, country, from, to)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(forecastText), nil
}

func generateForecast(ctx context.Context, city, country string, from, to time.Time) (string, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return "", errors.New("GOOGLE_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("error creating GenAI client: %v", err)
	}

	// Craft prompt for the generative model
	prompt := fmt.Sprintf(
		"Provide a concise weather forecast for %s, %s from %s to %s in Celsius. Include date, temperature, and a brief description.",
		city,
		country,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
	)

	parts := []*genai.Part{{Text: prompt}}
	contents := []*genai.Content{{Parts: parts}}

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, nil)
	if err != nil {
		return "", fmt.Errorf("error generating forecast: %v", err)
	}
	if len(resp.Candidates) == 0 {
		return "", errors.New("no candidates returned from GenAI")
	}

	// Concatenate parts of the first candidate
	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		sb.WriteString(part.Text)
	}
	return sb.String(), nil
}
