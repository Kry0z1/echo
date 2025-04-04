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
	Priority    int    `json:"priority"` // [0, 3]
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
	StartsAt    int    `json:"starts_at"`
	DueTo       int    `json:"due_to"`
	Done        bool   `json:"done"`
}

type TaskIn TaskBase

type TaskOut struct {
	TaskBase
	ID int `json:"id"`
}

type TaskStored struct {
	TaskOut
}

func (TaskStored) TableName() string {
	return "tasks"
}
