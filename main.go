package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username: "madmin",
		Password: "madmin",
	})
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cnf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnf()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("test")
	calendarsCol := database.Collection("calendars")

	dateString := "25-12-2022"
	holidayDate, err := time.Parse("02-01-2006", dateString)
	if err != nil {
		log.Fatalf("failed to format date: %s\n", err.Error())
	}
	filterDate := holidayDate.Format("20060102")

	filter := bson.M{
		"currency": "zar",
		"holiday":  filterDate,
	}
	cursor, err := calendarsCol.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var calendars []bson.M
	if err = cursor.All(ctx, &calendars); err != nil {
		log.Fatal(err)
	}
	fmt.Println(calendars)

	filter = bson.M{
		"currency": "blah",
	}
	test, err := calendarsCol.Distinct(ctx, "currency", filter)
	if err != nil {
		log.Fatal(err)
	}
	test1 := make([]string, len(test))
	for i, j := range test {
		test1[i] = fmt.Sprint(j)
	}
	fmt.Println(test1)

}
