package context

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/softmurata/nef/internal/logger"

	"github.com/softmurata/freeopenapi/models"
	"github.com/softmurata/nef/pkg/factory"
)

var nefContext = NEFContext{}

func init() {
	nefContext.NfId = uuid.New().String()

	nefContext.Name = "nef"

	nefContext.UriScheme = models.UriScheme_HTTPS
	nefContext.RegisterIPv4 = factory.NEF_DEFAULT_IPV4
	nefContext.SBIPort = factory.NEF_DEFAULT_PORT_INT

	serviceName := []models.ServiceName{
		models.ServiceName_3GPP_AS_SESSION_WITH_QOS,
		models.ServiceName_3GPP_SERVICE_PARAMETER,
	}

	nefContext.NfService = initNfService(serviceName, "1.0.0")

	nefContext.NrfUri = fmt.Sprintf("%s://%s:%d", models.UriScheme_HTTPS, nefContext.RegisterIPv4, nefContext.SBIPort)

}

type NEFContext struct {
	NfId         string
	Name         string
	UriScheme    models.UriScheme
	RegisterIPv4 string
	BindingIPv4  string
	SBIPort      int
	NfService    map[models.ServiceName]models.NfService
	NrfUri       string
}

// Initialize NEF context with configuration factory
func InitNefContext() {
	if !factory.Configured {
		logger.ContextLog.Warnf("NEF is not configured")
		return
	}
	nefConfig := factory.NefConfig

	if nefConfig.Configuration.NefName != "" {
		nefContext.Name = nefConfig.Configuration.NefName
	}

	nefContext.UriScheme = nefConfig.Configuration.Sbi.Scheme
	nefContext.RegisterIPv4 = nefConfig.Configuration.Sbi.RegisterIPv4
	nefContext.SBIPort = nefConfig.Configuration.Sbi.Port
	nefContext.BindingIPv4 = os.Getenv(nefConfig.Configuration.Sbi.BindingIPv4)

	if nefContext.BindingIPv4 != "" {
		logger.ContextLog.Info("Parsing ServerIPv4 address from ENV Variable.")
	} else {
		nefContext.BindingIPv4 = nefConfig.Configuration.Sbi.BindingIPv4
		if nefContext.BindingIPv4 == "" {
			logger.ContextLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
			nefContext.BindingIPv4 = "0.0.0.0"
		}
	}

	nefContext.NfService = initNfService(nefConfig.Configuration.ServiceNameList, nefConfig.Info.Version)

	if nefConfig.Configuration.NrfUri != "" {
		nefContext.NrfUri = nefConfig.Configuration.NrfUri
	} else {
		logger.InitLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		nefContext.NrfUri = fmt.Sprintf("%s://%s:%d", nefContext.UriScheme, "127.0.0.1", 29510)
	}

}

func initNfService(serviceName []models.ServiceName, version string) (nfService map[models.ServiceName]models.NfService) {
	versionUri := "v" + strings.Split(version, ".")[0]
	nfService = make(map[models.ServiceName]models.NfService)
	for idx, name := range serviceName {
		nfService[name] = models.NfService{
			ServiceInstanceId: strconv.Itoa(idx),
			ServiceName:       name,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          nefContext.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       GetIpv4Uri(),
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: nefContext.RegisterIPv4,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(nefContext.SBIPort),
				},
			},
		}
	}

	return
}

func GetIpv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", nefContext.UriScheme, nefContext.RegisterIPv4, nefContext.SBIPort)
}

func NEF_Self() *NEFContext {
	return &nefContext
}
