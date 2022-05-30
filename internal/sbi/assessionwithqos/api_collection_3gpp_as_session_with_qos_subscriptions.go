package assessionwithqos

import (
	"fmt"
	"net/http"

	"github.com/free5gc/openapi"
	"github.com/free5gc/util/httpwrapper"

	"github.com/gin-gonic/gin"
	"github.com/softmurata/freeopenapi/models"
	"github.com/softmurata/nef/internal/logger"
	"github.com/softmurata/nef/internal/sbi/producer"
)

func HTTP3GPPAsSessionWithQosSubscriptionsAsAndScsLevelGet(c *gin.Context) {
	scsAsId := c.Params.ByName("scsAsId")
	fmt.Println("ScsAsId:  ", scsAsId)

	rsp := producer.Handle3GPPAsSessionWithQosSubscriptionsAsAndScsLevelGet(scsAsId)

	for key, val := range rsp.Header {
		c.Header(key, val[0])
	}

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")

	if err != nil {
		logger.SessionQosLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}

	// c.JSON(http.StatusOK, gin.H{})
}

func HTTP3GPPAsSessionWithQosSubscriptionsPost(c *gin.Context) {
	var asSessionWithQosSubscriptionsRequest models.ScsAsDataContext

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.ServParamLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&asSessionWithQosSubscriptionsRequest, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.ServParamLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// scsAsId := c.Params.ByName("scsAsId")

	req := httpwrapper.NewRequest(c.Request, asSessionWithQosSubscriptionsRequest)
	req.Params["scsAsId"] = c.Params.ByName("scsAsId")

	rsp := producer.Handle3GPPAsSessionWithQosSubscriptionsPost(req)

	fmt.Println("response:", rsp)

	for key, val := range rsp.Header {
		c.Header(key, val[0])
	}

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")

	if err != nil {
		logger.SessionQosLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}
}
