package notifier

type Notifier interface {
	Success(msg string) error
	Failure(msg string) error
	Progress(msg string) error
}
