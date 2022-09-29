package http

import (
	"encoding/json"
	"net/http"
	"time"

	Stresser "github.com/RacoWireless/iot-gw-stresser/implementation/service"
	"github.com/RacoWireless/iot-gw-stresser/implementation/utils"
	"github.com/RacoWireless/iot-gw-stresser/model"
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
	var wsData = Stresser.SummaryChannel{}
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
			var slowestPerformance float64
			var fastestPerformance float64
			if wsData.SlowestPublishPerformance > summary.PublishPerformance[0] {
				slowestPerformance = summary.PublishPerformance[0]
			}
			if wsData.FastestPublishPerformance > summary.PublishPerformance[len(summary.PublishPerformance)-1] {
				fastestPerformance = summary.PublishPerformance[len(summary.PublishPerformance)-1]
			}
			wsData = Stresser.SummaryChannel{
				Clients:                   wsData.Clients + summary.Clients,
				TotalMessages:             wsData.TotalMessages + summary.TotalMessages,
				MessagesReceived:          wsData.MessagesReceived + summary.MessagesReceived,
				MessagesPublished:         wsData.MessagesPublished + summary.MessagesPublished,
				Errors:                    wsData.Errors + summary.Errors,
				Completed:                 wsData.Completed + summary.Completed,
				InProgress:                wsData.InProgress + summary.InProgress,
				ConnectFailed:             wsData.ConnectFailed + summary.ConnectFailed,
				SubscribeFailed:           wsData.SubscribeFailed + summary.SubscribeFailed,
				TimeoutExceeded:           wsData.TimeoutExceeded + summary.TimeoutExceeded,
				Aborted:                   wsData.Aborted + summary.Aborted,
				FastestPublishPerformance: fastestPerformance,
				SlowestPublishPerformance: slowestPerformance,
				PublishPerformanceMedian:  (wsData.PublishPerformanceMedian + summary.PublishPerformanceMedian) / 2,
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
