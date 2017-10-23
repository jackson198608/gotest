package task

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/tools"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

const dbAuth = "dog123:dog123"
const dbDsn = "192.168.86.193:3307"
const dbName = "new_dog123"
const mongoConn = "192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"

func newtask() (*Task, error) {
	//getXormEngine
	connStr := tools.GetMysqlDsn(dbAuth, dbDsn, dbName)
	engine, err := xorm.NewEngine("mysql", connStr)
	if err != nil {
		return nil, err
	}

	//get mongo session
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		return nil, err
	}

	t, err := NewTask("raw|club", engine, session)
	return t, err

}

func TestNewTask(t *testing.T) {

	task, err := newtask()
	if task == nil {
		t.Log("task create error", err)
		t.Fail()
	}

}

func TestDoTask(t *testing.T) {
	task, err := newtask()
	if task == nil {
		t.Log("task create error", err)
		t.Fail()
	}

	err = task.Do()
	if err != nil {
		t.Log("task do error", err)
		t.Fail()
	}

	closetask(task)

}

func closetask(t *Task) {
	t.MysqlXorm.Close()
	t.MongoConn.Close()
}