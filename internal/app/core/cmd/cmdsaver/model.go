package cmdsaver

type model struct {
	quitBroadcast chan<- struct{}
}

func New(quitBroadcast chan<- struct{}) model {
	return model{
		quitBroadcast: quitBroadcast,
	}
}
