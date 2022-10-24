//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"stopwatch/apiserver"
	"stopwatch/apiservices"
	"stopwatch/eliona"
	"stopwatch/stopwtch"
	"syscall"
	"time"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	MODULE          = "app"
	API_DATA_BUFFER = 20000
)

var (
	stopwatches []api.Asset
	ir          chan bool
	swMng       stopwtch.StopwatchManager
	apiServer   *http.Server
	osIr        bool
)

func actualizeStopwatches() {
	var err error

	stopwatches, err = eliona.GetStopwatches()
	log.Debug(MODULE, "got stopwatches %v with err %v", stopwatches, err)
}

func setupApp() {
	osIr = false

	log.Info(MODULE, "setup application")

	var (
		apiData chan api.Data
	)

	ir = make(chan bool)
	apiData = make(chan api.Data, API_DATA_BUFFER)

	actualizeStopwatches()

	swMng = *stopwtch.NewStopwatchManager(stopwatchCallbackHandler)

	go eliona.ListenHeapEvents(ir, apiData)
	go stopwatchEventCatcher(apiData)
	go catchAndExitOnOsSig()
}

func catchAndExitOnOsSig() {
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-osSig
	osIr = true
	log.Info("app", "Os Signal catched, exiting...")
	ir <- true
	swMng.StopAll()
	log.Info("app", "all stopwatches stopped")
	apiServer.Shutdown(context.Background())
	log.Info("app", "api server stopped")
}

// start / stop stopwatches
func stopwatchEventCatcher(apiData <-chan api.Data) {
	for data := range apiData {
		log.Debug(MODULE, "data from datalistener: %v", data)

		asset := getStopwatchByAssetId(data.AssetId)
		if asset != nil {
			processStartStopEvent(data.AssetId, data.Data)
		}
	}
	log.Warn(MODULE, "eventcatcher exited")
}

func processStartStopEvent(assetId int32, data map[string]interface{}) {
	start := data[eliona.ATTRIBUTE_START]
	stop := data[eliona.ATTRIBUTE_STOP]

	// prioritize stop signal
	if stop != nil && stop.(float64) >= 1 {
		lastTime := swMng.Stop(assetId)
		if lastTime.Seconds() > 0 {
			log.Info(MODULE, "timer stoppend %d @ %f s",
				assetId,
				lastTime.Seconds())
			eliona.UpdateLastTime(assetId, lastTime.Seconds())
		}
	} else if start != nil && start.(float64) >= 1 {
		log.Info(MODULE, "timer started %d", assetId)
		swMng.Start(assetId)
	}
}

// stopwatch ticks
func stopwatchCallbackHandler(assetId int32, t time.Duration) {
	log.Debug(MODULE, "stopwatch callback called %d,%v", assetId, t)

	eliona.UpdateTime(assetId, t.Seconds())
}

func getStopwatchByAssetId(assetId int32) *api.Asset {
	for _, sw := range stopwatches {
		if sw.GetId() == assetId {
			return &sw
		}
	}
	return nil
}

// listenApiRequests starts an API server and listen for API requests
// The API endpoints are defined in the openapi.yaml file
func listenApiRequests() {
	apiServer = &http.Server{
		Addr: ":" + common.Getenv("API_SERVER_PORT", "3000"),
		Handler: apiserver.NewRouter(
			apiserver.NewUtilsApiController(apiservices.NewUtilsApiService()),
		),
	}

	err := apiServer.ListenAndServe()

	if !osIr {
		log.Fatal(MODULE, "Error in API Server: %v", err)
	} else {
		log.Info(MODULE, "Api Server stopped todue an os interrupt")
	}
}
