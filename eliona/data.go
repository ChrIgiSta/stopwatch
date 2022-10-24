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

package eliona

import (
	"context"
	"encoding/json"
	"stopwatch/conn"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

const (
	ATTRIBUTE_START        = "start"
	ATTRIBUTE_STOP         = "stop"
	ATTRIBUTE_CURRENT_TIME = "current_time"
	ATTRIBUTE_LAST_TIME    = "last_time"
	SUBTYPE_ACTIONS        = api.SUBTYPE_OUTPUT
	SUBTYPE_VALUES         = api.SUBTYPE_INPUT
	SUBTYPE_LASTTIME       = api.SUBTYPE_INFO
	RX_BUFFER              = 20000
	WEBSOCKET_ENDPOINT     = "/data-listener"
)

func UpdateTime(assetId int32, newTime float64) error {

	dataData := make(map[string]interface{})
	dataData[ATTRIBUTE_CURRENT_TIME] = newTime

	data := api.Data{
		AssetId: assetId,
		Subtype: SUBTYPE_VALUES,
		// Timestamp:     *api.NewNullableTime(nil),
		Data: dataData,
	}
	resp, err := client.NewClient().DataApi.PutData(context.Background()).Data(data).Execute()

	log.Debug(MODULE, "get current time response: %v", resp)

	return err
}

func UpdateLastTime(assetId int32, newTime float64) error {

	dataData := make(map[string]interface{})
	dataData[ATTRIBUTE_LAST_TIME] = newTime

	data := api.Data{
		AssetId: assetId,
		Subtype: SUBTYPE_LASTTIME,
		// Timestamp:     *api.NewNullableTime(nil),
		Data: dataData,
	}
	resp, err := client.NewClient().DataApi.PutData(context.Background()).Data(data).Execute()

	log.Debug(MODULE, "get last time response: %v", resp)

	return err
}

func ListenHeapEvents(ir <-chan bool, rxApiData chan<- api.Data) {
	var (
		rx chan []byte
	)
	defer close(rxApiData)

	rx = make(chan []byte, RX_BUFFER)

	ws := conn.NewWebsocketClient(client.ApiEndpointString()+WEBSOCKET_ENDPOINT, false)
	go ws.ServeForever(rx, ir)
	for data := range rx {
		var apiData api.Data
		err := json.Unmarshal(data, &apiData)
		if err == nil {
			// prefilter data
			if apiData.Subtype == SUBTYPE_ACTIONS {
				rxApiData <- apiData
			}
		} else {
			log.Warn(MODULE, "cannot unmarshal ws data", err)
		}
	}
	log.Warn(MODULE, "rx channel from ws closed. exiting event listener...")
}
