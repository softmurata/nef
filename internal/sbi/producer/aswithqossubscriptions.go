package producer

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/free5gc/util/httpwrapper"
	"github.com/softmurata/freeopenapi/Nnrf_NFDiscovery"
	"github.com/softmurata/freeopenapi/Npcf_PolicyAuthorization"
	"github.com/softmurata/freeopenapi/models"
	nef_context "github.com/softmurata/nef/internal/context"
	"github.com/softmurata/nef/internal/logger"
	"github.com/softmurata/nef/internal/sbi/consumer"
)

// Get
func Handle3GPPAsSessionWithQosSubscriptionsAsAndScsLevelGet(scsAsId string) *httpwrapper.Response {

	logger.ProducerLog.Info("Handle 3GPPAsSessionWithQos Subscriptions AsAndScsLevel Get")

	response, _, _ := GetQosSubscriptionsAsAndScsLevelProcedure(scsAsId)

	return response
}

func GetQosSubscriptionsAsAndScsLevelProcedure(scsAsId string) (*httpwrapper.Response, string, *models.ProblemDetails) {
	// get nef context
	nefSelf := nef_context.NEF_Self()

	// find scs as data context from scs as id
	ScsAsDataContexts, err := nefSelf.FindScsAsDataContextsByScsAsId(scsAsId)

	if err != nil {
		logger.ProducerLog.Errorf("Cannot Find Scs As Data Contexts")

		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED",
		}

		response := httpwrapper.NewResponse(http.StatusInternalServerError, nil, problemDetails)

		return response, "", &problemDetails
	}

	locationHeader := "http://10.0.10.1:8000/3gpp-as-session-with-qos/v1/{nfId}/subscriptions/{subscriptionId}"

	// convert model
	var rspScsAsDataContexts []models.ScsAsDataContextResponse

	for _, sadContext := range ScsAsDataContexts {
		tmpSadContext := models.ScsAsDataContextResponse{
			Self:                    sadContext.ScsAsDataContext.Self,
			FlowInfos:               sadContext.ScsAsDataContext.FlowInfos,
			NotificationDestination: sadContext.ScsAsDataContext.NotificationDestination,
			EthFlowInfos:            sadContext.ScsAsDataContext.EthFlowInfos,
			QosReference:            sadContext.ScsAsDataContext.QosReference,
			UeIpv4Addr:              sadContext.ScsAsDataContext.UeIpv4Addr,
			UeIpv6Addr:              sadContext.ScsAsDataContext.UeIpv6Addr,
			MacAddr:                 sadContext.ScsAsDataContext.MacAddr,
			QosMonInfo:              sadContext.ScsAsDataContext.QosMonInfo,
			SupportedFeatures:       sadContext.ScsAsDataContext.SupportedFeatures,
			Supi:                    sadContext.ScsAsDataContext.Supi,
			Dnn:                     sadContext.ScsAsDataContext.Dnn,
		}

		rspScsAsDataContexts = append(rspScsAsDataContexts, tmpSadContext)
	}

	response := models.ScsAsDataContextsResponse{
		ScsAsDataContexts: rspScsAsDataContexts,
	}

	httpRsp := httpwrapper.NewResponse(http.StatusOK, nil, response)

	return httpRsp, locationHeader, nil

}

// Post
// models

