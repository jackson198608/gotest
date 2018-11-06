package cardFansPersons

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

func testConn() ([]*xorm.Engine, []*mgo.Session) {
	dbAuth := "dog123:dog123"
	dbDsn := "192.168.86.194:3307"
	// dbDsn := "210.14.154.117:33068"
	dbName := "new_dog123"
	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	dbName1 := "card"
	dataSourceName1 := dbAuth + "@tcp(" + dbDsn + ")/" + dbName1 + "?charset=utf8mb4"
	engine1, err := xorm.NewEngine("mysql", dataSourceName1)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	mongoConn := "192.168.86.80:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		fmt.Println("[error] connect mongodb err")
		return nil, nil
	}

	var engineAry []*xorm.Engine
	engineAry = append(engineAry, engine)
	engineAry = append(engineAry, engine1)
	var sessionAry []*mgo.Session
	sessionAry = append(sessionAry, session)
	Init()
	return engineAry, sessionAry
	// return engine, session
}

func jsonData() *job.FocusJsonColumn {
	var jsonData job.FocusJsonColumn
	jsonData.Uid = 2060500
	jsonData.TypeId = 30
	jsonData.Created = "2017-10-23 22:54:11"
	jsonData.Infoid = 56921
	jsonData.Title = ""
	jsonData.Content = "测首页"
	jsonData.Imagenums = 0
	jsonData.ImageInfo = "7916"
	jsonData.Source = 2
	jsonData.Status = -1
	jsonData.Action = -1
	jsonData.PetType = 1
	jsonData.PetId = 71
	jsonData.VideoUrl= ""
	jsonData.IsVideo =  0
	return &jsonData
}

var m map[int]bool

func Init() {
	m = make(map[int]bool)

	mongoConn := "192.168.86.80:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		// return m
	}

	var uids []int
	c := session.DB("ActiveUser").C("active_user")
	err = c.Find(nil).Distinct("uid", &uids)
	if err != nil {
		// panic(err)
		// return m
	}
	for i := 0; i < len(uids); i++ {
		m[uids[i]] = true
	}
}

func TestGetPersons(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()
	f := NewCardFansPersons(mysqlXorm, mongoConn, jsonData, &m)
	fmt.Println(f.getPersons(1))
}

func TestPushPerson(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()

	f := NewCardFansPersons(mysqlXorm, mongoConn, jsonData, &m)
	fmt.Println(f.pushPerson(881050))
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()

	f := NewCardFansPersons(mysqlXorm, mongoConn, jsonData, &m)
	fmt.Println(f.Do())
}