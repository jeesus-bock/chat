# API layer

This directory and file layouts are subject to change drastically as we've began to add functionality besides the websocket connection.
A practical example here being the voice messages that are and should not be delivered thru the WS but posted as a reqular HTTP POST message formdata.

## Endpoints

### GET
- /server - Gets all the server, room and user info. Because of the low amount of data this bundles all current data in the response.
```json
{
  "name":"server1",
  "type":"MASTER",
  "url":"127.0.0.1:9393",
  "voiceUrl":"127.0.0.1:9393",
  "rooms":[{"name":"main","topic":"Main chatroom","users":["test_user1"]},
           {"name":"testing","topic":"Testing grounds","users":[]},
           {"name":"test2","topic":"Testing grounds #2","users":[]}
  ],
  "users":[{"userId":"f13df928-10e8-459d-ba7e-e16164b4afed",
            "nick":"test_user1",
            "server":"127.0.0.1:9393",
            "rooms":["main"]}
  ]
}
```
- /uploads/{room name}/{filename} - Serves static uploaded files or a status 404 if the file doesn't exist. Filenames are a combination of user nick and unix timestamp.
- /rooms - Returns the internal map with room names as key and topic as value.
- / - All other paths lead to WS upgrading.
- /ws/{user id}/{room name} - Offers ws:// in localhost and wss:// on deployment. At the moment the user id is generated clientside and room is created if it doesn't already exist.
Websocket message example:
```json
{"type":"msg","from":"be71c8d4-e356-4ddd-827d-5b545656a034","to":"main","msg":"Hello!","ts":1659097919309}
```

  ### POST
  - /upload/{user id}/{room name} - Uploads a file as field named "document" in FormData. At the moment this is only implemented for audio files (Recoded voice messages), but ought to support other kinds of files, images for example.
