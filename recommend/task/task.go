package task

import (
	"database/sql"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/recommend/Pushdata"
	mgo "gopkg.in/mgo.v2"
	"strconv"
	// "strings"
)

type Task struct {
	loggerLevel int
	uid         int
	db          *sql.DB
	session     *mgo.Session
	session1    *mgo.Session
	// slave       *mgo.Session
	// event       *EventLogNew
}

func NewTask(loggerLevel int, redisStr string, db *sql.DB, session *mgo.Session, session1 *mgo.Session) *Task {
	if loggerLevel < 0 {
		loggerLevel = 0
	}
	logger.SetLevel(logger.LEVEL(loggerLevel))

	t := new(Task)
	t.session = session
	t.session1 = session1
	// t.slave = slave
	t.uid,_ = strconv.Atoi(redisStr)
	t.db = db
	return t
}

func (t *Task) Dopush(pustLimit string) {
	m := Pushdata.RecommendUser(t.loggerLevel, t.db, t.session, t.session1)
	m.PushActiveUserRecommendTask(t.uid, pustLimit)
}

func (t *Task) Dopushdog(pustLimit string) {
	m := Pushdata.RecommendUser(t.loggerLevel, t.db, t.session, t.session1)
	m.PushActiveUserDogRecommendTask(t.uid, pustLimit)
}