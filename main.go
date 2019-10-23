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

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "print incomplete images only",
		},
		cli.BoolFlag{
			Name:  "sniffing, s",
			Usage: "determine the image type of the first 32 bytes of data",
		},
	}
	app.Name = "image-checksum"
	app.Usage = "find out incomplete images in paths"
	app.UsageText = app.Name + " [paths...]"

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.ShowAppHelp(c)
		}

		ctx, cancel := context.WithCancel(context.Background())
		go interrupt(cancel)

		for _, f := range c.GlobalFlagNames() {
			ctx = context.WithValue(ctx, f, c.GlobalBool(f))
		}
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
			console(ctx, skip, i)
			continue
		}

		img := i
		eg.Go(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			ra, err := image.NewReaderAt(img)
			if err != nil {
				return err
			}
			defer ra.Close()

			if ctx.Value("sniffing") == true {
				checker, err = image.Sniff(ra)
				if err != nil {
					return fmt.Errorf("%s: %w", img, err)
				}
				if checker == nil {
					console(ctx, skip, i)
					return nil
				}
			}

			result, err := checker.Check(ra)
			if err != nil {
				return fmt.Errorf("%s: %w", img, err)
			}

			if result {
				console(ctx, success, img)
			} else {
				console(ctx, failure, img)
				fmt.Println(img)
			}
			return nil
		})
	}
	return eg.Wait()
}

func console(ctx context.Context, a ...interface{}) {
	if ctx.Value("quiet") == true {
		return
	}
	fmt.Fprintln(os.Stderr, a...)
}
