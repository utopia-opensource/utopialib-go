
![logo](https://github.com/Sagleft/utopialib-go/raw/master/logo.png)

Utopia Ecosystem API wrapper written in Golang

Docs
-----

[![GoDoc](https://godoc.org/github.com/sagleft/utopialib-go?status.svg)](https://godoc.org/gopkg.in/sagleft/utopialib-go.v1)
[![go-report](https://goreportcard.com/badge/github.com/Sagleft/utopialib-go)](https://goreportcard.com/report/github.com/Sagleft/utopialib-go)
[![Build Status](https://travis-ci.org/sagleft/utopialib-go.svg?branch=master)](https://travis-ci.org/sagleft/utopialib-go)

Install
-----

```bash
go get gopkg.in/sagleft/utopialib-go.v2
```

or

```go
import "gopkg.in/sagleft/utopialib-go.v2"
```

Usage
-----

```go
client := utopiago.UtopiaClient{
	protocol: "http",
	token:    "C17BF2E95821A6B545DC9A193CBB750B",
	host:     "127.0.0.1",
	port:     22791,
}

fmt.Println(client.GetSystemInfo())
```

How can this be used?
-----

* creating a web service that processes client requests;
* creation of a payment service;
* development of a bot for the channel;
* utility for working with uNS;
* experiments to explore web3.0;
