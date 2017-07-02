# Golang Kafka Example
A simple example to how to use Kafka with Go (golang)

# Dependencies
- Docker
- Docker Compose

# How to test

- Start the app: `make run`
- Make requests

    - Send sync message: ```curl -i -X POST -H 'Content-Type: application/json' -d '{"key": "User", "value": "user@email.com"}' http://localhost/producer/sync```
    - Send async message: ```curl -i -X POST -H 'Content-Type: application/json' -d '{"key": "User", "value": "user@email.com"}' http://localhost/producer/async```
