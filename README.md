# Go URL Shortener

A lightweight URL shortening service written in Go.  
The application provides both a CLI interface and a browser-based UI, with persistent JSON storage and an architecture designed for easy migration to databases like SQLite or PostgreSQL.

---

## Features

- Shorten long URLs into compact codes
- Redirect using browser-friendly routes
- Persistent storage in JSON
- Storage abstraction via interfaces
- Thread-safe operations
- List and delete stored URLs
- Simple, clean web interface
- Ready for future upgrades (SQL / Redis / Auth / Analytics)

---

## Architecture

The project follows a layered design:

---

## User Interface Design Screenshot
- Login/Register 
<img width="748" height="490" alt="image" src="https://github.com/user-attachments/assets/366b85b7-4b23-4633-9a71-9236c930546d" />

- User Home 
<img width="858" height="366" alt="image" src="https://github.com/user-attachments/assets/2e17decd-0895-4d21-b81b-c9bfe911b09a" />

- Output
<img width="811" height="420" alt="image" src="https://github.com/user-attachments/assets/3fb5b8d3-24b7-4f96-99a0-56c32e732d3b" />
