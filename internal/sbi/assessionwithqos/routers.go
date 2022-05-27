package assessionwithqos

import (
	"net/http"
	"strings"

	logger_util "github.com/free5gc/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/softmurata/nef/internal/logger"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)
	AddService(router)
	return router
}

func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("/3gpp-as-session-with-qos/v1")

	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.HandlerFunc)
		case "PATCH":
			group.PATCH(route.Pattern, route.HandlerFunc)
		}
	}

	return group
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/",
		Index,
	},
	{
		"3GPPAsSessionWithQosSubscriptionsAsAndScsLevelGet",
		strings.ToUpper("Get"),
		"/:scsAsId/subscriptions",
		HTTP3GPPAsSessionWithQosSubscriptionsAsAndScsLevelGet,
	},
	{
		"3GPPAsSessionWithQosSubscriptionsPost",
		strings.ToUpper("Post"),
		"/:scsAsId/subscriptions",
		HTTP3GPPAsSessionWithQosSubscriptionsPost,
	},

	{
		"3GPPAsSessionWithQosSubscriptionsSubscriptionIdGet",
		strings.ToUpper("Get"),
		"/:scsAsId/subscriptions/:subscriptionId",
		HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionIdGet,
	},

	{
		"3GPPAsSessionWithQosSubscriptionsSubscriptionIdPut",
		strings.ToUpper("Put"),
		"/:scsAsId/subscriptions/:subscriptionId",
		HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionIdPut,
	},

	{
		"3GPPAsSessionWithQosSubscriptionsSubscriptionIdPatch",
		strings.ToUpper("Patch"),
		"/:scsAsId/subscriptions/:subscriptionId",
		HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionIdPatch,
	},
	{
		"3GPPAsSessionWithQosSubscriptionsSubscriptionIdDelete",
		strings.ToUpper("Delete"),
		"/:scsAsId/subscriptions/:subscriptionId",
		HTTP3GPPAsSessionWithQosSubscriptionsSubscriptionIdDelete,
	},
}
