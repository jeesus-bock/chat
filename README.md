# Chat

A proof-of-concept implementation of a Gofiber websocket and HTTP APIs that use NATS on the backend.

The work is very much in progress but the basics of both NATS and websocket are in working state.

API ought to be defined better, among other things but as it is the emphasis of the project is to
have a fun simple project to work on and test all sorts of things. Basically everything is replaceable,
nothing's written in stone besides the fact that I'm not going work hard on anything.

Pullrequests and feature requests/issues and whatnot are **more than welcome** and can be written in english
or finnish, don't be shy.

# Running (on linux)
- Git clone repo
- go get
- go build
- Start NATS server, more info https://github.com/jeesus-bock/natikka
- ./chat
