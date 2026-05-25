# Boot.dev course: Build a Blog Aggregator

## Requirements

- Golang 1.26+
- Postgres

## Installation

Regular Go binary package:
```sh
go install github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator@latest
```

## Configuration

A JSON `gatorconfig.json` with the following format is expected:
```json
{
  "db_url": "postgres://<postgres_user>:<postgres_pswd>@<db_domain>:<db_port>/<db_name>?sslmode=disable",
  "current_user_name": "<username>"
}
```
- `db_name` is the database you created in Postgres.

`gatorconfig.json` file's location is `$XDG_CONFIG_HOME/bootdotdev-go-build-a-blog-aggregator/`


## Usage

This blog aggregator is a REPL program, you interact with it in the following format:
```sh
bootdotdev-go-build-a-blog-aggregator <command> [cmd_arguments]
```

Supported commands include:
- `login <username>`
- `users`
- `addfeed <feed_name> <feed_url>`
- `feeds`
- `agg <scrape_frequency>`
- `browse [nb_feed_limit_to_print]`
  - default: two posts
