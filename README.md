# GoProxy


[GoDoc](https://github.com/autom8ter/goproxy/blob/master/GODOC.md)

    go get github.com/autom8ter/goproxy/...


## Overview
GoProxy is a lightweight reverse proxy server written in Golang

It registers target urls and appends basic authentication, headers, form values, and more to the inbound request so that you may regulate and modify requests to 
many different API endpoints from a single gateway.

run:
    goproxy

output:

```text
              (                         
 (            )\ )                      
 )\ )        (()/( (            )  (    
(()/(     (   /(_)))(    (   ( /(  )\ ) 
 /(_))_   )\ (_)) (()\   )\  )\())(()/( 
(_)) __| ((_)| _ \ ((_) ((_)((_)\  )(_))
  | (_ |/ _ \|  _/| '_|/ _ \\ \ / | || |
   \___|\___/|_|  |_|  \___//_\_\  \_, |
                                   |__/

Current Config:
map[]

Usage:
  GoProxy [command]

Available Commands:
  help        Help about any command
  serve       start the GoProxy server

Flags:
  -a, --addr string     address to run server on (default ":8080")
  -c, --config string   relative path to file containing proxy configuration (default "config.yaml")
  -h, --help            help for GoProxy
      --version         version for GoProxy

Use "GoProxy [command] --help" for more information about a command.

```

run:

    goproxy serve

output:

    starting GoProxy server: :8080



## Example Config File

**COMING SOON**


## TODO

- [ ] Add JWT/OAuth2 based authentication