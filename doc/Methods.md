## Handlers

#### Registration
* Request
```
curl '\n' -X POST 'localhost:1323/v1/register'  -d '{"name": "vika", "password":"qwerty","email":"qweqwe@mail.ru"}'
```
* Response
```
{"code":"200","message":"You are registered. Welcome vika"}
```
#### Login
* Request
```
curl '\n' -X POST 'localhost:1323/v1/login'  -d '{"name": "pasha", "password":"qwerty","email":"qweqwe@mail.ru"}'
```
* Response
```
{"code":"200","message":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicGFzaGEiLCJwYXNzd29yZCI6InF3ZXJ0eSIsImV4cCI6MTUwNjEwNzM5N30.MtvT1t6DYOXJNgi9IQBGmRExxy7XOm3XMqdGnWdtLi0"}
```
#### Ping
* Request
```
curl localhost:1323/v1/ping -H "Authorization: Bearer eyJuYW1lIjoicGFzaGEiLCJwYXNzd29yZCI6InF3ZXJ0eSIsImV4cCI6MTUwNjEwNzQ3N30.-jUuyH9jfbLCJgUXVlxPijvhoWFXh6vfxTS33HnJsaw"
```
* Response
```
{"code":"200","message":"ECHO_PING"}
```
#### PingDb
* Request
```
curl localhost:1323/v1/ping_db -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicGFzaGEiLCJwYXNzd29yZCI6InF3ZXJ0eSIsImV4cCI6MTUwNjEwNzM5N30.MtvT1t6DYOXJNgi9IQBGmRExxy7XOm3XMqdGnWdtLi0"
```
* Response
```
{"code":"200","message":"1"}
```

