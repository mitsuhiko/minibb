# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MiniBB is a small bulletin board system inspired by phpBB and 4chan with no user authentication but tripcode support for identity verification. It's a full-stack application with a Go backend and React frontend.

## Development Commands

```bash
make dev          # Development (starts both frontend and backend.  This autoreloads and auto compiles.  Don't ever stop the server)
make build        # Production build (we rarely need this)
make check        # Linting and type checking (for backend and frontend)
make format       # Code formatting (for backend and frontend)
make clean        # Clean build artifacts (we rarely need this)
make tail-log     # Reads the current log file (last 100 lines of code)
```

**IMPORTANT:**

* The server and the frontend log everything into the dev.log file.
* Use the `make tail-log` command to read the log file. 
* Never stop the server! It keeps running.  It auto compiles and auto reloads.  It does log to `dev.log`
* If you fail to run the Makefile, you have to remember that you have to run it from the top-level directory.

## Architecture

**Backend (Go)**
- Single binary using Chi router with standard HTTP stack
- SQLite database using `modernc.org/sqlite`
- Embedded frontend files for production builds
- API endpoints under `/api/`
- Rate limiting by IP address

**Frontend (React/TypeScript)**
- We always use `npm` as a package manager
- React with Vite build system
- TanStack Query for data fetching
- TanStack Router for routing  
- Tailwind CSS v4 for styling
- Development proxy to backend on port 8080

**Development Setup**
- Uses Procfile with custom shoreman.sh script
- Frontend dev server proxies `/api` requests to backend
- Watchexec for auto-reloading Go backend on file changes

## Database Schema

- To see the database schema that is currently used, just use `sqlite3 minibb.db` and run `.schema`

## Key Features

- **Tripcode Authentication**: Uses 4chan newstyle tripcode algorithm
- **Read Status Tracking**: Client-side localStorage tracking highest read post ID per topic
- **Pagination**: Cursor-based (not offset-based)
- **Markdown Support**: Uses goldmark for rendering

## Development Notes

- Backend runs on port 8080 by default
- Frontend dev server proxies to backend and runs on port 5173
- Production embeds frontend files in Go binary using fs.embed
- Environment variable PORT can override default port