# pop

A point-of-presence server

## Usage

This app has a sample origin server in `server.js`. Start this server by running

```bash
node server.js
```

Then start the POP server by running

```bash
go run main.go
```

The point-of-presence (pop) server acts as a caching reverse-proxy for the origin server. Pop receives a request first, if a corresponding cache object is found, it is returned immediately. Else, a proxy request is sent to the origin server, the response received is cached and forwarded to the client.

This pop server can be deployed at multiple locations and their IPs can be added to the DNS so that requests are routed to these pop servers first.

The cache is a simple in-memory implementation for now.
