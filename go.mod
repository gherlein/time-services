module time-services

go 1.22.3

replace github.com/gherlein/time-services/time-service-pb => ./timeservice-pb

require (
	github.com/gherlein/time-services/time-service-pb v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.64.1
)

require (
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
