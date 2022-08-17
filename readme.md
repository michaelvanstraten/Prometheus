# Prometheus

I started writing this web framework for the future online shop of my school. This was my first real go project and i really enjoined the typing experience of the language, at first. I would run into problems here and there which i could overlook or hack around, something i couldn't work around was the support for Generics. This is why i abandoned this project and switched to rust-lang.

## Development Process

I started looking into NodeJS frameworks like Express and NestJS but i found them to be to abstracted away, i couldn't really see what was going on under the hood and i wanted something were i knew what was going on. So i decided to write my own.

I chose go because i recently tried it for the first time and i really liked the syntax. Further more i heard that google was using it to parse html documents for there search engine and i wanted something fast so i figured, let it be go.

### The Router
I started by writing my own Router seen here:

```go
type Router struct {
	Routes 			map[string]*Service
	DefaultRoute 	string
}
```
it consisted of Routes which where mapped to a Pointer to Service. This way you could have multiple routes pointing to one service, for example "/hello" -> *HelloWorldService, "/hello-world" -> *HelloWorldService and the service would respone if the url path was prefix with the path of the route so "/hello/some-random-string" would still be mapped to the service HelloWorldService. The default route field was, if set, written out to the response writer if the requested route could not be found in the map. By this point it was a pretty primitive router but i wanted to add the option to chain services together. 