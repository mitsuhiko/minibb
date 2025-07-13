# MiniBB

MiniBB is a small bulletin board modeled after phpBB and 4chan. It doesn't use any user authentication but it has a way to prevent impersonation.

## Features

* Users can use tripcode to authenticate
* Admins can create multiple boards (admin permissions are hardcoded in an env var)
* Users can start multiple topics in those boards. 
* Each topic can have multiple replies. 
* There are about 20-30 posts per page, and we use pagination. 
* And we allow basic Markdown for formatting.
* Users can see which topics have new posts, and when they read a post, it's marked as read. 

## Tech Stack

* Single binary Go back-end (chi router with std http stack)
* We use TypeScript for the frontend and use prettier for formatting
* The front-end is built with React, Vite, Tailwind CSS4, and the tanstack query and router system.
* During development, we want to use the vite server
* We want to use SQLite and we want to use modernc.org/sqlite as package
* We put all our tools into a makefile (`dev` (starts frontend + backend), `check` (runs lint), `format` (formats)`, `build` (builds the production build))

## Tech Decisions

* Tripcode:
  * Use the 4chan newstyle tripcode algorithm
* Pagination:
  * Use cursors and not offsets
* Markdown:
  * use goldmark
* Read status:
  * we want to track the read status purely in the browser with local storage

## Deployments

* prod:
  * we use fs embed to embed the production build of the frontend app and serve it from the backend
* dev:
  * we use the vite server to proxy to the backend service
  * we rebuild the thing all the time on changes

## Database Schema

* boards
  * id (primary key)
  * slug (short code used in the url)
  * description (markdown summary)
* topics
  * id (primary key)
  * board_id (where they are stored in)
  * pub_date (timestamp of when the topic was created)
  * title (name of the topic, just text)
  * status (open, locked, might need more later)
  * author (name of the author + tripcode if used), denormalized
  * last_post_id (helps the read status tracking)
  * post_count (number of replies + initial post)
* posts
  * id (primary key)
  * topic_id (topic it belongs to)
  * pub_date (when it was written)
  * author (name of the author + tripcode if used)
  * content (markdown formatted)

## API Interface

* All api endpoints are in `/api/`
* We always send JSON in and out
* We use a rate limiter by IP address to prevent flooding the forums unnecessarily

## Read Status Tracking

Instead of storing every read post ID:

// Bad: huge storage
localStorage = {
  "read_posts": [1,2,3,4,5,6,7,8,9,10,11,12,...]
}

// Good: minimal storage
localStorage = {
  "read_topics": {
    "123": 567,  // topic 123: read up to post 567
    "124": 89    // topic 124: read up to post 89
  }
}

Client Logic

- When user reads a topic, store the highest post ID they've seen
- To check if topic has unread posts: topic.last_post_id > localStorage.read_topics[topic.id]
- Show unread count: topic.post_count - (posts_read_up_to_stored_id)
