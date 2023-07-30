package redis

import "strconv"

const (
	likeSuffix  = ":like"
	likedSuffix = ":liked"
)

type (
	Favorite struct{}
)

func (f Favorite) AddLike(userId, videoId int64) {
	add(rdbFavorite, strconv.FormatInt(userId, 10)+likeSuffix, videoId)
}

func (f Favorite) AddLiked(userId, videoId int64) {
	add(rdbFavorite, strconv.FormatInt(videoId, 10)+likedSuffix, userId)
}

func (f Favorite) DelLike(userId, videoId int64) {
	del(rdbFavorite, strconv.FormatInt(userId, 10)+likeSuffix, videoId)
}

func (f Favorite) DelLiked(userId, videoId int64) {
	del(rdbFavorite, strconv.FormatInt(videoId, 10)+likedSuffix, userId)
}

func (f Favorite) CheckLike(userId int64) bool {
	return check(rdbFavorite, strconv.FormatInt(userId, 10)+likeSuffix)
}

func (f Favorite) CheckLiked(videoId int64) bool {
	return check(rdbFavorite, strconv.FormatInt(videoId, 10)+likedSuffix)
}

func (f Favorite) ExistLike(userId, videoId int64) bool {
	return exist(rdbFavorite, strconv.FormatInt(userId, 10)+likeSuffix, videoId)
}

func (f Favorite) ExistLiked(userId, videoId int64) bool {
	return exist(rdbFavorite, strconv.FormatInt(videoId, 10)+likedSuffix, userId)
}

func (f Favorite) CountLike(userId int64) (int64, error) {
	return count(rdbFavorite, strconv.FormatInt(userId, 10)+likeSuffix)
}

func (f Favorite) CountLiked(videoId int64) (int64, error) {
	return count(rdbFavorite, strconv.FormatInt(videoId, 10)+likedSuffix)
}

func (f Favorite) GetLike(userId int64) []int64 {
	return get(rdbFavorite, strconv.FormatInt(userId, 10)+likeSuffix)
}

func (f Favorite) GetLiked(videoId int64) []int64 {
	return get(rdbFavorite, strconv.FormatInt(videoId, 10)+likedSuffix)
}
