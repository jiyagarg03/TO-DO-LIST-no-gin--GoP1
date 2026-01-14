# Go Todo List API(net/http)

A simple **Todo List REST API** built using **Go (net/http)**.  
This project demonstrates core backend concepts like routing, JSON handling, concurrency, and mutex usage.
No framework(gin) used. **Each line** of code includes a comment explaining what it does, making the project easy to understand for both me and the reader.

---

## Features

- Create a todo
- Get all todos
- Update a todo (mark as done)
- Delete a todo
- In-memory storage
- Thread-safe using `sync.Mutex`
- JSON based REST API

---

## Tech Stack

- Go
- net/http
- encoding/json
- sync.Mutex
