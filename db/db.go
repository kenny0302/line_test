package db

import (
	"context"
	"fmt"
	proto "main/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, host, port string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + host + ":" + port))
	if err != nil {
		fmt.Errorf("client establish failed. err: %v", err)
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		fmt.Errorf("client establish failed. err: %v", err)
		return nil, err
	}
	return client, nil
}

func GetUser(host, port, db, col string, filter interface{}) ([]proto.Output, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var output []proto.Output
	cli, err := Connect(ctx, host, port)
	if err != nil {
		return nil, err
	}

	collection := cli.Database(db).Collection(col)
	queryResult, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for queryResult.Next(context.TODO()) {
		var lr proto.Output
		if queryResult.Decode(&lr) != nil {
			fmt.Print(err)
		}
		output = append(output, proto.Output{lr.UserId, lr.DisplayName, lr.PictureUrl})
	}

	defer func() {
		if err = cli.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	return output, nil
}

func SetUser(host, port, db, col string, filter interface{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli, err := Connect(ctx, host, port)
	if err != nil {
		return err
	}

	collection := cli.Database(db).Collection(col)
	_, err = collection.InsertOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func SetMessage(host, port, db, col string, filter interface{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli, err := Connect(ctx, host, port)
	if err != nil {
		return err
	}

	collection := cli.Database(db).Collection(col)
	_, err = collection.InsertOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
