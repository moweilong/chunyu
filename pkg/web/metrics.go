package web

import (
	"expvar"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moweilong/chunyu/pkg/queue"
)

// 您可能想了解:
// 1. 应用程序使用了多少内存? 使用率是如何随着时间变化的?
// 2. 目前有多少个 Goroutine 正在使用?
// 3. 有多少个数据库连接正在使用中，有多少个处于空闲状态?
// 4. HTTP 响应成功和错误的比率是多少?
// 深入了解以上内容有助于把控程序，并得到预警。

// Metrics ...
func Metrics() gin.HandlerFunc {
	request := expvar.NewInt("request")
	totalRequests := expvar.NewInt("requests")
	totalResponses := expvar.NewInt("responses")
	urls := expvar.NewMap("requestURLs")
	statusCodes := expvar.NewMap("statusCodes")
	// 协程数量
	// go func() {
	// queue.NewCirQueue[expvar.Int](10)
	// queue.CirQueue
	// }()

	return func(c *gin.Context) {
		totalRequests.Add(1)
		request.Add(1)
		c.Next()
		request.Add(-1)
		totalResponses.Add(1)

		status := c.Writer.Status()
		if status != 404 {
			urls.Add(c.Request.Method+" "+c.FullPath(), 1)
		}
		statusCodes.Add(strconv.Itoa(status), 1)
	}
}

type GoroutineNum struct {
	Time string `json:"time"`
	Num  int    `json:"num"`
}

// CountGoroutines 协程数量，间隔 duration 记录一次
func CountGoroutines(d time.Duration, num uint8) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	goroutine := queue.NewCirQueue[GoroutineNum](num)

	expvar.Publish("goroutine_num", expvar.Func(func() any {
		return goroutine.Range()
	}))

	for {
		goroutine.Push(GoroutineNum{
			Time: time.Now().Format(time.DateTime),
			Num:  runtime.NumGoroutine(),
		})
		<-ticker.C
	}
}
