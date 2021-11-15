package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Set struct {
	db *redis.Client
	setName string
}

type SetError struct {
	setName string
	identifier string
}

func (se *SetError) Error() string {
	return "identifier: " + fmt.Sprint(se.identifier) + " is part of set " + se.setName + "."
}

func (s Set) Cardinality() (int64, error) {
	var context = context.Background()
	return s.db.ZCard(context, s.setName).Result()
}

func (s Set) AddMembers(Identifiers ...string) []error {
	var context = context.Background()
	var errs []error
	var SetCardinality, err = s.db.ZCard(context, s.setName).Result()
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	for i := 0; i < len(Identifiers); i++ {
		var newMember = &redis.Z{
			Score: float64(SetCardinality) + float64(i),
			Member: Identifiers[i],
		}
		var code, err = s.db.ZAdd(context, s.setName, newMember).Result()
		errs = append(errs, err)
		if code == 0 {
			errs = append(errs, &SetError{setName: s.setName,identifier: Identifiers[i],})
		}
	}
	return errs
}

func (s Set) AddMember(Identifier interface{}) (int64, error) {
	var context = context.Background()
	var SetCardinality, err = s.db.ZCard(context, s.setName).Result()
	if err != nil {
		return -1, err
	}
	var newMember = &redis.Z{
		Score: float64(SetCardinality),
		Member: Identifier,
	}
	return s.db.ZAdd(context, s.setName, newMember).Result()
}

func (s Set) IncrementScore(Identifier string, Score ...float64) error {
	var context = context.Background()
	if totalScore := sum(Score); totalScore == 0 {
		return s.db.ZIncrBy(context, s.setName, totalScore, Identifier).Err()
	} else {
		return s.db.ZIncrBy(context, s.setName, 1, Identifier).Err()
	}
}

// func (s Set) SetScore(Identifier string, Score float64) error {
// 	var context = context.Background()
// 	var currentscore, err = s.db.ZScore()
// }

func (s Set) GetMembers(Count, Offset int64) ([]string, error) {
	var context = context.Background()
	return s.db.ZRange(context, s.setName, Offset, Offset + Count).Result()
}

func (s Set) RemoveMembers(Identifier ...string) error {
	var context = context.Background()
	return s.db.ZRem(context, s.setName, Identifier).Err()
}
 
func (s Set) GetMembersByIntersect(Setnames ...string) ([]string, error) {
	var context = context.Background()
	var store = redis.ZStore{
		Keys: Setnames,
	}
	return s.db.ZInter(context, &store).Result()
}

func sum(In []float64) float64 {
	var result float64 = 0
	for i:= 0; i < len(In); i++ {
		result += In[i]
	}
	return result
}