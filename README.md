# prometheus
application implements "Hello world!" web-server on port 9000  
 It has two paths: 
 + / (root) - displays "Hello world!" after a random time interval (0...10 sec)
 + /metrics - displays metrics collected by prometheus exporter ("github.com/prometheus/client_golang/prometheus")

user's metrics are:
 + http_response_time_seconds - histogram that shows the amount of requests served during specified time interval
 + http_simplehandler_counter - counts the amount of requests

Both: the web application and Prometheus server - are docerized.  
Prometheus UI is available on http://localhost:9090 
