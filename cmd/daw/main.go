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
	iamSQL "github.com/adoublef/daw/internal/iam/sql"
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
	iamDSN      = "iam"
	iamDSNHelp  = "iam datasource name"
	dawDSN      = "daw"
	dawDSNShort = 'd'
	dawDSNHelp  = "daw datasource name"
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
		s.Flag(dawDSN, dawDSNHelp).Short(dawDSNShort).StringVar(&v.dawDSN)
		s.Flag(iamDSN, iamDSNHelp).StringVar(&v.iamDSN)
	}
	// migrate
	{
		v := &migrate{}
		m := f.Command(migrateName, migrateHelp).Action(v.migrate)
		m.Flag(dawDSN, dawDSNHelp).Short(dawDSNShort).StringVar(&v.dawDSN)
		m.Flag(iamDSN, iamDSNHelp).StringVar(&v.iamDSN)
	}
	// parse flags
	if _, err := f.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

// serve
type serve struct {
	addr   string
	dawDSN string
	iamDSN string
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
	if c.dawDSN == _EMPTY_ {
		c.dawDSN = defaultDSN
	}
	if c.iamDSN == _EMPTY_ {
		c.iamDSN = defaultDSN
	}

	dawDB, err := sql.Open(c.dawDSN)
	if err != nil {
		return err
	}
	defer dawDB.Close()

	iamDB, err := sql.Open(c.iamDSN)
	if err != nil {
		return err
	}
	defer iamDB.Close()

	// run a ping to check it created a

	s := server.New(c.addr, iamDB, dawDB)

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
	dawDSN string
	iamDSN string
}

func (c *migrate) migrate(_ *fisk.ParseContext) error {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		eg          = errgroup.New(ctx)
	)
	defer cancel()

	if c.dawDSN == _EMPTY_ {
		c.dawDSN = defaultDSN
	}

	dawDB, err := sql.Open(c.dawDSN)
	if err != nil {
		return err
	}
	defer dawDB.Close()

	iamDB, err := sql.Open(c.iamDSN)
	if err != nil {
		return err
	}
	defer iamDB.Close()

	eg.Go(func(ctx context.Context) error {
		return dawSQL.Up(ctx, dawDB)
	})

	eg.Go(func(ctx context.Context) error {
		return iamSQL.Up(ctx, iamDB)
	})

	return eg.Wait()
}