/*
type FlowInfo struct {
	FlowDescriptions []string `json:"flowDescriptions" yaml:"flowDescriptions" bson:"flowDescriptions" mapstructure:"FlowDescriptions"`
	FlowId           int32    `json:"flowId" yaml:"flowId" bson:"flowId" mapstructure:"FlowId"`
}

type EthFlowInfo struct {
	DestMacAddr    string
	EthType        string
	FDescs         string
	FDir           []string
	SourceMacAddr  string
	VlanTags       []string
	SrcMacAddrEnd  string
	DestMacAddrEnd string
}

type QosMonitoringInformation struct {
	ReqQosMonParams []RequestedQosMonitoringParameter `json:"reqQosMonParams" yaml:"reqQosMonParams" bson:"reqQosMonParams" mapstructure:"ReqQosMonParams"`
	RepFreqs        []ReportingFrequency              `json:"repFreqs" yaml:"repFreqs" bson:"repFreqs" mapstructure:"RepFreqs"`
	RepThreshDl     int                               `json:"repThreshDl,omitempty" yaml:"repThreshDl" bson:"repThreshDl" mapstructure:"RepThreshDl"`
	RepThreshUl     int                               `json:"repThreshUl,omitempty" yaml:"repThreshUl" bson:"repThreshUl" mapstructure:"RepThreshUl"`
	RepThreshRp     int                               `json:"repThreshRp,omitempty" yaml:"repThreshRp" bson:"repThreshRp" mapstructure:"RepThreshRp"`
	RepPeriod       int                               `json:"repPeriod,omitempty" yaml:"repPeriod" bson:"repPeriod" mapstructure:"RepPeriod"`
	WaitTime        int                               `json:"waitTime" yaml:"waitTime" bson:"waitTime" mapstructure:"WaitTime"`
}

// Requested Qos Monitoring Parameter.go
type RequestedQosMonitoringParameter string

const (
	RequestedQosMonitoringParameter_DOWNLINK   RequestedQosMonitoringParameter = "DOWNLINK"
	RequestedQosMonitoringParameter_UPLINK     RequestedQosMonitoringParameter = "UPLINK"
	RequestedQosMonitoringParameter_ROUND_TRIP RequestedQosMonitoringParameter = "ROUND_TRIP"
)

// Reporting Frequency.go
type ReportingFrequency string

const (
	ReportingFrequency_EVENT_TRIGGERED ReportingFrequency = "EVENT_TRIGGERED"
	ReportingFrequency_PERIODIC        ReportingFrequency = "PERIODIC"
	ReportingFrequency_SESSION_RELEASE ReportingFrequency = "SESSION_RELEASE"
)

type AsSessionWithQosSubscriptionsRequest struct {
	FlowInfos               []FlowInfo               `json:"flowInfos,omitempty" yaml:"flowInfos" bson:"flowInfos" mapstructure:"FlowInfos"`
	NotificationDestination string                   `json:"notificationDestination" yaml:"notificationDestination" bson:"notificationDestination" mapstructure:"NotificationDestination"`
	EthFlowInfos            []EthFlowInfo            `json:"ethFlowInfos,omitempty" yaml:"ethFlowInfos" bson:"ethFlowInfos" mapstructure:"EthFlowInfos"`
	QosReference            string                   `json:"qosReference,omitempty" yaml:"qosReference" bson:"qosReference" mapstructure:"QosReference"`
	UeIpv4Addr              string                   `json:"ueIpv4Addr,omitempty" yaml:"ueIpv4Addr" bson:"ueIpv4Addr" mapstructure:"UeIpv4Addr"`
	UeIpv6Addr              string                   `json:"ueIpv6Addr,omitempty" yaml:"ueIpv6Addr" bson:"ueIpv6Addr" mapstructure:"UeIpv6Addr"`
	MacAddr                 string                   `json:"macAddr,omitempty" yaml:"macAddr" bson:"macAddr" mapstructure:"MacAddr"`
	QosMonInfo              QosMonitoringInformation `json:"qosMonInfo,omitempty" yaml:"qosMonInfo" bson:"qosMonInfo" mapstructure:"QosMonInfo"`
	SupportedFeatures       string                   `json:"supportedFeatures,omitempty" yaml:"supportedFeatures" bson:"supportedFeatures" mapstructure:"SupportedFeatures"`
	Supi                    string                   `json:"supi,omitempty" yaml:"supi" bson:"supi" mapstructure:"Supi"` // for loadcore
	Dnn                     string                   `json:"dnn,omitempty" yaml:"dnn" bson:"dnn" mapstructure:"Dnn"`     // for free5gc
}


type AsSessionWithQosSubscriptionsResponse struct {
	Self string
}
*/

