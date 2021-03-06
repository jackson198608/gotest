package main

import (
	"github.com/jackson198608/goProject/go_spider/core/common/page"
	"strconv"
	"strings"
	//"time"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

func getDetailId(tag string) (int, bool) {
	tags := strings.Split(tag, "|")
	shopDetailIdStr := tags[1]
	shopDetailId, err := strconv.Atoi(shopDetailIdStr)
	if err != nil {
		logger.Println("[error]invaild shop id ", tag)
		return 0, false
	}
	return shopDetailId, true

}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
	logger.Println("[info]in Process ")
	//time.Sleep(1 * time.Second)
	if !p.IsSucc() {
		logger.Println("[Error]not 200: ", p.GetRequest().Url)
		return
	}

	tag := p.GetUrlTag()

	if tag == "shopDetail" {
		logger.Println("[info]shop detail page:", p.GetRequest().Url)
		//save shop into mysql
		shopDetailId, success := saveShopDetail(p)
		if success {
			logger.Println("[info]Get detail id", shopDetailId)
			//get all query
			qShopDetail(p, shopDetailId)
		}
	} else if tag == "DogCateList" {
		qShopCateList(p)

	} else if tag == "shopList" {
		qShopList(p)
	} else if tag == "sogou" {
		//@todo for sogou weixin
		logger.Println("sogou tag callback")
		qSogouHome(p)

	} else if tag == "weixinAccountHome" {
		//@todo for sogou weixin
		logger.Println("weixinAccountHome tag callback")
		qWeixinAccountHome(p)

	} else if tag == "shopImage" {
		saveImage(p)

	} else if strings.Contains(tag, "shopImage") {
		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		saveShopImagePath(p, int64(shopDetailId))

	} else if strings.Contains(tag, "shopCommentList") {
		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		saveShopCommentList(p, int64(shopDetailId))
	} else if strings.Contains(tag, "goodsPrice") {
		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		saveGoodsPrice(p, int64(shopDetailId))
	} else if strings.Contains(tag, "goodsCommentNumScore") {
		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		updateSuccess := saveGoodsCommentNumAndScore(p, int64(shopDetailId))
		if !updateSuccess {
			logger.Println("[error] update goods comment num and score fail ", tag, p.GetRequest().Url)
			return
		}
		qGoodsCommentList(p, int64(shopDetailId))
	} else if strings.Contains(tag, "goodsDescImage") {
		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		qGoodsDescImage(p, int64(shopDetailId))
	} else if strings.Contains(tag, "skuPrice") {
		logger.Println("[info] get sku price by tag : ", tag)
		qSkuPrice(p)
	} else if strings.Contains(tag, "skuCommentNum") {
		logger.Println("[info] get sku comment num by tag : ", tag)
		qSkuCommentNum(p)
	}
}

func (this *MyPageProcesser) Finish() {
}
