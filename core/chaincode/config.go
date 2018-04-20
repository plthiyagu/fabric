/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

const (
	defaultExecutionTimeout = 30 * time.Second
	minimumStartupTimeout   = 5 * time.Second
)

type Config struct {
	PeerNetworkID  string
	PeerID         string
	TLSEnabled     bool
	Keepalive      time.Duration
	ExecuteTimeout time.Duration
	StartupTimeout time.Duration
	LogFormat      string
	LogLevel       string
	ShimLogLevel   string
}

func GlobalConfig() *Config {
	c := &Config{}
	c.load()
	return c
}

func (c *Config) load() {
	viper.SetEnvPrefix("CORE")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	c.PeerNetworkID = viper.GetString("peer.networkId")
	c.PeerID = viper.GetString("peer.id")
	c.TLSEnabled = viper.GetBool("peer.tls.enabled")

	c.Keepalive = toSeconds(viper.GetString("chaincode.keepalive"), 0)
	c.ExecuteTimeout = viper.GetDuration("chaincode.executetimeout")
	if c.ExecuteTimeout < time.Second {
		c.ExecuteTimeout = defaultExecutionTimeout
	}
	c.StartupTimeout = viper.GetDuration("chaincode.startuptimeout")
	if c.StartupTimeout < minimumStartupTimeout {
		c.StartupTimeout = 5 * time.Second
	}

	c.LogFormat = viper.GetString("chaincode.logging.format")
	c.LogLevel = getLogLevelFromViper("chaincode.logging.level")
	c.ShimLogLevel = getLogLevelFromViper("chaincode.logging.shim")
}

func toSeconds(s string, def int) time.Duration {
	seconds, err := strconv.Atoi(s)
	if err != nil {
		return time.Duration(def) * time.Second
	}

	return time.Duration(seconds) * time.Second
}

// getLogLevelFromViper gets the chaincode container log levels from viper
func getLogLevelFromViper(key string) string {
	levelString := viper.GetString(key)
	_, err := logging.LogLevel(levelString)
	if err != nil {
		chaincodeLogger.Warningf("%s has invalid log level %s. defaulting to %s", key, levelString, flogging.DefaultLevel())
		levelString = flogging.DefaultLevel()
	}

	return levelString
}
