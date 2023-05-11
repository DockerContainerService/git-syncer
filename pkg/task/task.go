package task

type Task interface {
	Run() error
	GetTitle() string
	GetSrcRepo() string
	GetDstRepo() string
	SetError(err error)
	GetError() error
}