func Handle3GPPAsSessionWithQosSubscriptionsPost(request *httpwrapper.Request) *httpwrapper.Response {
	logger.SessionQosLog.Infof("Handle 3GPP As Session With Qos Subscriptions Post")
	scsAsId := request.Params["scsAsId"]
	asSessionWithQosSubscriptionsRequest := request.Body.(models.ScsAsDataContext)

	response, locationHeader, problemDetails := PostQosSubscriptionsProcedure(scsAsId, asSessionWithQosSubscriptionsRequest)

	headers := http.Header{
		"Location": {locationHeader},
	}

	if response != nil {
		return httpwrapper.NewResponse(http.StatusCreated, headers, response)
	} else if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
}

func PostQosSubscriptionsProcedure(scsAsId string, req models.ScsAsDataContext) (*models.ScsAsDataContextResponse, string, *models.ProblemDetails) {

	// error handling for request body
	if req.UeIpv4Addr == "" && req.UeIpv6Addr == "" && req.MacAddr == "" {
		problemDetails := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Detail: "Ue UeIpv4 and UeIpv6 and UeMac are all empty",
			Cause:  "ERROR_REQUEST_PARAMETERS",
		}

		return nil, "", &problemDetails
	}

	// get nef context
	nefSelf := nef_context.NEF_Self()

	// find scs as data context from scs as id
	tempScsAsDataContexts, _ := nefSelf.FindScsAsDataContextsByScsAsId(scsAsId)

	if tempScsAsDataContexts != nil {
		logger.ProducerLog.Warnf("Existed Scs As Data Context: %+v", tempScsAsDataContexts)
		// problemDetails := models.ProblemDetails{
		// 	Status: http.StatusInternalServerError,
		// 	Detail: "Existed Scs As Data Context",
		// 	Cause:  "FAILED_NEW_SCS_AS_DATA_CONTEXT_CREATED",
		// }

		// return nil, "", &problemDetails
	}

	serviceName := "npcf-policyauthorization"

	// get from request body
	supi := req.Supi                  // "imsi-208930000000003"
	suppFeat := req.SupportedFeatures // "f0"
	ueIpv4 := req.UeIpv4Addr          // "10.60.0.1"
	dnn := req.Dnn                    // "internet"

	fmt.Println("  supi:", supi)
	fmt.Println(" suppFeat:", suppFeat)
	fmt.Println(" ueIpv4: ", ueIpv4)

	// call npcf-authorization/v1/app-sessions POST
	// request body
	notifUri := nefSelf.GetIpv4Uri() + "/" + serviceName
	medComponents := make(map[string]models.MediaComponent)
	medSubComps := make(map[string]models.MediaSubComponent)

	for idx, info := range req.FlowInfos {
		// fmt.Println("idx: ", idx)
		// fmt.Println("info: ", info)

		stridx := strconv.Itoa(idx + 1)
		medSubComps[stridx] = models.MediaSubComponent{
			FDescs: info.FlowDescriptions,
			FNum:   info.FlowId,
		}

		medComponents[stridx] = models.MediaComponent{
			FStatus:      models.FlowStatus_ENABLED,
			MedCompN:     int32(idx + 1),
			MedSubComps:  medSubComps,
			QosReference: req.QosReference,
		}

	}

	/*
		// for test
		fdescs := []string{"permit out ip from any to 10.60.0.1", "permit out 6 from 10.60.0.0/16 to 10.60.0.1"}
		medSubComps["1"] = models.MediaSubComponent{
			FDescs: fdescs,
			FNum:   1,
		}

		medComponents["1"] = models.MediaComponent{
			FStatus:      models.FlowStatus_ENABLED,
			MedCompN:     1,
			MedSubComps:  medSubComps,
			QosReference: "ref2",
		}
	*/

	ascReqData := &models.AppSessionContextReqData{
		EvSubsc: &models.EventsSubscReqData{
			Events: []models.AfEventSubscription{
				{
					Event:       models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION,
					NotifMethod: models.AfNotifMethod_EVENT_DETECTION,
				},
				{
					Event:       models.AfEvent_FAILED_RESOURCES_ALLOCATION,
					NotifMethod: models.AfNotifMethod_EVENT_DETECTION,
				},
			},
			NotifUri: notifUri,
		},
		MedComponents: medComponents,
		NotifUri:      notifUri,
		Supi:          supi,
		SuppFeat:      suppFeat,
		UeIpv4:        ueIpv4,
		Dnn:           dnn,
	}

	appSessionContextReq := models.AppSessionContext{
		AscReqData: ascReqData,
	}

	fmt.Println(" appSessionContextReq: ", appSessionContextReq)

	// call POST npcf-authorization/v1/app-sessions
	pcfUri := ""
	targetNfType := models.NfType_PCF
	requestNfType := models.NfType_NEF
	localVarOptions := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}

	result, localErr := consumer.SendSearchNFInstances(nefSelf.NrfUri, targetNfType, requestNfType, localVarOptions)

	if localErr != nil {
		logger.ConsumerLog.Errorf("Search Nf Instances response error: %+v", localErr)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: localErr.Error(),
			Cause:  "FAILED_SEARCH_PCF",
		}

		return nil, "", &problemDetails
	}

	nfProfile := result.NfInstances[0] // ToDo: implement algorithm for selecting suitable pcf @rofile
	for _, nfService := range *nfProfile.NfServices {
		if nfService.ServiceName == models.ServiceName_NPCF_POLICYAUTHORIZATION && nfService.NfServiceStatus == models.NfServiceStatus_REGISTERED {
			fmt.Println("nfservice: ", nfService) // http://127.0.0.7:8000/npcf-smpolicycontrol/v1/
			pcfPrefix := nfService.ApiPrefix      // http://127.0.0.7:8000

			pcfUri = pcfPrefix
		}

	}

	if pcfUri == "" {
		logger.ProducerLog.Errorln(" cannot find pcf")

		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: localErr.Error(),
			Cause:  "FAILED_SEARCH_PCF",
		}

		return nil, "", &problemDetails

	}

	configuration := Npcf_PolicyAuthorization.NewConfiguration()
	configuration.SetBasePath(pcfUri)
	client := Npcf_PolicyAuthorization.NewAPIClient(configuration)

	fmt.Println(" API Client: ", client)

	rspData, _, err := client.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), appSessionContextReq)

	if err != nil {
		logger.ProducerLog.Errorf("Failed to call Policy authorization API app session context error: %+v", err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: localErr.Error(),
			Cause:  "FAILED_CALL_API",
		}

		return nil, "", &problemDetails
	}

	fmt.Println(" Policy Authorization rspData:", rspData)

	// create response for 3gpp-as-session-with-qos api
	// ToDo: get api prefix and service name and change subscription id and scsasid
	selfUri := "http://10.0.10.1:8000/3gpp-as-session-with-qos/v1/{nfId}/subscriptions/{subscriptionId}"

	response := models.ScsAsDataContextResponse{
		Self:                    selfUri,
		FlowInfos:               req.FlowInfos,
		NotificationDestination: req.NotificationDestination,
		EthFlowInfos:            req.EthFlowInfos,
		QosReference:            req.QosReference,
		UeIpv4Addr:              req.UeIpv4Addr,
		UeIpv6Addr:              req.UeIpv6Addr,
		MacAddr:                 req.MacAddr,
		QosMonInfo:              req.QosMonInfo,
		SupportedFeatures:       req.SupportedFeatures,
		Supi:                    req.Supi,
		Dnn:                     req.Dnn,
	}

	// create scs as data context
	subscriptionId := strconv.Itoa(nefSelf.SubscriptionIdGenerator) // ToDo: create id generator

	storeDataSubsc := nef_context.SubscriptionScsAsDataContext{
		SubscriptionId:   subscriptionId,
		ScsAsId:          scsAsId,
		ScsAsDataContext: &response,
	}
	subscPool := nefSelf.SubscScsAsPool
	subscPool[subscriptionId] = &storeDataSubsc

	var storeData []nef_context.SubscriptionScsAsDataContext

	if tempScsAsDataContexts != nil {
		storeData = tempScsAsDataContexts
		storeData = append(storeData, storeDataSubsc)
	} else {
		storeData = []nef_context.SubscriptionScsAsDataContext{}
		storeData = append(storeData, storeDataSubsc)
	}

	// storeData := make(map[string]*nef_context.SubscriptionScsAsDataContext)
	// storeData[subscriptionId] = &storeDataSubsc

	pool := nefSelf.ScsAsPool

	// fmt.Println(" scsas Pool:", pool)
	pool[scsAsId] = &nef_context.ScsAsData{
		ScsAsId:           scsAsId,
		ScsAsDataContexts: storeData,
	}

	nefSelf.SubscriptionIdGenerator += 1 // ToDo: check how to sum up subscription id generator

	/*
		// check data context by scsas id
		checkScsAsDataContexts, err := nefSelf.FindScsAsDataContextsByScsAsId(scsAsId)

		if err != nil {
			fmt.Println(" check scs as data contexts error: ", err)
		}
		// fmt.Println("  check scs as data contexts: ", checkScsAsDataContexts)
		for _, cont := range checkScsAsDataContexts {
			// fmt.Println("scsAsId: ", saId)
			// subscId := cont.SubscriptionId
			// subscCont := cont.ScsAsDataContext
			// fmt.Println("subscId: ", subscId)
			// fmt.Println("subscCont: ", subscCont)
			fmt.Println("check context:", cont)
		}

		// check data context by subscriptionId
		checkSubscScsAsDataContext, err := nefSelf.FindScsAsDataContextByScsAsIdAndSubscriptionId(scsAsId, subscriptionId)
		if err != nil {
			fmt.Println(" check scs as and subscription data context error: ", err)
		}

		fmt.Println("check subscription context:", checkSubscScsAsDataContext)
	*/

	return &response, selfUri, nil
}

