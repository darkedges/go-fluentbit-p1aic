package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/fluent/fluent-bit-go/input"
	"github.com/imroc/req/v3"
)

// Values needed to make this work
var previousBeginTime string
var beginTime string
var parsedTime time.Time
var forgerockIdCloudTenant string
var apiKeyId string
var apiKeySecret string
var logSources string
var logFilter string
var client *req.Client

const MonitoringLogsTemplate = "%s/monitoring/logs"

// This holds the result from the log request
type MonitoringLogsResponse struct {
	Result []struct {
		Payload   map[string]interface{} `json:"payload"`
		Source    string                 `json:"source"`
		Timestamp time.Time              `json:"timestamp"`
		Type      string                 `json:"type"`
	} `json:"result"`
	ResultCount             int    `json:"resultCount"`
	PagedResultsCookie      string `json:"pagedResultsCookie"`
	TotalPagedResultsPolicy string `json:"totalPagedResultsPolicy"`
	TotalPagedResults       int    `json:"totalPagedResults"`
	RemainingPagedResults   int    `json:"remainingPagedResults"`
}

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return input.FLBPluginRegister(def, "p1aic", "p1aic GO!")
}

// (fluentbit will call this)
// plugin (context) pointer to fluentbit context (state/ c code)
//
//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	// Example to retrieve an optional configuration parameter
	forgerockIdCloudTenant = input.FLBPluginConfigKey(plugin, "p1aic_id_cloud_tenant")
	if strings.HasPrefix(forgerockIdCloudTenant, "http://") {
		forgerockIdCloudTenant = strings.Replace(forgerockIdCloudTenant, "http://", "https://", 1)
	}
	apiKeyId = input.FLBPluginConfigKey(plugin, "api_key_id")
	apiKeySecret = input.FLBPluginConfigKey(plugin, "api_key_secret")
	logSources = input.FLBPluginConfigKey(plugin, "log_sources")
	logFilter = input.FLBPluginConfigKey(plugin, "log_filter")
	if logSources == "" {
		logSources = "am-authentication,am-access,am-config,idm-activity"
	}

	// Attempt to parse previousBeginTime
	parsedTime, err := time.Parse("2006-01-02T15:04:05Z", previousBeginTime)
	previousBeginTime = parsedTime.String()
	if err != nil {
		previousBeginTime = ""
	}
	if previousBeginTime == "" {
		fmt.Println("No previous beginTime saved so backdating to 1 minute ago")
		beginTime = time.Now().UTC().Add(-120 * time.Second).Format("2006-01-02T15:04:05Z")
	} else if time.Since(parsedTime).Seconds() > 3600 {
		fmt.Printf("Previous saved beginTime %s too old (> 1 hour) so backdating to 1 minute ago\n", parsedTime.Format("2006-01-02T15:04:05Z"))
		beginTime = time.Now().UTC().Add(-120 * time.Second).Format("2006-01-02T15:04:05Z")
	} else {
		beginTime = previousBeginTime
	}
	client = req.C().
		SetUserAgent("go-fluentbit-p1aic")
	return input.FLB_OK
}

//export FLBPluginInputCallback
func FLBPluginInputCallback(data *unsafe.Pointer, size *C.size_t) int {
	// Calculate the end time
	endTime := time.Now().UTC().Add(-60 * time.Second).Format("2006-01-02T15:04:05Z")
	//Store it for future lookup
	saveCheckPoint("beginTime", endTime)

	buf := bytes.NewBuffer([]byte{})
	pagedResultsCookie := ""
	urlString := fmt.Sprintf(MonitoringLogsTemplate, forgerockIdCloudTenant)
	//Loop through the results
	for {
		// Configure the Client Request
		queryParams := map[string]string{
			"_pageSize": "10",
			"source":    logSources,
			"beginTime": beginTime,
			"endTime":   endTime,
		}
		if pagedResultsCookie != "" {
			queryParams["_pagedResultsCookie"] = pagedResultsCookie
		}
		if logFilter != "" {
			queryParams["_queryFilter"] = logFilter
		}
		var monitoringLogsResponse MonitoringLogsResponse

		// Get the log results
		resp, err := client.R().
			SetQueryParams(queryParams).
			SetSuccessResult(&monitoringLogsResponse).
			SetHeaders(map[string]string{
				"x-api-key":    apiKeyId,
				"x-api-secret": apiKeySecret,
			}).
			Get(urlString)
		if err != nil { // Error handling.
			return input.FLB_ERROR
		}
		if resp.IsSuccessState() { // Status code is between 200 and 299.
			pagedResultsCookie = monitoringLogsResponse.PagedResultsCookie

			for _, element := range monitoringLogsResponse.Result {
				flb_time := input.FLBTime{element.Timestamp}
				entry := []interface{}{flb_time, element.Payload}
				enc := input.NewEncoder()
				packed, err := enc.Encode(entry)
				if err != nil {
					fmt.Fprintf(os.Stderr, "msgpack marshal: %s\n", err)
					return input.FLB_ERROR
				}
				buf.Grow(len(packed))
				buf.Write(packed)
			}
			// No more entries to read so break if that is the case
			if pagedResultsCookie == "" {
				break
			}
		}
	}
	if buf.Len() > 0 {
		b := buf.Bytes()
		cdata := C.CBytes(b)
		*data = cdata
		if size != nil {
			*size = C.size_t(len(b))
		}
	}
	// Update beginTime to be the current endTime
	beginTime = endTime
	// For emitting interval adjustment.
	// todo Need to work out how to do an Interval, b ut emit straightaway
	time.Sleep(60 * time.Second)

	return input.FLB_OK
}

//export FLBPluginInputCleanupCallback
func FLBPluginInputCleanupCallback(data unsafe.Pointer) int {
	return input.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return input.FLB_OK
}

func saveCheckPoint(key string, value string) {
	// todo
}

func main() {
}