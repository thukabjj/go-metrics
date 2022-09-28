# Go with Gin Framework and OpenTelemetry 
The goal of this project is show a Simple Example of how configure the Distributed Tracing Using OpenTelemetry library for Golang and also using one of the most popular Golang Framework Gin Gonic.

## How to run
```
    docker compose up --force-recreate --build 
```

## Payloads
### service A:
```curl
   curl --location --request GET 'http://localhost:8080/ping' \
    --header 'X-B3-TraceId: 463ac35c9f6413ad48485a3953bb6124' \
    --header 'X-B3-SpanId: a2fb4a1d1a96d312' \
    --header 'X-B3-ParentSpanId: 0020000000000001' \
    --header 'X-B3-Sampled: 1'
```
### service B:
```curl
    curl --location --request GET 'http://localhost:8081/pong' \
    --header 'X-B3-TraceId: 463ac35c9f6413ad48485a3953bb6124' \
    --header 'X-B3-SpanId: a2fb4a1d1a96d312' \
    --header 'X-B3-ParentSpanId: 0020000000000001' \
    --header 'X-B3-Sampled: 1'
```

### You can also use the with [ApacheBench](https://httpd.apache.org/docs/2.4/programs/ab.html) in Linux or macOS.
```
ab -n 30 -c 10 http://localhost:8080/ping    
```

## References
 - [Zipkin - b3-propagation](https://github.com/openzipkin/b3-propagation)
 - [Open Telemetry](https://opentelemetry.io/docs/collector/getting-started/)
 - [jaeger-go-example](https://github.com/albertteoh/jaeger-go-example)
 - [Implementing OpenTelemetry and Jaeger tracing in Golang HTTP API](http://www.inanzzz.com/index.php/post/4qes/implementing-opentelemetry-and-jaeger-tracing-in-golang-http-api)
 - [OpenTelemetry Go: All you need to know](https://lightstep.com/blog/opentelemetry-go-all-you-need-to-know)