# ChatApp

This is a chatting app, which consists of React on frontend and Golang on the backend.

## Used tech
- Go
- React
- Gin Framework
- GORM
- Websockets with gorilla package
- AWS S3 Storage

## Implemented functionality
- Sending messages via websocket connection
- Authentication using jwt
- REST API for creating, updating and deleting groups
- Adding and removing members from group
- Changing members rights in a group (to delete and add members)
- Updating websocket connection when adding or deleting from group
- Handling messages from other groups on frontend
- Profile pictures and pictures for groups stored in AWS S3 bucket
