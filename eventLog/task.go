package main

import (
    "database/sql"
    "errors"
    "github.com/donnie4w/go-logger/logger"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jackson198608/squirrel"
    "strconv"
    "fmt"
)

// // import (
// //     // "strconv"
// //     "fmt"
// // )
// // func NewTask(taskNum int) {
// //     var startId int = 1
// //     var endId int = 10
// //     var limit int = 10
// //     var offset int = 0

// //     // count := getPostCount(startPid,endPid)
// //     for {
// //         task := getEventLogData(startId,endId,limit,offset)
// //         if len(task) == 0 {
// //             fmt.Println("task data is empty")
// //             return
// //         }
// //         for _,v := range task {
// //             insertEventLog(v.typeId,v.uid,v.info,v.created,v.infoid,v.status,v.tid)
// //         }
// //         break;
// //         // offset += limit
// //         // if offset > count {
// //         //     break;
// //         // }
// //         // fmt.Println(task)
// //     }
// // }


type Task struct {
    id      int64
    logLevel int
    dbAuth   string
    dbDsn    string
    dbName   string
    con      *sql.DB
    dbCache  squirrel.DBProxyBeginner
    ids     []int64
}

func NewTask(logLevel int, idStr string, args []string) *Task {

    logger.SetLevel(logger.LEVEL(logLevel))

    //check the string
    if len(args) != 3 {
        logger.Error("there is not enough args to start")
        return nil
    }

    //set value
    id, err := strconv.Atoi(idStr)
    if err != nil {
        logger.Error("tid error")
        return nil
    }

    t := new(Task)
    t.id = int64(id)
    t.dbAuth = args[0]
    t.dbDsn = args[1]
    t.dbName = args[2]
    t.logLevel = logLevel
    //make db comon value and check error
    err = t.getDbCache()
    if err != nil {
        logger.Error("get id list error", err, t.id)
        return nil
    }

    // err = t.getPids()
    // if err != nil {
    //     logger.Error("get id list error", err, t.id)
    //     t.con.Close()
    //     return nil
    // }

    // //if have no pid,no task
    // if len(t.ids) == 0 {
    //     logger.Error("this id have no pid ,so pass", t.id)
    //     t.con.Close()
    //     return nil
    // }

    return t
}

func (t *Task) Over() {
    if t.con != nil {
        t.con.Close()
    }
}
func (t *Task) getDbCache() error {
    con, err := sql.Open("mysql", t.dbAuth+"@tcp("+t.dbDsn+")/"+t.dbName+"?charset=utf8")
    if err != nil {
        logger.Error("connect err", t.dbDsn, t.dbAuth, t.dbName)
        return errors.New("connect db error")

    }
    // Third, we wrap in a prepared statement cache for better performance.
    cache := squirrel.NewStmtCacheProxy(con)
    t.dbCache = cache
    t.con = con
    return nil
}

func (t *Task) Do() error {
    if t.id == 0 {
        return errors.New("id is nil")
    }
    fmt.Println(t.id)
    // idlen := len(t.ids)
    // for i := 0; i < idlen; i++ {
        t.handleId(t.id)
    // }
    return nil
}

func (t *Task) handleId(id int64) error {
    if id <= 0 {
        logger.Error("id <0")
        return errors.New("id<0")
    }
    event := NewEvent(t.logLevel, t.dbCache, "mysql", id, false)
    logger.Info("doing for this event id", event.id, id)
    // exist := event.IdExists()
    // fmt.Println(event.created)
    // if exist {
    //     logger.Info("doing for this event id", event.id, id)
        result := event.MoveToSplit()
        if !result {
            logger.Error("change event id error", id)
            return errors.New("change event id error")
        }
    // } else {
    //     logger.Info("not exsit", id)
    //     return nil
    // }

    return nil

}


