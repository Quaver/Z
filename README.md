# Z

New Quaver realtime multiplayer & chat server.

## Setup

* Install Go version 1.19.
* Copy `config.example.json` and make a file named `config.json`.
* Fill out the config with the appropriate details.
  * If you do not have a Steam Publisher account, you can set `bypass_steam_login` to `true`. This should NOT be used in production, however.
* Start the server with `go run .` or build and run the executable.