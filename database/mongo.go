package database

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/xulei324/spiderlab-go-sdk/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func GetDataSourceCol(host string, port string, username string, password string, authSource string, database string, col string) (*mongo.Client, *mongo.Collection, error) {
	//timeout := time.Second * 10
	//
	//dialInfo := mgo.DialInfo{
	//	Addrs:         []string{net.JoinHostPort(host, port)},
	//	Timeout:       timeout,
	//	Database:      database,
	//	PoolLimit:     100,
	//	PoolTimeout:   timeout,
	//	ReadTimeout:   timeout,
	//	WriteTimeout:  timeout,
	//	AppName:       "crawlab",
	//	FailFast:      true,
	//	MinPoolSize:   10,
	//	MaxIdleTimeMS: 1000 * 30,
	//}
	//if username != "" {
	//	dialInfo.Username = username
	//	dialInfo.Password = password
	//	dialInfo.Source = authSource
	//}
	//if Session == nil {
	//	s, err := mgo.DialWithInfo(&dialInfo)
	//	if err != nil {
	//		log.Errorf("dial mongo error: " + err.Error())
	//		debug.PrintStack()
	//		return nil, nil, err
	//	}
	//	Session = s
	//}
	//db := Session.DB(database)
	//return Session, db.C(col), nil

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", username, password, host, port, database, authSource)))
	if err != nil {
		log.Errorf("error:", err)
		return nil, nil, err
	}
	collection := client.Database(database).Collection(col)

	return client, collection, nil
}

func GetMongoCol(ds entity.DataSource) (*mongo.Client, *mongo.Collection, error) {
	if ds.Type == "" {
		return GetDataSourceCol(
			os.Getenv("CRAWLAB_MONGO_HOST"),
			os.Getenv("CRAWLAB_MONGO_PORT"),
			os.Getenv("CRAWLAB_MONGO_USERNAME"),
			os.Getenv("CRAWLAB_MONGO_PASSWORD"),
			os.Getenv("CRAWLAB_MONGO_AUTHSOURCE"),
			os.Getenv("CRAWLAB_MONGO_DB"),
			os.Getenv("CRAWLAB_COLLECTION"),
		)
	}
	return GetDataSourceCol(
		ds.Host,
		ds.Port,
		ds.Username,
		ds.Password,
		ds.AuthSource,
		ds.Database,
		os.Getenv("CRAWLAB_COLLECTION"),
	)
}
