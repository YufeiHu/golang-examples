package main

import (
	"context"
	"fmt"
	"time"
)

func task1(pctx context.Context, pCancelFunc context.CancelFunc) {
	go func() {
		i := 0
		for {
			fmt.Println(i)
			time.Sleep(1*time.Second)
			i++
		}
	}()

	select {
	case <-pctx.Done():
	}

	fmt.Println("Done")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	go func() {
		time.Sleep(3*time.Second)
		cancel()
	}()

	task1(ctx, cancel)
}
