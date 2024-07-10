package inject_tomongo

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Natannegara/inject_tomongo/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnyData interface {
	GetId() string
}

func CreateId() string {
	year, mon, _ := time.Now().Date()
	monStr := strconv.Itoa(int(mon))
	yearStr := strconv.Itoa(year)
	return string(monStr + yearStr)
}

func Controller(data AnyData, command string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := mongodb.Connect()
	dbCollection := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))
	switch command {
	case "add":
		// result := checkDuplicate(ctx, dbCollection, data.GetId())
		// if result != nil {
		// 	fmt.Println("you have generated data for this month, want to recreate?")
		// } else {
		// }
		insertData(ctx, dbCollection, data)
	case "read":
		result := readData(ctx, dbCollection)
		for _, res := range result {
			fmt.Println(res)
		}
	default:
		fmt.Println("please define your argument")
	}
}

func checkDuplicate(ctx context.Context, collection *mongo.Collection, id string) interface{} {
	var result interface{}
	if err := collection.FindOne(ctx, bson.D{{"id", id}}).Decode(&result); err != nil {
		fmt.Println(err)
	}
	return result
}

func insertData(ctx context.Context, collection *mongo.Collection, data interface{}) {
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		panic(err)
	}
	fmt.Println("succesfuly added", res.InsertedID)
}

func readData(ctx context.Context, collection *mongo.Collection) []interface{} {
	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		panic(err)
	}

	var result []interface{}
	if err = cur.All(ctx, &result); err != nil {
		panic(err)
	}
	fmt.Println("going to print all datas")
	return result
}
