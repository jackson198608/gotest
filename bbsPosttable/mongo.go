package main

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "strconv"
  "fmt"
)


type AppMessage struct {
    Id    bson.ObjectId `bson:"_id"`
    Info  string        `bson:"info"`
    Tid  string        `bson:"tid"` //bson:"name" 表示mongodb数据库中对应的字段名称
    Type string         `bson:"type"`
}

const URL = "192.168.86.68:27017" //mongodb连接字符串

var (
    mgoSession *mgo.Session
    dataBase   = "AppMessage"
)
/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
    if mgoSession == nil {
        var err error
        mgoSession, err = mgo.Dial(URL)
        if err != nil {
            panic(err) //直接终止程序运行
        }
    }
    //最大连接池默认为4096
    return mgoSession.Clone()
}
//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
    session := getSession()
    defer session.Close()
    c := session.DB(dataBase).C(collection)
    return s(c)
}

/**
 * 添加person对象
 */
func AddPerson(tableid int,p AppMessage) string {
    p.Id = bson.NewObjectId()
    query := func(c *mgo.Collection) error {
        return c.Insert(p)
    }
    tablename := "app_message_"+strconv.Itoa(tableid)
    err := witchCollection(tablename, query)
    if err != nil {
        return "false"
    }
    return p.Id.Hex()
}

/**
 * 获取一条记录通过objectid
 */
func GetAppMessageById(tableid int,id string) *AppMessage {
    objid := bson.ObjectIdHex(id)
    message := new(AppMessage)
    query := func(c *mgo.Collection) error {
        return c.FindId(objid).One(&message)
    }
    tablename := "app_message_"+strconv.Itoa(tableid)
    witchCollection(tablename, query)
    return message
}

/**
 * 获取一条记录通过objectid
 */
func GetAppMessageCount(tableid int) (int){
    var count int
    where := bson.M{"type":1}
    // where := bson.M{"type": bson.M{"$in":{23, 26, 32}}}//find({"age":{"$in":(23, 26, 32)}})
    query := func(c *mgo.Collection) error {
        count,_ := c.Find(where).Count()
        return int(count)
    }
    tablename := "app_message_"+strconv.Itoa(tableid)
    witchCollection(tablename, query)
    return count
}


//获取所有的person数据
func PageAppMessage(tableid int) []AppMessage {
    var messages []AppMessage
    // where := bson.M{"type": {"$in":{23, 26, 32}}}
    where := bson.M{"type":10}
    query := func(c *mgo.Collection) error {
        return c.Find(where).All(&messages)
    }
    tablename := "app_message_"+strconv.Itoa(tableid)
    err := witchCollection(tablename, query)
    if err != nil {
        return messages
    }
    return messages
}


//更新person数据
func UpdatePerson(tableid int,query bson.M, change bson.M) string {
    exop := func(c *mgo.Collection) error {
        return c.Update(query, change)
    }
    tablename := "app_message_"+strconv.Itoa(tableid)
    err := witchCollection(tablename, exop)
    if err != nil {
        return "true"
    }
    return "false"
}

/**
 * 执行查询，此方法可拆分做为公共方法
 * [SearchPerson description]
 * @param {[type]} collectionName string [description]
 * @param {[type]} query          bson.M [description]
 * @param {[type]} sort           bson.M [description]
 * @param {[type]} fields         bson.M [description]
 * @param {[type]} skip           int    [description]
 * @param {[type]} limit          int)   (results      []interface{}, err error [description]
 */
func SearchPerson(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
    exop := func(c *mgo.Collection) error {
        return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
    }
    err = witchCollection(collectionName, exop)
    return
}
type AppMessages struct {
    AppMessages []AppMessage
}
func test(){
    // var redisString string = `{"uid":1895167,"type":1,"info":281,"isnew":0,"from":0,"created":"1483951718","modified":"1484032727"}`
    // tasks := PageAppMessage(50)
    message := GetAppMessageById(50,"58734e66bdd4fb1e004175e1")
    fmt.Println(message.Id)

    // messages := PageAppMessage(50)
    // fmt.Println(messages)
    // collection.Update(bson.M{"name": "ddd"}, bson.M{"$set": bson.M{"phone": "12345678"}})
    id := bson.M{"_id": message.Id}
    // message := bson.M{"_id": bson.ObjectIdHex("58734e66bdd4fb1e004175e1")}
    change := bson.M{"$set":bson.M{"tid":122211}}
    UpdatePerson(50,id,change)
    fmt.Println(message)

}

func task() {
    count := GetAppMessageCount(50)
    fmt.Println(count)
    result := PageAppMessage(50)
    for _,v := range result {
        fmt.Println(v.Id,v.Info)
    }
    // fmt.Println(result)
}