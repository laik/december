package december

import "testing"
import "fmt"
import "math/rand"

var dispatcher *Dispatcher = SharedDispatcher()

func BenchmarkReactor(t *testing.B) {

	var _cb func(e Event) = func(e Event) {
		for key, _ := range e.Params {
			switch key {
			case "start":
				//			fmt.Fprintf(os.Stdout, "i receive start ..........tututu %v\n", value)
			case "stop":
				//			fmt.Fprintf(os.Stdout, "i receive stop ..........dongdongdong %v\n", value)
			}
		}
	}

	var cb *callback = &callback{&_cb}

	dispatcher.AddEventListener("/access/dafeiji1", cb)

	go func() {
		for i := 0; i < t.N; i++ {
			args := make(Parameters, 0)
			args["start"] = i
			args["stop"] = i
			dispatcher.DispatchEvent(&Event{"/access/dafeiji1", args})
		}
	}()
}

type _userResuestStreamBenchmark struct {
	user   int
	result string
}

type _streamsBenchmark chan *_userResuestStreamBenchmark

func (self _streamsBenchmark) Handle(event Event) {
	for _, params := range event.Params {
		self <- &_userResuestStreamBenchmark{params.(int), fmt.Sprintf("receive user id %d request result data %d", params, params)}
	}
}

func (self _streamsBenchmark) request(_route string, _args Parameters, dispatcher *Dispatcher) {
	dispatcher.DispatchEvent(&Event{_route, _args})
}

func (self _streamsBenchmark) response(t *testing.B) {
	for {
		select {
		case x, ok := <-self:
			if ok {
				t.Logf("%s", x.result)
			}
			if !ok {
				break
			}
		}
	}
}

func BenchmarkStream(t *testing.B) {

	var __streams _streamsBenchmark = make(chan *_userResuestStreamBenchmark, 0)

	var _route string = "/access/redisMonitor"

	dispatcher := SharedDispatcher()

	dispatcher.AddEventListener(_route, __streams)

	go __streams.response(t)

	__streams.request(_route, map[string]interface{}{"test": rand.Int()}, dispatcher)

}
