package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	// var c = context.Background()
	// 新建一个fanout 对象 名称为cache 名称主要用来上报监控和打日志使用 最好不要重复
	// (可选参数) worker数量为1 表示后台只有1个线程在工作
	// (可选参数) buffer 为1024 表示缓存chan长度为1024 如果chan慢了 再调用Do方法就会报错 设定长度主要为了防止OOM
	// cache := fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024))
	// // 需要异步执行的方法
	// // 这里传进来的c里面的meta信息会被复制 超时会忽略 addCache拿到的context已经没有超时信息了
	// // cache.Do(c, func(c context.Context) { addCache(c, 0, 0) })
	// // 程序结束的时候关闭fanout 会等待后台线程完成后返回
	// cache.Close()
	// cache.Close()

	// fmt.Println(cache)

	app := cli.NewApp()
	app.Name = "image-checksum"
	app.Usage = "检查给定路径内所有图片的完整性"
	app.UsageText = app.Name + " [paths...]"

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.ShowAppHelp(c)
		}
		fmt.Printf("Hello %q\n", c.Args().Get(0))
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
