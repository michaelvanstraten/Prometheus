package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type CollectionMemberType int

const (
	RedisString = iota
	RedisList = iota
	RedisHashmap = iota
	RedisSet = iota
	RedisSortedSet = iota
	RedisBitmap = iota
)

type CollectionOptions struct {
	collectionName string
	redisClient *redis.Client
}

type Collection struct {
	db *redis.Client
	collectionName string
	Sets map[string]Set
}

func NewColletion(ColletionName string, RedisClient *redis.Client) *Collection {
	var newColletion = Collection{
		db : RedisClient,
		collectionName: ColletionName,
		Sets: make(map[string]Set),
	}
	newColletion.Sets["MASTER"] = Set{db: RedisClient, setName: ColletionName  + "s"}
	return &newColletion
}

func (c *Collection) AddSubset(SetName string) {
	c.Sets[SetName] = Set{
		db : c.db,
		setName : c.Sets["MASTER"].setName + ":" + SetName,
	}
}

func (c *Collection) GetMember(Identifier string, MemberInterface interface{}, Fields ...string) error {
	var context = context.Background()
	return c.db.HMGet(context, c.collectionName + ":" + Identifier, Fields...).Scan(MemberInterface)
}

func (c *Collection) GetMembers(MemberInterfaces []interface{}, Set Set, Offset int64, Fields ...string) []error {
	var context = context.Background()
	var errs []error
	var Identifiers, err = Set.GetMembers(int64(len(MemberInterfaces)), Offset)
 	if err != nil {
		errs = append(errs, err)
		return errs
	}
	var li, lm = len(Identifiers), len(MemberInterfaces)
	if len(Fields) > 0 && li <= lm {
		for i := 0; i < li && i < lm ; i++ {
			errs = append(errs, c.db.HMGet(context, c.collectionName + ":" + Identifiers[i], Fields...).Scan(MemberInterfaces[i]))
		}
	} else if li <= lm {
		for i := 0; i < li && i < lm; i++ {
			errs = append(errs, c.db.HGetAll(context, c.collectionName + ":" + Identifiers[i]).Scan(MemberInterfaces[i]))
		}
	}
	return errs
}

func (c *Collection) AddMemberWithCustomIdentifier(Identifier interface{}, Values ...interface{}) error {
	var context = context.Background()
	var addItem = func(tx *redis.Tx) error {
		if returnCode, err := c.Sets["MASTER"].AddMember(fmt.Sprint(Identifier)); err == nil && returnCode == 1 {
			err = tx.HSet(context, c.collectionName + ":" + fmt.Sprint(Identifier), Values...).Err()
			if err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}
		return errors.New("identifier: " + fmt.Sprint(Identifier) + " is part of collection " + c.collectionName + "s.")
	}
	for {
		var err = c.db.Watch(context, addItem, c.collectionName + "s")
		if err == redis.TxFailedErr {
			continue
		} else {
			return err
		}
	}
}

func (c *Collection) AddMember(Values ...interface{}) (string, error) {
	var context = context.Background()
	var Identifier string
	var addItem = func(tx *redis.Tx) error {
		var Identifier, err = tx.ZCard(context, c.Sets["MASTER"].setName).Result()
		if err != nil {
			return err
		}
		err = tx.HSet(context, c.collectionName + ":" + fmt.Sprint(Identifier), Values...).Err()
		if err != nil {
			return err
		} 
		return nil
	}
	for {
		var err = c.db.Watch(context, addItem, c.Sets["MASTER"].setName)
		if err == redis.TxFailedErr {
			continue
		} else {
			return Identifier, err
		}
	}
}

func (c *Collection) SetValueOnMember(Identifier string, Values ...interface{}) error {
	var context = context.Background()
	return c.db.HMSet(context, c.collectionName + Identifier, Values...).Err()
}

func (c *Collection) RemoveMembers(Identifiers ...string) []error {
	var context = context.Background()
	var errs []error
	for i := 0; i < len(Identifiers); i++ {
		errs = append(errs, c.db.HDel(context, c.collectionName + ":" + Identifiers[i]).Err()) 
	}
	for _, set := range c.Sets {
		errs = append(errs, set.RemoveMembers(Identifiers...))
	}
	return errs
}