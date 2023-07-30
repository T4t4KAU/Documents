package constants

const (
	MySQLDSN      = "root:123456@tcp(127.0.0.1:3306)/douyin?charset=utf8&parseTime=true"
	RedisAddr     = "127.0.0.1:6379"
	RedisPassword = ""
)

const (
	UserTableName      = "users"
	FollowsTableName   = "follows"
	VideosTableName    = "videos"
	MessageTableName   = "messages"
	FavoritesTableName = "favorites"
	CommentTableName   = "comments"

	VideoFeedCount       = 30
	FavoriteActionType   = 1
	UnFavoriteActionType = 2
)

const (
	TestAva        = "test.jpg"
	TestBackground = "back.jpg"
)

const (
	MinioEndPoint        = "127.0.0.1:18001"
	MinioAccessKeyId     = "minio"
	MinioSecretAccessKey = "12345678"
	MinioUseSSL          = false

	MinioVideoBucketName = "videobucket"
	MinioImageBucketName = "imagebucket"
)
