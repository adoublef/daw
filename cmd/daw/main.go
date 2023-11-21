package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adoublef/daw/cmd/daw/server"
	"github.com/adoublef/daw/errgroup"
	dawSQL "github.com/adoublef/daw/internal/daw/sql"
	"github.com/adoublef/daw/sql"
	"github.com/choria-io/fisk"
)

const (
	appName     = "daw"
	appHelp     = ""
	serveName   = "serve"
	serveHelp   = "run application"
	migrateName = "migrate"
	migrateHelp = "run database migrations"
	addrName    = "addr"
	addrHelp    = "listen address"
	defaultAddr = ":8080"
	dsnName     = "dsn"
	dsnShort    = 'd'
	dsnHelp     = "datasource name"
	defaultDSN  = ":memory:"
	_EMPTY_     = ""
)

func main() {
	f := fisk.New(appName, appHelp)
	// serve
	{
		v := &serve{}
		s := f.Command(serveName, serveHelp).Action(v.serve)
		s.Flag(addrName, addrHelp).StringVar(&v.addr)
		s.Flag(dsnName, dsnHelp).Short(dsnShort).StringVar(&v.dsn)
	}
	// migrate
	{
		v := &migrate{}
		m := f.Command(migrateName, migrateHelp).Action(v.migrate)
		m.Flag(dsnName, dsnHelp).Short(dsnShort).StringVar(&v.dsn)
	}
	// parse flags
	if _, err := f.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

// serve
type serve struct {
	addr string
	dsn  string
}

func (c *serve) serve(_ *fisk.ParseContext) error {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		eg          = errgroup.New(ctx)
	)
	defer cancel()

	if c.addr == _EMPTY_ {
		c.addr = defaultAddr
	}
	if c.dsn == _EMPTY_ {
		c.dsn = defaultDSN
	}

	db, err := sql.Open(c.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// run a ping to check it created a

	s := server.New(c.addr)

	eg.Go(func(ctx context.Context) error {
		return s.ListenAndServe()
	})

	eg.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return s.Shutdown()
	})

	return eg.Wait()
}

// migrate
type migrate struct {
	dsn string
}

func (c *migrate) migrate(_ *fisk.ParseContext) error {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	)
	defer cancel()

	if c.dsn == _EMPTY_ {
		c.dsn = defaultDSN
	}

	db, err := sql.Open(c.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return dawSQL.Up(ctx, db)
}
