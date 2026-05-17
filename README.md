# My First Go API 

Hey! This is a backend API for a simple Blog system that I built in order to practice my Go (Golang) skills. It handles user registration, login with JWT, and full CRUD operations for blog posts with some neat filtering features.

I also containerized the whole app using Docker Compose so it runs anywhere without messing up with local database installations.

---

## 🛠️ What I Used

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Containerization:** Docker & Docker Compose

---

## 💡 Features I Implemented

- **User System:** People can sign up and log in. Passwords are safely hashed using `bcrypt`.
- **JWT Protection:** Sensitive actions like creating, updating, or deleting blogs require a valid token.
- **Blog Operations:** Users can create posts, and authors can modify/delete their own posts.
- **Query Filters:** Added pagination (`page`, `limit`), global text `search`, dynamic `field` filtering, and `sort` (asc/desc).

---

## 🚀 How to Use

### 1. Requirements
You just need to have `Docker Desktop` installed and running on your computer. `No need to install Go` locally if you use Docker, as the Go runtime comes built-in with the container.

### 2. Setup Env Variables
Create a file named `.env` in the root folder and add your local configs:

````env
DB_HOST=localhost
DB_PORT=5432
DB_USER=whoami
DB_PASSWORD=my_super_secret_password
DB_NAME=demodb
JWT_SECRET_KEY=my_local_secret_key
````

### 3. Start everyting
Open your terminal in the project and run this command:

````bash
docker-compose up --build # Use --build if running for the first time or after code changes
# To close everything, just hit [Ctrl + C]
# To close everything cleanly, press [CTRL + C] and then run: docker-compose down
````

---
## 📝 Blog Management (CRUD)

### 🔌 API Endpoints Reference

| Method | Endpoint | Middleware / Auth | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/register` | None | Register a new user account |
| `POST` | `/login` | None | Login and receive a signed JWT token |
| `GET` | `/blogs` | None | List all public blogs (Supports pagination, search, sorting) |
| `GET` | `/blogs/{id}` | Optional Auth | View a single blog post details |
| `POST` | `/blogs` | Required Auth | Create a new blog post |
| `PATCH` | `/blogs/{id}` | Required Auth | Update a blog post (Author only) |
| `DELETE` | `/blogs/{id}` | Required Auth | Delete a blog post (Author only) |
| `GET` | `/users/{id}/blogs` | Optional Auth | View all blogs written by a specific user |
| `PATCH` | `/users/{id}` | Required Auth | Update user profile details (Owner only) |
| `DELETE` | `/users/{id}` | Required Auth | Delete user account and all their blogs (Owner only) |

## 🔍 Query Parameters for GET Requests

Here is how you can use query parameters.

#### 1. Pagination
Specifies which page of the results to retrieve.
- Example Request: `http://localhost:8000/blogs?page=2`

#### 2. Limit
Controls the maximum number of blog posts returned in a single page.
- Example Request: `http://localhost:8000/blogs?limit=5`

#### 3. Global Search
Searches for a specific keyword across the blog title and content fields.
- Example Request: `http://localhost:8000/blogs?search=golang`

#### 4. Dynamic Field Filtering
Filters the results by specifying columns and their matching values.
- Example Request: `http://localhost:8000/blogs?field=title,content`

#### 5. Sorting
Sorts the blog posts in ascending (`asc`) or descending (`desc`) order based on their creation time.
- Example Request: `http://localhost:8000/blogs?sort=desc`

#### 6. Mixing Everything Together
You can combine multiple query parameters using the `&` symbol to build advanced queries.
- Example Request: `http://localhost:8000/blogs?page=1&limit=10&search=docker&sort=desc`

----

## 📦 API Payloads & Responses

````request
$ http://localhost:8000/users/5/blogs

~ If the request is made by the blog owner, hidden blogs (isPublic: false) will also be returned

{
  "total_count": 3,
  "response": [
    {
      "id": 1,
      "author_id": 5,
      "created_at": "2026-05-16T20:56:54.931799Z",
      "updated_at": "2026-05-16T20:56:54.931799Z",
      "title": "golang",
      "content": "pointers",
      "is_public": true
    },
    {
      "id": 2,
      "author_id": 5,
      "created_at": "2026-05-16T20:56:59.031547Z",
      "updated_at": "2026-05-16T20:56:59.031547Z",
      "title": "golang",
      "content": "pointers 1",
      "is_public": true
    },
    {
      "id": 4,
      "author_id": 5,
      "created_at": "2026-05-16T20:57:14.988785Z",
      "updated_at": "2026-05-16T20:57:14.988785Z",
      "title": "python",
      "content": "error handling",
      "is_public": true
    }
  ]
}
````

````request
$ http://localhost:8000/blogs?page=1&limit=10&field=title,is_public&search=go

~ The filtered_count key represents the number of blog posts returned based on the applied query parameters.

{
  "total_count": 5,
  "filtered_count": 3,
  "response": [
    {
      "title": "golang",
      "is_public": true
    },
    {
      "title": "golang",
      "is_public": true
    }
  ]
}
````