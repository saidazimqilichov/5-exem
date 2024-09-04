package kafka

import (
	"context"
	"fmt"
	"user-service/protos/genuser"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/encoding/protojson"
)

type Kafka struct {
	Client *kgo.Client
}

func ConnectKafka(kurl string) (*Kafka, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(kurl),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed kafka Client error: %v", err)
	}

	err = client.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed kafka connection: %v", err)
	}
	return &Kafka{Client: client}, nil
}

func (k *Kafka) ProduceRegistrationEmail(user *genuser.NewUser) error {
	data, err := protojson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal userinfo in KAFKA: %v", err)
	}

	record := kgo.Record{
		Topic: "user-registration",
		Value: data,
	}
	err = k.Client.ProduceSync(context.Background(), &record).FirstErr()
	if err != nil {
		return fmt.Errorf("failed to produce message on Registration: %v", err)
	}
	return nil
}