// subscription id get
func Handle3GPPAsSessionWithQosSubscriptionsSubscriptionIdGet(scsAsId string, subscriptionId string) *httpwrapper.Response {
	logger.ProducerLog.Infof(" Handle 3GPPAsSessionWithQos Subscriptions SubscriptionId Get")

	response, _, _ := GetQosSubscriptionsSubscriptionIdProcedure(scsAsId, subscriptionId)

	return response
}

func GetQosSubscriptionsSubscriptionIdProcedure(scsAsId string, subscriptionId string) (*httpwrapper.Response, string, *models.ProblemDetails) {
	// get nef context
	nefSelf := nef_context.NEF_Self()

	// find scs as data context from scs as id and subscription id
	scsAsDataContext, err := nefSelf.FindScsAsDataContextByScsAsIdAndSubscriptionId(scsAsId, subscriptionId)

	if err != nil {
		logger.ProducerLog.Errorf("Cannot Find Scs As Data Context by subscription id and scs as id")

		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED",
		}

		response := httpwrapper.NewResponse(http.StatusInternalServerError, nil, problemDetails)

		return response, "", &problemDetails
	}

	if scsAsDataContext == nil {

		logger.ProducerLog.Errorf("Cannot Find Scs As Data Context by subscription id and scs as id")

		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED",
		}

		response := httpwrapper.NewResponse(http.StatusInternalServerError, nil, problemDetails)

		return response, "", &problemDetails

	}

	httpRsp := httpwrapper.NewResponse(http.StatusOK, nil, scsAsDataContext)

	return httpRsp, "", nil
}

