package helpers

import (
	"context"
	"errors"
	"net"
)

func DomainVerificationHost(domain string) string {
	return "_grovehub." + domain
}

func DomainVerificationValue(token string) string {
	return "grovehub-verification=" + token
}

func VerifyDomainTXT(ctx context.Context, domain string, token string) (bool, error) {
	values, err := net.DefaultResolver.LookupTXT(ctx, DomainVerificationHost(domain))
	if err != nil {
		var dnsError *net.DNSError
		if errors.As(err, &dnsError) && dnsError.IsNotFound {
			return false, nil
		}
		return false, err
	}

	expected := DomainVerificationValue(token)
	for _, value := range values {
		if value == expected {
			return true, nil
		}
	}
	return false, nil
}
