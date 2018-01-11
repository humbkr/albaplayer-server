# Alba Player


## Installation


Features


### Developement

BACKEND  
Golang version: 1.9  
Dependencies:  
This project uses golang/dep to manage its dependencies. Download it (https://github.com/golang/dep) and and it to your path.

#### Docker
##### Prerequisites
- docker set up on your machine
- ``make`` command available

##### Set up
From the project root, run ``make up``  
You can then log into the container by running ``make ssh``

To start the application in watch mode, inside the container run ``fresh`` and wait for the process to finish.

To access the application in a browser, get the container port from docker: ``docker ps``, then go to your browser and
access localhost:<port>.

Available endpoints:
- /graphql : graphql server
- /graphiql : graphiql client for testing






--------------------------------------------------------------------------
FRONT  
Set up IDE
disable safe write in intelliJ settings so the hot reloading can work.
Install React Developer Tools for FF or Chrome.


Tool to get a graphql schema from an endpoint:
npm install -g get-graphql-schema

get-graphql-schema http://localhost:8888/graphql > ./schema.graphql
Then remove all the ' @deprecated' occurences



TODO
- Create an icon component to mask the fontello stuff
