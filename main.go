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
	success = "\033[32m✔\033[0m"
	failure = "\033[31m✘\033[0m"
	skip    = "\033[33m⚠\033[0m"
)

var quiet bool

func console(a ...interface{}) {
	if quiet {
		return
	}
	fmt.Fprintln(os.Stderr, a...)
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "print incomplete only",
		},
	}
	app.Name = "image-checksum"
	app.Usage = "检查给定路径内所有图片的完整性"
	app.UsageText = app.Name + " [roots...]"

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.ShowAppHelp(c)
		}

		quiet = c.Bool("quiet")
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
			console(skip, i)
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
			if err != nil {
				console(failure, img)
				if err == image.Incomplete {
					fmt.Println(img)
					return nil
				}
				return fmt.Errorf("%s: %w", img, err)
			}
			console(success, img)
			return nil
		})
	}
	return eg.Wait()
}
