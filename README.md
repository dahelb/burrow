# burrow ― A Gopher Server in Go

This is a learner project that implements the [Gopher protocol](https://datatracker.ietf.org/doc/html/rfc1436) in Go.

## Usage

Starting the server will serve a directory (by default the current one).

$ burrow

```
Usage of ./burrow:
  -bind string
    	bind to this address (default "127.0.0.1")
  -host string
    	Hostname for gopher menus (default "127.0.0.1")
  -port int
    	The port to listen on. (default 7000)
  -root string
    	The root directory to serve (default ".")
```

Assuming that there is a folder `gopherhole` with the following content:

```
gopherhole/
├── about.txt
└── foo
    └── bar.txt
```

Start the server using `burrow -root gopherhole` will listen on port 7000.

Now, connect using a gopher client like lynx or [ncgopher](https://github.com/jansc/ncgopher).

