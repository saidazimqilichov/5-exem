package main

import (
	"notif-service/kafka"
)

func main() {
	kafka.UserConsumer("broker:29092", "user-registration")
}
