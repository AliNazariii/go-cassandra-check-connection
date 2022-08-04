package database

import (
	"github.com/gocql/gocql"
	"strings"
	"cassandra_connection_check/pkg/config"
	"cassandra_connection_check/pkg/log"
	"time"
)

type CassandraDB struct {
	log      *log.Logger
	Session  *gocql.Session
	KeySpace string
}

func NewCassandraDB(log *log.Logger, conf *config.SectionCassandra) *CassandraDB {
	log.Debugln(conf.Consistency, conf.Hosts, conf.KeySpace, conf.User, conf.Password, conf.PageSize, conf.Port, conf.Timeout, conf.DataCenter)
	db := CassandraDB{KeySpace: conf.KeySpace, log: log}

	consistency, err := gocql.ParseConsistencyWrapper(conf.Consistency)
	if err != nil {
		consistency = gocql.LocalOne
		log.Infof("Error in consistency, so set to %s", consistency)
	}

	cluster := gocql.NewCluster(conf.Hosts...)
	cluster.Keyspace = conf.KeySpace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.User,
		Password: conf.Password}
	cluster.PageSize = conf.PageSize
	cluster.Port = conf.Port
	cluster.Consistency = consistency
	cluster.Timeout = time.Duration(conf.Timeout) * time.Millisecond
	cluster.WriteCoalesceWaitTime = time.Millisecond
	cluster.NumConns = 10
	cluster.ProtoVersion = 3
	//cluster.CQLVersion = "3.4.5"
	if conf.DataCenter != "" {
		cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy(conf.DataCenter))
		cluster.HostFilter = gocql.DataCentreHostFilter(conf.DataCenter)
	}

	db.Session, err = cluster.CreateSession()

	if err != nil {
		log.Fatalf("Unable to create cassandra session. %s", err)
	}
	return &db
}

func (c *CassandraDB) CreateTable(table string, createQueryString string) error {
	keySpace, err := c.Session.KeyspaceMetadata(c.KeySpace)
	if err != nil {
		c.log.Errorf("Error in getting keyspace metadata. Error: %s", err)
		return err
	}
	if _, exists := keySpace.Tables[table]; exists != true {
		c.log.Infof("Creating table %s", table)
		queries := strings.Split(createQueryString, ";")
		for _, query := range queries {
			if query != "" {
				if err := c.Session.Query(query).Exec(); err != nil {
					c.log.Errorf("Error in running query in creating table. Error: %s", err)
					return err
				}
			}
		}
	}
	return nil
}
