package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/rahulbharuka/dns-record-manager/model"
	"github.com/rahulbharuka/dns-record-manager/storage"
)

var (
	// initOnce protects the following
	initServerRepoOnce  sync.Once
	singletonServerRepo *serverRepoImpl
)

// Server is a storage object for server table.
type Server struct {
	ID              uint64 `json:"id"`
	Name            string `json:"name"`
	ClusterID       uint64 `json:"cluster_id"`
	IP              string `json:"ip"`
	AddedToRotation bool   `json:"added_to_rotation"`
}

// TableName returns table name.
func (s Server) TableName() string {
	return "server"
}

// serverRepoImpl ...
type serverRepoImpl struct {
	db storage.Handler
}

// ServerRepo implements following methods.
// go:generate mockery -inpkg -case underscore -name ClusterRepo
type ServerRepo interface {
	ListAll() ([]*model.Server, error)
	ListAllWithClusterInfo() ([]*model.Server, error)
	AddToRotation(id uint64) error
	RemoveFromRotation(id uint64) error
	LoadByID(id uint64) (*model.Server, error)
	FindByIPs(ips []string) ([]*model.Server, error)
}

// NewServerRepo returns the ServerRepo handler.
func NewServerRepo() ServerRepo {
	initServerRepoOnce.Do(func() {
		singletonServerRepo = &serverRepoImpl{
			db: storage.NewDBHandler(),
		}
	})
	return singletonServerRepo
}

// ListAll lists all servers in db.
func (r *serverRepoImpl) ListAll() ([]*model.Server, error) {
	query := fmt.Sprintf(selectAllQuery, "server")
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to list all servers, err: %v", err)
		return nil, err
	}

	return r.rowsToServerModel(rows)
}

// ListAllWithClusterInfo lists all servers with corresponding cluster metadata.
func (r *serverRepoImpl) ListAllWithClusterInfo() ([]*model.Server, error) {
	query := fmt.Sprintf("select s.id, s.name, s.ip, s.added_to_rotation, c.name, c.subdomain from server s inner join cluster c where s.cluster_id = c.id")
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to list all servers with cluster info, err: %v", err)
		return nil, err
	}

	res := []*model.Server{}
	for rows.Next() {
		entity := &model.Server{}
		err := rows.Scan(&entity.ID, &entity.Name, &entity.IP, &entity.AddedToRotation, &entity.ClusterName, &entity.Subdomain)
		if err != nil {
			log.Printf("failed to convert sql.Row to server model, err: %v", err)
			return nil, err
		}
		res = append(res, entity)
	}
	return res, nil
}

// AddToRotation marks the server as added to rotation.
func (r *serverRepoImpl) AddToRotation(id uint64) error {
	query := fmt.Sprintf(updateQuery, "server", "added_to_rotation", 1, id)
	_, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to add server %v to rotation, err: %v", id, err)
		return err
	}
	return nil
}

// RemoveFromRotation marks the server as not added to rotation.
func (r *serverRepoImpl) RemoveFromRotation(id uint64) error {
	query := fmt.Sprintf(updateQuery, "server", "added_to_rotation", 0, id)
	_, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to remove server %v from rotation, err: %v", id, err)
		return err
	}
	return nil
}

// LoadByID loads a server by ID.
func (r *serverRepoImpl) LoadByID(id uint64) (*model.Server, error) {
	query := fmt.Sprintf(selectByIDQuery, "server", id)
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to load server %v, err: %v", id, err)
		return nil, err
	}

	servers, err := r.rowsToServerModel(rows)
	if err != nil {
		return nil, err
	}

	return servers[0], nil
}

// FindByIPs finds servers for given IPs.
func (r *serverRepoImpl) FindByIPs(ips []string) ([]*model.Server, error) {
	var serverIPs []string
	for _, ip := range ips {
		serverIPs = append(serverIPs, fmt.Sprintf("'%v'", ip))
	}
	query := fmt.Sprintf(findByQuery, "server", "ip", strings.Join(serverIPs, ","))
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to find servers for %v ips, err: %v", ips, err)
		return nil, err
	}

	servers, err := r.rowsToServerModel(rows)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

// rowsToServerModel converts sql.Rows to server model object.
func (r *serverRepoImpl) rowsToServerModel(rows *sql.Rows) ([]*model.Server, error) {
	res := []*model.Server{}
	for rows.Next() {
		entity := &model.Server{}
		err := rows.Scan(&entity.ID, &entity.Name, &entity.ClusterID, &entity.IP, &entity.AddedToRotation)
		if err != nil {
			log.Printf("failed to convert sql.Rows to server model, err: %v", err)
			return nil, err
		}
		res = append(res, entity)
	}

	return res, nil
}
