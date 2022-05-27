package assessionwithqos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softmurata/nef/internal/logger"
)

func HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionIdGet(c *gin.Context) {
	logger.SessionQosLog.Info("HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionId Get Method")
	scsAsId := c.Params.ByName("scsAsId")
	subscriptionId := c.Params.ByName("subscriptionId")

	fmt.Println("scsAsId: ", scsAsId)
	fmt.Println("subscriptionId: ", subscriptionId)

	c.JSON(http.StatusOK, gin.H{})

}

func HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionIdPut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})

}

func HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionIdPatch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})

}

func HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionIdDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})

}
