package database

import "github.com/Kry0z1/echo/internal"

func CreateUser(user *UserStored) error {
	return db.Create(user).Error
}

func CreateUserWithUsernameAndPassword(user UserIn) (UserStored, error) {
	hashedPassword, err := internal.GetPasswordHash(user.Password)
	if err != nil {
		return UserStored{}, err
	}

	userDB := UserStored{
		UserOut{user.UserBase, 0},
		hashedPassword,
	}

	err = CreateUser(&userDB)

	if err != nil {
		return UserStored{}, err
	}

	return userDB, nil
}

func CreateTask(task TaskIn) (TaskStored, error) {
	taskStored := TaskStored{TaskOut{TaskBase(task), 0}}

	err := db.Create(&taskStored).Error

	return taskStored, err
}

func GetUser(id int) (UserStored, error) {
	var user UserStored
	user.ID = id
	err := db.First(&user).Error
	return user, err
}

func GetUserByUsername(username string) (UserStored, error) {
	var user UserStored
	err := db.First(&user, "username=?", username).Error
	return user, err
}
