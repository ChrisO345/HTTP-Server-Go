# HTTP Server in GO

Allows for the creation of a simple HTTP server that can serve files and accept POST requests to write to files.

Made following the CodeCrafter's HTTP Server challenge.

Note that this is a volatile server, and will not save files between runs.

### To Use (on Windows)

1. Run the Server File
   ```bash
   go run app\server.go
   ```
2. Use `curl` to interact with the server
   ```bash
   curl -d "Hello World" localhost:4221/files/hello.txt
   ```
3. Use `curl` to read the file
   ```bash
   curl localhost:4221/files/hello.txt
   ```