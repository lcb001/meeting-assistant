/**
 * client.ts
 * 
 * This file implements a test client for the Todo MCP server.
 * It demonstrates how to connect to the server, call various tools,
 * and handle the responses.
 * 
 * WHY HAVE A TEST CLIENT?
 * - Validates that the server works correctly
 * - Provides a working example of how to use the MCP client SDK
 * - Makes it easy to test changes without needing an LLM
 * - Serves as documentation for how to interact with the server
 */
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";

/**
 * Response content type definition
 * 
 * The MCP protocol returns content as an array of typed objects.
 * This interface defines the structure of text content items.
 */
interface ContentText {
  type: "text";
  text: string;
}

/**
 * Main function that runs the test client
 * 
 * This function:
 * 1. Connects to the Todo MCP server
 * 2. Demonstrates all the available tools
 * 3. Creates, updates, completes, and deletes a test todo
 */
async function main() {
  console.log("Starting Todo MCP Test Client...");

  try {
    /**
     * Create a client transport to the server
     * 
     * The StdioClientTransport launches the server as a child process
     * and communicates with it via standard input/output.
     * 
     * WHY STDIO TRANSPORT?
     * - Simple to set up and use
     * - Works well for local testing
     * - Doesn't require network configuration
     * - Similar to how Claude Desktop launches MCP servers
     */
    const transport = new StdioClientTransport({
      command: "node",
      args: ["dist/index.js"],
    });

    /**
     * Create and connect the client
     * 
     * We configure the client with basic identity information
     * and the capabilities it needs (tools in this case).
     */
    const client = new Client(
      {
        name: "todo-test-client",
        version: "1.0.0",
      },
      {
        capabilities: {
          tools: {}
        }
      }
    );

    // Connect to the server through the transport
    await client.connect(transport);
    console.log("Connected to Todo MCP Server");

    /**
     * List available tools
     * 
     * This demonstrates how to query what tools the server provides,
     * which is useful for discovery and documentation.
     */
    const toolsResult = await client.listTools();
    console.log("\nAvailable tools:", toolsResult.tools.map(tool => tool.name));

    /**
     * Create a test todo
     * 
     * This demonstrates the create-todo tool, which takes a title
     * and markdown description as arguments.
     */
    console.log("\nCreating a test todo...");
    const createTodoResult = await client.callTool({
      name: "create-todo",
      arguments: {
        title: "Learn about MCP",
        description: "# Model Context Protocol\n\n- Understand core concepts\n- Build a simple server\n- Test with Claude"
      }
    });
    
    // Type assertion to access the content
    const createContent = createTodoResult.content as ContentText[];
    console.log(createContent[0].text);

    /**
     * Extract the todo ID from the response
     * 
     * We use a simple regex to parse the ID from the formatted response.
     * In a real application, you might want a more structured response format.
     */
    const idMatch = createContent[0].text.match(/ID: ([0-9a-f-]+)/);
    const todoId = idMatch ? idMatch[1] : null;

    // Only proceed if we successfully created a todo and extracted its ID
    if (todoId) {
      /**
       * List all todos
       * 
       * This demonstrates the list-todos tool, which takes no arguments
       * and returns a formatted list of all todos.
       */
      console.log("\nListing all todos...");
      const listTodosResult = await client.callTool({
        name: "list-todos",
        arguments: {}
      });
      const listContent = listTodosResult.content as ContentText[];
      console.log(listContent[0].text);

      /**
       * Update the todo
       * 
       * This demonstrates the update-todo tool, which takes an ID
       * and optional title/description fields to update.
       */
      console.log("\nUpdating the test todo...");
      const updateTodoResult = await client.callTool({
        name: "update-todo",
        arguments: {
          id: todoId,
          description: "# Updated MCP Learning Plan\n\n- Learn MCP core concepts\n- Build a server with tools\n- Connect to Claude\n- Create amazing AI experiences"
        }
      });
      const updateContent = updateTodoResult.content as ContentText[];
      console.log(updateContent[0].text);

      /**
       * Mark todo as completed
       * 
       * This demonstrates the complete-todo tool, which takes an ID
       * and marks the corresponding todo as completed.
       */
      console.log("\nCompleting the test todo...");
      const completeTodoResult = await client.callTool({
        name: "complete-todo",
        arguments: {
          id: todoId
        }
      });
      const completeContent = completeTodoResult.content as ContentText[];
      console.log(completeContent[0].text);

      /**
       * Summarize active todos
       * 
       * This demonstrates the summarize-active-todos tool, which
       * generates a summary of all non-completed todos.
       */
      console.log("\nSummarizing active todos...");
      const summaryResult = await client.callTool({
        name: "summarize-active-todos",
        arguments: {}
      });
      const summaryContent = summaryResult.content as ContentText[];
      console.log(summaryContent[0].text);

      /**
       * Delete the todo
       * 
       * This demonstrates the delete-todo tool, which permanently
       * removes a todo from the database.
       */
      console.log("\nDeleting the test todo...");
      const deleteTodoResult = await client.callTool({
        name: "delete-todo",
        arguments: {
          id: todoId
        }
      });
      const deleteContent = deleteTodoResult.content as ContentText[];
      console.log(deleteContent[0].text);
    }

    // Close the client connection
    await client.close();
    console.log("\nTest completed successfully!");
  } catch (error) {
    console.error("Error in test client:", error);
    process.exit(1);
  }
}

// Start the test client
main(); 