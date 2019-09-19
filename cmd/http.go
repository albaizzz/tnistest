package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tnistest/config"
	"github.com/tnistest/internal/http/handlers"
	"github.com/tnistest/internal/repositories"
	"github.com/tnistest/internal/services"

	internalHttp "github.com/tnistest/internal/http"

	"github.com/tnistest/pkg/mysql"
)

type httpCmd struct {
	stop     <-chan bool
	BaseCmd  *cobra.Command
	Config   *config.Config
	filename string
}

func NewHttpCmd(config *config.Config) *httpCmd {
	return NewHttpCmdSignaled(config, nil)
}

func NewHttpCmdSignaled(config *config.Config, stop <-chan bool) *httpCmd {

	cc := &httpCmd{stop: stop}
	cc.Config = config
	cc.BaseCmd = &cobra.Command{
		Use:   "http",
		Short: "Used to run http",
		RunE:  cc.server,
	}
	fs := pflag.NewFlagSet("root", pflag.ContinueOnError)
	fs.StringVarP(&cc.filename, "file", "f", "", "Custom configuration filename")
	cc.BaseCmd.Flags().AddFlagSet(fs)
	return cc
}

func (h *httpCmd) server(cmd *cobra.Command, args []string) (err error) {
	if len(h.filename) > 1 {
		h.Config = config.New(h.filename,
			"./config",
			"../config",
			"../../config",
			"../../../config")
	}

	// init sql
	conn := fmt.Sprintf(
		mysql.MysqlDataSourceFormat,
		h.Config.MariaDB.User,
		h.Config.MariaDB.Password,
		h.Config.MariaDB.Host,
		h.Config.MariaDB.Port,
		h.Config.MariaDB.DbName,
		h.Config.MariaDB.Charset,
	)

	//db.OpenConnection(conn, injector.Config)
	db := mysql.NewMySQL()
	db.OpenConnection(conn, h.Config)
	db.SetConnMaxLifetime(h.Config.MariaDB.MaxLifeTime)
	db.SetMaxIdleConn(h.Config.MariaDB.MaxIdleConnection)
	db.SetMaxOpenConn(h.Config.MariaDB.MaxOpenConnection)

	// init repo
	accountRepo := repositories.NewAccountRepository(db)
	balanceRepo := repositories.NewBalanceRepository(db)
	transactionRepo := repositories.NewTransactionHistory(db)

	accountSvc := services.NewAccountService(accountRepo, balanceRepo)
	balanceSvc := services.NewBalanceService(balanceRepo, transactionRepo)
	transactionSvc := services.NewTransactionHistory(transactionRepo)

	accountHandler := handlers.NewAccountHandler(accountSvc)
	balanceHandler := handlers.NewBalanceHandler(balanceSvc)
	transactionHandler := handlers.NewTransactionHandler(transactionSvc)

	// inject routes
	route := &internalHttp.Routes{
		Config:             h.Config,
		AccountHandler:     accountHandler,
		BalanceHandler:     balanceHandler,
		TransactionHandler: transactionHandler,
	}
	router := route.NewRoutes()

	// Description µ micro service
	fmt.Println(
		fmt.Sprintf(
			WelkomText,
			h.Config.App.Port,
			strings.Join([]string{
				h.Config.Log.Dir,
				h.Config.Log.Filename}, "/"),
		))

	//tableRoute(router) // Prettier Route Pattern

	return h.serve(router)
}

func (h *httpCmd) serve(router http.Handler) error {
	errCh := make(chan error, 1)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	addr := net.JoinHostPort("",
		strconv.Itoa(h.Config.App.Port))
	s := StartWebServer(
		addr,
		h.Config.App.ReadTimeout,
		h.Config.App.WriteTimeout,
		router,
	)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logrus.Info(
				"Server gracefully ListenAndServe",
				logrus.Fields{})
			errCh <- err
		}
		<-h.stop
	}()

	if h.stop != nil {
		select {
		case err := <-errCh:
			logrus.Info(
				"Server gracefully h stop stopped",
				logrus.Fields{})
			return err
		case <-h.stop:
		case <-quit:
		}
	} else {
		select {
		case err := <-errCh:
			logrus.Info(
				"Server gracefully stopped",
				logrus.Fields{})
			return err
		case <-quit:
		}
	}
	return nil
}

// StartWebServer starts a web server
func StartWebServer(addr string, readTimeout, writeTimeout int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}
}
