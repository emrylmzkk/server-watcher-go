package generic

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"

	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func InitFirebaseMessaging() *messaging.Client {

	opt := option.WithAuthCredentialsFile(option.ServiceAccount, "serviceAccountKey.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatalf("Firebase service could not be started %v", err)
	}

	client, err := app.Messaging(context.Background())

	if err != nil {
		log.Fatalf("Firebase Messaging Client could not be started %v", err)
	}

	return client

}
