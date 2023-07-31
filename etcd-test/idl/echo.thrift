namespace go api

struct EchoRequest {
    1: string message
}

struct EchoResponse {
    2: string message
}

service EchoService {
    EchoResponse Echo(1:EchoRequest req)
}