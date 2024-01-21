package internal

type Tracker interface {
	Setup() error
	Start(description string) (string, error)
	Stop() (string, error)
}
