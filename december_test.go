package december

import "testing"
import "time"
import "log"

func EventTest(t *testing.T) {

	params := make(map[string]interface{}, 0)
	params["ID"] = "dsadsadsa"
	tmpBuffer := make([]string, 10)

	main := CreateEvent("main", params)

	var rawParse EventHandler = func(event *Event) {
		ch := make(chan *Event)
		go func() { ch <- event }()
		for {
			select {
			case tmpEvent, ok := <-ch:
				if ok {
					tmpBuffer = append(tmpBuffer, tmpEvent.Params["Data"].(string))
				}
			}
		}
	}

	main.AddEventDependency(CreateEvent("raw", params), &rawParse)

	if eventList, ok := main.EventDependency["raw"]; ok {
		for _, event := range eventList {
			t.Log("add event dependency sucess , event name %s", &(<-event).Name)
		}
	}

	main.DeleteEventDependency(CreateEvent("raw", params))
	if _, ok := main.EventDependency["raw"]; !ok {
		t.Log("delete event dependency sucess")
	}

	t.Log("sucesses")
	t.Logf("%v", tmpBuffer)

	var logger EventHandler = func(event *Event) {
		ch := make(chan *Event)
		ch <- event
		log.Println(<-ch)
	}
	dispatcher := NewDispatcher()
	dispatcher.AddListener("main", &logger)
	dispatcher.DispatchEvent(main)
	time.Sleep(1 * time.Second)
}
