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
