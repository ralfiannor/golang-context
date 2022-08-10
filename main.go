package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Millisecond*2000))
	defer cancel()

	timeNow := time.Now()

	err := callWithoutTimeoutParent(ctx)
	if err != nil {
		fmt.Println("callWithoutTimeoutParent() error ", err.Error())
	}

	err = callWithTimeoutParent(ctx)
	if err != nil {
		fmt.Println("callWithTimeoutParent() error ", err.Error())
	}
	fmt.Println("=== Stopped Application ===")
	fmt.Println("=== Duration: ", time.Since(timeNow))
}

// callWithoutTimeoutParent will call the HTTP Request and create new context if parent context have a Deadline
func callWithoutTimeoutParent(ctx context.Context) (err error) {
	req, err := http.NewRequest(http.MethodGet, "http://deelay.me/2000/http://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		return
	}

	if _, isHaveDeadline := ctx.Deadline(); isHaveDeadline {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*9000)
	defer cancel()

	req = req.WithContext(ctx)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println(string(out))
	return
}

// callWithTimeoutParent will call the HTTP Request and use existing context
// if we set context.WithTimeout() in child func is not affected because ctx parent already set it
func callWithTimeoutParent(ctx context.Context) (err error) {
	req, err := http.NewRequest(http.MethodGet, "http://deelay.me/50/http://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		return
	}

	// try to create new deadline with 10 seconds
	// but it fail because parent have existing deadline
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Millisecond*10000))
	defer cancel()

	req = req.WithContext(ctx)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println(string(out))
	return
}
