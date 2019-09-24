package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"

	"github.com/bilibili/kratos/pkg/sync/errgroup"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "image-checksum"
	app.Usage = "检查给定路径内所有图片的完整性"
	app.UsageText = app.Name + " [roots...]"

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.ShowAppHelp(c)
		}

		ctx, cancel := context.WithCancel(context.Background())
		go interrupt(cancel)
		return start(ctx, c.Args())
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func interrupt(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
	time.Sleep(time.Second)
}

func start(ctx context.Context, roots []string) error {
	eg := errgroup.WithContext(ctx)

	eg.GOMAXPROCS(runtime.NumCPU())
	for _, v := range roots {
		root := v
		eg.Go(func(ctx context.Context) error {
			return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					fmt.Println(info.Name())
				}
				return nil
			})
		})
	}

	return eg.Wait()
}
