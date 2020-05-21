package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/rahulbharuka/dns-record-manager/storage"
)

var (
	// initOnce protects the following
	initClusterRepoOnce  sync.Once
	singletonClusterRepo *clusterRepoImpl
)

// Cluster is a storage object for cluster table.
type Cluster struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
}

// TableName returns table name.
func (c Cluster) TableName() string {
	return "cluster"
}

// clusterRepoImpl ...
type clusterRepoImpl struct {
	db storage.Handler
}

// ClusterRepo implements following methods.
// go:generate mockery -inpkg -case underscore -name ClusterRepo
type ClusterRepo interface {
	ListAll() ([]*Cluster, error)
	FindByIDs(ids []uint64) ([]*Cluster, error)
	LoadByID(id uint64) (*Cluster, error)
	FindBySubdomains(subdomains []string) ([]*Cluster, error)
}

// NewClusterRepo returns the ClusterRepo handler.
func NewClusterRepo() ClusterRepo {
	initClusterRepoOnce.Do(func() {
		singletonClusterRepo = &clusterRepoImpl{
			db: storage.NewDBHandler(),
		}
	})
	return singletonClusterRepo
}

// ListAll lists all clusters
func (r *clusterRepoImpl) ListAll() ([]*Cluster, error) {
	query := fmt.Sprintf(selectAllQuery, "cluster")
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to list all clusters, err: %v", err)
		return nil, err
	}

	return r.rowsToStorageObjects(rows)
}

// Find finds all clusters by given IDs.
func (r *clusterRepoImpl) FindByIDs(ids []uint64) ([]*Cluster, error) {
	var clusterIDs []string
	for _, id := range ids {
		clusterIDs = append(clusterIDs, strconv.FormatUint(id, 10))
	}
	query := fmt.Sprintf(findByQuery, "cluster", "id", strings.Join(clusterIDs, ","))
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to find clusters for IDs %v, err: %v", clusterIDs, err)
		return nil, err
	}

	return r.rowsToStorageObjects(rows)
}

// LoadByID loads a cluster by ID.
func (r *clusterRepoImpl) LoadByID(id uint64) (*Cluster, error) {
	query := fmt.Sprintf(selectByIDQuery, "cluster", id)
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to load cluster ID %v, err: %v", id, err)
		return nil, err
	}

	clusters, err := r.rowsToStorageObjects(rows)
	if err != nil {
		return nil, err
	}

	return clusters[0], nil
}

// FindBySubdomains finds all clusters by subdomains.
func (r *clusterRepoImpl) FindBySubdomains(subdomains []string) ([]*Cluster, error) {
	var arr []string
	for _, sd := range subdomains {
		arr = append(arr, fmt.Sprintf("'%v'", sd))
	}
	query := fmt.Sprintf(findByQuery, "cluster", "subdomain", strings.Join(arr, ","))
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("failed to load clusters for subdomains %v, err: %v", subdomains, err)
		return nil, err
	}

	return r.rowsToStorageObjects(rows)
}

// rowsToStorageObjects converts a sql.Rows to storage object format.
func (r *clusterRepoImpl) rowsToStorageObjects(rows *sql.Rows) ([]*Cluster, error) {
	res := []*Cluster{}
	for rows.Next() {
		entity := &Cluster{}
		err := rows.Scan(&entity.ID, &entity.Name, &entity.Subdomain)
		if err != nil {
			log.Printf("failed to convert sql.Rows to Cluster object, err: %v", err)
			return nil, err
		}
		res = append(res, entity)
	}
	return res, nil
}
