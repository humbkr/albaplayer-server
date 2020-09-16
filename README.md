# Alba Player (server)
> Minimalistic audio library web player.

## About
Alba player is an audio library web-based player optimized for desktop. I was tired of all the desktop audio players available on linux 
and macos with questionable user interfaces and/or library management, and I wanted to learn golang and react, so I decided
to build my own.   
This player is by purpose very limited in terms of functionalities, as its main goal is to just play music.  

This repository is only for the server part of the application. You can find the client part at [https://github.com/humbkr/albaplayer-client](https://github.com/humbkr/albaplayer-client)   
Although fully operationnal, this project is still under heavy development. Contributions and remarks on the existing 
code are welcomed (new react dev here)!

WEBSITE (with demo): [https://albaplayer.com](https://albaplayer.com)

### Basic features

- HTML5 audio player with basic controls (play / pause / previous / next / timeline progress / volume / random / repeat)
- Main playing queue à la Winamp (playback doesn't magically change to the songs you are browsing in the library)
- Library browser with Artists / Albums / Tracks views that actually manage compilations properly
- Now playing screen with current song info and buttons to google lyrics and guitar tabs
- Client / server app, so can be installed on a server to access a music library remotely
- Can manage huge libraries (tested with 30000+ songs)

**Note:** this player is not adapted for mobile or tablet use. A good mobile UI would be completely different from the
desktop one, so I focused on the desktop first, as there are already a lot of good mobile players app.

## Installation

Grab the archive corresponding to your system on the [official website](https://albaplayer.com), unzip it somewhere, tinker with the alba.yml
configuration file and run the alba executable from the command line.

## Developement

**Tech stack:**
- Golang 1.9
- GraphQL API
- SQLite

**Dependencies:**   

This project uses go modules.

**Code organization:**   

This project mostly follows [https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout) for its repository structure.   
In terms of go packages structure, it follows the concepts of Clean Architecture: 
[http://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/](http://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/)

#### Docker

A docker image is provided for developement purposes, if you don't want to install the dev stack on your machine.   
(Note that you will unfortunately still have to install Golang if you want IDEs like JetBrains' Goland to work 
properly.)

##### Prerequisites
- docker set up on your machine
- ``make`` command available (windows users)

##### Set up
From the project root, cd into /docker then run ``make up``  
Once the container is mounted, log into it by running ``make ssh``
From inside the container, install the dependencies by running ```go dep ensure``` from the project's root

##### Use
- To start the application in watch mode, from inside the container run ``fresh`` and wait for the process to finish.
- To access the application in a browser, get the container port from docker: ``docker ps`` (image name is "alba_server"), then go to your browser and
access localhost:<port>.

Available endpoints:
- /graphql : graphql server
- /graphiql : graphiql client for testing (when dev mode is enabled only, see alba.yml)

Note that you need to build the [client app](https://github.com/humbkr/albaplayer-client) separately to access the user interface. By running only the server part
in this repository you will only have access to the two endpoints mentionned earlier.   

#### Include the client app in the build
To build the client app, follow the readme at [https://github.com/humbkr/albaplayer-client](https://github.com/humbkr/albaplayer-client).   
Once the client prod build is generated, dump its contents in the /web directory of this project.

Alternatively you can also get the end-user build on the [official website](https://albaplayer.com) and copy-paste the web directory contents into /web.

#### Test
From the project root run ``go test ./...``
