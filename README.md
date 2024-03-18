# Tune Tracker

## Problem
[Spotify](https://open.spotify.com/) is a great source for DJs to find new records. However, getting copies of these tracks to be used in mixing sets is a manual task which includes but is not limited to the following steps:

- Buying / Downloading tracks. Features idle time while waiting for a download to complete.
- Arranging files by folders / labelling track files to simplify mxing in record box or another mixing program. 

## Solution
Tune Tracker (name is not final) is a web application that aims to solve this problem giving the users an opportunity to provide a spotify playlist URL to monitor for new tracks. Locating, purchasing and downloading of the music will be handled automatically.

## Architecture

TBD

## Running the App Locally
The application relies uses docker to run an instance of [Postgres](https://www.postgresql.org/). Therefore, before running the application you need to make sure that [docker](https://www.docker.com/) is installed on your system. Once this requirement is met, complete the following steps to launch a local instance of Tune Tracker:

Clone the repostiory & install the app's dependencies:
```bash
git clone https://github.com/FirstDayAtWork/spfy-app-pages
cd spfy-app-pages
go get .
go mod download
```

Install [modd](https://github.com/cortesi/modd), it is needed for live-reloading of the application. Helps with the development process A LOT:
```bash
go install github.com/cortesi/modd/cmd/modd@latest
```

In the root folder of the project, create a `.env` file. Ask project owners about the contents.

Start docker app aka docker deamon. Without it the following steps will not work.

### Instructions for Windows
Launch the app's database with Docker:
```bash
docker compose -f dev_compose.yaml up --build
```

Start the application using modd:
```bash
modd -f dev_modd.conf
```


## Testing

TBD