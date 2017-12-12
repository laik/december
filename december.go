// Copyright 2017 The laik Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations

// under the License.

package december

import "unsafe"
import "fmt"

type EventHandler func(*Event)

type Event struct {
	Name            string
	Params          map[string]interface{}
	EventDependency map[string][]chan *Event
	Handler         map[string][]*EventHandler
}

func (e *Event) AddEventDependency(event *Event, handler *EventHandler) {
	c := make(chan *Event)
	c <- event
	e.EventDependency[e.Name] = append(e.EventDependency[e.Name], c)
	e.Handler[e.Name] = append(e.Handler[e.Name], handler)
}

func (e *Event) DeleteEventDependency(event *Event) {
	for name, _ := range e.EventDependency {
		if event.Name == name {
			delete(e.EventDependency, name)
			fmt.Printf("delete event name %s master event name %s", event.Name, e.Name)
		}
	}
}

func (e *Event) DeleteEventDependencyByName(event string) {
	for name, _ := range e.EventDependency {
		if event == name {
			delete(e.EventDependency, name)
			fmt.Printf("delete event name %s master event name %s", event, e.Name)
		}
	}
}

func CreateEvent(name string, params map[string]interface{}) *Event {
	return &Event{Name: name, Params: params, EventDependency: make(map[string][]chan *Event, 0), Handler: make(map[string][]*EventHandler, 0)}
}

type Dispatcher struct {
	Listeners map[string]*Event
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{Listeners: make(map[string]*Event, 0)}
}

func (d *Dispatcher) AddListener(eventName string, handler *EventHandler) {
	event, ok := d.Listeners[eventName]
	if !ok {
		d.Listeners[eventName] = CreateEvent(eventName, make(map[string]interface{}))
	}

	for _, item := range event.Handler[eventName] {
		if *(*int)(unsafe.Pointer(item)) == *(*int)(unsafe.Pointer(handler)) {
			return
		}
	}

	ch := make(chan *Event)
	event.EventDependency[eventName] = append(event.EventDependency[eventName][:], ch)
	event.Handler[eventName] = append(event.Handler[eventName][:], handler)
	go d.depthHandler(eventName, ch, handler)
}

func (d *Dispatcher) depthHandler(eventName string, ch chan *Event, handler *EventHandler) {
	for {
		event, ok := <-ch
		if !ok || event == nil {
			break
		}
		go (*handler)(event)
	}
}

func (d *Dispatcher) DispatchEvent(event *Event) {
	tmp, ok := d.Listeners[event.Name]
	if ok {
		for _, tmpEvent := range tmp.EventDependency[event.Name] {
			tmpEvent <- event
		}
	}
}
