#  TaskFlow
**TaskFlow is a Task Management application.**

TaskFlow is a production-ready REST API built with the Gin framework, focusing on project organization, task tracking, and secure JWT-based authentication.

---

##  Tech Stack
* **Language:**  Go (Gin Framework)
* **Database:** PostgreSQL with GORM
* **Migrations:** golang-migrate
* **DevOps:** Docker & Docker Compose
* **Auth:** JWT (JSON Web Tokens)
* **API:** REST API

---

# Architecture Decisions
##  **Structure**
**I followed a layered architecture:**

* **Controllers:**  handle HTTP requests, validation, and responses
* **Services:**  contain business logic and database operations
* **Routes:**  define API endpoints and group them logically
* **Middleware:**  handle authentication (JWT)
* **Config:**  manage database connection and environment setup

**This separation keeps the codebase clean, testable, and easier to scale.**

## Why This Approach?
* I used Gin because it is lightweight, fast, and widely used in Go backend systems.
* I separated controllers and services to avoid mixing HTTP logic with business logic.
* I used JWT authentication because it is stateless and easy to integrate with APIs.
* I used manual migrations (golang-migrate) instead of GORM auto-migrate to have full control over schema changes.
  
## Tradeoffs I Made
* I did not implement a refresh token mechanism → kept auth simple (only access token with 24h expiry).
* I used basic validation and error handling → not fully standardized across all endpoints.
* I focused more on functionality than optimization → queries and indexing can be improved.
* I used UUIDs for IDs which are good for distributed systems, but slightly heavier than integers.

---

##  Getting Started

### Prerequisites
* Docker Desktop installed required.

### Installation & Run
1. **Clone the repository**
   ```bash
   git clone https://github.com/madan-kumar-tm/taskflow-madan-kumar.git
   cd taskflow-madan-kumar/backend
   ```
2. **Setup Environment(Update DB_PASSWORD and JWT_SECRET field)**
   ```bash
   cp .env.example .env
   ```
3. **Spin up containers**
   ```bash
   docker compose up --build
   ```
The API will be live at `http://localhost:8080`.

---

## ER Diagram
  ![alt text](image.png)

##  Authentication
**Default Test User:**
* **Email:** `test@example.com`
* **Password:** `password123`

---
##  API Reference
**Import collection**
   ```bash
    https://github.com/madan-kumar-tm/taskflow-madan-kumar/tree/main/postman/collections
   ```
* **Steps:**
* **1.** Import collection into Postman
* **2.** Set environment variable base_url
* **3.** Run Login API to generate token
* **4.** Use token for all protected routes


| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/auth/register` | Create a new account |
| `POST` | `/auth/login` | Receive JWT access token |

### Projects
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/projects` | List all projects |
| `POST` | `/projects` | Create a new project |
| `GET` | `/projects/:id` | Get project details & tasks |
| `PATCH` | `/projects/:id` | Update project info |
| `DELETE` | `/projects/:id` | Remove project |

### Tasks
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/projects/:id/tasks` | List tasks (Supports `status` & `assignee` filters) |
| `POST` | `/projects/:id/tasks` | Create task within a project |
| `PATCH` | `/tasks/:id` | Update task status or details |
| `DELETE` | `/tasks/:id` | Remove task |

---

## 🛠 Development Notes

### Database Migrations
Migrations run automatically on startup. To check migration status:
```bash
docker logs taskflow-migrate
```

### What I'D Do With More Time
* [ ] Pagination for project and task lists.
* [ ] Comprehensive Integration Test Suite
* [ ] Add refresh token system for better auth.
* [ ] Add frontend (React dashboard).
* [ ] Improve error handling consistency

---

**Author:** [Madan Kumar T M](https://github.com/madan-kumar-tm)
