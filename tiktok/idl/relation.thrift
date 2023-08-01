namespace go relation

include "common.thrift"

struct RelationActionRequest {
    1: required i64 current_user_id         // 用户鉴权token
    2: required i64 to_user_id       // 对方用户id
    3: required i32 action_type = 3; // 1-关注 2-取消关注
}

struct RelationActionResponse {
    1: required i32 status_code      // 状态码 0-成功 other-失败
    2: required string status_msg    // 状态描述
}

struct RelationFollowListRequest {
    1: required i64 user_id
    2: required i64 current_user_id
}

struct RelationFollowListResponse {
    1: required i32 status_code
    2: required string status_msg
    3: list<common.User> user_list   // 用户列表
}

struct RelationFollowerListRequest {
    1: required i64 user_id
    2: required i64 current_user_id
}

struct RelationFollowerListResponse {
    1: required i32 status_code
    2: optional string status_msg
    3: list<common.User> user_list
}

struct RelationFriendListRequest {
    1: i64 user_id;
    2: i64 current_user_id;
}

struct RelationFriendListResponse {
    1: i32 status_code
    2: required string status_msg
    3: list<FriendUser> user_list   // 用户列表
}

struct FriendUser {
    1: optional string message
    2: required i64 msgType
}

service RelationService {
    RelationActionResponse RelationAction(1: RelationActionRequest req)
    RelationFollowListResponse RelationFollowList(1: RelationFollowListRequest req)
    RelationFollowerListResponse RelationFollowerList(1: RelationFollowerListRequest req)
    RelationFriendListResponse RelationFriendList(1: RelationFriendListRequest req)
}