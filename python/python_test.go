package python

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"github.com/chuwt/zing/object"
	"github.com/chuwt/zing/python/lib"
)

func TestPython(t *testing.T) {
	pe := NewPyEngine(
		"/Volumes/hdd1000gb/workspace/src/github.com/chuwt/zing/python/github.com/chuwt/zing/strategies",
		"/Volumes/hdd1000gb/workspace/src/github.com/chuwt/zing/python/github.com/chuwt/zing")
	if err := pe.Init(); err != nil {
		fmt.Println(err.Error())
		return
	}

	wg := sync.WaitGroup{}
	strategies := make([]*lib.PyObject, 0)
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			strategy := pe.NewStrategyInstance(
				"TestStrategy",
				object.StrategyId(i),
				object.VtSymbol{
					GatewayName: object.Gateway(fmt.Sprintf("huobi%d", i)),
					Symbol:      "btcusdt",
				},
				"")
			strategies = append(strategies, strategy)
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("init ok")

	for i := 0; i < 2; i++ {
		swg := sync.WaitGroup{}
		for _, s := range strategies {
			swg.Add(1)
			s := s
			go func() {
				pe.ObjectCallFunc(s, "test", "2")
				swg.Done()
			}()
		}
		swg.Wait()
		<-time.After(time.Second)
	}
	fmt.Println("run ok")
	gil := lib.PyGILState_Ensure()
	o := strategies[0].GetAttrString("count")
	fmt.Println(lib.PyUnicode_AsUTF8(o.Repr()))
	lib.PyGILState_Release(gil)
	_ = pe.Close()
}
