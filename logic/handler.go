package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/dns-record-manager/repository"
)

// Handler is the logic handler interface
type Handler interface {
	ListServers(ctx *gin.Context)
	AddServer(ctx *gin.Context)
	RemoveServer(ctx *gin.Context)
	ListDNSARecords(ctx *gin.Context)
}

// handlerImpl is a implementation of Handler interface
type handlerImpl struct {
	clusterRepo repository.ClusterRepo
	serverRepo  repository.ServerRepo
}

// GetHandler initializes and returns the logic layer handler.
func GetHandler() Handler {
	return &handlerImpl{
		clusterRepo: repository.NewClusterRepo(),
		serverRepo:  repository.NewServerRepo(),
	}
}

// handlerError is a helper function to return JSON error.
func handlerError(ctx *gin.Context, errCode int, err error) {
	ctx.JSON(errCode, gin.H{
		"message": err.Error(),
	})
}
