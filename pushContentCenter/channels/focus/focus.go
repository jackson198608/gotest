package focus

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mysql/new_dog123"
	"math"
	"strings"
)

type Focus struct {
	mysqlXorm *xorm.Engine
	mongoConn *mgo.Session
	jobstr    string
	jsonData  *jsonColumn
}

//json column
type jsonColumn struct {
	TypeId    int
	Uid       int
	Created   int
	Infoid    int
	Status    int
	Tid       int
	Bid       int
	Fid       int
	Content   string
	Title     string
	Imagenums int
	Forum     string
	Tag       int
	Qsttype   int
	Source    int
}

type EventLogX struct {
	Id        bson.ObjectId "_id"
	TypeId    int           "type"
	Uid       int           "uid"
	Fuid      int           "fuid" //fans id
	Created   string        "created"
	Infoid    int           "infoid"
	Status    int           "status"
	Tid       int           "tid"
	Bid       int           "bid"
	Content   string        "content"
	Title     string        "title"
	Imagenums int           "image_num"
	Forum     string        "forum"
	Tag       int           "tag"
	Qsttype   int           "qst_type"
	Source    int           "source"
}

const count = 100

func NewFocus(mysqlXorm *xorm.Engine, mongoConn *mgo.Session, jobStr string) *Focus {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobStr == "") {
		return nil
	}

	f := new(Focus)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jobstr = jobStr

	//@todo pass params
	jsonColumn, err := f.parseJson()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	f.jsonData = jsonColumn

	return f

}

func (f *Focus) Do() error {
	page, page1 := f.getPersionsPageNum()
	if f.jsonData.TypeId == 1 {
		//帖子 粉丝和俱乐部
		if (page <= 0) && (page1 <= 0) {
			return errors.New("you have  no person to push " + f.jobstr)
		}
	}

	if page <= 0 {
		return errors.New("you have  no person to push " + f.jobstr)
	}

	var startId int
	var endId int
	for i := 1; i <= page; i++ {
		currentPersionList := f.getPersons(page, startId, endId)
		f.pushPersons(currentPersionList)

	}
	return nil
}

//change json colum to object private member
func (f *Focus) parseJson() (*jsonColumn, error) {
	var jsonC jsonColumn

	jobs := strings.Split(f.jobstr, "|")
	if len(jobs) <= 1 {
		return &jsonC, errors.New("you have no job")
	}

	jsonStr := jobs[0]
	js, err := simplejson.NewJson([]byte(jsonStr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.Uid, _ = js.Get("uid").Int()
	jsonC.TypeId, _ = js.Get("event_type").Int()
	jsonC.Created, _ = js.Get("time").Int()
	jsonC.Tid, _ = js.Get("tid").Int()
	jsonC.Bid, _ = js.Get("event_info").Get("bid").Int()
	jsonC.Infoid, _ = js.Get("event_info").Get("infoid").Int()
	jsonC.Title, _ = js.Get("event_info").Get("title").String()
	jsonC.Content, _ = js.Get("event_info").Get("content").String()
	jsonC.Forum, _ = js.Get("event_info").Get("forum").String()
	jsonC.Imagenums, _ = js.Get("event_info").Get("image_num").Int()
	jsonC.Tag, _ = js.Get("event_info").Get("tag").Int()
	jsonC.Qsttype, _ = js.Get("event_info").Get("qst_type").Int()
	jsonC.Fid, _ = js.Get("event_info").Get("fid").Int()
	jsonC.Source, _ = js.Get("source").Int()
	jsonC.Status, _ = js.Get("status").Int()

	return &jsonC, nil
}

func (f *Focus) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("you have no person to push " + f.jobstr)
	}

	for _, person := range persons {
		err := f.pushPerson(person)
		if err != nil {
			//@todo if err times < 5 ,just print log
			//      if err times > 5 ,return err
		}
	}
	return nil
}

func (f *Focus) pushPerson(person int) error {

	return nil
}

func (f *Focus) pushData(fuid int) *EventLogX {
	var data EventLogX
	return &data
}

//@todo how to remove duplicate uid from to lists
func (f *Focus) getPersons(page int, startId int, endId int) []int {
	var uid []int
	typeId := f.jsonData.TypeId
	if typeId == 1 {
		//帖子 推所有活跃粉丝 + 相同俱乐部的活跃用户
		clubUids := f.getClubPersons(page)
		fansUids := f.getFansPersons(startId, endId)
		uid = MergePersons(fansUids, clubUids)
	} else if typeId == 6 {
		// 视频 推活跃粉丝
		uid = f.getFansPersons(startId, endId)
	} else if typeId == 8 {
		//问答 推相同犬种的活跃用户
		uid = f.getBreedPersons(page)
	} else if ((typeId == 9) || (typeId == 15)) && (f.jsonData.Source == 1) {
		//人工小编 推全部活跃用户
		uid = f.getActivePersons(page)
	} else {
		//推全部活跃用户
		uid = f.getActivePersons(page)
	}
	fmt.Println(uid)
	return uid
}

