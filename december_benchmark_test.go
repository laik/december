package december

import "testing"

//import "fmt"
//import "os"

var dispatcher *Dispatcher = SharedDispatcher()

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

func BenchmarkReactor(t *testing.B) {
	dispatcher.AddEventListener("/access/dafeiji1", cb)

	for i := 0; i < t.N; i++ {
		args := make(Parameters, 0)
		args["start"] = i
		args["stop"] = i
		dispatcher.DispatchEvent(&Event{"/access/dafeiji1", args})
	}
}
