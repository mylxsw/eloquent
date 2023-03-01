package event

import (
	"reflect"
)

// Listener is a event listener
type Listener interface{}

// EventStore is a interface for event store
type EventStore interface {
	Listen(eventName string, listener Listener)
	Publish(eventName string, evt interface{})
	SetManager(*EventManager)
}

// EventManager is a manager for event dispatch
type EventManager struct {
	store EventStore
}

// NewEventManager create a eventManager
func NewEventManager(store EventStore) *EventManager {
	manager := &EventManager{
		store: store,
	}

	store.SetManager(manager)

	return manager
}

// Listen create a relation from event to listners
func (em *EventManager) Listen(listeners ...Listener) {
	for _, listener := range listeners {
		listenerType := reflect.TypeOf(listener)
		if listenerType.Kind() != reflect.Func {
			panic("listener must be a function")
		}

		if listenerType.NumIn() != 1 {
			panic("listener must be a function with only one arguemnt")
		}

		if listenerType.In(0).Kind() != reflect.Struct {
			panic("listener must be a function with only on argument of type struct")
		}

		em.store.Listen(listenerType.In(0).String(), listener)
	}
}

// Publish a event
func (em *EventManager) Publish(evt interface{}) {
	em.store.Publish(reflect.TypeOf(evt).String(), evt)
}

// Call trigger listener to execute
func (em *EventManager) Call(evt interface{}, listener Listener) {
	reflect.ValueOf(listener).Call([]reflect.Value{reflect.ValueOf(evt)})
}

// MemoryEventStore is a event store for sync operations
type MemoryEventStore struct {
	listeners map[string][]Listener
	manager   *EventManager
}

// NewMemoryEventStore create a sync event store
func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		listeners: make(map[string][]Listener),
	}
}

// Listen add a listener to a event
func (eventStore *MemoryEventStore) Listen(evtType string, listener Listener) {
	if _, ok := eventStore.listeners[evtType]; !ok {
		eventStore.listeners[evtType] = make([]Listener, 0)
	}

	eventStore.listeners[evtType] = append(eventStore.listeners[evtType], listener)
}

// Publish publish a event
func (eventStore *MemoryEventStore) Publish(evtType string, evt interface{}) {
	if listeners, ok := eventStore.listeners[evtType]; ok {
		for _, listener := range listeners {
			eventStore.manager.Call(evt, listener)
		}
	}
}

// SetManager event manager
func (eventStore *MemoryEventStore) SetManager(manager *EventManager) {
	eventStore.manager = manager
}
