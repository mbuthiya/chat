## Chat application

A simple chat application that uses web sockets to communicate with the server. 


## How to run 

To the application just clone the repo and run the following in the terminal

**Build Executable**
```bash
go build -o chat
```

**Start server**
```bash
./chat -addr=":8080"
```

The *addr* flag allows you to specify which port you want to run the application