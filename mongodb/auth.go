package mongodb

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func auth() options.Credential {
	credential := options.Credential{
		Username: os.Getenv("MONGODB_USER"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	}

	return credential
}
