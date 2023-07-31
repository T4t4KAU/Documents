namespace go api

struct AddRequest {
    1: i32 first
    2: i32 second
}

struct AddResponse {
    1: i32 sum
}

service AddService {
    AddResponse Add(AddRequest req)
}