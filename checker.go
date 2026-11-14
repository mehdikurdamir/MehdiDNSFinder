package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type Result struct {
	Name    string   `json:"name"`
	DNS     string   `json:"dns"`
	Group   string   `json:"group"`
	Status  string   `json:"status"` // ok | warn | noconn | blocked | nxdomain | down
	Message string   `json:"message"`
	IPs     []string `json:"ips"`
	DNSms   int64    `json:"dnsMs"`
	ConnMs  int64    `json:"connMs"`
}

type Target struct {
	Host string
	Port string
}

func parseTarget(in string) (Target, error) {
	in = strings.TrimSpace(in)
	if in == "" {
		return Target{}, errors.New("آدرس خالی است")
	}
	if !strings.Contains(in, "://") {
		in = "https://" + in
	}
	u, err := url.Parse(in)
	if err != nil || u.Hostname() == "" {
		return Target{}, errors.New("آدرس نامعتبر است")
	}
	port := u.Port()
	if port == "" {
		if u.Scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}
	return Target{Host: strings.ToLower(u.Hostname()), Port: port}, nil
}

func isFakeIP(ip net.IP) (bool, string) {
	if ip.IsUnspecified() || ip.IsLoopback() {
		return true, "آی‌پی بی‌معنی (۰.۰.۰.۰ / لوپ‌بک)"
	}
	if ip4 := ip.To4(); ip4 != nil && ip4[0] == 10 && ip4[1] == 10 && ip4[2] == 34 {
		return true, "صفحه فیلترینگ (10.10.34.x)"
	}
	if ip.IsPrivate() || ip.IsLinkLocalUnicast() {
		return true, "آی‌پی خصوصی (احتمالاً مسدودسازی)"
	}
	return false, ""
}

func queryDNS(ctx context.Context, server, host string) ([]net.IP, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 3 * time.Second}
			return d.DialContext(ctx, "udp", net.JoinHostPort(server, "53"))
		},
	}
	return r.LookupIP(ctx, "ip4", host)
}

func checkOne(s DNSServer, t Target) Result {
	res := Result{Name: s.Name, DNS: s.IP, Group: s.Group, IPs: []string{}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	start := time.Now()
	ips, err := queryDNS(ctx, s.IP, t.Host)
	res.DNSms = time.Since(start).Milliseconds()

	if err != nil || len(ips) == 0 {
		var de *net.DNSError
		if errors.As(err, &de) && de.IsNotFound {
			res.Status = "nxdomain"
			res.Message = "دامنه پیدا نشد (NXDOMAIN)"
		} else {
			res.Status = "down"
			res.Message = "سرور DNS پاسخ نداد"
		}
		return res
	}

	good := []net.IP{}
	reason := ""
	for _, ip := range ips {
		res.IPs = append(res.IPs, ip.String())
		if fake, why := isFakeIP(ip); fake {
			reason = why
		} else {
			good = append(good, ip)
		}
	}
	if len(good) == 0 {
		res.Status = "blocked"
		res.Message = reason
		return res
	}

	// تست اتصال واقعی به اولین آی‌پی‌های سالم
	var lastErr error
	limit := len(good)
	if limit > 3 {
		limit = 3
	}
	for _, ip := range good[:limit] {
		cs := time.Now()
		d := net.Dialer{Timeout: 4 * time.Second}
		conn, err := d.Dial("tcp", net.JoinHostPort(ip.String(), t.Port))
		if err != nil {
			lastErr = err
			continue
		}
		res.ConnMs = time.Since(cs).Milliseconds()
		if t.Port == "80" {
			conn.Close()
			res.Status = "ok"
			res.Message = "دسترسی برقرار است"
			return res
		}
		// TLS با تأیید گواهی
		tc := tls.Client(conn, &tls.Config{ServerName: t.Host})
		_ = tc.SetDeadline(time.Now().Add(5 * time.Second))
		herr := tc.Handshake()
		tc.Close()
		if herr == nil {
			res.Status = "ok"
			res.Message = "دسترسی برقرار است"
			return res
		}
		// دوباره بدون تأیید گواهی تا بفهمیم فقط گواهی مشکل دارد یا کل اتصال
		conn2, err2 := d.Dial("tcp", net.JoinHostPort(ip.String(), t.Port))
		if err2 == nil {
			tc2 := tls.Client(conn2, &tls.Config{ServerName: t.Host, InsecureSkipVerify: true})
			_ = tc2.SetDeadline(time.Now().Add(5 * time.Second))
			if tc2.Handshake() == nil {
				tc2.Close()
				res.Status = "warn"
				res.Message = "متصل شد ولی گواهی معتبر نیست (احتمال دستکاری)"
				return res
			}
			tc2.Close()
		}
		lastErr = herr
	}
	res.Status = "noconn"
	if lastErr != nil && strings.Contains(lastErr.Error(), "timeout") {
		res.Message = "آی‌پی دریافت شد ولی اتصال برقرار نشد (تایم‌اوت)"
	} else {
		res.Message = fmt.Sprintf("آی‌پی دریافت شد ولی اتصال برقرار نشد")
	}
	return res
}
