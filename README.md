This is a CLI application that helps follow RSS feeds.
Requirements:
  - PostgresSQL 16.8
  - Go 1.24.2
  - Goose: github.com/pressly/goose

Installation steps:
  - install the binary with [go install]
  - create .gatorconfig.json file at home dir with the following contents:
    {
      "Db_url": [postgres connection string],
      "Current_user_name": ""
    }
  - dir into schema folder: /sql/schema
  - migrate neccessary tables with [goose up] (after setting up environment variables for goose)

Supported commands (call the app with 'bootdev-gator [command]'):
  - reset: reset the database
  - users: list all registered users
  - register [username]: register user and login
  - login [username]: login with username
  - addfeed [name] [url]: add a RSS feed
  - feeds: list all feeds
  - follow [url]: follow a specific feed (for current logged in user)
  - following: list all following feeds of current logged in user
  - unfollow [url]: unfollow a specific feed (for currenet logged in user)
  - browse [limit]: print [limit] posts from following feeds
  - agg [duration] (e.g '1h2m3s'): fetch posts from registered feeds every [duration] (only 1 feed is processed per [duration])
    + WARNING: this command will initiate a infinite loop, so remember to run it in a seperate terminal.

Disclaimer: This is a demo Golang project guided by bootdev.
