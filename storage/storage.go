package storage

type Storage interface {
	Flush() error
	Init() error

	UserAdd(User) (int, error)
	UserGetData(string) (User, error)
}
