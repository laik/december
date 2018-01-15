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

// Inspiron by events
type Parameters map[string]interface{}

// Event define
type Event struct {
	Name   string
	Params Parameters
}

func (self *Event) String() string { return self.Name }

// Listener defines event handler interface
type Listener interface {
	Handle(Event)
}

// Stream implements Listener interface on channel
type Stream chan Event

// Handle Listener
func (stream Stream) Handle(event Event) { stream <- event }

// Callback implements Listener interface on function
func Callback(function func(Event)) Listener { return callback{&function} }

// Callback Prototype
type callback struct{ function *func(Event) }

// Handle Listener
func (callback callback) Handle(event Event) { (*callback.function)(event) }

// Event has one or multi listener
type Listeners struct {
	events    []chan *Event
	listeners []Listener
}

func createListeners() *Listeners {
	return &Listeners{events: []chan *Event{}, listeners: []Listener{}}
}

// Shard instance global singleton
var _instance *Dispatcher

func SharedDispatcher() *Dispatcher {
	if _instance == nil {
		_instance = &Dispatcher{}
		_instance.Init()
	}
	return _instance
}

// Dispatcher stores event(Pointer) and one or multi subscribe
type Dispatcher struct{ listeners map[string]*Listeners }

// Initial Dispatcher
func (self *Dispatcher) Init() { self.listeners = make(map[string]*Listeners, 0) }

// Add *Event & on or multi Listener(implement Handle Method) &struct{}
func (self *Dispatcher) AddEventListener(event string, listener Listener) {
	_, ok := self.listeners[event]
	if !ok {
		self.listeners[event] = createListeners()
	}

	//..
	for _, item := range self.listeners[event].listeners {
		a := *(*int)(unsafe.Pointer(&item))
		b := *(*int)(unsafe.Pointer(&listener))
		if a == b {
			return
		}
	}

	//  Make a channel recevie outer event
	_ch := make(chan *Event, 0)
	self.listeners[event].events = append(self.listeners[event].events[:], _ch)
	//  Add new listener to listeners
	self.listeners[event].listeners = append(self.listeners[event].listeners[:], listener)
	//  Async call back
	go self._handler(event, _ch, listener)
}

func (self *Dispatcher) _handler(event string, ch chan *Event, listener Listener) {
	for {
		event := <-ch
		if event == nil {
			break
		}
		go listener.Handle(*event)
	}
}

func (self *Dispatcher) RemoveEventListener(event string, listener Listener) {
	lsncrt, ok := self.listeners[event]
	if !ok {
		return
	}

	var (
		ch    chan *Event
		key   int = 0
		exist bool
	)

	for k, item := range self.listeners[event].listeners {
		a := *(*int)(unsafe.Pointer(&item))
		b := *(*int)(unsafe.Pointer(&listener))
		if a == b {
			exist = true
			ch = lsncrt.events[k]
			key = k
			break
		}
	}

	if exist {
		ch <- nil
		lsncrt.events = append(lsncrt.events[:key], lsncrt.events[key+1:]...)
		lsncrt.listeners = append(lsncrt.listeners[:key], lsncrt.listeners[key+1:]...)
	}
}

func (self *Dispatcher) DispatchEvent(event *Event) {
	events, ok := self.listeners[event.Name]
	if ok {
		for _, chEvent := range events.events {
			chEvent <- event
		}
	}
}
