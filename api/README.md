# API layer

This directory and file layouts are subject to change drastically as we've began to add functionality besides the websocket connection.
A practical example here being the voice messages that are and should not be delivered thru the WS but posted as a reqular HTTP POST message formdata.

## Endpoints

### GET
- /server - Gets all the server info. Because of the low amount of data this bundles all current info
```json
{"name":"server1","type":"MASTER","url":"127.0.0.1:9393","voiceUrl":"127.0.0.1:9393","rooms":[{"name":"main","topic":"Main chatroom","users":["test_user1"]},{"name":"testing","topic":"Testing grounds","users":[]},{"name":"test2","topic":"Testing grounds #2","users":[]}],"users":[{"userId":"f13df928-10e8-459d-ba7e-e16164b4afed","nick":"test_user1","server":"127.0.0.1:9393","rooms":["main"]}]}
```

