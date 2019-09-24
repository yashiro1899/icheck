package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"

	"imagechecksum/image"

	"github.com/bilibili/kratos/pkg/sync/errgroup"
	"github.com/urfave/cli"
)

const (
	success = "✔"
	failure = "✘"
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
	eg := errgroup.WithContext(ctx)
	eg.GOMAXPROCS(runtime.NumCPU())

	images := make(chan string)
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
	return check(ctx, images)
}

func check(ctx context.Context, images <-chan string) error {
	eg := errgroup.WithCancel(ctx)
	eg.GOMAXPROCS(runtime.NumCPU())

	for i := range images {
		checker := image.Get(path.Ext(i))

		if checker == nil {
			fmt.Fprintf(os.Stderr, "No checker for %q\n", i)
			continue
		}

		img := i
		eg.Go(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			file, err := os.Open(img)
			if err != nil {
				return err
			}
			defer file.Close()

			err = checker.Check(&ReaderAt{file})
			fmt.Println(err)
			return nil
		})
	}
	return eg.Wait()
}
