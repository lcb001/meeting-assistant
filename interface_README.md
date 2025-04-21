# Meeting API Documentation

This document provides detailed information about the Meeting API endpoints and how to use them.

## API Endpoints

### 1. Create Meeting
Creates a new meeting and returns a meeting ID.

**Endpoint:** `POST /meeting`

**Request Body:**
```json
{
  "title": "Team Weekly Sync",
  "description": "Weekly team sync meeting",
  "participants": ["john@example.com", "jane@example.com"]
}
```

**Response:**
```json
{
  "id": "meeting_123abc"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8888/meeting \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Weekly Sync",
    "description": "Weekly team sync meeting",
    "participants": ["john@example.com", "jane@example.com"]
  }'
```

### 2. List Meetings
Retrieves a list of all meetings.

**Endpoint:** `GET /meeting`

**Response:**
```json
{
  "meetings": [
    {
      "id": "meeting_123abc",
      "content": {
        "title": "Team Weekly Sync",
        "description": "Weekly team sync meeting",
        "participants": ["john@example.com", "jane@example.com"]
      }
    }
  ]
}
```

**Curl Example:**
```bash
curl -X GET http://localhost:8888/meeting
```

### 3. Get Meeting Summary
Retrieves the summary of a specific meeting.

**Endpoint:** `GET /summary`

**Query Parameters:**
- `meeting_id` (required): The ID of the meeting

**Response:**
```json
{
  "summary": "Meeting discussion points and conclusions..."
}
```

**Curl Example:**
```bash
curl -X GET "http://localhost:8888/summary?meeting_id=meeting_123abc"
```

### 4. Start Chat Session
Initiates a Server-Sent Events (SSE) connection for real-time chat updates.

**Endpoint:** `GET /chat`

**Query Parameters:**
- `meeting_id` (required): The ID of the meeting
- `session_id` (required): The ID of the chat session

**Response:**
Server-Sent Events stream with messages in the following format:
```json
{
  "data": {
    "message": "Chat message content",
    "timestamp": "2024-03-21T10:00:00Z",
    "sender": "john@example.com"
  }
}
```

**Curl Example:**
```bash
curl -X GET "http://localhost:8888/chat?meeting_id=meeting_123abc&session_id=session_xyz789&message=Hello"
```


## Content Types

- All regular endpoints use `application/json` for request and response bodies
- The chat endpoint uses `text/event-stream` for Server-Sent Events streaming 