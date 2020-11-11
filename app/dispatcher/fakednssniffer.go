// +build !confonly

package dispatcher

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/features/dns"
)

// newFakeDNSSniffer Create a Fake DNS metadata sniffer
func newFakeDNSSniffer(ctx context.Context) (protocolSnifferWithMetadata, error) {
	var fakeDNSEngine dns.FakeDNSEngine
	err := core.RequireFeatures(ctx, func(fdns dns.FakeDNSEngine) {
		fakeDNSEngine = fdns
	})
	if err != nil {
		return protocolSnifferWithMetadata{}, err
	}
	return protocolSnifferWithMetadata{protocolSniffer: func(bytes []byte, ctx context.Context) (SniffResult, error) {
		Target := session.OutboundFromContext(ctx).Target
		if Target.Network == net.Network_TCP || Target.Network == net.Network_UDP {
			domainFromFakeDNS := fakeDNSEngine.GetDomainFromFakeDNS(Target.Address)
			if domainFromFakeDNS != "" {
				newError("fake dns got domain: ", domainFromFakeDNS, " for ip: ", Target.Address.String()).WriteToLog(session.ExportIDToError(ctx))
				return &fakeDNSSniffResult{domainName: domainFromFakeDNS}, nil
			}
		}
		return nil, common.ErrNoClue
	}, metadataSniffer: true}, nil
}

type fakeDNSSniffResult struct {
	domainName string
}

func (f fakeDNSSniffResult) Protocol() string {
	return "fakedns"
}

func (f fakeDNSSniffResult) Domain() string {
	return f.domainName
}
