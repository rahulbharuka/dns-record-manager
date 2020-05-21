package route53

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var (
	// client for AWS Route53 service.
	client *route53.Route53
	// Domain is the domain under management
	Domain string
	// hostedZoneID is the hosted zone created for the domain.
	hostedZoneID string
	// ErrNotFound is resource not found error
	ErrNotFound = errors.New("resource not found")
)

// Init initializes the AWS Route53 service client.
func Init() {
	Domain = os.Getenv("DOMAIN")
	hostedZoneID = os.Getenv("HOSTED_ZONE_ID")

	sess, err := session.NewSession()
	if err != nil {
		log.Printf("failed to initialize AWS Route53 service client, err: %v", err)
		return
	}
	client = route53.New(sess)
}

// DNSRecord is an object for DNS record.
type DNSRecord struct {
	FQDN string `json:"fqdn"`
	IP   string `json:"ip"`
}

// ListARecords list all DNS A records for given subdomain.
func ListARecords(ctx context.Context, subdomain string) ([]*DNSRecord, error) {
	rrSets, err := client.ListResourceRecordSetsWithContext(ctx,
		&route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(hostedZoneID),
			StartRecordName: aws.String(subdomain + Domain),
			StartRecordType: aws.String(route53.RRTypeA),
		})
	if err != nil {
		log.Printf("failed to retreive A records for subdomain %v, err: %+v", subdomain, err)
		return nil, err
	}

	records := []*DNSRecord{}
	for _, rrset := range rrSets.ResourceRecordSets {
		if *rrset.Type != route53.RRTypeA {
			continue
		}
		for _, rr := range rrset.ResourceRecords {
			rec := &DNSRecord{
				FQDN: *rrset.Name,
				IP:   *rr.Value,
			}
			records = append(records, rec)
		}
	}
	return records, nil
}

// AddServer adds a server(IP) to Route53 DNS A record of specified subdomain.
// It creates a new DNS A record for the subdomain if it does not exist.
// If the server is already added, its a NO-OP.
func AddServer(ctx context.Context, subdomain string, ip string) error {
	updatedResourceRecords := []*route53.ResourceRecord{{Value: aws.String(ip)}}

	// First check whether the subdomain DNS entry already exist
	rrSets, err := client.ListResourceRecordSetsWithContext(ctx,
		&route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(hostedZoneID),
			StartRecordName: aws.String(subdomain + Domain),
			StartRecordType: aws.String(route53.RRTypeA),
			MaxItems:        aws.String(strconv.Itoa(1)),
		})
	if err != nil {
		log.Printf("failed to retreive A record for subdomain %v, err: %+v", subdomain, err)
		return err
	}

	// If subdomain DNS A record exists, check whether the server is already added.
	if len(rrSets.ResourceRecordSets) > 0 && *rrSets.ResourceRecordSets[0].Name == subdomain+Domain {
		for _, rr := range rrSets.ResourceRecordSets[0].ResourceRecords {
			if *rr.Value == ip {
				log.Printf("Server %v is already added for subdomain %v", ip, subdomain)
				return nil
			}
		}
		// The server is not added for the subdmoain. So prepare to add it.
		updatedResourceRecords = append(updatedResourceRecords, rrSets.ResourceRecordSets[0].ResourceRecords...)
	}

	// add new server to the rotation.
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            aws.String(subdomain + Domain),
						ResourceRecords: updatedResourceRecords,
						TTL:             aws.Int64(300),
						Type:            aws.String("A"),
					},
				},
			},
		},
		HostedZoneId: aws.String(hostedZoneID),
	}

	_, err = client.ChangeResourceRecordSets(input)
	if err != nil {
		return err
	}
	return nil
}

// RemoveServer removes a server(IP) from Route53 DNS A record of given subdomain.
// It also deletes the DNS A record if this was the last server.
// If the server does not exist, its a NO-OP.
func RemoveServer(ctx context.Context, subdomain string, ip string) error {
	// First check whether the subdomain DNS entry exist.
	rrSets, err := client.ListResourceRecordSetsWithContext(ctx,
		&route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(hostedZoneID),
			StartRecordName: aws.String(subdomain + Domain),
			StartRecordType: aws.String(route53.RRTypeA),
			MaxItems:        aws.String(strconv.Itoa(1)),
		})
	if err != nil {
		log.Printf("failed to retreive DNS A record for subdomain %v, err: %+v", subdomain, err)
		return err
	}

	if len(rrSets.ResourceRecordSets) == 0 || *rrSets.ResourceRecordSets[0].Name != subdomain+Domain {
		log.Printf("No DNS A record found for subdomain %v", subdomain)
		return nil
	}

	// If DNS entry exists, check whether the server is added.
	idx := -1
	for i, rr := range rrSets.ResourceRecordSets[0].ResourceRecords {
		if *rr.Value == ip {
			idx = i
			break
		}
	}

	if idx == -1 {
		log.Printf("The server %v is not added for subdomain %v", ip, subdomain)
		return ErrNotFound
	}

	// Remove the server from subdomain server list.
	updatedResourceRecords := []*route53.ResourceRecord{}
	updatedResourceRecords = append(updatedResourceRecords, rrSets.ResourceRecordSets[0].ResourceRecords[:idx]...)
	if idx < len(rrSets.ResourceRecordSets[0].ResourceRecords)-1 {
		updatedResourceRecords = append(updatedResourceRecords, rrSets.ResourceRecordSets[0].ResourceRecords[idx+1:]...)
	}

	// First, delete the DNS A record for the subdomain.
	dInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            aws.String(subdomain + Domain),
						ResourceRecords: rrSets.ResourceRecordSets[0].ResourceRecords,
						TTL:             aws.Int64(300),
						Type:            aws.String("A"),
					},
				},
			},
		},
		HostedZoneId: aws.String(hostedZoneID),
	}

	_, err = client.ChangeResourceRecordSets(dInput)
	if err != nil {
		log.Printf("failed to delete DNS A record for subdomain %v, err: %v", subdomain, err)
		return err
	}

	// if this was the only server for the subdomain, no need to add back the record. So just return.
	if len(updatedResourceRecords) == 0 {
		return nil
	}

	// Now, add back the subdomain DNS A record for rest of the servers as it has multiple servers.
	cInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("CREATE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            aws.String(subdomain + Domain),
						ResourceRecords: updatedResourceRecords,
						TTL:             aws.Int64(300),
						Type:            aws.String("A"),
					},
				},
			},
		},
		HostedZoneId: aws.String(hostedZoneID),
	}

	_, err = client.ChangeResourceRecordSets(cInput)
	if err != nil {
		log.Printf("failed to add back DNS A record for subdomain %v, err: %v", subdomain, err)
		return err
	}

	return nil
}
