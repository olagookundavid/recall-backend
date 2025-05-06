package notification

import (
	"context"
	"fmt"
	"recall-app/internal/logger"

	"firebase.google.com/go/v4/messaging"
)

func SendNotification(
	fcmClient *messaging.Client, ctx context.Context, tokens []string, userId string, notification messaging.Notification) error {
	//Send to One Token
	// resp, err := fcmClient.Send(ctx, &messaging.Message{
	// 	Token: tokens[0],
	// 	Data: map[string]string{
	// 		message: message,
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }

	//Send to Multiple Tokens
	res, err := fcmClient.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Notification: &notification,
		Tokens:       tokens,
	})
	logger := logger.GetLogger(logger.Options{})

	logger.Info(fmt.Sprint("Notification sent for : ", userId), nil)
	logger.Info(fmt.Sprint("Response success count : ", res.SuccessCount), nil)
	logger.Info(fmt.Sprint("Response failure count : ", res.FailureCount), nil)
	return err
}

func SendTopicNotification(
	fcmClient *messaging.Client, ctx context.Context, topic string, notification messaging.Notification) error {

	//Send to Multiple Tokens
	res, err := fcmClient.Send(ctx, &messaging.Message{
		Notification: &notification,
		Topic:        topic,
	})
	logger := logger.GetLogger(logger.Options{})
	logger.Info(fmt.Sprint("Report status : ", res), nil)
	return err
}
