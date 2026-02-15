# Go URL Shortener

A full-stack URL shortener built with Go, featuring user authentication, SQLite storage, and a modern dark-themed web interface. Each user can create, manage, and organize their own shortened links.

---
# Screenshots

### Register
<img width="625" height="483" alt="image" src="https://github.com/user-attachments/assets/ec207a97-1098-4009-a72b-184da5312cab" />

### Login
<img width="569" height="475" alt="image" src="https://github.com/user-attachments/assets/da30a601-e776-47d0-9fbd-0835a95cf0fd" />

### Home
<img width="1125" height="628" alt="image" src="https://github.com/user-attachments/assets/a9877d8a-8a1a-46dd-8e77-784ce00af5cd" />

### URL list
<img width="1102" height="453" alt="image" src="https://github.com/user-attachments/assets/7cca40df-6698-48ec-99b2-71dac231dff4" />

---

## Features

### Core Features
- Shorten long URLs
- Custom short codes
- Redirect via `/r/{code}`
- Per-user URL management
- Delete links

### Authentication
- User registration and login
- Password hashing using bcrypt
- Session-based authentication
- Logout functionality

### Web Interface
- Dark theme UI (black + purple)
- Centered card layout
- Copy-to-clipboard button
- Search/filter user links
- Responsive and clean design

### Storage
- SQLite database (`data.db`)
- User data stored securely
- URLs stored per user

---

## Project Structure

```
go-url-shortener/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── actions/
│   │   └── url.go
│   │
│   ├── handlers/
│   │   └── api.go
│   │
│   ├── middlewares/
│   │   └── auth.go
│   │
│   ├── store/
│   │   ├── url_store.go
│   │   └── user_store.go
│   │
│   ├── sqlite/
│   │   ├── url_store.go
│   │   └── user_store.go
│   │
│   ├── templates/
│   │   ├── index.html
│   │   ├── list.html
│   │   ├── login.html
│   │   └── register.html
│   │
│   └── static/
│       └── css/
│           └── style.css
│
├── go.mod
└── README.md
```

---

## Architecture

The project follows a layered architecture:

```
HTTP Handlers
      ↓
Business Logic (Actions)
      ↓
Store Interfaces
      ↓
SQLite Implementation
```

Benefits:
- Separation of concerns  
- Easy to switch storage (JSON, SQLite, PostgreSQL, etc.)  
- Follows Open/Closed Principle  

---

## Requirements

- Go 1.20 or higher  
- SQLite driver  
- On Windows: GCC installed (required for CGO if using `github.com/mattn/go-sqlite3`)

---

## Installation

Clone the repository:

```bash
git clone https://github.com/hyarisara/go-url-shortener.git
cd go-url-shortener
```

Install dependencies:

```bash
go mod tidy
```

---

## Running the Application

Start the server:

```bash
go run ./cmd/server
```

Open your browser:

```
http://localhost:8080
```

---

## How to Use

1. Register a new account
2. Login
3. Enter a long URL
4. (Optional) Add a custom short code
5. Click **Shorten**
6. Copy or open the generated short link
7. View all links from **My URLs**

Example:

Original:
```
https://example.com/very/long/url
```

Short:
```
http://localhost:8080/r/abc123
```

---

## Database

SQLite file created automatically:

```
data.db
```

Tables:
- **users** → username, password hash
- **urls** → stored as `username::code → original URL`

---

## Deployment

Build the application:

```bash
go build -o app ./cmd/server
./app
```

The project can be deployed on:
- Render
- Railway
- Fly.io
- VPS (Linux)

Make sure the server has write permission for the SQLite file.

---

## Future Improvements

- Click analytics (visit count)
- Link expiration
- Rate limiting
- REST API endpoints
- Custom domain support
- Docker support
- PostgreSQL support

---

## Author

Sara Alhyari  
Computer Science Student | Data Analyst 

GitHub:  
https://github.com/hyarisara

---

## License

MIT License
