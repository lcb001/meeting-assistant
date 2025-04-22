/**
 * formatters.ts
 * 
 * This file contains utility functions for formatting data in the application.
 * These utilities handle the transformation of internal data structures into
 * human-readable formats appropriate for display to LLMs and users.
 * 
 * WHY SEPARATE FORMATTERS?
 * - Keeps formatting logic separate from business logic
 * - Allows consistent formatting across the application
 * - Makes it easier to change display formats without affecting core functionality
 * - Centralizes presentation concerns in one place
 */
import { Todo } from "../models/Todo.js";

/**
 * Format a todo item to a readable string representation
 * 
 * This formatter converts a Todo object into a markdown-formatted string
 * with clear visual indicators for completion status (emojis).
 * 
 * WHY USE MARKDOWN?
 * - Provides structured, readable output
 * - Works well with LLMs which understand markdown syntax
 * - Allows rich formatting like headers, lists, and emphasis
 * - Can be displayed directly in many UI contexts
 * 
 * @param todo The Todo object to format
 * @returns A markdown-formatted string representation
 */
export function formatTodo(todo: Todo): string {
  return `
## ${todo.title} ${todo.completed ? '✅' : '⏳'}

ID: ${todo.id}
Created: ${new Date(todo.createdAt).toLocaleString()}
Updated: ${new Date(todo.updatedAt).toLocaleString()}

${todo.description}
  `.trim();
}

/**
 * Format a list of todos to a readable string representation
 * 
 * This formatter takes an array of Todo objects and creates a complete
 * markdown document with a title and formatted entries.
 * 
 * @param todos Array of Todo objects to format
 * @returns A markdown-formatted string with the complete list
 */
export function formatTodoList(todos: Todo[]): string {
  if (todos.length === 0) {
    return "No todos found.";
  }

  const todoItems = todos.map(formatTodo).join('\n\n---\n\n');
  return `# Todo List (${todos.length} items)\n\n${todoItems}`;
}

/**
 * Create success response for MCP tool calls
 * 
 * This utility formats successful responses according to the MCP protocol.
 * It wraps the message in the expected content structure.
 * 
 * WHY THIS FORMAT?
 * - Follows the MCP protocol's expected response structure
 * - Allows the message to be properly displayed by MCP clients
 * - Clearly indicates success status
 * 
 * @param message The success message to include
 * @returns A properly formatted MCP response object
 */
export function createSuccessResponse(message: string) {
  return {
    content: [
      {
        type: "text" as const,
        text: message,
      },
    ],
  };
}

/**
 * Create error response for MCP tool calls
 * 
 * This utility formats error responses according to the MCP protocol.
 * It includes the isError flag to indicate failure.
 * 
 * @param message The error message to include
 * @returns A properly formatted MCP error response object
 */
export function createErrorResponse(message: string) {
  return {
    content: [
      {
        type: "text" as const,
        text: message,
      },
    ],
    isError: true,
  };
} 