// subscription id delete
func Handle3GPPAsSessionWithQosSubscriptionsSubscriptionIdDelete(scsAsId string, subscriptionId string) *httpwrapper.Response {
	problemDetails := Delete3GPPAsSessionWithQosSubscriptionsSubscriptionIdProcedure(scsAsId, subscriptionId)

	if problemDetails == nil {
		return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
	} else {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func Delete3GPPAsSessionWithQosSubscriptionsSubscriptionIdProcedure(scsAsId string, subscriptionId string) *models.ProblemDetails {
	nefSelf := nef_context.NEF_Self()

	// find scs as data context from scs as id and subscription id
	scsAsDataContext, err := nefSelf.FindScsAsDataContextByScsAsIdAndSubscriptionId(scsAsId, subscriptionId)

	// error handling
	if err != nil {
		logger.ProducerLog.Errorf("Cannot Find Scs As Data Context by subscription id and scs as id")

		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "FAILED_SCS_AS_DATA_CONTEXT",
		}

		return &problemDetails

	}

	if scsAsDataContext == nil {
		logger.ProducerLog.Errorf("Cannot Find Scs As Data Context by subscription id and scs as id")

		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "FAILED_SCS_AS_DATA_CONTEXT",
		}

		return &problemDetails

	}

	delete(nefSelf.SubscScsAsPool, subscriptionId)

	return nil
}

