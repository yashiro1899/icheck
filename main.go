package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

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
}

func start(ctx context.Context, roots []string) error {
	images := make(chan string)
	eg := errgroup.WithContext(ctx)
	eg.GOMAXPROCS(runtime.NumCPU())

	for _, v := range roots {
		root := v
		eg.Go(func(ctx context.Context) error {
			return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.Mode().IsRegular() {
					return nil
				}
				select {
				case images <- path:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			})
		})
	}

	go func() {
		eg.Wait()
		close(images)
	}()

	for i := range images {
		fmt.Println(i)
	}

	return nil
}
