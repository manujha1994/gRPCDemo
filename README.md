# gRPCDemo

This project is intended to give us a basic idea how we can implement gRPC APIs using GoLang. 
There are 4 ways in which we can implement API in gRPC. They are as following :
- Unary API
- Server side streaming API
- Client side Streaming API
- Bi Directional Streaming API

In this Demo I have all the 4 APIs using calculation of sum of numbers between a given range. The aim was to keep the example as simple as possible for purpose of understanding.

## gRPC Advantages

- gRPC uses HTTP/2 under the covers
- Single protocol buffer file can be shared by various language
- In comparison to JSON, size of protocol buffer is less
- gRPC uses a binary payload that is efficient to create and to parse, and it exploits HTTP/2 for efficient management of connections
- gRPC provide set of responses based on the type of result/error
- By proper planning architecture, the number of API hits from client to server can be reduced via streaming.
