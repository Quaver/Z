package sessions

import (
	"example.com/Quaver/Z/db"
	"strconv"
)

// UpdateRedisOnlineUserCount Updates the online user count in Redis
func UpdateRedisOnlineUserCount() error {
	_, err := db.Redis.Set(db.RedisCtx, "quaver:server:online_users", GetOnlineUserCount(), 0).Result()

	if err != nil {
		return err
	}

	return nil
}

// Adds a user's session token to redis
func addUserTokenToRedis(user *User) error {
	_, err := db.Redis.Set(db.RedisCtx, user.getRedisSessionKey(), strconv.Itoa(user.Info.Id), 0).Result()

	if err != nil {
		return err
	}

	return nil
}

// Removes a user's session token from redis
func removeUserTokenFromRedis(user *User) error {
	_, err := db.Redis.Del(db.RedisCtx, user.getRedisSessionKey()).Result()

	if err != nil {
		return err
	}

	return nil
}
