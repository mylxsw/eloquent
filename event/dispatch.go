package event

// Dispatcher is an interface for event dispatcher
type Dispatcher interface {
	Publish(evt interface{})
}

var dispatcher Dispatcher

// SetDispatcher set a event dispatcher
func SetDispatcher(disp Dispatcher) {
	dispatcher = disp
}

// Dispatch an event to Dispatcher
func Dispatch(evt interface{}) {
	if dispatcher == nil {
		return
	}

	dispatcher.Publish(evt)
}
