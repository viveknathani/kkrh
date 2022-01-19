# kkrh

kkrh is a web application for tracking my activities. I use it in my daily life to understand my habits better and get insights out of them. Hence the name is an abbreviation of the Hindi sentence "kya kar raha hai".

## why not just use something that already exists?

- I want control over my data. The logs stored in tracking apps describe a lot about the users.
- I get to choose and find the statistics I care about and not be limited to what an application provides.
- It is fun to spin up your own web application.

## tech

- Go
- PostgreSQL
- Redis

## source description

- [cache](./cache): provides an API to work with the cache 
- [database](./database): provides an API to work with the database 
- [entity](./entity): provides the `struct` to represent the entities in the system 
- [repository](./repository): provides the interface implemented by `database` package. 
- [scripts](./scripts): for things like setting up PostgreSQL one time
- [server](./server): provides the HTTP API with some handy middleware
- [service](./service): the core logic of this system, to be used by `server`
- [shared](./shared): for things like managing `context` across `server` and `service` 
- [web](./web): the frontend. 

## what is the framework for frontend?

I chose to write it in vanilla javascript this time. Something like React was not a necessity here and writing in vanilla javascript is fun. 

## kkrh-meta

`kkrh` outputs some logs to `stdout` and `stderr`. I wanted to capture and process them which is why I wrote [kkrh-meta](https://github.com/viveknathani/kkrh-meta/).

## cannot create an account?

- I have disabled the signup functionality for now. This service exists for my own usage with limited use of cloud resources. But you are free to host your own version of kkrh.

## host your own version

- clone this repository.
- create a new app in Heroku with the PostgreSQL addon.
- setup a new Redis server from [here](https://app.redislabs.com/).
- setup the environment variables in Heroku as per the [.env.example](./.env.example) file.
- push the repository to heroku.
- have fun.

## contributing

- `make build`, `make test`, and `make run` will be some handy commands for you. 
- It would be cool if you write an issue before working on a PR. 

## license

[MIT](./LICENSE)

