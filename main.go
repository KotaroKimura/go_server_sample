package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"golang.org/x/sync/errgroup"
	_ "github.com/go-sql-driver/mysql"

	"github.com/KotaroKimura/server_sample/handlers"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	session, err := newDBSession()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	eg.Go(func() error {
			return runServer(ctx, session)
	})
	eg.Go(func() error {
			<-ctx.Done()
			return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
			fmt.Println(err)
			return 1
	}
	return 0
}

func runServer(ctx context.Context, session *sql.DB) error {
	s := &http.Server{
			Addr:    ":8888",
			Handler: &handlers.RootHandlerFunc{},
	}

	fmt.Println(session)

	errCh := make(chan error)
	go func() {
			defer close(errCh)
			if err := s.ListenAndServe(); err != nil {
					errCh <- err
			}
	}()

	select {
	case <-ctx.Done():
			return s.Shutdown(ctx)
	case err := <-errCh:
			return err
	}
}

func newDBSession() (*sql.DB, error) {
	s, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test_20191010?parseTime=true")
	if err != nil {
		return nil, err
	}
	return s, nil
}
