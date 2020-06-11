package mongodb

import (
	"time"

	"github.com/globalsign/mgo"
	envconfig "github.com/praslar/to-do-list-micro/pkg/env"
	"github.com/sirupsen/logrus"
)

type (
	// Config hold MongoDB configuration information
	Config struct {
		Addrs    []string      `envconfig:"MONGODB_ADDRS" default:"localhost:27017"`
		Database string        `envconfig:"MONGODB_DATABASE" default:"gogodo"`
		Username string        `envconfig:"MONGODB_USERNAME"`
		Password string        `envconfig:"MONGODB_PASSWORD"`
		Timeout  time.Duration `envconfig:"MONGODB_TIMEOUT" default:"10s"`
	}
)

// LoadConfigFromEnv load mongodb configurations from environments
func LoadConfigFromEnv() *Config {
	var conf Config
	envconfig.Load(&conf)
	return &conf
}

// Dial dial to target server with Monotonic mode
func Dial(conf *Config) (*mgo.Session, error) {
	logrus.Info(conf.Addrs)
	logrus.Infof("Dialing to target MongoDB at: %v, Database: %v", conf.Addrs, conf.Database)
	ms, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    conf.Addrs,
		Database: conf.Database,
		Username: conf.Username,
		Password: conf.Password,
		Timeout:  conf.Timeout,
	})
	if err != nil {
		return nil, err
	}
	ms.SetMode(mgo.Monotonic, true)
	logrus.Infof("successfully dialing to MongoDB at %v", conf.Addrs)
	return ms, nil
}

// DialInfo return dial info from config
func (conf *Config) DialInfo() *mgo.DialInfo {
	return &mgo.DialInfo{
		Addrs:    conf.Addrs,
		Database: conf.Database,
		Username: conf.Username,
		Password: conf.Password,
		Timeout:  conf.Timeout,
	}
}
