package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/xinyu/infra/inventory/host"
	"github.com/xinyu/infra/inventory/service"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var hostInfo host.Host
	{
		hostInfo = host.NewInmemHost()
		hostInfo = host.LoggingMiddleware(logger)(hostInfo)
	}

	var serviceInfo service.Service
	{
		serviceInfo = service.NewInmemService()
		serviceInfo = service.LoggingMiddleware(logger)(serviceInfo)
		serviceInfo = service.HostMiddleware(hostInfo)(serviceInfo)
	}

	mux := http.NewServeMux()
	mux.Handle("/host/v1/", host.MakeHTTPHandler(hostInfo, log.With(logger, "component", "HTTP")))
	mux.Handle("/service/v1/", service.MakeHTTPHandler(serviceInfo, log.With(logger, "component", "HTTP")))

	http.Handle("/", accessControl(mux))

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()

	logger.Log("exit", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
