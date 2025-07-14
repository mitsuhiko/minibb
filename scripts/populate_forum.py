#!/usr/bin/env python3
"""
Forum Population Script for MiniBB

This script generates forum conversation data using Claude Code CLI and populates
the MiniBB database with realistic test data including topics and posts
across multiple boards.

Usage:
    uv run scripts/populate_forum.py <num_threads> <posts_per_thread>
    
Examples:
    uv run scripts/populate_forum.py 10 5     # Generate 10 threads with 5 posts each
    uv run scripts/populate_forum.py 20 8     # Generate 20 threads with 8 posts each
"""

import json
import sqlite3
import sys
import subprocess
from datetime import datetime, timedelta
import random
from pathlib import Path
import argparse

def generate_conversations_with_claude(num_threads, posts_per_thread):
    """Generate conversation data using Claude Code CLI."""
    prompt = f"""Generate realistic forum conversations in JSON format for a bulletin board system. Create exactly {num_threads} conversations with exactly {posts_per_thread} posts each.

Requirements:
- Use boards: "general", "watercooler", "test"
- Each conversation should have a realistic title and engaging discussion
- Authors should use realistic usernames, some with tripcodes (format: username!tripcode)
- Posts should feel natural and conversational
- Topics can be about programming, technology, hobbies, current events, or general discussion
- Vary the writing styles and personalities

Output ONLY valid JSON in this exact format (NO MARKDOWN WRAPPER OR ADDITIONAL TEXT):

[
  {{
    "board": "general",
    "title": "Topic title here",
    "posts": [
      {{
        "author": "username",
        "content": "Post content here"
      }},
      {{
        "author": "username!tripcode",
        "content": "Reply content here"
      }}
    ]
  }}
]

IMPORTANT: Do not use any tools or external resources. Generate the complete response directly without any tool invocations. Output only the JSON array, no explanations or additional text."""

    try:
        print(f"Generating {num_threads} conversations with {posts_per_thread} posts each using Claude Code...")
        
        # Call Claude Code CLI in direct mode
        result = subprocess.run([
            'claude', 
            '--model',
            'sonnet',
            "-p",
            prompt
        ], capture_output=True, text=True, check=True)
        
        # Parse the JSON response
        try:
            conversations = json.loads(result.stdout.strip())
            print(f"Successfully generated {len(conversations)} conversations")
            return conversations
        except json.JSONDecodeError as e:
            print(f"Error: Claude Code returned invalid JSON: {e}")
            print("Raw output:")
            print(result.stdout)
            sys.exit(1)
            
    except subprocess.CalledProcessError as e:
        print(f"Error calling Claude Code CLI: {e}")
        print("Stderr:", e.stderr)
        sys.exit(1)
    except FileNotFoundError:
        print("Error: Claude Code CLI not found. Please install it first.")
        print("Visit: https://claude.ai/code for installation instructions")
        sys.exit(1)

def get_board_id(cursor, board_slug):
    """Get board ID by slug, or None if not found."""
    cursor.execute("SELECT id FROM boards WHERE slug = ?", (board_slug,))
    result = cursor.fetchone()
    return result[0] if result else None

def create_topic(cursor, board_id, title, author, pub_date):
    """Create a new topic and return its ID."""
    cursor.execute("""
        INSERT INTO topics (board_id, title, author, pub_date, post_count)
        VALUES (?, ?, ?, ?, 0)
    """, (board_id, title, author, pub_date))
    return cursor.lastrowid

def create_post(cursor, topic_id, author, content, pub_date):
    """Create a new post and return its ID."""
    cursor.execute("""
        INSERT INTO posts (topic_id, author, content, pub_date)
        VALUES (?, ?, ?, ?)
    """, (topic_id, author, content, pub_date))
    return cursor.lastrowid