// subscription id put
// if you designate the subscription id, you should write specified rule such as ind-1
// This method has route for calling policy authorization api
func Handle3GPPAsSessionWithQosSubscriptionsSubscriptionIdPut(request *httpwrapper.Request) *httpwrapper.Response {
	scsAsId := request.Params["scsAsId"]
	subscriptionId := request.Params["subscriptionId"]
	asSessionWithQosSubscriptionsRequest := request.Body.(models.ScsAsDataContext)

	response, locationHeader, problemDetails := PutQosSubscriptionsSubscriptionIdProcedure(scsAsId, subscriptionId, asSessionWithQosSubscriptionsRequest)

	headers := http.Header{
		"Location": {locationHeader},
	}

	if response != nil {
		return httpwrapper.NewResponse(http.StatusCreated, headers, response)
	} else if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)

}

func PutQosSubscriptionsSubscriptionIdProcedure(scsAsId string, subscriptionId string, req models.ScsAsDataContext) (*models.ScsAsDataContextResponse, string, *models.ProblemDetails) {
	// error handling for request body
	if req.UeIpv4Addr == "" && req.UeIpv6Addr == "" && req.MacAddr == "" {
		problemDetails := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Detail: "Ue UeIpv4 and UeIpv6 and UeMac are all empty",
			Cause:  "ERROR_REQUEST_PARAMETERS",
		}

		return nil, "", &problemDetails
	}

	nefSelf := nef_context.NEF_Self()

	// find scs as data context from scs as id and subscription id
	scsAsDataContext, err := nefSelf.FindScsAsDataContextByScsAsIdAndSubscriptionId(scsAsId, subscriptionId)

	if err != nil {
		logger.ProducerLog.Errorf("Cannot Find Scs As Data Context by subscription id and scs as id")

		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED",
		}

		return nil, "", &problemDetails
	}

	if scsAsDataContext != nil {

		logger.ProducerLog.Errorf("Can Find Scs As Data Context by subscription id and scs as id")

		problemDetails := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "UNSPECIFIED",
		}

		return nil, "", &problemDetails

	}

	// create response for 3gpp-as-session-with-qos api
	// ToDo: get api prefix and service name and change subscription id and scsasid
	selfUri := "http://10.0.10.1:8000/3gpp-as-session-with-qos/v1/{nfId}/subscriptions/{subscriptionId}"

	response := models.ScsAsDataContextResponse{
		Self:                    selfUri,
		FlowInfos:               req.FlowInfos,
		NotificationDestination: req.NotificationDestination,
		EthFlowInfos:            req.EthFlowInfos,
		QosReference:            req.QosReference,
		UeIpv4Addr:              req.UeIpv4Addr,
		UeIpv6Addr:              req.UeIpv6Addr,
		MacAddr:                 req.MacAddr,
		QosMonInfo:              req.QosMonInfo,
		SupportedFeatures:       req.SupportedFeatures,
		Supi:                    req.Supi,
		Dnn:                     req.Dnn,
	}

	storeDataSubsc := nef_context.SubscriptionScsAsDataContext{
		SubscriptionId:   subscriptionId,
		ScsAsId:          scsAsId,
		ScsAsDataContext: &response,
	}

	subscPool := nefSelf.SubscScsAsPool
	subscPool[subscriptionId] = &storeDataSubsc

	nefSelf.SubscriptionIdGenerator += 1 // ToDo: need process?

	/*
		// Problem: change scsasId subscriptions Pool?
		var storeData []nef_context.SubscriptionScsAsDataContext

		if tempScsAsDataContexts != nil {
			storeData = tempScsAsDataContexts
			storeData = append(storeData, storeDataSubsc)
		} else {
			storeData = []nef_context.SubscriptionScsAsDataContext{}
			storeData = append(storeData, storeDataSubsc)
		}
	*/

	return &response, selfUri, nil

}

