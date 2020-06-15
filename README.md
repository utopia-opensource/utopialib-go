
![logo](https://github.com/Sagleft/utopialib-go/raw/master/logo.png)

Utopia Ecosystem API wrapper written in Golang

Usage
-----

```go
client := UtopiaClient{
	protocol: "http",
	token:    "C17BF2E95821A6B545DC9A193CBB750B",
	host:     "127.0.0.1",
	port:     22791,
}

fmt.Println(client.GetSystemInfo())
```
