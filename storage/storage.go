package storage

type Storage interface {
	Flush() error
	Init() error

	UserAdd(User) (int, error)
	UserGetData(string, int) (User, error)
	UserUpdateActivity(int) error
	UserGetParam(string, string) (string, error)
}
