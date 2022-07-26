# Models

Abstractions for different types of data that the backend uses for DAO and DTO purposes.

- Room - A chatroom with name and topic
- User - Chat user with id, nick, servery they're on and list of roomnames
- Config - Configuration of the server, rather short list at the moment
- ServerData - dynamic information about the server, will become more relevant when multiple servers are implmented.
- Msg - DTO message for both the NATS and Websocket communication
