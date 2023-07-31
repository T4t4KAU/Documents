namespace go user

include "common.thrift"

struct UserRegisterRequest {
    1: string username
    2: string password
}

struct UserRegisterResponse {
    1: i32 status_code
    2: string status_msg
    3: i64 user_id
    4: string token
}

struct UserLoginRequest {
    1: string username
    2: string password

}

struct UserLoginResponse {
    1: i32 status_code
    2: string status_msg
    3: i64 user_id
    4: string token
}

service UserService {
    UserRegisterResponse UserRegister(1: UserRegisterRequest req)
    UserLoginResponse UserLogin(1: UserLoginRequest req)
}