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

func Controller(data AnyData, command string, isTrash bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := mongodb.Connect()
	var dbCollection *mongo.Collection

	dbCollection = client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))

	if isTrash == true {
		dbCollection = client.Database(os.Getenv("DATABASE")).Collection("trash")
	}

	switch command {
	case "add":
		// result := CheckDuplicate(ctx, dbCollection, data.GetId())
		result := CheckDuplicate(ctx, dbCollection, CreateId())
		if result != nil {
			fmt.Println("you have generated data for this month, want to recreate?")
			break
		} else {
			insertData(ctx, dbCollection, data)
			fmt.Println("added succesfuly")
		}
	case "read":
		result := readData(ctx, dbCollection)
		fmt.Println(result...)
	default:
		fmt.Println("please define your argument")
	}
}

func CheckDuplicate(ctx context.Context, collection *mongo.Collection, id string) interface{} {
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

func GetAllData(result interface{}, timeFilter string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := mongodb.Connect()
	collection := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))

	if timeFilter == "now" {
		timeFilter = CreateId()
	}

	cur, err := collection.Find(ctx, bson.D{{"id", timeFilter}})
	if err != nil {
		panic(err)
	}

	defer cur.Close(ctx)
	if err := cur.All(ctx, result); err != nil {
		return err
	}
	return nil
}

func GetTrashData(result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := mongodb.Connect()
	collection := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}

	defer cur.Close(ctx)
	if err := cur.All(ctx, result); err != nil {
		return err
	}
	return nil
}

func DeleteData() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := mongodb.Connect()
	collection := client.Database(os.Getenv("DATABASE")).Collection("trash")

	collection.Drop(ctx)
}
