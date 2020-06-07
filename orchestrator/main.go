package orchestrator

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"html/template"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const Name = "geoffrey-orchestrator"

func CreateCLI() cli.App {
	a := cli.App{
		Name:   Name,
		Action: Main,
		Usage:  "geoffrey, the general automation tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				DefaultText: "/etc/geoffrey/config.yaml",
				Usage:       "use `FILE` for the config for the orchestrator",
				EnvVars:     []string{"GEOFFREY_CONFIG_O"},
			},
		},
	}
	return a
}

func Main(ctx *cli.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go httpServer(context.Background(), wg.Done)

	wg.Wait()
	return nil
}

func httpServer(ctx context.Context, done func()) {
	defer done()
	logrus.Info("Starting server")
	s := http.Server{
		Addr: ":3000",
	}
	mux := http.DefaultServeMux
	s.Handler = mux
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		rand.Seed(time.Now().Unix())
		requestId := rand.Int63()
		log := logrus.WithFields(logrus.Fields{
			"pathRequested": request.RequestURI,
			"ip":            request.RemoteAddr,
			"id":            requestId,
		})
		log.Info("start request")
		defer func() {
			log.WithField("duration", time.Now().Sub(now).String()).Info("end request")
		}()
		t, err := template.New("homepage").Parse(homepage)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		writer.WriteHeader(http.StatusOK)
		err = t.Execute(writer, nil)
		if err != nil {
			log.Error(err)
		}

	})
	go func() {
		for {
			select {
			case <-ctx.Done():

			}
		}
	}()


	err := s.ListenAndServe()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Server done")
}

var homepage = `{{/* for formatting*/ -}}
<!DOCTYPE html>
<html>
	<head>
	</head>
	<body>
		<h1>Hello</h1>
	</body>
</html>
`
