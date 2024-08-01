# Z

> 🌎 The new real-time login, multiplayer, & chat server for Quaver.

**Z** is one of many services that power Quaver and is the successor to [Albatross](https://github.com/Quaver/Albatross). 

This server handles any and all real-time game events such as:

* Client login & user sessions
* Multiplayer
* Spectator
* Chat
* Twitch Song Requests
* & any other real-time events

**This application is being developed for internal use. As such, no support or proper documentation will be provided for the usage of this software.**

## Requirements

* MariaDB / MySQL
* Redis
* Steam API Key
* Steam Publisher API Key

## Setup

1. Install `Go 1.22` or later.
2. Clone the repository.
3. Copy `config.example.json` and make a file named `config.json`
4. Fill out the config file with the appropriate details.
5. If you do not have a **Steam Publisher Account** (you are not a Quaver developer), you can set `bypass_steam_login` to `true` in the config file. This should **NOT** be used in a production environment.
6. Navigate to the `/cmd/server/` directory
7. Start the server with `go run .` or your method of choice.
8. The server is now available at `ws://localhost:3000`.

## LICENSE

This software is licensed under the GNU Affero General Public License v3.0. Please see the LICENSE file for more information.
