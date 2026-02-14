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

