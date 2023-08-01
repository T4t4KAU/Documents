namespace go user

include "common.thrift"

struct UserRegisterRequest {
    1: required string username
    2: required string password
}

struct UserRegisterResponse {
    1: required i32 status_code
    2: required string status_msg
    3: required i64 user_id
    4: required string token
}

struct UserLoginRequest {
    1: required string username
    2: required string password

}

struct UserLoginResponse {
    1: required i32 status_code
    2: required string status_msg
    3: required i64 user_id
    4: required string token
}

service UserService {
    UserRegisterResponse UserRegister(1: UserRegisterRequest req)
    UserLoginResponse UserLogin(1: UserLoginRequest req)
}