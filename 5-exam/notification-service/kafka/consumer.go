package kafka

import (
	"context"
	"fmt"
	"log"

	"notif-service/email"
	"notif-service/protos/protos"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/encoding/protojson"
)

func UserConsumer(kurl, topic string) {
	client, err := ConnectKafka(kurl, topic)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("User-Registration: started Consuming messages.....")
	for {
		fetches := client.PollFetches(context.Background())
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatal(errs)
		}
		fetches.EachPartition(func(ftp kgo.FetchTopicPartition) {
			for _, record := range ftp.Records {
			
				var s protos.NewUser
				err := protojson.Unmarshal(record.Value, &s)
				if err != nil {
					log.Fatalf("failed to unmarshal user information: %v", err)
				}
				email.SendEmail(s.Email, s.Name, s.Email)
				log.Println("user malumot", s.Email, s.Name, s.UserId)
			}
		})
	}
}

func ConnectKafka(kurl string, topic string) (*kgo.Client, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(kurl),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup("my-group"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed consumer client:%v", err)
	}
	return client, nil
}
