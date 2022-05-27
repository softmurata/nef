package assessionwithqos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HTTP3GPPAsSessionWithQosSubscriptionsAsAndScsLevelGet(c *gin.Context) {
	scsAsId := c.Params.ByName("scsAsId")
	fmt.Println("ScsAsId:  ", scsAsId)
	c.JSON(http.StatusOK, gin.H{})
}

func HTTP3GPPAsSessionWithQosSubscriptionsPost(c *gin.Context) {
	requestBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}
	fmt.Println("request Body:", requestBody)
	c.JSON(http.StatusOK, gin.H{})
}
