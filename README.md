# Homemade Dropbox
Synchronizes the file system in a client directory with a server directory over TCP. 

## How to run the client and server
1. Start the server with the server directory as the argument:
	```bash
	$ go run main.go /path/to/server_directory
	```

2. Start the client with the client directory as the argument:
	```bash
	$ go run main.go /path/to/client_directory 
	```
