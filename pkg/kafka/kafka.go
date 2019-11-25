package kafka_operator

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/wuxiaoxiaoshen/go-anything/configs"
)

type (
	kafkaAction struct {
		producer sarama.AsyncProducer
	}
	kafkaSettings struct {
		broker        string
		topic         string
		consumerGroup string
	}
)

var (
	DefaultAsyncProducer kafkaAction
	Topic                string
	settings             kafkaSettings
)

func init() {
	r := configs.DefaultConfigs.LoadConfigs("kafka")
	a := r.(map[string]interface{})
	log.Println(fmt.Sprintf("Keys: Kafka: %#v", r))
	settings = kafkaSettings{
		broker:        a["broker"].(string),
		topic:         a["topic"].(string),
		consumerGroup: a["consumergroup"].(string),
	}

}

func KafkaInit() {
	broker := settings.broker
	Topic = settings.topic
	DefaultAsyncProducer.producer = newProducer([]string{broker})
}

func newProducer(broker []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Version = sarama.V2_2_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	producer, err := sarama.NewAsyncProducer(broker, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer1:", err)
	}
	return producer
}

func (K *kafkaAction) Close() {
	defer K.producer.AsyncClose()
}

func (K *kafkaAction) Run(topic string, v interface{}) {
	message := v.(Message)
	log.Println(fmt.Sprintf("KafkaAction: Send Message: %+v", message))
	K.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: &message,
	}
}
