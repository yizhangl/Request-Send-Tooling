package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"time"
)

func main() {
	buildEnv("cs-dev")
	syncProducer()
}

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

// data prep
func syncProducer() {
	config.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Println("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()

	topic := "target_topic"
	value :=  // json file
	listForProduct := []string{"A", "B", "C"}
	listForDataSource := []string{"a", "b", "c"}

	sendMockData(listForProduct, listForDataSource, topic, value, p, "9034")

	value2 :=  // json file
	listForProduct2 := []string{"A"}
	listForDataSource2 := []string{"a"}

	sendMockData(listForProduct2, listForDataSource2, topic, value2, p, "7009")

	value3 :=  // json file
	listForProduct3 := []string{"A", "B"}
	listForDataSource3 := []string{"a", "b"}

	sendMockData(listForProduct3, listForDataSource3, topic, value3, p, "9034")
}

func sendMockData(listForProduct []string, listForDataSource []string, topic string, value string, p sarama.SyncProducer, tenant string) {
	for _, product := range listForProduct {
		for _, dataSource := range listForDataSource {
			//value := fmt.Sprintf(srcValue, i)

			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(value),
				Headers: []sarama.RecordHeader{
					{Key: []byte("Product"), Value: []byte(product)},
					{Key: []byte("Datasource"), Value: []byte(dataSource)},
					{Key: []byte("TenantID"), Value: []byte(tenant)},
				},
			}
			part, offset, err := p.SendMessage(msg)
			if err != nil {
				log.Printf("send message(%s) err=%s \n", value, err)
			} else {
				fmt.Fprintf(os.Stdout, value+"Send request successfully，partition=%d, offset=%d \n", part, offset)
			}
			//time.Sleep(2 * time.Second)
		}
	}
}