// subscription id PATCH
func Handle3GPPAsSessionWithQosSubscriptionsSubscriptionIdPatch(request *httpwrapper.Request) *httpwrapper.Response {
	scsAsId := request.Params["scsAsId"]
	subscriptionId := request.Params["subscriptionId"]
	asSessionWithQosSubscriptionsRequest := request.Body.(models.ScsAsDataContext)

	response, locationHeader, problemDetails := PatchQosSubscriptionsSubscriptionIdProcedure(scsAsId, subscriptionId, asSessionWithQosSubscriptionsRequest)

	headers := http.Header{
		"Location": {locationHeader},
	}

	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, headers, response)
	} else if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)

}

func PatchQosSubscriptionsSubscriptionIdProcedure(scsAsId string, subscriptionId string, req models.ScsAsDataContext) (*models.ScsAsDataContextResponse, string, *models.ProblemDetails) {
	// error handling for request body
	if req.UeIpv4Addr == "" && req.UeIpv6Addr == "" && req.MacAddr == "" {
		problemDetails := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Detail: "Ue UeIpv4 and UeIpv6 and UeMac are all empty",
			Cause:  "ERROR_REQUEST_PARAMETERS",
		}

		return nil, "", &problemDetails
	}

	nefSelf := nef_context.NEF_Self()

	// find scs as data context from scs as id and subscription id
	scsAsDataContext, err := nefSelf.FindScsAsDataContextByScsAsIdAndSubscriptionId(scsAsId, subscriptionId)

	if err != nil {
		logger.ProducerLog.Errorf("Cannot Find Scs As Data Context by subscription id and scs as id")

		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED",
		}

		return nil, "", &problemDetails
	}

	if scsAsDataContext == nil {

		logger.ProducerLog.Errorf("Cannot modify Scs As Data Context by subscription id and scs as id")

		problemDetails := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "UNSPECIFIED",
		}

		return nil, "", &problemDetails

	}

	// create response for 3gpp-as-session-with-qos api
	// ToDo: get api prefix and service name and change subscription id and scsasid
	selfUri := "http://10.0.10.1:8000/3gpp-as-session-with-qos/v1/{nfId}/subscriptions/{subscriptionId}"

	// modify store data
	response := models.ScsAsDataContextResponse{
		Self:                    selfUri,
		FlowInfos:               req.FlowInfos,
		NotificationDestination: req.NotificationDestination,
		EthFlowInfos:            req.EthFlowInfos,
		QosReference:            req.QosReference,
		UeIpv4Addr:              req.UeIpv4Addr,
		UeIpv6Addr:              req.UeIpv6Addr,
		MacAddr:                 req.MacAddr,
		QosMonInfo:              req.QosMonInfo,
		SupportedFeatures:       req.SupportedFeatures,
		Supi:                    req.Supi,
		Dnn:                     req.Dnn,
	}

	storeDataSubsc := nef_context.SubscriptionScsAsDataContext{
		SubscriptionId:   subscriptionId,
		ScsAsId:          scsAsId,
		ScsAsDataContext: &response,
	}

	subscPool := nefSelf.SubscScsAsPool
	subscPool[subscriptionId] = &storeDataSubsc

	// need scsaspool?

	return &response, selfUri, nil

}