func (f *Focus) getPersionsPageNum() (int, int) {
	typeId := f.jsonData.TypeId
	if typeId == 1 {
		//帖子 推所有活跃粉丝 + 相同俱乐部的活跃用户

	} else if typeId == 6 {
		// 视频 推活跃粉丝

	} else if typeId == 8 {
		//问答 推相同犬种的活跃用户

	} else if ((typeId == 9) || (typeId == 15)) && (f.jsonData.Source == 1) {
		//人工小编 推全部活跃用户
	} else {
		//推全部活跃用户
	}

	return 0, 0
}

//合并俱乐部和粉丝数据
func MergePersons(fansuids []int, clubuids []int) []int {
	var alluids []int

	//@todo 发帖用户的所有活跃粉丝

	//@todo 所有加入该帖子俱乐部的活跃用户

	return alluids
}

//获取相同犬种的活跃用户数
func (f *Focus) getBreedPersonsPagNum() int {
	Bid := f.jsonData.Bid
	if Bid == 0 {
		return 0
	}

	c := f.mongoConn.DB("ActiveUser").C("active_breed_user")
	countNum, err := c.Find(&bson.M{"breed_id": Bid}).Count()
	if err != nil {
		panic(err)
	}
	page := int(math.Ceil(float64(countNum) / float64(count)))
	return page
}

//获取相同犬种的活跃用户
//@todo 使用id范围分页查询
func (f *Focus) getBreedPersons(page int) []int {
	var uids []int
	Bid := f.jsonData.Bid
	if Bid == 0 {
		return uids
	}

	c := f.mongoConn.DB("ActiveUser").C("active_breed_user")
	err := c.Find(&bson.M{"breed_id": Bid}).
		Skip((page-1)*count).
		Limit(count).
		Distinct("uid", &uids)
	if err != nil {
		panic(err)
	}
	return uids
}

func (f *Focus) getClubPersonPageNum() int {
	fid := f.jsonData.Fid
	c := f.mongoConn.DB("ActiveUser").C("active_forum_user")
	countNum, err := c.Find(&bson.M{"forum_id": fid}).Count()
	if err != nil {
		panic(err)
	}
	page := int(math.Ceil(float64(countNum) / float64(count)))

	return page
}

//获取相同俱乐部的活跃用户
//@todo 使用id范围分页查询
func (f *Focus) getClubPersons(page int) []int {
	var uids []int
	fid := f.jsonData.Fid
	c := f.mongoConn.DB("ActiveUser").C("active_forum_user")
	err := c.Find(&bson.M{"forum_id": fid}).
		Skip((page-1)*count).
		Limit(count).
		Distinct("uid", &uids)
	if err != nil {
		panic(err)
	}
	return uids
}

//获取活跃用户第一个ID
func (f *Focus) getFansPersonFirstId() int {
	uid := f.jsonData.Uid
	var follows []new_dog123.Follow
	err := f.mysqlXorm.Where("user_id=? and fans_active=1", uid).Asc("id").Limit(1).Find(&follows)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	if len(follows) == 0 {
		return 0
	}

	id := follows[0].Id
	return id
}

//获取活跃用户最后一个ID
func (f *Focus) getFansPersonLastId() int {
	uid := f.jsonData.Uid
	var follows []new_dog123.Follow
	err := f.mysqlXorm.Where("user_id=? and fans_active=1", uid).Desc("id").Limit(1).Find(&follows)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if len(follows) == 0 {
		return 0
	}

	id := follows[0].Id

	return id
}

//获取活跃粉丝用户
func (f *Focus) getFansPersons(startId int, endId int) []int {
	var uids []int
	uid := f.jsonData.Uid
	var follows []new_dog123.Follow
	err := f.mysqlXorm.Where("user_id=? and id>=? and id<=? and fans_active=1", uid, startId, endId).Cols("follow_id").Find(&follows)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, v := range follows {
		uids = append(uids, v.FollowId)
	}
	return uids
}

func (f *Focus) getActivePersonPageNum() int {
	c := f.mongoConn.DB("ActiveUser").C("active_user")
	countNum, err := c.Find(nil).Count()
	if err != nil {
		panic(err)
	}
	page := int(math.Ceil(float64(countNum) / float64(count)))

	return page
}

//获取所有活跃用户
//@todo 使用id范围分页查询
func (f *Focus) getActivePersons(page int) []int {
	var uids []int
	c := f.mongoConn.DB("ActiveUser").C("active_user")
	err := c.Find(nil).Skip((page-1)*count).Limit(count).Distinct("uid", &uids)
	if err != nil {
		panic(err)
	}
	return uids
}
