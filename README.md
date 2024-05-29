# Learning by Doing:  Simple gRPC in Go (and using ChatGPT)

This is an example repo just for me to learn how to build gPRC programs in go.  I specifically wanted to learn certain things:

* how to generate go code from .proto files
* how to properly include that code in a working go client and server
* how to import local go modules
* how to use ChatGPT to help craft the protobufs

## TLDR and Conclusion

ChatGPT won't be taking any golang programmers job anytime soon.  It made several key mistakes (discussed below).  Using the code it generated actually required modern knowledge of the go ecosystem and language to make it work.

However, in a few hours of wall-clock time I have a working demo and I can expand it to explore the specifics of using protobufs.  

## This Repo

This example is a "time server" that accepts a request for the time and replies with the current time.  

## Background

It's been a few years since I wrote anything in golang and a lot has changed.  Specifically the $GOPATH is no longer used.  So I had a few things to refresh on.

Part of what made my learning difficult is that a lot of the guidance pages assume knowledge that I did not have.  I missed the transition from $GOPATH to modules so [pages](https://protobuf.dev/getting-started/gotutorial/) that describe how to do things didn't make a lot of sense to me.  Specifically, one cannot use relative paths in the import block anymore.  So I HAD to learn the new way.  And learning is always good.

Caveat:  my web searches didn't find easy answers for me, but you may be much better and find better pages right off the bat.  I'm mostly writing this so that I know that I actually understand it, and so that if it comes to pass that I don't get to do much golang for awhile again that I can refresh my memory.  Of course, if I wait too long then it may all change again anyway.  Oh well, welcome to the Internet.

## Basic Setup

The very first thing I had to re-learn was how the modern go modules work.  I particularly wanted to understand how to use "local modules."  I wanted to develop modules somewhat in a monorepo style, where I have a cli client and a server in the same repo.  I wanted to have a common protobuf definition that both used.  I wanted to autogenerate go code from the protobuf and then use that go code in both the client and server.  

This [article](https://brokencode.io/how-to-use-local-go-modules-with-golang-with-examples/) helped me the most.  The key things you have to do:

* the root folder can have a shared go.mod for each of the discrete executables, but it needs a "replace" section (see below)
* anything that will compile to be a discrete executable needs to be in a folder of it's own and declare "package main"
* any shared code needs to be in a folder of it's own and it needs a go.mod file in that folder
    * that go.mod file serves to name the package that golang will treat that code as
    * I went ahead and used a fully qualified github path name to preserve the option later of splitting it out to a discrete, non-monorepo package
* the root folder go.mod file is where the magic happens

## ChatGPT

I asked ChatGPT to do this:

```text
write a server and client program in golang that uses grpc.  The server should expose an API that allows a user to request the time.
```

See below for what it spit out.


## gPRC Stuff

I split out the protofuf file into a folder "time-services-pb" and defined new targets in my Makefile to generate the go code:

```
TSPB := ./timeservice-pb


proto-clean:
        -rm ${TSPB}/*.go
        ls ${TSPB}



proto: proto-clean
        @echo "creating go files from .proto files"
        protoc -I${TSPB} --go_out=${TSPB} --go-grpc_out=${TSPB} --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ${TSPB}/timeservice.proto 
        ls ${TSPB}
```

The ChatGPT solution (below) was not wrong, but it would not work.  And putting the go files into the same folder sent me down the wild goose chase of finding out how to declare the package so that the generated code would be included.

The [Protobuif Docs page](https://protobuf.dev/reference/go/go-generated/) was the place I found the information letting me specify exactly what command line I needed.  Note that I preferred the "declaring it within the .proto file" approach.

## Solution

### Main Code "import" Section

I split out the client and server code into their own folders since they both are part of package "main" and they both have a main() function.  This means that they both cannot be in the same folder and do a clean build.  They both will inherit the go.mod from the parent folder.

###  Module "go.mod" File

Since this example added a local module, I had to add a go.mod file into the folder holding the proto definition (timeservice-pb):


```
module github.com/gherlein/time-services/time-service-pb
```

This is the defined module name.  Note that it's just a string.  It does happen to map to the folder in the git repo, but it being a repo is not at all important.  What is important is that this is the name of the module.

### Main "go.mod" File

You are no longer allowed to directly put relative paths in go.mod files.  But, using the "replace" section you can accomplish the same thing.  This allows you to have the module included in the main repo instead of having it broken out into another repo.  Of course, if you wanted to you could keep it in another repo and include it here as a sub-module.

But there's a chicken before the egg problem here.  First step is:

```bash
go mod init time-services
go: creating new go.mod: module time-services
go: to add module requirements and sums:
	go mod tidy
```

But if you immediately follow the instructions they fail:

```bash
go mod tidy
go: finding module for package google.golang.org/grpc
go: finding module for package github.com/gherlein/time-services/time-service-pb
go: found google.golang.org/grpc in google.golang.org/grpc v1.64.0
go: finding module for package github.com/gherlein/time-services/time-service-pb
go: time-services/client imports
	github.com/gherlein/time-services/time-service-pb: cannot find module providing package github.com/gherlein/time-services/time-service-pb: module github.com/gherlein/time-services/time-service-pb: git ls-remote -q origin in /home/gherlein/go/pkg/mod/cache/vcs/ed1fcb598d0c41f0acbd4fda33252d6e5c60f87cf6ab9a5030f796f3fd8220a2: exit status 128:
	fatal: could not read Username for 'https://github.com': terminal prompts disabled
Confirm the import path was entered correctly.
If this is a private repository, see https://golang.org/doc/faq#git_https for additional information.
```

The problem is that our module is local, not in a repo.  So we need to add a "replace" section:

```
module time-service

go 1.22.3

replace (
        github.com/gherlein/time-services/time-service-pb  => ./timeservice-pb
        
)
```

This is the magic.  This tells golang that the module is local, not a separate go module.  NOW you can do "go mod tidy" like usual.

```bash
go mod tidy
go: finding module for package google.golang.org/grpc
go: found github.com/gherlein/time-services/time-service-pb in github.com/gherlein/time-services/time-service-pb v0.0.0-00010101000000-000000000000
go: found google.golang.org/grpc in google.golang.org/grpc v1.64.0
```

and the new go.mod file:

```
module time-services

go 1.22.3

replace github.com/gherlein/time-services/time-service-pb => ./timeservice-pb

require (
	github.com/gherlein/time-services/time-service-pb v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.64.0
)

require (
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
```

### Output

Run the server:

```bash
make run-server
```

And the client in another terminal:

```bash
make run-client
```

It works!


## What ChatGPT Said

ChatGPT left out some important things.  Specifically:

* putting the generated go code into it's own module.  It said "Be sure to replace "path/to/your/package/timepb" with the actual path to your generated timepb package." But that is not enough information at all to make things work.

* putting things into a local module and using the replace directive to properly include it.

* it made a mistake in the client code:  it used res.GetCurrent_time when the function is actually res.GetCurrentTime()

It's clear to me that ChatGPT won't be taking a golang programmers job anytime soon.


### Define the gRPC Service - time.proto

```
syntax = "proto3";

package timepb;

service TimeService {
  rpc GetTime (TimeRequest) returns (TimeResponse);
}

message TimeRequest {}

message TimeResponse {
  string current_time = 1;
}
```

### Generate the gRPC Code 

```
protoc --go_out=. --go-grpc_out=. time.proto
```

### Implement the Server - server.go

```
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"

    "google.golang.org/grpc"
    "path/to/your/package/timepb"
)

type server struct {
    timepb.UnimplementedTimeServiceServer
}

func (s *server) GetTime(ctx context.Context, req *timepb.TimeRequest) (*timepb.TimeResponse, error) {
    currentTime := time.Now().Format(time.RFC3339)
    res := &timepb.TimeResponse{
        Current_time: currentTime,
    }
    return res, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
    timepb.RegisterTimeServiceServer(s, &server{})

    fmt.Println("Server is running on port 50051...")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

###  Implement the Client - client.go

```
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
    "path/to/your/package/timepb"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := timepb.NewTimeServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    res, err := client.GetTime(ctx, &timepb.TimeRequest{})
    if err != nil {
        log.Fatalf("Error calling GetTime: %v", err)
    }

    fmt.Printf("Current Time: %s\n", res.GetCurrent_time())
}
```

## Next Steps

I really want to explore debugging gRPC applications.  Here's some things I plan next:

* [grpc-tools](https://github.com/bradleyjkemp/grpc-tools) - Like Fiddler/Charles but for gRPC!
* [Analyzing gRPC messages using Wireshark](https://grpc.io/blog/wireshark/)
* [grpcdebug](https://github.com/grpc-ecosystem/grpcdebug) - fetches the internal states of the gRPC library from the application via gRPC protocol 

I'll do more examples of each as I work my way through.

