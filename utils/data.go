package utils

import (
	"context"
	"github.com/apex/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/xulei324/spiderlab-go-sdk/constants"
	"github.com/xulei324/spiderlab-go-sdk/database"
	"github.com/xulei324/spiderlab-go-sdk/entity"
	"runtime/debug"
)

func SaveItem(item entity.Item) (err error) {
	dsType := GetDataSourceType()
	if dsType == constants.DataSourceTypeMongo {
		if err := SaveItemMongo(item); err != nil {
			return err
		}
	} else if dsType == constants.DataSourceTypeKafka {
		if err := SaveItemKafka(item); err != nil {
			return err
		}
	} else if dsType == constants.DataSourceTypeElasticSearch {
		if err := SaveItemElasticSearch(item); err != nil {
			return err
		}
	} else {
		if err := SaveItemSql(item); err != nil {
			return err
		}
	}
	return nil
}

func SaveItemMongo(item entity.Item) (err error) {
	ds := GetDataSource()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	_, c, err := database.GetMongoCol(ds)
	item["task_id"] = GetTaskId()

	isDedup := GetIsDedup()

	if isDedup == "1" {
		// 去重
		dedupField := GetDedupField()
		dedupMethod := GetDedupMethod()
		if dedupMethod == constants.DedupMethodOverwrite {
			// 覆盖
			//var res interface{}

			if r := c.FindOne(ctx, bson.M{dedupField: item[dedupField]}); r.Err() != nil {
				if err == mongo.ErrNoDocuments {
					// 不存在
					if _, err := c.InsertOne(ctx, item); err != nil {
						log.Errorf("save item error: " + err.Error())
						debug.PrintStack()
						return err
					}
				} else {
					log.Errorf("find item error: " + err.Error())
					debug.PrintStack()
					return err
				}
			} else {
				// 已存在

				if _, err := c.UpdateOne(ctx, bson.M{dedupField: item[dedupField]}, item); err != nil {
					log.Errorf("update item error: " + err.Error())
					debug.PrintStack()
					return err
				}
			}
		} else if dedupMethod == constants.DedupMethodIgnore {
			// 忽略
			if _, err := c.InsertOne(ctx, item); err != nil {
				log.Errorf("save item error: " + err.Error())
				debug.PrintStack()
				return err
			}
		} else {
			// 其他
			if _, err := c.InsertOne(ctx, item); err != nil {
				log.Errorf("save item error: " + err.Error())
				debug.PrintStack()
				return err
			}
		}
	} else {
		// 不去重
		if _, err := c.InsertOne(ctx, item); err != nil {
			log.Errorf("save item error: " + err.Error())
			debug.PrintStack()
			return err
		}
	}
	return nil
}

func SaveItemSql(item entity.Item) error {
	ds := GetDataSource()
	item["task_id"] = GetTaskId()

	isDedup := GetIsDedup()

	if isDedup == "1" {
		// 去重
		dedupField := GetDedupField()
		dedupMethod := GetDedupMethod()
		if dedupMethod == constants.DedupMethodOverwrite {
			// 覆盖
			_item, _ := database.GetItem(ds, dedupField, item[dedupField].(string))
			if _item == nil {
				// 不存在
				if err := database.InsertItem(ds, item); err != nil {
					log.Errorf("save item error: " + err.Error())
					debug.PrintStack()
					return err
				}
			} else {
				// 已存在
				if err := database.UpdateItem(ds, item, dedupField); err != nil {
					log.Errorf("save item error: " + err.Error())
					debug.PrintStack()
					return err
				}
			}
		} else if dedupMethod == constants.DedupMethodIgnore {
			// 忽略
			if err := database.InsertItem(ds, item); err != nil {
				log.Errorf("save item error: " + err.Error())
				debug.PrintStack()
				return err
			}
		}
	} else {
		// 不去重
		if err := database.InsertItem(ds, item); err != nil {
			log.Errorf("save item error: " + err.Error())
			debug.PrintStack()
			return err
		}
	}

	return nil
}

func SaveItemKafka(item entity.Item) error {
	ds := GetDataSource()
	if err := database.SendKafkaMsg(ds, item); err != nil {
		return err
	}
	return nil
}

func SaveItemElasticSearch(item entity.Item) error {
	ds := GetDataSource()
	if err := database.IndexItem(ds, item); err != nil {
		return err
	}
	return nil
}
