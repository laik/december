package main

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/laik/december"
)

type _userResuestStreamBenchmark struct {
	user   string
	result string
}

type _streamsBenchmark chan *_userResuestStreamBenchmark

func (self _streamsBenchmark) Handle(event Event) {
	for key, params := range event.Params {
		self <- &_userResuestStreamBenchmark{key, fmt.Sprintf("xdsa request result data %s", params)}
	}
}

func (self _streamsBenchmark) request(_route string, _args Parameters, dispatcher *Dispatcher) {
	dispatcher.DispatchEvent(&Event{_route, _args})
}

func (self _streamsBenchmark) response() string { return (<-self).result }

var dispatcher *Dispatcher = SharedDispatcher()

var __streams _streamsBenchmark = make(chan *_userResuestStreamBenchmark, 0)

var _route string = "/access/redisMonitor"

func monitorFlush(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	dispatcher.AddEventListener(_route, __streams)

	__streams.request(_route, map[string]interface{}{"test_user_1": r.Form}, dispatcher)

	_str := (<-__streams).result

	log.Printf("%s", _str)

	w.Write([]byte(_str))
}

func main() {

	http.HandleFunc("/", monitorFlush) //设置访问的路由

	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	//http://127.0.0.1:9090/?xx=88&ss=00&wocao=3838438
}
