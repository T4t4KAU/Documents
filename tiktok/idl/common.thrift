namespace go common

struct User {
  1: required i64 id; // user id
  2: required string name; // user name
  3: optional i64 follow_count; // total number of people the user follows
  4: optional i64 follower_count; // total number of fans
  5: required bool is_follow; // whether the currently logged-in user follows this user
  6: optional string avatar; // user avatar URL
  7: optional string background_image; // image at the top of the user's personal page
  8: optional string signature; // user profile
  9: optional i64 total_favorited; // number of likes for videos published by user
  10: optional i64 work_count; // number of videos published by user
  11: optional i64 favorite_count; // number of likes by this user
}