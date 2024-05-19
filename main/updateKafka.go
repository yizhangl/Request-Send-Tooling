package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

var (
	brokers = make([]string, 1)
	config  = sarama.NewConfig()
	admin   sarama.ClusterAdmin
)

// Credentials
func buildEnv(env string) {
	if env == "cs-dev" {
		//admin, _ = sarama.NewClusterAdmin(brokers, config)
	}
	if env != "cs-dev" && env != "cs-sit" {
		panic("build env error")
	}
}

func main() {

	buildEnv("cs-dev")
	// get info of Kafka topic
	getDetails("")

	defer admin.Close()
}

func create(topicName, mmb string) {
	err := admin.CreateTopic(topicName, &detail, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("create topic ", topicName, " successfully.")

}

func getDetails(topicName string) {
	topics, _ := admin.ListTopics()
	for k, topic := range topics {

		if k == topicName {
			fmt.Println("topic name: ", k)

			cfg := topic.ConfigEntries
			for cfgK, cfgV := range cfg {
				if cfgK == "max.message.bytes" {
					fmt.Println("max.message.bytes: ", *cfgV, " bytes")
				}
			}
		}
	}
}

func update(topicName, mmb string) {
	err := admin.AlterConfig(sarama.TopicResource, topicName, configEntries, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(topicName, " Topic configuration updated successfully.")
}

func delete(topicName string) {
	err := admin.DeleteTopic(topicName)
	if err != nil {
		panic(err)
	}

	fmt.Println("delete topic ", topicName, " successfully.")
}