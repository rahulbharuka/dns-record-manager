package logic

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/dns-record-manager/external/route53"
)

// ListServers lists all servers along with corresponding cluster metadata.
func (h *handlerImpl) ListServers(ctx *gin.Context) {
	servers, err := h.serverRepo.ListAllWithClusterInfo()
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "servers.html", gin.H{
		"servers": servers,
		"domain":  ".domain.com",
	})
}

// AddServer adds DNA A record for the server.
func (h *handlerImpl) AddServer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("failed to add server ID %v, err: %v", id, err)
		handlerError(ctx, http.StatusBadRequest, err)
		return
	}

	server, err := h.serverRepo.LoadByID(uint64(id))
	if err != nil {
		log.Printf("failed to load server ID %v, err: %v", id, err)
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}
	if server.AddedToRotation {
		ctx.String(http.StatusOK, "Server already added to rotation")
		return
	}

	cluster, err := h.clusterRepo.LoadByID(uint64(server.ClusterID))
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = route53.AddServer(ctx, cluster.Subdomain, server.IP)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.serverRepo.AddToRotation(uint64(id))
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.String(200, "Successfully added server to rotation")
}

// RemoveServer removes DNS A record for the server.
func (h *handlerImpl) RemoveServer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	server, err := h.serverRepo.LoadByID(uint64(id))
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}
	if !server.AddedToRotation {
		ctx.String(200, "Server is currently NOT added to rotation")
		return
	}

	cluster, err := h.clusterRepo.LoadByID(uint64(server.ClusterID))
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = route53.RemoveServer(ctx, cluster.Subdomain, server.IP)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.serverRepo.RemoveFromRotation(uint64(id))
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.String(http.StatusOK, "Successfully removed server from rotation")
}
