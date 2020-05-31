package boot

var (
	Data chan interface{}//数据通道
	//Ticker *time.Ticker//定时器
)

/*const (
	MEMINFO = "/proc/meminfo"
	GAPTIME = time.Microsecond * 100)*/

func init() {
	Data = make(chan interface{}, 10)
	//run()
	//Ticker = time.NewTicker(GAPTIME)
}

//func Run(){
//	p := &modules.Pack{}
//	mem := modules.Module(&mem.Stats{})
//	p.SetModule(&mem)
//	//go p.Run()
//	p.Run()
//}