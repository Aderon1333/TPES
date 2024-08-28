package producer

import (
	"encoding/json"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"

	"github.com/Aderon1333/TPES/internal/models"
)

type ProducerKafka struct {
	prod sarama.SyncProducer
}

func NewProducer(prod sarama.SyncProducer) ProducerKafka {
	return ProducerKafka{prod: prod}
}

func (p *ProducerKafka) PushReqToQueue(c *gin.Context, topic string, body []byte) error {
	// Create new Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(body),
	}
	//Send message
	_, _, err := p.prod.SendMessage(msg)
	if err != nil {
		return err
	}
	c.JSON(200, gin.H{
		"Info": "",
	})
	return nil
}

func (p *ProducerKafka) PlaceReq(c *gin.Context) {
	// Parse request body
	var (
		taskDecode models.Task
	)

	err := json.NewDecoder(c.Request.Body).Decode(&taskDecode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to decode request body",
		})
		return
	}
	// convert req body into bytes
	taskInBytes, err := json.Marshal(taskDecode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to convert order into bytes",
		})
		return
	}

	// send bytes to Kafka
	err = p.PushReqToQueue(c, "tasks", taskInBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to push task to queue",
		})
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"INFO": c.Request.Method + " method successfully pushed to Kafka queue",
	})

}
