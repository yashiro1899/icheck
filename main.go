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

	"github.com/go-kratos/kratos/pkg/sync/errgroup"
	"github.com/urfave/cli"
)

type Message struct {
	Error string
	Out   string
}

const (
	success = "\033[32m✔\033[0m"
	failure = "\033[31m✘\033[0m"
	skip    = "\033[33m⚠\033[0m"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "show more",
		},
	}
	app.Name = "image-checksum"
	app.Usage = "find out incomplete images in paths"
	app.UsageText = app.Name + " [paths...]"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.ShowAppHelp(c)
		}

		ctx, cancel := context.WithCancel(context.Background())
		go interrupt(cancel)

		for _, fn := range c.GlobalFlagNames() {
			ctx = context.WithValue(ctx, fn, c.GlobalBool(fn))
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

	images := make(chan string, 32)
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

	lines := make(chan Message, 32)
	go console(ctx, lines)

	for i := range images {
		m := Message{}
		m.Error = failure + " "

		if chk := image.Get(path.Ext(i)); chk == nil {
			m.Error = fmt.Sprintf("%s %s\n", skip, i)
			lines <- m
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

			checker, err := image.Sniff(ra)
			if err != nil {
				return fmt.Errorf("%s: %w", img, err)
			}
			if checker == nil {
				m.Out = img
				lines <- m
				return nil
			}

			result, err := checker.Check(ra)
			if err != nil {
				return fmt.Errorf("%s: %w", img, err)
			}
			if result {
				m.Error = fmt.Sprintf("%s %s\n", success, img)
			} else {
				m.Out = img
			}
			lines <- m
			return nil
		})
	}

	err := eg.Wait()
	close(lines)
	return err
}

func console(ctx context.Context, lines <-chan Message) {
	verbose := ctx.Value("verbose").(bool)
	for m := range lines {
		if !verbose {
			continue
		}
		fmt.Fprint(os.Stderr, m.Error)
		if m.Out != "" {
			fmt.Println(m.Out)
		}
	}
}
