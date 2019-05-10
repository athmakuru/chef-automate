package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	secrets "github.com/chef/automate/api/external/secrets"
	cfgmgmtRequest "github.com/chef/automate/api/interservice/cfgmgmt/request"
	cfgmgmtResponse "github.com/chef/automate/api/interservice/cfgmgmt/response"
	cfgmgmt "github.com/chef/automate/api/interservice/cfgmgmt/service"
	"github.com/chef/automate/components/data-feed-service/config"
	notifications "github.com/chef/automate/components/notifications-client/api"
	"github.com/chef/automate/lib/grpc/secureconn"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type assetNotification struct {
	username   string
	password   string
	url        string
	attributes []attributesMessage
}

type attributesMessage struct {
	//Automatic attributes               `json:"automatic"`
	Automatic map[string]interface{} 	`json:"automatic"`
	LastRun   *cfgmgmtResponse.Run `json:"last_run"`
}


type chef struct {
	ChefRoot string `json:"chef_root"`
	Version string `json:"version"`
}

type ohai struct {
	OhaiRoot string `json:"ohai_root"`
	Version string `json:"version"`
}

type ChefPackages struct {
	Chef chef `json:"chef"`
	Ohai ohai `json:"ohai"`
}

type attributes struct {
	ChefEnvironment string `json:"chef_environment"`
	ChefGuid string `json:"chef_guid"`
	//ChefPackages ChefPackages `json:"chef_packages"`
	// Cloud
	// Command
	// Cookbooks
	Cookbooks map[string]interface{} `json:"cookbooks"`
	// Counters
	// Cpu
	CurrentUser string `json:"current_user"`
	// Dmi
	Domain string `json:"domain"`
	// cloud provider e.g. ec2
	// etc - possibly linux only
	// Expanded run list
	// filesystem - differs by os, what is common?
	// fips
	Fqdn string `json:"fqdn"`
	Hostname string `json:"hostname"`
	IpAddress string `json:"ipaddress"`
	Ip6Address string `json:"ip6address"`
	// json?
	// kernel differs by os
	// keys
	// languages format differ by os
	MacAddress string `json:"macaddress"`
	MachineId string `json:"machine_id"`
	MachineName string `json:"machinename"`
	// memory differes by os version
	Name string `json:"name"`
	// Network differs by os
	// nginx? if present would we bother?
	OhaiTime float64 `json:"ohai_time"`
	Os string `json:"os"`
	OsVersion string `json:"os_version"`
	// packages
	Platform string `json:"platform"`
	PlatformFamily string `json:"platform_family"`
	PlatformVersion string `json:"platform_version"`
	// recipes
	// roles
	RootGroup string `json:"root_group"`
	ShardSeed float64 `json:"shard_seed"`
	// tags
	// time
	Uptime string `json:"uptime"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	// virtualization

	// known os specifics
	// Ubuntu
	// shells
	// site
	// sysconf
	// sytemd_paths

	// windows
	// system_enclosure

}

type serviceClients struct {
	notifications notifications.NotificationsClient
	cfgMgmt       cfgmgmt.CfgMgmtClient
	secrets       secrets.SecretsServiceClient
}

func getConnection(connectionFactory *secureconn.Factory, service string, target string) *grpc.ClientConn {
	connection, err := connectionFactory.Dial(service, target, grpc.WithBlock())
	if err != nil {
		log.Fatalf("Error getting connection to %v: %v", service, err)
	}
	return connection
}

func initServiceClients(dataFeedConfig *config.DataFeedConfig) *serviceClients {
	clients := &serviceClients{}
	connectionFactory := secureconn.NewFactory(*dataFeedConfig.ServiceCerts)
	connection := getConnection(connectionFactory, "notifications-service", dataFeedConfig.NotificationsConfig.Target)
	clients.notifications = notifications.NewNotificationsClient(connection)
	log.Debugf("NotificationsClient created")

	connection = getConnection(connectionFactory, "config-mgmt-service", dataFeedConfig.CfgmgmtConfig.Target)
	clients.cfgMgmt = cfgmgmt.NewCfgMgmtClient(connection)
	log.Debugf("CfgMgmtClient created")

	connection = getConnection(connectionFactory, "secrets-service", dataFeedConfig.SecretsConfig.Target)
	clients.secrets = secrets.NewSecretsServiceClient(connection)
	log.Debugf("Secrets created")
	return clients
}

func Start(dataFeedConfig *config.DataFeedConfig) {
	log.Debugf("data-feed-service start")

	serviceClients := initServiceClients(dataFeedConfig)
	var retry = 0
	for {
		listResponse, err := serviceClients.notifications.ListRules(context.Background(), &notifications.Empty{})
		if err != nil {
			// this is currently inadequate, retry period must be calculated so it does not overlap with the interval
			// bring the service down if we can't get a response in 50 seconds since the smallest interval we allow is 1 minute
			if retry == 5 {
				log.Fatalf("Error calling list rules on notifications-service attempt %v: %v", retry, err)
			}
			log.Errorf("Error calling list rules on notifications-service attempt %v: %v", retry, err)
			retry++
			wait, _ := time.ParseDuration("10s")
			time.Sleep(wait)
			continue
		}
		responseCode := listResponse.Code
		log.Debugf("ListRules responseCode %v", responseCode)
		messages := listResponse.Messages
		// get the rules here - service now rules only
		log.Debugf("ListRules messages %v", messages)

		assetRules := make([]notifications.Rule, 0)
		rules := listResponse.Rules
		for rule := range rules {
			event := rules[rule].Event
			if event == notifications.Rule_Assets {
				assetRules = append(assetRules, *rules[rule])
			}
		}
		now := time.Now()
		feedStartTime, feedEndTime := getFeedTimes(dataFeedConfig, now)
		if len(assetRules) > 0 {
			attributes := buildDatafeed(serviceClients, dataFeedConfig, feedStartTime, feedEndTime)
			if len(attributes) > 0 {
				for rule := range assetRules {
					log.Debugf("Rule id %v", assetRules[rule].Id)
					log.Debugf("Rule Name %v", assetRules[rule].Name)
					log.Debugf("Rule Event %v", assetRules[rule].Event)
					log.Debugf("Rule Action %v", assetRules[rule].Action)
					switch action := rules[rule].Action.(type) {
					case *notifications.Rule_ServiceNowAlert:
						log.Debugf("Service now alert URL  %v", action.ServiceNowAlert.Url)
						log.Debugf("Service now alert secret  %v", action.ServiceNowAlert.SecretId)
						username, password, err := getCredentials(serviceClients, dataFeedConfig, action.ServiceNowAlert.SecretId)
						if err != nil {
							log.Errorf("Error retrieving credentials, cannot send asset notification: %v", err)
						} else {
							dataFeedNotification := &assetNotification{username: username, password: password, url: action.ServiceNowAlert.Url, attributes: attributes}
							sendNotification(dataFeedConfig, dataFeedNotification)
						}
					}
				}
			}
		}

		waitForInterval(dataFeedConfig.ServiceConfig.FeedInterval, feedEndTime, now)
	}
}

func waitForInterval(feedInterval time.Duration, feedEndTime time.Time, now time.Time) {

	nextFeedEndTime := feedEndTime.Add(feedInterval)
	/*
	 * Calculate the sleep duration as the first sleep may be in the middle of an interval.
	 * E.g. for first wait now may be 1:55pm, waiting an hour to get data from a rounded
	 * down window of 12pm-1pm will introduce a 55min lag by waking at 2:55pm
	 *
	 * Calulating the sleepDuration will reduce the value
	 * e.g nextFeedEndTime 2pm - 1.55pm = 5min sleep
	 * waking at 2pm rather than 2:55pm to get 1pm-2pm data witout lag.
	 */
	sleepDuration := nextFeedEndTime.Sub(now)
	log.Debugf("Sleeping for %v", sleepDuration)
	time.Sleep(sleepDuration)
}

func getFeedTimes(dataFeedConfig *config.DataFeedConfig, now time.Time) (time.Time, time.Time) {
	var feedEndTime time.Time
	var feedStartTime time.Time

	feedInterval := dataFeedConfig.ServiceConfig.FeedInterval

	feedEndTime = getFeedEndTime(feedInterval, now)
	feedStartTime = feedEndTime.Add(-feedInterval)

	log.Debugf("Feed interval start, end: %s, %s", feedStartTime.Format("15:04:05"), feedEndTime.Format("15:04:05"))
	return feedStartTime, feedEndTime
}

func getFeedEndTime(feedInterval time.Duration, now time.Time) time.Time {
	log.Debugf("Time Now: %s", now.Format("15:04:05"))
	/*
	 * We round the current time down based on the interval duration to get the end of the last interval.
	 *
	 * Round will round either up or down to the nearest value of the inteval duration.
	 * e.g 1:20pm rounds to 1pm, 1:40pm rounds to 2pm
	 *
	 * If we have rounded down that will be our feed end time. The end of a clock interval
	 * rather than current time e.g. 1pm
	 */
	feedEndTime := now.Round(feedInterval)
	log.Debugf("FeedInterval/Units: %s", feedInterval)
	/*
	 * If we have rounded up we subtract the interval to effectively round down
	 */
	if feedEndTime.After(now) {
		feedEndTime = feedEndTime.Add(-feedInterval)
		log.Debugf("feedEndTime after: %s", feedEndTime.Format("15:04:05"))
	}
	log.Debugf("feedEndTime: %s", feedEndTime.Format("15:04:05"))
	return feedEndTime
}

func getCredentials(serviceClients *serviceClients, dataFeedConfig *config.DataFeedConfig, secretId string) (string, string, error) {

	secret, err := serviceClients.secrets.Read(context.Background(), &secrets.Id{Id: secretId})

	username := ""
	password := ""
	if err == nil {
		data := secret.GetData()

		for kv := range data {
			if data[kv].Key == "username" {
				username = data[kv].Value
			} else if data[kv].Key == "password" {
				password = data[kv].Value
			}
		}
	}

	return username, password, err
}

func sendNotification(dataFeedConfig *config.DataFeedConfig, notification *assetNotification) {

	messageBytes, err := json.Marshal(notification.attributes)
	log.Debugf("sendNotification bytes length %v", len(messageBytes))

	if err != nil {
		log.Errorf("Error creating json bytes %v", err)
	}
	request, err := http.NewRequest("POST", notification.url, bytes.NewBuffer(messageBytes))
	if err != nil {
		log.Error("Error creating request")
	}
	request.SetBasicAuth(notification.username, notification.password)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	log.Debugf("request %v", request)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Errorf("Error sending test message %v", err)
	} else {
		log.Infof("Asset data posted to %v, Status %v", notification.url, response.Status)
	}
}

func buildDatafeed(serviceClients *serviceClients, dataFeedConfig *config.DataFeedConfig, feedStartTime time.Time, feedEndTime time.Time) []attributesMessage {

	feedStart := float64(feedStartTime.Unix())
	feedEnd := float64(feedEndTime.Unix())
	log.Info("Building data feed...")

	nodeFilters := &cfgmgmtRequest.Nodes{}
	nodes, err := serviceClients.cfgMgmt.GetNodes(context.Background(), nodeFilters)
	if err != nil {
		log.Errorf("Error getting cfgmgmt/nodes %v", err)
	}

	messages := make([]attributesMessage, 0)

	for _, node := range nodes.Values {
		nodeStruct := node.GetStructValue()
		id := nodeStruct.Fields["id"].GetStringValue()
		lastRunId := nodeStruct.Fields["latest_run_id"].GetStringValue()
		nodeAttributes, err := serviceClients.cfgMgmt.GetAttributes(context.Background(), &cfgmgmtRequest.Node{NodeId: id})
		if err != nil {
			log.Errorf("Error getting attributes %v", err)
		}
		var automaticJson map[string]interface{}
		err = json.Unmarshal([]byte(nodeAttributes.Automatic), &automaticJson)
		if err != nil {
			log.Errorf("Could not parse automatic attributes from json: %v", err)
			// continue to the next node
			continue
		}
		attributesJson := buildDynamicJson(automaticJson)
		ohaiTime := automaticJson["ohai_time"].(float64)
		log.Debugf("feedStartTime %v, feedEndTime %v, ohai_time %v", feedStart, feedEnd, ohaiTime)
		if ohaiTime > feedStart && ohaiTime < feedEnd {
			// get the latest node run
			lastRun, error := serviceClients.cfgMgmt.GetNodeRun(context.Background(), &cfgmgmtRequest.NodeRun{NodeId: id, RunId: lastRunId})
			// check time? for paranoia's sake?
			if error != nil {
				log.Errorf("Error getting node run %v", error)
			} else {
				log.Debugf("Last run\n %v", lastRun)
			}

			message := attributesMessage{Automatic: attributesJson, LastRun: lastRun}
			log.Debugf("Message: %v", message)
			messages = append(messages, message)
		}
	}
	log.Debugf("%v node attribute messages retrieved in interval", len(messages))
	return messages
}

func buildAttributesJson(automaticJson map[string]interface{}) attributes {

	attributesJson := attributes{
		ChefEnvironment: automaticJson["chef_environment"].(string),
		ChefGuid: automaticJson["chef_guid"].(string),
		Cookbooks: automaticJson["cookbooks"].(map[string]interface{}),
		//ChefPackages: automaticJson["chef_packages"],
		CurrentUser: getSafeString(automaticJson, "current_user"),
		Domain: getSafeString(automaticJson, "domain"),
		Fqdn: automaticJson["fqdn"].(string),
		Hostname: automaticJson["hostname"].(string),
		IpAddress: automaticJson["ipaddress"].(string),
		Ip6Address: automaticJson["ip6address"].(string),
		MacAddress: automaticJson["macaddress"].(string),
		MachineId: getSafeString(automaticJson, "machine_id"),
		MachineName: automaticJson["machinename"].(string),
		Name: automaticJson["name"].(string),
		OhaiTime: automaticJson["ohai_time"].(float64),
		Os: automaticJson["os"].(string),
		OsVersion: automaticJson["os_version"].(string),
		Platform: automaticJson["platform"].(string),
		PlatformFamily: automaticJson["platform_family"].(string),
		PlatformVersion: automaticJson["platform_version"].(string),
		RootGroup: automaticJson["root_group"].(string),
		ShardSeed: automaticJson["shard_seed"].(float64),
		Uptime: automaticJson["uptime"].(string),
		UptimeSeconds: automaticJson["uptime_seconds"].(float64)}

	return attributesJson
}

func buildDynamicJson(automaticJson map[string]interface{}) map[string]interface{} {
	// check what format he json is e.g. ubuntu/ windows
	// remove specific sections
	// insert normalised
	var dynamicJson map[string]interface{}
	if automaticJson["lsb"] != nil {
		dynamicJson = buildUbuntuJson(automaticJson)
	} else if automaticJson["os"] == "windows" {
		dynamicJson = buildWindowsJson(automaticJson)
	}

	return dynamicJson
}


func buildUbuntuJson(ubuntuJson map[string]interface{}) map[string]interface{} {
	// check what format he json is e.g. ubuntu/ windows
	// remove specific sections
	// insert normalised

	lsb := ubuntuJson["lsb"].(map[string]interface{})
	ubuntuJson["description"] = lsb["description"]

	dmi := ubuntuJson["dmi"].(map[string]interface{})
	system := dmi["system"].(map[string]interface{})
	ubuntuJson["serial_number"] = system["serial_number"]

	//delete(ubuntuJson, "lsb")
	//delete(ubuntuJson, "system")

	return ubuntuJson
}


func buildWindowsJson(windowsJson map[string]interface{}) map[string]interface{} {
	// check what format he json is e.g. ubuntu/ windows
	// remove specific sections
	// insert normalised

	kernel := windowsJson["kernel"].(map[string]interface{})
	windowsJson["description"] = kernel["name"]
	osInfo := kernel["os_info"].(map[string]interface{})
	windowsJson["serial_number"] = osInfo["serial_number"]
	servicePackMajorVersion := fmt.Sprintf("%g", osInfo["service_pack_major_version"].(float64))
	servicePackMinorVersion := fmt.Sprintf("%g", osInfo["service_pack_minor_version"].(float64))
	servicePack := strings.Join([]string{servicePackMajorVersion, servicePackMinorVersion}, ".")
	windowsJson["os_service_pack"] = servicePack

	//delete(windowsJson, "kernel")

	return windowsJson
}

func getSafeString(automaticJson map[string]interface{}, s string) string {
	if automaticJson[s] == nil {
		return ""
	}
	return automaticJson[s].(string)
}