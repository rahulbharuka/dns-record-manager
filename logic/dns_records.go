package logic

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rahulbharuka/dns-record-manager/external/route53"
	"github.com/rahulbharuka/dns-record-manager/model"
	"github.com/rahulbharuka/dns-record-manager/repository"
)

// ListDNSARecords lists all DNA A records published along with corresponding server metadata.
func (h *handlerImpl) ListDNSARecords(ctx *gin.Context) {
	// fetch the list of all DNS A records
	records, err := route53.ListARecords(ctx, "*")
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	// create temporary maps to find corresponding server and cluster metadata.
	ips := []string{}
	subdomains := []string{}
	clusterMap := map[string]*repository.Cluster{}
	for _, r := range records {
		ips = append(ips, r.IP)
		sd := strings.TrimSuffix(r.FQDN, route53.Domain)
		if _, ok := clusterMap[sd]; !ok {
			subdomains = append(subdomains, sd)
			clusterMap[sd] = nil
		}
	}

	// find corresponding servers
	servers, err := h.serverRepo.FindByIPs(ips)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	serverMap := map[string]*model.Server{}
	for _, s := range servers {
		serverMap[s.IP] = s
	}

	// find corresponding clusters
	clusters, err := h.clusterRepo.FindBySubdomains(subdomains)
	if err != nil {
		handlerError(ctx, http.StatusInternalServerError, err)
		return
	}

	for _, c := range clusters {
		clusterMap[c.Subdomain] = c
	}

	// populate metadata for each DNS record.
	dnsRecords := []*model.DNSRecord{}
	for _, r := range records {
		serverName := "not found"
		clusterName := "N/A"
		if server, ok := serverMap[r.IP]; ok {
			serverName = server.Name
		}
		if cluster, ok := clusterMap[strings.TrimSuffix(r.FQDN, route53.Domain)]; ok && cluster != nil {
			clusterName = cluster.Name
		}

		dnsRecords = append(dnsRecords,
			&model.DNSRecord{
				FQDN:        r.FQDN,
				IP:          r.IP,
				ServerName:  serverName,
				ClusterName: clusterName,
			})
	}

	ctx.HTML(http.StatusOK, "dns_records.html", gin.H{
		"dns_records": dnsRecords,
	})
}
