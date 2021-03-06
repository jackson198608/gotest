package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn           string
	dbName          string
	dbAuth          string
	numloops        int
	redisConn       string
	queueName       string
	saveDir         string
	logFile         string
	logLevel        int
	mongoConn       string
	h5templatefile  string
	miptemplatefile string
	tidStart        string
	tidEnd          string
	domain          string
}

func loadConfig(args []string) {
	var config *jconfig.Config

	if len(args) >= 3 {
		config = jconfig.LoadConfig(args[2])
	} else {
		config = jconfig.LoadConfig("/etc/goumincms.json")
	}
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.numloops = config.GetInt("numloops")
	c.redisConn = config.GetString("redisConn")
	c.queueName = config.GetString("queueName")
	c.saveDir = config.GetString("saveDir")
	c.logFile = config.GetString("logFile")
	c.logLevel = config.GetInt("logLevel")
	c.mongoConn = config.GetString("mongoConn")
	c.h5templatefile = config.GetString("h5templatefile")
	c.miptemplatefile = config.GetString("miptemplatefile")
	c.tidStart = config.GetString("tidStart")
	c.tidEnd = config.GetString("tidEnd")
	c.domain = config.GetString("domain")
}
