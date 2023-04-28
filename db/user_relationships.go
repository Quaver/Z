package db

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
