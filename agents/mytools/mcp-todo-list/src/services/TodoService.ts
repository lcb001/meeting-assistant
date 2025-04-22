/**
 * TodoService.ts
 *
 * This service implements the core business logic for managing todos.
 * It acts as an intermediary between the data model and the database,
 * handling all CRUD operations and search functionality.
 *
 * WHY A SERVICE LAYER?
 * - Separates business logic from database operations
 * - Provides a clean API for the application to work with
 * - Makes it easier to change the database implementation later
 * - Encapsulates complex operations into simple method calls
 */
import { Todo, createTodo, CreateTodoSchema, UpdateTodoSchema } from '../models/Todo.js';
import { z } from 'zod';
import { databaseService } from './DatabaseService.js';

/**
 * TodoService Class
 *
 * This service follows the repository pattern to provide a clean
 * interface for working with todos. It encapsulates all database
 * operations and business logic in one place.
 */
class TodoService {
  /**
   * Create a new todo
   *
   * This method:
   * 1. Uses the factory function to create a new Todo object
   * 2. Persists it to the database
   * 3. Returns the created Todo
   *
   * @param data Validated input data (title and description)
   * @returns The newly created Todo
   */
  createTodo(data: z.infer<typeof CreateTodoSchema>): Todo {
    // Use the factory function to create a Todo with proper defaults
    const todo = createTodo(data);

    // Get the database instance
    const db = databaseService.getDb();

    // Prepare the SQL statement for inserting a new todo
    const stmt = db.prepare(`
      INSERT INTO todos (id, meetingID, title, description, completedAt, createdAt, updatedAt, list, assignee)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `);

    // Execute the statement with the todo's data
    stmt.run(
      todo.id,
        todo.meetingID,
      todo.title,
      todo.description,
      todo.completedAt,
      todo.createdAt,
      todo.updatedAt,
        todo.list,
        todo.assignee,
    );

    // Return the created todo
    return todo;
  }

  /**
   * Get a todo by ID
   *
   * This method:
   * 1. Queries the database for a todo with the given ID
   * 2. Converts the database row to a Todo object if found
   *
   * @param id The UUID of the todo to retrieve
   * @returns The Todo if found, undefined otherwise
   */
  getTodo(id: string): Todo | undefined {
    const db = databaseService.getDb();

    // Use parameterized query to prevent SQL injection
    const stmt = db.prepare('SELECT * FROM todos WHERE id = ?');
    const row = stmt.get(id) as any;

    // Return undefined if no todo was found
    if (!row) return undefined;

    // Convert the database row to a Todo object
    return this.rowToTodo(row);
  }

  /**
   * Get all todos
   *
   * This method returns all todos in the database without filtering.
   *
   * @returns Array of all Todos
   */
  getAllTodos(): Todo[] {
    const db = databaseService.getDb();
    const stmt = db.prepare('SELECT * FROM todos');
    const rows = stmt.all() as any[];

    // Convert each database row to a Todo object
    return rows.map(row => this.rowToTodo(row));
  }

  /**
   * Get all active (non-completed) todos
   *
   * This method returns only todos that haven't been marked as completed.
   * A todo is considered active when its completedAt field is NULL.
   *
   * @returns Array of active Todos
   */
  getActiveTodos(): Todo[] {
    const db = databaseService.getDb();
    const stmt = db.prepare('SELECT * FROM todos WHERE completedAt IS NULL');
    const rows = stmt.all() as any[];

    // Convert each database row to a Todo object
    return rows.map(row => this.rowToTodo(row));
  }

  /**
   * Update a todo
   *
   * This method:
   * 1. Checks if the todo exists
   * 2. Updates the specified fields
   * 3. Returns the updated todo
   *
   * @param data The update data (id required, title/description optional)
   * @returns The updated Todo if found, undefined otherwise
   */
  updateTodo(data: z.infer<typeof UpdateTodoSchema>): Todo | undefined {
    // First check if the todo exists
    const todo = this.getTodo(data.id);
    if (!todo) return undefined;

    // Create a timestamp for the update
    const updatedAt = new Date().toISOString();

    const db = databaseService.getDb();
    const stmt = db.prepare(`
      UPDATE todos
      SET title = ?, description = ?, updatedAt = ?
      WHERE id = ?
    `);

    // Update with new values or keep existing ones if not provided
    stmt.run(
      data.title || todo.title,
      data.description || todo.description,
      updatedAt,
      todo.id
    );

    // Return the updated todo
    return this.getTodo(todo.id);
  }

