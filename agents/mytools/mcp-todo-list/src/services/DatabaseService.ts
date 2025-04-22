/**
 * DatabaseService.ts
 *
 * This file implements a lightweight SQLite database service for the Todo application.
 *
 * WHY SQLITE?
 * - SQLite is perfect for small to medium applications like this one
 * - Requires no separate database server (file-based)
 * - ACID compliant and reliable
 * - Minimal configuration required
 * - Easy to install with minimal dependencies
 */
import Database from 'better-sqlite3';
import { config, ensureDbFolder } from '../config.js';

/**
 * DatabaseService Class
 *
 * This service manages the SQLite database connection and schema.
 * It follows the singleton pattern to ensure only one database connection exists.
 *
 * WHY SINGLETON PATTERN?
 * - Prevents multiple database connections which could lead to conflicts
 * - Provides a central access point to the database throughout the application
 * - Makes it easier to manage connection lifecycle (open/close)
 */
class DatabaseService {
  private db: Database.Database;

  constructor() {
    // Ensure the database folder exists before trying to create the database
    ensureDbFolder();

    // Initialize the database with the configured path
    this.db = new Database(config.db.path);

    /**
     * Set pragmas for performance and safety:
     * - WAL (Write-Ahead Logging): Improves concurrent access performance
     * - foreign_keys: Ensures referential integrity (useful for future expansion)
     */
    this.db.pragma('journal_mode = WAL');
    this.db.pragma('foreign_keys = ON');

    // Initialize the database schema when service is created
    this.initSchema();
  }

  /**
   * Initialize the database schema
   *
   * This creates the todos table if it doesn't already exist.
   * The schema design incorporates:
   * - TEXT primary key for UUID compatibility
   * - NULL completedAt to represent incomplete todos
   * - Timestamp fields for tracking creation and updates
   */
  private initSchema(): void {
    // Create todos table if it doesn't exist
    this.db.exec(`
      CREATE TABLE IF NOT EXISTS todos (
        id TEXT PRIMARY KEY,
        meetingID TEXT NOT NULL,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        completedAt TEXT NULL, -- ISO timestamp, NULL if not completed
        createdAt TEXT NOT NULL,
        updatedAt TEXT NOT NULL,
        list TEXT NOT NULL,
        assignee TEXT NOT NULL
      )
    `);
  }

  /**
   * Get the database instance
   *
   * This allows other services to access the database for operations.
   *
   * @returns The SQLite database instance
   */
  getDb(): Database.Database {
    return this.db;
  }

  /**
   * Close the database connection
   *
   * This should be called when shutting down the application to ensure
   * all data is properly saved and resources are released.
   */
  close(): void {
    this.db.close();
  }
}

// Create a singleton instance that will be used throughout the application
export const databaseService = new DatabaseService();
