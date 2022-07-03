package batch

import (
	"context"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

var appendMutex sync.Mutex

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	errG, _ := errgroup.WithContext(context.Background())
	errG.SetLimit(int(pool))

	var i int64
	for i = 0; i < n; i++ {
		func(i int64) {
			errG.Go(func() error {
				user := getOne(i)

				appendMutex.Lock()
				defer appendMutex.Unlock()
				res = append(res, user)

				return nil
			})
		}(i)
	}

	err := errG.Wait()

	if err != nil {
		return nil
	}

	return res
}
