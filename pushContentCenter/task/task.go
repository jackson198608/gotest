package task

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/club"
	"github.com/jackson198608/goProject/pushContentCenter/channels/focus"
	"github.com/jackson198608/goProject/pushContentCenter/channels/recommend"
	"gopkg.in/mgo.v2"
	"strings"
	"github.com/olivere/elastic"
)

type Task struct {
	Raw       string         //the data get from redis queue
	MysqlXorm []*xorm.Engine //mysql single instance
	MongoConn []*mgo.Session //mongo single instance
	Jobstr    string         //private member parse from raw
	JobType   string         //private membe parse from raw jobType: focus|club
	esConn  *elastic.Client
}

//job: redisQueue pop string
//taskarg: mongoHost,mongoDatabase,mongoReplicaSetName
func NewTask(raw string, mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, esConn *elastic.Client) (*Task, error) {
	//check prams
	if (raw == "") || (mysqlXorm == nil) || (esConn == nil) {
		return nil, errors.New("params can not be null")
	}

	t := new(Task)
	if t == nil {
		return nil, errors.New("there is no space to create struct")
	}

	//pass params
	t.Raw = raw
	t.MysqlXorm = mysqlXorm
	t.MongoConn = mongoConn

	//create private member
	err := t.parseRaw()
	if err != nil {
		return nil, errors.New("raw format error ,can not find jobstr and jobtype " + raw)
	}

	t.esConn = esConn
	return t, nil

}

// public interface for task
// if you have New channles you need to add logic here
func (t *Task) Do() error {
	switch t.JobType {
	case "club":
		err := t.ChannelClub()
		if err != nil {
			return err
		} else {
			return nil
		}
	case "focus":
		err := t.ChannelFocus()
		if err != nil {
			return err
		} else {
			return nil
		}
	//case "recommend":
	//	err := t.ChannelRecommend()
	//	if err != nil {
	//		return err
	//	} else {
	//		return nil
	//	}

	}
	return nil
}

// focus channel's invoke function
func (t *Task) ChannelFocus() error {
	c := focus.NewFocus(t.MysqlXorm, t.MongoConn, t.Jobstr, t.esConn)
	err := c.Do()
	if err != nil {
		return err
	}
	return nil

}

// club channel's invoke function
func (t *Task) ChannelClub() error {
	c := club.NewClub(t.MysqlXorm, t.MongoConn, t.Jobstr)
	err := c.Do()
	if err != nil {
		return err
	}
	return nil
}

// recommend channel's invoke function
func (t *Task) ChannelRecommend() error {
	c := recommend.NewRecommend(t.MysqlXorm, t.MongoConn, t.esConn,t.Jobstr)
	err := c.Do()
	if err != nil {
		return err
	}
	return nil
}

// this function parase raw to judge jobstr and job type
// sep string : '|'
//return:
//         jobstr
//	       type
//		   error
func (t *Task) parseRaw() error {
	rawSlice := []byte(t.Raw)
	rawLen := len(rawSlice)
	lastIndex := strings.LastIndex(t.Raw, "|")
	t.Jobstr = string(rawSlice[0:lastIndex])
	t.JobType = string(rawSlice[lastIndex+1 : rawLen])

	return nil

}
