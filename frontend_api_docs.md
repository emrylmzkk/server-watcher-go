# Server Watcher API Documentation for Frontend Development

This documentation provides details about the API endpoints, request/response structures (DTOs), and enumerations required for developing the frontend application.

## API Common Response Structure
All responses are wrapped in this structure:
```json
{
  "success": boolean,
  "message": "string (optional)",
  "data": object | array | null,
  "error": object | string | null
}
```


### Login
*   **Endpoint:** `POST /login`
*   **Request DTO:**
    ```json
    {
      "username": "string",
      "password": "string"
    }
    ```
*   **Response Data:**
    ```json
    {
      "access_token": "string",
      "refresh_token": "string"
    }
    ```

### Refresh Token
*   **Endpoint:** `POST /refresh-token`
*   **Request DTO:**
    ```json
    {
      "refresh_token": "string"
    }
    ```
*   **Response Data:**
    ```json
    {
      "access_token": "string",
      "refresh_token": "string"
    }
    ```

### Get Current User Information
*   **Endpoint:** `GET /me`
*   **Middleware:** `AuthToken` Required
*   **Request DTO:** None
*   **Response Data:**
    ```json
    {
      "id": number,
      "name": "string",
      "surname": "string",
      "user_role": number (See Enums)
    }
    ```

---

## 2. PM2 Management (`/api/v1/pm2`)
*All require `AuthToken`*

### Get All Projects
*   **Endpoint:** `GET /projects`
*   **Response Data:** `[]Pm2ProjectResponseDTO`
    ```json
    [
      {
        "id": number,
        "external_id": "string",
        "name": "string",
        "project_type": number,
        "status": "string",
        "project_path": "string",
        "project_start_command": "string",
        "project_runtime_tpye": number
      }
    ]
    ```

### Create Project
*   **Endpoint:** `POST /create`
*   **Request DTO:**
    ```json
    {
      "external_id": "string",
      "name": "string",
      "project_path": "string",
      "project_start_command": "string",
      "project_runtime_tpye": number
    }
    ```
*   **Response Data:** `boolean`

### Update Project
*   **Endpoint:** `PUT /:id`
*   **Request DTO:**
    ```json
    {
      "name": "string",
      "project_path": "string (optional)",
      "project_start_command": "string (optional)",
      "project_runtime_tpye": number (optional)
    }
    ```
*   **Response Data:** `boolean`

### Delete Project
*   **Endpoint:** `DELETE /:id`
*   **Response Data:** `boolean`

### Start/Stop/Reset
*   **Endpoints:** 
    *   `POST /start/:id`
    *   `POST /stop/:id`
    *   `POST /reset`
*   **Response Data:** `boolean`

### Get PM2 Real-time List
*   **Endpoint:** `GET /inside-list`
*   **Response Data:** `[]Pm2InsideListResponseDTO`
    ```json
    [
      {
        "pm_id": number,
        "name": "string",
        "pm2_env": { 
          "status": "string",
          "pm_cwd": "string",
          "pm_exec_path": "string",
          "interpreter": "string",
          "restart_time": number
        },
        "monit": { "cpu": float, "memory": number }
      }
    ]
    ```

### Sync PM2 Projects
*   **Endpoint:** `POST /sync-projects`
*   **Description:** Synchronizes local PM2 processes with the database. Adds new processes, updates status, and marks missing ones as "Deleted".
*   **Response Data:** `boolean`

### Get PM2 Project by ID
*   **Endpoint:** `GET /:id`
*   **Response Data:** `Pm2ProjectResponseDTO`
    ```json
    {
      "id": number,
      "external_id": "string",
      "name": "string",
      "project_type": number,
      "status": "string",
      "project_path": "string",
      "project_start_command": "string",
      "project_runtime_tpye": number
    }
    ```

### Clear and Delete PM2 Project
*   **Endpoint:** `POST /:id/clear-project`
*   **Description:** Runs `pm2 delete` on the server and removes the project from the database.
*   **Response Data:** `boolean`

---

## 3. Docker Management (`/api/v1/docker`)
*All require `AuthToken`*

### Get Containers
*   **Endpoint:** `GET /containers`
*   **Response Data:**
    ```json
    [
      {
        "name": "string",
        "isRunning": boolean,
        "uptime": "string"
      }
    ]
    ```

### Get Containers with Stats
*   **Endpoint:** `GET /containers-stats`
*   **Response Data:** `[]DockerStats`
    ```json
    [
      {
        "name": "string",
        "is_running": boolean,
        "uptime": "string",
        "cpu": float,
        "memory_mb": float,
        "network_rx": number,
        "network_tx": number
      }
    ]
    ```

### Start/Stop
*   **Endpoints:**
    *   `POST /start/:id`
    *   `POST /stop/:id`
*   **Response Data:** `string` (e.g., "Container started")

---

## 4. Container Stats (`/api/v1/docker-stats`)

### Get Stats by Name
*   **Endpoint:** `GET /:containerName`
*   **Response Data:** `DockerStats` (See structure in Docker Management)

### Get Paginated Stats
*   **Endpoint:** `GET /:containerName/paginate`
*   **Query Params:** `page`, `limit`, `sort_by`, `sort_order`
*   **Response Data:** `PaginatedResult<DockerStats>`
    ```json
    {
      "data": [ ...DockerStats ],
      "page": number,
      "limit": number,
      "total": number,
      "total_pages": number
    }
    ```

---

## 5. Server General Stats (`/api/v1/server-general`)

| Method | Endpoint | Response Data (In `data` field) |
| :--- | :--- | :--- |
| GET | `/cpu` | `float64` (Percentage) |
| GET | `/ram` | `{ "used": float, "total": float, "percent": float }` |
| GET | `/disk` | `{ "used": float, "total": float, "percent": float }` |
| GET | `/stats` | `object` (General System Stats) |

---

## 6. Enumerations

### `ProcessType` (int)
- `1`: Docker, `2`: PM2

### `ProjectRuntimeType` (int)
- `1`: Go, `2`: Python, `3`: React, `4`: Expo

### `ProjectStatus` (string)
- `"exited"`, `"closed"`, `"dead"`

### `UserRole` (int)
- `95`: Admin, `2`: Member
