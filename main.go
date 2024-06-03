package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/segmentio/kafka-go"
)

// Define the Kafka writer with your broker addresses and topic
var kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
    Brokers: []string{"broker1:9092", "broker2:9092", "broker3:9092"},
    Topic:   "your-topic",
})

// SampleRecord represents the structure of your Kafka message
type SampleRecord struct {
    AgentId        string  `json:"agentId"`
    CustomerId     float64 `json:"customerId"` // Use float64 for numbers in Go
    MessageId      int     `json:"messageId"`
    MessageText    string  `json:"messageText"`
    SentimentScore string  `json:"sentimentScore"`
}

// Handler for your API that takes the message and sentiment
func evaluateHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the incoming JSON request
    var record SampleRecord
    if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Convert the record to a JSON string
    recordBytes, err := json.Marshal(record)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Create a Kafka message
    msg := kafka.Message{
        Key:   []byte(fmt.Sprintf("%v", record.MessageId)), // Use MessageId as the key
        Value: recordBytes,
    }

    // Write the message using the Kafka writer
    if err := kafkaWriter.WriteMessages(context.Background(), msg); err != nil {
        http.Error(w, fmt.Sprintf("Failed to write message: %s", err), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "Message sent successfully")
}

func main() {
    // Set up your HTTP server and route
    http.HandleFunc("/evaluate", evaluateHandler)
    http.ListenAndServe(":8080", nil)
}
