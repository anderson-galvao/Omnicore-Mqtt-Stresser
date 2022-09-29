package http

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	Stresser "github.com/RacoWireless/iot-gw-mqtt-stresser/implementation/service"
	"github.com/RacoWireless/iot-gw-mqtt-stresser/implementation/utils"
	"github.com/RacoWireless/iot-gw-mqtt-stresser/model"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Stress godoc
// @Summary      Stress Mqtt
// @Description  Api For Stressing
// @Tags         stress
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.GenericResponse
// @Failure      400  {object}  model.GenericResponse
// @Failure      404  {object}  model.GenericResponse
// @Failure      500  {object}  model.GenericResponse
// @Router       /stress [post]
func (r *Handler) ExecuteStresser(c echo.Context) error {
	//ctx := c.Request().Context()
	req := new(model.Stresser)
	if err := c.Bind(req); err != nil {
		log.Error().Err(err).Msg("Error in Binding Request")
		r := utils.FrameGenericResponse(400, model.INVALIDJSON, err.Error())
		return c.JSON(r.StatusCode, r.Message)
	}
	if err := c.Validate(req); err != nil {
		r := utils.FrameGenericResponse(400, model.INVALIDJSON, err.Error())
		return c.JSON(http.StatusBadRequest, r.Message)
	}
	mResponse := utils.FrameGenericResponse(200, "Success", "")
	c.JSON(http.StatusOK, mResponse.Message)
	go func() {
		_, _ = r.Usecase.ExecuteStresser(*req)
	}()
	return nil
}

var singleData = Stresser.SummaryChannel{FastestPublishPerformance: 0, SlowestPublishPerformance: 9999}

var slowestPerformance float64
var fastestPerformance float64
var tenants = map[string]int{"Pepsi": 0, "Cola": 1}
var wsData = make([]Stresser.SummaryChannel, 0, 10)

func formBaseArray() {
	for _, _ = range tenants {
		wsData = append(wsData, singleData)
	}
}

func (r *Handler) StreamResults(c echo.Context) error {
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	var quit bool = false
	formBaseArray()
	channel := make(chan bool)
	go func() {
		for {
			_, _, errRead := ws.ReadMessage()
			if errRead != nil {
				channel <- true
			}
		}
	}()
	for {

		select {
		case <-time.After(time.Second):

		case <-channel:
			quit = true
		case summary := <-Stresser.SummaryChannelData:
			id := tenants[summary.Tenant]
			if wsData[id].SlowestPublishPerformance > summary.PublishPerformance[0] {
				slowestPerformance = summary.PublishPerformance[0]
			}
			if wsData[id].FastestPublishPerformance < summary.PublishPerformance[len(summary.PublishPerformance)-1] {
				fastestPerformance = summary.PublishPerformance[len(summary.PublishPerformance)-1]
			}
			wsData[id] = Stresser.SummaryChannel{
				Tenant:                    summary.Tenant,
				Clients:                   wsData[id].Clients + summary.Clients,
				TotalMessages:             wsData[id].TotalMessages + summary.TotalMessages,
				MessagesReceived:          wsData[id].MessagesReceived + summary.MessagesReceived,
				MessagesPublished:         wsData[id].MessagesPublished + summary.MessagesPublished,
				Errors:                    wsData[id].Errors + summary.Errors,
				Completed:                 wsData[id].Completed + summary.Completed,
				InProgress:                wsData[id].InProgress + summary.InProgress,
				ConnectFailed:             wsData[id].ConnectFailed + summary.ConnectFailed,
				SubscribeFailed:           wsData[id].SubscribeFailed + summary.SubscribeFailed,
				TimeoutExceeded:           wsData[id].TimeoutExceeded + summary.TimeoutExceeded,
				Aborted:                   wsData[id].Aborted + summary.Aborted,
				FastestPublishPerformance: math.Round(fastestPerformance),
				SlowestPublishPerformance: math.Round(slowestPerformance),
				PublishPerformanceMedian:  (wsData[id].PublishPerformanceMedian + summary.PublishPerformanceMedian) / 2,
			}

		}
		summaryJson, err := json.Marshal(wsData)
		if err != nil {
			log.Error().Err(err).Msg("Json Marshal Error")
			break
		}
		err = ws.WriteMessage(websocket.TextMessage, summaryJson)
		if err != nil {
			log.Error().Err(err).Msg("Unable TO Write Message to websocket")
			break
		}
		if quit == true {
			break
		}

	}
	return nil

}
