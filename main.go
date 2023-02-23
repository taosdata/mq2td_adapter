package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	_ "github.com/taosdata/driver-go/v2/taosSql"
	"github.com/taosdata/mq2td_adapter/config"
	"github.com/taosdata/mq2td_adapter/log"
	"github.com/taosdata/mq2td_adapter/mqtt"
	"github.com/taosdata/mq2td_adapter/pool"
	"github.com/taosdata/mq2td_adapter/rule"
)

func main() {
	config.Init()
	log.ConfigLog()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Close(ctx)
	}()
	manager, err := rule.NewRuleManage(config.Conf.RulePath)
	if err != nil {
		panic(err)
	}
	logger := log.GetLogger("main")
	initDB(&config.Conf.TDengine, logger)
	pool.InitGoPool(logger, 50000)
	dsn := fmt.Sprintf("%s:%s/tcp(%s:%d)/%s", config.Conf.TDengine.User, config.Conf.TDengine.Password, config.Conf.TDengine.Host, config.Conf.TDengine.Port, config.Conf.TDengine.DB)
	db, err := sql.Open("taosSql", dsn)
	createSql := manager.GenerateCreateSql()
	for _, s := range createSql {
		if config.Conf.ShowSql {
			logger.Info(createSql)
		}
		_, err = db.Exec(s)
		if err != nil {
			logger.WithError(err).WithField("sqls", createSql).Panic("exec create sql error")
		}
	}
	mqttLogger := log.GetLogger("mqtt_connect").WithField("addr", config.Conf.MQTT.Address)
	connected := make(chan struct{})
	connector := mqtt.NewConnector(config.Conf.MQTT, pool.GoPool, mqttLogger, func() {
		connected <- struct{}{}
	})
	<-connected
	logger.Info("mqtt server connected")
	connector.SubscribeWithReceiveTime("#", 0, func(topic string, msg []byte, t time.Time) {
		if manager.RuleExist(topic) {
			r, err := manager.Parse(topic, msg)
			if err != nil {
				logger.WithError(err).WithField("topic", topic).WithField("msg", msg).Error("parse message error")
				return
			}
			s := r.ToSql()
			if config.Conf.ShowSql {
				logger.WithField("topic", topic).Infof("insert sql %s", s)
			}
			_, err = db.Exec(s)
			if err != nil {
				logger.WithError(err).WithField("topic", topic).WithField("msg", msg).WithField("sql", s).Error("execute insert error")
			}
		}
	})
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	logger.Println("stop server")
}

func initDB(c *config.TDengine, logger *logrus.Entry) {
	dsnWithoutDB := fmt.Sprintf("%s:%s/tcp(%s:%d)/", c.User, c.Password, c.Host, c.Port)
	db, err := sql.Open("taosSql", dsnWithoutDB)
	if err != nil {
		logger.WithError(err).Panic("connect TDengine error")
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("create database if not exists %s", c.DB))
	if err != nil {
		logger.WithError(err).Panic("execute create db error")
	}
}
