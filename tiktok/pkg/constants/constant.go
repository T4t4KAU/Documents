package constants

const (
	MySQLDSN      = "root:123456@tcp(127.0.0.1:3306)/tiktok?charset=utf8&parseTime=true"
	RedisAddr     = "127.0.0.1:6379"
	RedisPassword = ""
	EtcdAddress   = "127.0.0.1:2379"

	SecretKey   = "tiktok"
	IdentityKey = "user_id"
)

const (
	UserTableName     = "users"
	RelationTableName = "relations"
)

const (
	UserServiceName     = "user"
	RelationServiceName = "relation"
	ApiServiceName      = "api"
)

const (
	TestAva        = "test.jpg"
	TestBackground = "test.jpg"
)
