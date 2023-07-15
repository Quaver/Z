package db

import "database/sql"

type UserRelationship struct {
	Id           int `db:"id"`
	UserId       int `db:"user_id"`
	TargetUserId int `db:"target_user_id"`
	Relationship int `db:"relationship"`
}

// GetUserFriendsList Retrieves a slice of user ids that a given user is friends with
func GetUserFriendsList(userId int) ([]int, error) {
	const query string = "SELECT target_user_id FROM user_relationships WHERE user_id = ? AND (relationship & 1) != 0"

	relationships := make([]int, 0)

	err := SQL.Select(&relationships, query, userId)

	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// GetUserRelationship Gets a relationship with a user
func GetUserRelationship(userId int, targetUserId int) (*UserRelationship, error) {
	const query string = "SELECT * FROM user_relationships WHERE user_id = ? AND target_user_id = ? LIMIT 1"

	var relationship UserRelationship
	err := SQL.Get(&relationship, query, userId, targetUserId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &relationship, nil
}

// AddFriend Adds a player to a user's friends list
func AddFriend(userId int, targetUserId int) error {
	const query string = "INSERT INTO user_relationships (user_id, target_user_id, relationship) VALUES (?, ?, ?)"

	_, err := SQL.Exec(query, userId, targetUserId, 1)

	if err != nil {
		return err
	}

	return nil
}

// RemoveFriend Remvoes a player from a user's friends list
func RemoveFriend(userId int, targetUserId int) error {
	const query string = "DELETE FROM user_relationships WHERE user_id = ? AND target_user_id = ?"

	_, err := SQL.Exec(query, userId, targetUserId)

	if err != nil {
		return err
	}

	return nil
}
