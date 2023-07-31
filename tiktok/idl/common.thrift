namespace go common

struct User {
  1: i64 id; // user id
  2: string name; // user name
  3: i64 follow_count; // total number of people the user follows
  4: i64 follower_count; // total number of fans
  5: bool is_follow; // whether the currently logged-in user follows this user
  6: string avatar; // user avatar URL
  7: string background_image; // image at the top of the user's personal page
  8: string signature; // user profile
  9: i64 total_favorited; // number of likes for videos published by user
  10: i64 work_count; // number of videos published by user
  11: i64 favorite_count; // number of likes by this user
}