def update_topic_stats(cursor, topic_id, last_post_id, post_count):
    """Update topic's last post ID and post count."""
    cursor.execute("""
        UPDATE topics 
        SET last_post_id = ?, post_count = ?
        WHERE id = ?
    """, (last_post_id, post_count, topic_id))

def generate_realistic_timestamps(post_count, start_days_ago=30):
    """Generate realistic timestamps for posts in a conversation."""
    # Start the conversation some time in the past
    start_time = datetime.now() - timedelta(days=random.randint(1, start_days_ago))
    
    timestamps = [start_time]
    current_time = start_time
    
    # Generate subsequent post times with realistic gaps
    for i in range(1, post_count):
        # Posts can be minutes to hours apart, with some clustering
        if random.random() < 0.3:  # 30% chance of quick reply (within an hour)
            gap = timedelta(minutes=random.randint(2, 60))
        elif random.random() < 0.6:  # 60% chance of same-day reply
            gap = timedelta(hours=random.randint(1, 12))
        else:  # 10% chance of next-day or later reply
            gap = timedelta(days=random.randint(1, 3), hours=random.randint(0, 12))
        
        current_time += gap
        timestamps.append(current_time)
    
    return timestamps

def populate_database(db_path, conversations):
    """Populate the database with conversation data."""
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        total_posts = 0
        total_topics = 0
        
        print(f"Populating database with {len(conversations)} conversations...")
        
        for i, conversation in enumerate(conversations, 1):
            board_slug = conversation['board']
            title = conversation['title']
            posts = conversation['posts']
            
            # Get board ID
            board_id = get_board_id(cursor, board_slug)
            if board_id is None:
                print(f"Warning: Board '{board_slug}' not found, skipping conversation '{title}'")
                continue
            
            # Generate realistic timestamps for this conversation
            timestamps = generate_realistic_timestamps(len(posts))
            
            # Create the topic with the first post's author and timestamp
            first_post = posts[0]
            topic_id = create_topic(cursor, board_id, title, first_post['author'], timestamps[0])
            total_topics += 1
            
            # Create all posts for this topic
            last_post_id = None
            for j, post in enumerate(posts):
                post_id = create_post(cursor, topic_id, post['author'], post['content'], timestamps[j])
                last_post_id = post_id
                total_posts += 1
            
            # Update topic statistics
            update_topic_stats(cursor, topic_id, last_post_id, len(posts))
            
            print(f"  Created topic {i}/{len(conversations)}: '{title}' with {len(posts)} posts")
        
        conn.commit()
        print(f"\nSuccessfully populated database:")
        print(f"  Topics created: {total_topics}")
        print(f"  Posts created: {total_posts}")
        
    except sqlite3.Error as e:
        print(f"Database error: {e}")
        sys.exit(1)
    finally:
        if conn:
            conn.close()

def main():
    """Main function."""
    # Parse command line arguments
    parser = argparse.ArgumentParser(
        description="Generate and populate MiniBB forum with realistic conversations using Claude Code CLI"
    )
    parser.add_argument("num_threads", type=int, help="Number of forum threads to generate")
    parser.add_argument("posts_per_thread", type=int, help="Number of posts per thread")
    
    args = parser.parse_args()
    
    if args.num_threads <= 0 or args.posts_per_thread <= 0:
        print("Error: Number of threads and posts per thread must be positive integers")
        sys.exit(1)
    
    # Set up paths
    script_dir = Path(__file__).parent
    project_root = script_dir.parent
    db_path = project_root / "minibb.db"
    
    # Check if database exists
    if not db_path.exists():
        print(f"Error: Database file '{db_path}' not found")
        print("Please run this script from the project root or ensure the database exists")
        sys.exit(1)
    
    # Generate conversations using Claude Code CLI
    conversations = generate_conversations_with_claude(args.num_threads, args.posts_per_thread)
    
    # Populate database
    populate_database(db_path, conversations)
    
    print("\nDatabase population complete!")
    print("You can now view the populated forum in your MiniBB application.")

if __name__ == "__main__":
    main()