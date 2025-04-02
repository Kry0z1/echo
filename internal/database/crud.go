package database

func GetUserByUsername(username string) (UserStored, error) {
	var user UserStored
	err := db.First(&user, "username=?", username).Error
	return user, err
}
