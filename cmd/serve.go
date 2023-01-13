/*
Copyright © 2022 Tom Siewert <tom@siewert.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tomsiewert/fast-dhcp/pkg/model"
	fprometheus "github.com/tomsiewert/fast-dhcp/pkg/prometheus"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the DHCP server",
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" {
			logrus.Fatalf("No config file defined")
		}
		config := readConfig(cfgFile)

		if config.PProfListen != "" {
			go func() {
				logrus.Print("Start pprof server on ", config.PProfListen)
				if err := http.ListenAndServe(config.PProfListen, nil); err != nil {
					logrus.Fatal("Error starting pprof server: ", err)
				}
			}()
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				logrus.Println(sig)
				os.Exit(0)
			}
		}()

		initSentry(config)
		serveDHCP(config)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func initSentry(config *model.Config) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   config.SentryDSN,
		Debug: true,
	})

	if err != nil {
		logrus.Fatalf("Sentry Init failed. %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}

func serveDHCP(config *model.Config) {
	go func() {
		if err := fprometheus.ServePrometheusHandler(&config.Prometheus); err != nil {
			logrus.Fatalf("Prometheus HTTP failed. %s", err)
		}
	}()
}
