# Gator

## Introduction
Welcome to the Gator Feed Aggregator! This is a CLI tool which allows users to:
- Add RSS feeds from the internet for collection
- Store posts in a PostgreSQL Database
- Follow and unfollow RSS feeds other users have added
- View summaries of the posts in the terminal with links to the posts

There are no servers involved besides the one storing the database, this is intended primarily for local use.
To view a list of commands, type gator help in the terminal. Then you can go from there to conduct all your feed aggregator needs.
To get started you can use the register command followed by a name to register a new user.
You can use the addfeed command followed by a name for the feed and it's url to add a feed for the current user.

## Install libraries
To run the program, you will need to have go and postgres installed.

You can install go using the webi installer. Run this in your terminal:
```bash
curl -sS https://webi.sh/golang | sh
```

You can install PostgreSQL on macOS with brew:
```bash
brew install postgresql@15
```
Or on windows or linux:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Ensure the installation worked with
```bash
psql --version
```
## Config
The program looks for a .gatorconfig.json file in the home directory. Create that file now, and store this code in it:
```
{"db_url":"url","current_user_name":"username"}
```
The username will be overwritten when you run the register command.
The db_url will be the one for your postgres system, and will look like: postgres://username:@localhost:5432/gator.
If you have a password, it will go after the username: and before the @
Add the following flags to the end: ?sslmode=disable, making it look like: postgres://username:@localhost:5432/gator?sslmode=disable.

## Clone
To clone the repo, enter this into the terminal
```bash
git clone https://github.com/jacobhuneke/gator.git
cd gator
```

## Database Migrations
To create the necessary tables, install [goose](https://github.com/pressly/goose) and run the following command from the `sql/schema` directory:
```bash
goose postgres <your_db_url> up
```

## Install Gator
You can install Gator to your system so you dont have to have the file running to use it's functions. You can call them from anywhere in the terminal.
You can do this by running this prompt in the terminal:
```bash
go install github.com/jacobhuneke/gator@latest
```

