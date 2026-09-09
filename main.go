package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	NWSAPIBase = "https://api.weather.gov"
	UserAgent  = "weather-app/1.0"
)

func main() {
	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-go-starter-kit",
		Version: "1.0.0",
	}, nil)

	// Add get_forecast tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_forecast",
		Description: "Get weather forecast for a location",
	}, GetForecast)

	// Add get_alerts tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_alerts",
		Description: "Get weather alerts for a US state",
	}, GetAlerts)

	// Run server on stdio transport
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
