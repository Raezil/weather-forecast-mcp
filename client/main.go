package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// 1. Launch stdio transport, pointing to your MCP server binary
	//    Adjust the path to where your server executable resides
	stdio := transport.NewStdio("../server/my_mcp_server", nil)

	// 2. Create the MCP client
	cli := client.NewClient(stdio)
	defer cli.Close()

	ctx := context.Background()

	// 3. Start communication
	if err := cli.Start(ctx); err != nil {
		log.Fatalf("failed to start client: %v", err)
	}

	// 4. Perform handshake & capability negotiation
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = "0.4.0"
	initReq.Params.Capabilities = cli.GetClientCapabilities()
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "weather-client",
		Version: "v1.0.0",
	}
	if _, err := cli.Initialize(ctx, initReq); err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	// 5. Prepare the CallToolRequest for the "weather" tool
	req := mcp.CallToolRequest{}
	req.Params.Name = "weather"
	req.Params.Arguments = map[string]interface{}{
		"city":     "Warsaw",
		"country":  "Poland",
		"fromDate": "2025-05-02",
		"toDate":   "2025-05-15",
	}

	// 6. Call the tool
	resp, err := cli.CallTool(ctx, req)
	if err != nil {
		log.Fatalf("tool call failed: %v", err)
	}

	// 7. Extract and print the forecast text
	for _, content := range resp.Content {
		if tc, ok := mcp.AsTextContent(content); ok {
			fmt.Println(tc.Text)
		}
	}
}
