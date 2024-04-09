package cassandra

import "github.com/gocql/gocql"

type CassandraDBConnection struct {
	Cluster *gocql.ClusterConfig
	Session *gocql.Session
}

func NewCassandraDBConection() (*CassandraDBConnection, error) {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = "chatapp"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &CassandraDBConnection{
		Cluster: cluster,
		Session: session,
	}, nil
}