  /**
   * Mark a todo as completed
   *
   * This method:
   * 1. Checks if the todo exists
   * 2. Sets the completedAt timestamp to the current time
   * 3. Returns the updated todo
   *
   * @param id The UUID of the todo to complete
   * @returns The updated Todo if found, undefined otherwise
   */
  completeTodo(id: string): Todo | undefined {
    // First check if the todo exists
    const todo = this.getTodo(id);
    if (!todo) return undefined;

    // Create a timestamp for the completion and update
    const now = new Date().toISOString();

    const db = databaseService.getDb();
    const stmt = db.prepare(`
      UPDATE todos
      SET completedAt = ?, updatedAt = ?
      WHERE id = ?
    `);

    // Set the completedAt timestamp
    stmt.run(now, now, id);

    // Return the updated todo
    return this.getTodo(id);
  }

  /**
   * Delete a todo
   *
   * This method removes a todo from the database permanently.
   *
   * @param id The UUID of the todo to delete
   * @returns true if deleted, false if not found or not deleted
   */
  deleteTodo(id: string): boolean {
    const db = databaseService.getDb();
    const stmt = db.prepare('DELETE FROM todos WHERE id = ?');
    const result = stmt.run(id);

    // Check if any rows were affected
    return result.changes > 0;
  }

  /**
   * Search todos by title
   *
   * This method performs a case-insensitive partial match search
   * on todo titles.
   *
   * @param title The search term to look for in titles
   * @returns Array of matching Todos
   */
  searchByTitle(title: string): Todo[] {
    // Add wildcards to the search term for partial matching
    const searchTerm = `%${title}%`;

    const db = databaseService.getDb();

    // COLLATE NOCASE makes the search case-insensitive
    const stmt = db.prepare('SELECT * FROM todos WHERE title LIKE ? COLLATE NOCASE');
    const rows = stmt.all(searchTerm) as any[];

    return rows.map(row => this.rowToTodo(row));
  }

  /**
   * Search todos by date
   *
   * This method finds todos created on a specific date.
   * It matches the start of the ISO string with the given date.
   *
   * @param dateStr The date to search for in YYYY-MM-DD format
   * @returns Array of matching Todos
   */
  searchByDate(dateStr: string): Todo[] {
    // Add wildcard to match the time portion of ISO string
    const datePattern = `${dateStr}%`;

    const db = databaseService.getDb();
    const stmt = db.prepare('SELECT * FROM todos WHERE createdAt LIKE ?');
    const rows = stmt.all(datePattern) as any[];

    return rows.map(row => this.rowToTodo(row));
  }

  /**
   * Generate a summary of active todos
   *
   * This method creates a markdown-formatted summary of all active todos.
   *
   * WHY RETURN FORMATTED STRING?
   * - Provides ready-to-display content for the MCP client
   * - Encapsulates formatting logic in the service
   * - Makes it easy for LLMs to present a readable summary
   *
   * @returns Markdown-formatted summary string
   */
  summarizeActiveTodos(): string {
    const activeTodos = this.getActiveTodos();

    // Handle the case when there are no active todos
    if (activeTodos.length === 0) {
      return "No active todos found.";
    }

    // Create a bulleted list of todo titles
    const summary = activeTodos.map(todo => `- ${todo.title}`).join('\n');
    return `# Active Todos Summary\n\nThere are ${activeTodos.length} active todos:\n\n${summary}`;
  }

  /**
   * Helper to convert a database row to a Todo object
   *
   * This private method handles the conversion between the database
   * representation and the application model.
   *
   * WHY SEPARATE THIS LOGIC?
   * - Avoids repeating the conversion code in multiple methods
   * - Creates a single place to update if the model changes
   * - Isolates database-specific knowledge from the rest of the code
   *
   * @param row The database row data
   * @returns A properly formatted Todo object
   */
  private rowToTodo(row: any): Todo {
    return {
      id: row.id,
      meetingID: row.meetingID,
      title: row.title,
      description: row.description,
      completedAt: row.completedAt,
      completed: row.completedAt !== null, // Computed from completedAt
      createdAt: row.createdAt,
      updatedAt: row.updatedAt,
      list: row.list,
      assignee: row.assignee,
    };
  }
}

// Create a singleton instance for use throughout the application
export const todoService = new TodoService();
