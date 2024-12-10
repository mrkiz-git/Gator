# Gator

Gator is a multi-user CLI application for managing and browsing RSS feeds locally. It's designed for single-device use but supports multiple user profiles.

## Prerequisites

Before you begin, ensure you have met the following requirements:
* You have installed the latest version of [Go](https://golang.org/doc/install)
* You have a PostgreSQL database set up and running

## Installing Gator

To install Gator, follow these steps:

```
go install github.com/mrkiz-git/gator@latest
```

## Configuration

Create a JSON configuration file named `config.json` in the same directory as the Gator executable. The file should have the following structure:

```json
{
  "db_url": "your_postgresql_connection_string",
  "current_user_name": "username_goes_here"
}
```

Replace `your_postgresql_connection_string` with your actual PostgreSQL connection string.

## Using Gator

Gator provides several commands for managing users, feeds, and posts:

* `gator login`: Log in as a user
* `gator register`: Register a new user
* `gator reset`: Reset the database
* `gator users`: List all users
* `gator agg`: Fetch RSS feeds
* `gator addfeed`: Add a new RSS feed (requires login)
* `gator feeds`: List all feeds
* `gator follow`: Follow a feed (requires login)
* `gator following`: List followed feeds (requires login)
* `gator unfollow`: Unfollow a feed (requires login)
* `gator brows`: Browse posts

To use a command, run:

```
gator [command]
```

For commands that require login, ensure you're logged in first using the `login` command.

## Note on Security

This application doesn't include user-based authorization. Anyone with database credentials can act as any user. This design is intentional for learning purposes, focusing on SQL, CLIs, and long-running services.

## License

This project is licensed under [specify your license here].

Sources
