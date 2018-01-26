package december

import "testing"
import "fmt"
import "os"

func TestEvent(t *testing.T) {

	dispatcher := SharedDispatcher()

	args := make(Parameters, 0)

	args["start"] = "lualu"
	args["stop"] = "leipale"

	event1 := &Event{"/access/dafeiji", args}

	_cb := func(e Event) {
		for key, value := range e.Params {
			switch key {
			case "start":
				fmt.Fprintf(os.Stdout, "i receive start ..........tututu %s\n", value)
			case "stop":
				fmt.Fprintf(os.Stdout, "i receive stop ..........dongdongdong %s\n", value)
			}
		}
	}

	cb := &callback{&_cb}

	dispatcher.AddEventListener(event1.String(), cb)

	dispatcher.DispatchEvent(event1)

}

type _userResuestStream struct {
	user   int
	result string
}

type _streams chan *_userResuestStream

func (self _streams) Handle(event Event) {
	for _, params := range event.Params {
		self <- &_userResuestStream{params.(int), fmt.Sprintf("receive user id %d request result data %d", params, params)}
	}
}

func (self _streams) request(_route string, _args Parameters, dispatcher *Dispatcher) {
	dispatcher.DispatchEvent(&Event{_route, _args})
}

func (self _streams) response(t *testing.T) {
	for {
		select {
		case x, ok := <-self:
			t.Logf("%s", x.result)
			if !ok {
				break
			}
		}
	}
}

func TestReactorStream(t *testing.T) {

	var __streams _streams = make(chan *_userResuestStream, 0)

	var _route string = "/access/papapa"

	dispatcher := SharedDispatcher()

	dispatcher.AddEventListener(_route, __streams)

	go __streams.response(t)

	for i := 0; i < 10; i++ {
		__streams.request(_route, map[string]interface{}{"test": i}, dispatcher)
	}
}
