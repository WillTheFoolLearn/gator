# Gator
Gator RSS Blog Feed Aggregator project\

## Pre-requisites
PostgreSQL and Go need to be installed\

## Build Gator
Run `go install` to install the Gator CLI from the Gator directory in the terminal.\

## Config file
A config file needs to be created in your home directory. Name the file .gatorconfig.json containing the following:\

`{
  "db_url": "postgres://example"
}`

## Running Gator
After installing the Gator CLI per instructions above, the following commands can be used to interact with it:\

`gator login <username>` Logs in as given username\
`gator register <username>` Adds a new user to Gator\
`gator reset` Removes all users\
`gator users` Lists all users\
`gator agg <time duration>` Scrapes the RSS feeds to collect the posts every `<time duration>` provided. Examples of time duration are 5s, 1m, 2h, etc. (Please don't make it too short)\
`gator addfeed <"feed name"> <feed url>` Adds a new RSS feed\
`gator feeds` Lists all feeds\
`gator follow <url>` Adds an RSS feed to the logged in user's follow list\
`gator following` Lists all RSS feeds the current user is following\
`gator unfollow <url>` Unfollows an RSS feed that the logged in user is following\
`gator browse <count>` Shows the `<count>` most recent posts from RSS feeds that have been scraped that the logged in user is following\
