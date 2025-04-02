package database

type UserBase struct {
	Username string       `json:"username"`
	Tasks    []TaskStored `json:"tasks" gorm:"foreignKey:UserID"`
}

type UserIn struct {
	UserBase
	Password string `json:"password"`
}

type UserOut struct {
	UserBase
	ID int
}

type UserStored struct {
	UserOut
	HashedPassword string
}

func (UserStored) TableName() string {
	return "users"
}

type TaskBase struct {
	Title       string `json:"title"`
	Priority    int    `json:"priority"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
	StartsAt    int    `json:"starts_at"`
	DueTo       int    `json:"due_to"`
}

type TaskIn TaskBase

type TaskOut struct {
	TaskBase
	ID int
}

type TaskStored TaskOut

func (TaskStored) TableName() string {
	return "tasks"
}
