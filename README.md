Run the server with:

`go run server.go --address=0.0.0.0:8090`

Run the client with:

`go run client.go --sender=<set your username> --receiver=<username of the person whom you want to chat with> --server-address=<server-ip:port>`

Ex.
So if your username you want to set as Whitebox and you know that your friend's username is Blackbox:
`go run client.go --sender=Whitebox --receiver=Blackbox --server-address=35.228.52.141:30542`

And your friend have to start her client as:
`go run client.go --sender=Blackbox --receiver=Whitebox --server-address=35.228.52.141:30542`

Run the client with docker:
docker build -t simplechat-server -f Dockerfile .
`docker run -d -p 10001:10001 simplechat-server`
