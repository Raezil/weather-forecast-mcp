# Weather Forecast 🚀

Welcome to the **Weather Forecast** MCP tool! 🎉 This small but powerful service lets you fetch concise weather forecasts for any city and country over a date range, powered by Google’s GenAI Gemini model.

## 🚀 Features

* **Simple interface:** Call a single `weather` tool with `city`, `country`, and optional date range.
* **GenAI-powered:** Leverages Google Gemini to generate human-friendly, contextual forecasts.
* **Customizable dates:** Default to today, or specify `fromDate` and `toDate` in `YYYY-MM-DD` format.
* **Extensible via MCP:** Built on the `mark3labs/mcp-go` framework for easy integration and expansion.

## 📦 Prerequisites

* Go **1.18+**
* A valid **Google Cloud API key** with access to the Gemini API

## 🛠 Installation

1. Clone this repository:

   ```bash
   git clone https://github.com/Raezil/weather-forecast-mcp
   cd weather-forecast-tool
   ```
2. Install dependencies:

   ```bash
   go mod download
   ```

## 🔧 Configuration

Set the following environment variable before running the server:

```bash
export GOOGLE_API_KEY="YOUR_GOOGLE_CLOUD_API_KEY"
```

*Tip:* Store your API key securely (e.g., in a secrets manager or `.env` file).

## ▶️ Running the Server

Start the MCP server over stdio:

```bash
go run main.go
```

You should see output indicating that the `weather` tool is registered:

```
Weather Forecast 🚀 v1.0.0 listening on stdio...
```

## 💬 Using the Client

Below is an example Go client that demonstrates calling the `weather` tool:

```go
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

```

We’ve covered missing/invalid parameters, date parsing errors, and environment checks. Contributions of additional test scenarios are encouraged!

## 🤝 Contributing

We’d love your help to make this even better:

* Open issues for bugs or feature requests
* Submit pull requests for new functionality or improvements
* Add support for more transport layers or AI backends

## 📝 License

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for details.

---

Thanks for checking out **Weather Forecast**! Let’s build something amazing together. 🌟
