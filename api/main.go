package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// getenvDefault returns env var or default if unset
func getenvDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

// DNS server configuration
type DNSServer struct {
	Name string
	IP   string
}

// BlocklistResult represents the result of checking a domain against a DNS server
type BlocklistResult struct {
	Server    string `json:"server"`
	ServerIP  string `json:"server_ip"`
	IsBlocked bool   `json:"is_blocked"`
}

// BlocklistInfo represents blocklist checking information
type BlocklistInfo struct {
	Domain  string            `json:"domain"`
	Results []BlocklistResult `json:"results"`
	Error   string            `json:"error,omitempty"`
}

// SSLInfo represents SSL certificate information
type SSLInfo struct {
	Domain          string    `json:"domain"`
	Valid           bool      `json:"valid"`
	Issuer          string    `json:"issuer"`
	Subject         string    `json:"subject"`
	NotBefore       time.Time `json:"not_before"`
	NotAfter        time.Time `json:"not_after"`
	DaysUntilExpiry int       `json:"days_until_expiry"`
	SerialNumber    string    `json:"serial_number"`
	SignatureAlg    string    `json:"signature_algorithm"`
	PublicKeyAlg    string    `json:"public_key_algorithm"`
	KeySize         int       `json:"key_size"`
	Error           string    `json:"error,omitempty"`
}

// HTTP3Info represents HTTP/3 detection information
type HTTP3Info struct {
	Domain    string `json:"domain"`
	Supported bool   `json:"supported"`
	Protocol  string `json:"protocol"`
	Status    int    `json:"status"`
	Details   string `json:"details"`
	Error     string `json:"error,omitempty"`
}

// DNSInfo represents DNS resolution information
type DNSInfo struct {
	Domain string   `json:"domain"`
	IPv4   []string `json:"ipv4"`
	IPv6   []string `json:"ipv6"`
	CNAME  []string `json:"cname"`
	MX     []string `json:"mx"`
	TXT    []string `json:"txt"`
	NS     []string `json:"ns"`
	Error  string   `json:"error,omitempty"`
}

// IPInfo represents IP address information
type IPInfo struct {
	Input        string   `json:"input"`
	IsDomain     bool     `json:"is_domain"`
	ResolvedIPs  []string `json:"resolved_ips,omitempty"`
	IP           string   `json:"ip,omitempty"`
	Country      string   `json:"country"`
	Region       string   `json:"region"`
	City         string   `json:"city"`
	ISP          string   `json:"isp"`
	Organization string   `json:"organization"`
	Timezone     string   `json:"timezone"`
	Error        string   `json:"error,omitempty"`
}

// HSTSInfo represents HTTP Strict Transport Security information
type HSTSInfo struct {
	Enabled           bool   `json:"enabled"`
	MaxAge            int    `json:"max_age"`
	IncludeSubDomains bool   `json:"include_subdomains"`
	Preload           bool   `json:"preload"`
	Directive         string `json:"directive"`
	Details           string `json:"details"`
}

// RobotsTxtInfo represents robots.txt information
type RobotsTxtInfo struct {
	Domain     string   `json:"domain"`
	Exists     bool     `json:"exists"`
	Status     string   `json:"status"`
	Content    string   `json:"content"`
	Lines      []string `json:"lines"`
	UserAgents []string `json:"user_agents"`
	Disallowed []string `json:"disallowed"`
	Allowed    []string `json:"allowed"`
	Sitemaps   []string `json:"sitemaps"`
	CrawlDelay string   `json:"crawl_delay"`
	Error      string   `json:"error,omitempty"`
}

// SitemapURL represents a URL entry in a sitemap
type SitemapURL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

// SitemapIndex represents a sitemap index file
type SitemapIndex struct {
	XMLName  xml.Name `xml:"sitemapindex"`
	Sitemaps []struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	} `xml:"sitemap"`
}

// Sitemap represents a regular sitemap file
type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapInfo represents sitemap information
type SitemapInfo struct {
	Domain         string   `json:"domain"`
	SitemapURL     string   `json:"sitemap_url"`
	Exists         bool     `json:"exists"`
	Status         string   `json:"status"`
	IsSitemapIndex bool     `json:"is_sitemap_index"`
	URLCount       int      `json:"url_count"`
	SubSitemaps    []string `json:"sub_sitemaps,omitempty"`
	SampleURLs     []string `json:"sample_urls,omitempty"`
	LastModified   []string `json:"last_modified,omitempty"`
	Error          string   `json:"error,omitempty"`
}

// WebSettingsInfo represents web server settings and headers
type WebSettingsInfo struct {
	Domain        string            `json:"domain"`
	StatusCode    int               `json:"status_code"`
	Headers       map[string]string `json:"headers"`
	Server        string            `json:"server"`
	ContentType   string            `json:"content_type"`
	ContentLength int64             `json:"content_length"`
	LastModified  string            `json:"last_modified"`
	ETag          string            `json:"etag"`
	RedirectURL   string            `json:"redirect_url,omitempty"`
	HSTS          HSTSInfo          `json:"hsts"`
	ResponseTime  int64             `json:"response_time_ms"`
	Error         string            `json:"error,omitempty"`
}

// EmailConfigInfo represents email authentication configuration
type EmailConfigInfo struct {
	Domain string    `json:"domain"`
	SPF    SPFInfo   `json:"spf"`
	DKIM   DKIMInfo  `json:"dkim"`
	DMARC  DMARCInfo `json:"dmarc"`
	BIMI   BIMIInfo  `json:"bimi"`
	Error  string    `json:"error,omitempty"`
}

// SPFInfo represents SPF (Sender Policy Framework) information
type SPFInfo struct {
	Configured bool   `json:"configured"`
	Record     string `json:"record,omitempty"`
	Valid      bool   `json:"valid"`
	Details    string `json:"details,omitempty"`
	Error      string `json:"error,omitempty"`
}

// DKIMInfo represents DKIM (DomainKeys Identified Mail) information
type DKIMInfo struct {
	Configured bool     `json:"configured"`
	Selectors  []string `json:"selectors"`
	Valid      bool     `json:"valid"`
	Details    string   `json:"details,omitempty"`
	Error      string   `json:"error,omitempty"`
}

// DMARCInfo represents DMARC (Domain-based Message Authentication) information
type DMARCInfo struct {
	Configured bool   `json:"configured"`
	Record     string `json:"record,omitempty"`
	Policy     string `json:"policy,omitempty"` // none, quarantine, reject
	Valid      bool   `json:"valid"`
	Details    string `json:"details,omitempty"`
	Error      string `json:"error,omitempty"`
}

// BIMIInfo represents BIMI (Brand Indicators for Message Identification) information
type BIMIInfo struct {
	Configured bool   `json:"configured"`
	Record     string `json:"record,omitempty"`
	LogoURL    string `json:"logo_url,omitempty"`
	Valid      bool   `json:"valid"`
	Details    string `json:"details,omitempty"`
	Error      string `json:"error,omitempty"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// NetChecker handles all network checking operations
type NetChecker struct {
	httpClient  *http.Client
	http3Client *http.Client
}

// -----------------------------
// Rate limiting (per-IP, per-route)
// -----------------------------

type tokenBucket struct {
	capacity   int
	tokens     float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket // key: ip|route
	// per-route overrides: path pattern -> requests per minute
	perRoute   map[string]int
	defaultRpm int
}

func NewRateLimiter(defaultRpm int, perRoute map[string]int) *RateLimiter {
	return &RateLimiter{
		buckets:    make(map[string]*tokenBucket),
		perRoute:   perRoute,
		defaultRpm: defaultRpm,
	}
}

func (rl *RateLimiter) allow(ip string, route string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rpm := rl.defaultRpm
	if v, ok := rl.perRoute[route]; ok {
		rpm = v
	}
	if rpm <= 0 { // disabled
		return true
	}

	key := ip + "|" + route
	b, ok := rl.buckets[key]
	if !ok {
		b = &tokenBucket{
			capacity:   rpm,
			tokens:     float64(rpm),
			refillRate: float64(rpm) / 60.0, // per second
			lastRefill: time.Now(),
		}
		rl.buckets[key] = b
	}

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	if elapsed > 0 {
		b.tokens = minFloat(float64(b.capacity), b.tokens+elapsed*b.refillRate)
		b.lastRefill = now
	}
	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true
	}
	return false
}

// getRealClientIP extracts the real client IP from request headers
// Priority: X-Forwarded-For (first IP) > X-Real-IP > CF-Connecting-IP > RemoteAddr
func getRealClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (format: "client, proxy1, proxy2")
	// This is set by the frontend proxy with the real client IP
	xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Take the first IP in the chain (the original client)
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			// Validate it's a valid IP address
			if ip != "" && net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xRealIP := c.Request.Header.Get("X-Real-IP")
	if xRealIP != "" {
		ip := strings.TrimSpace(xRealIP)
		if ip != "" && net.ParseIP(ip) != nil {
			return ip
		}
	}

	// Check CF-Connecting-IP (Cloudflare)
	cfIP := c.Request.Header.Get("CF-Connecting-IP")
	if cfIP != "" {
		ip := strings.TrimSpace(cfIP)
		if ip != "" && net.ParseIP(ip) != nil {
			return ip
		}
	}

	// Fallback to RemoteAddr (removes port if present)
	remoteAddr := c.Request.RemoteAddr
	if remoteAddr != "" {
		// Remove port if present (format: "ip:port")
		if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
			return remoteAddr[:idx]
		}
		return remoteAddr
	}

	// Last resort: use Gin's ClientIP() which should handle trusted proxies
	return c.ClientIP()
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getRealClientIP(c)
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		if !rl.allow(ip, route) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, APIResponse{
				Success: false,
				Error:   "Too many requests",
			})
			return
		}
		c.Next()
	}
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// -----------------------------
// Simple in-memory TTL cache
// -----------------------------

type cacheItem struct {
	value  interface{}
	expiry time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewCache() *Cache {
	return &Cache{items: make(map[string]cacheItem)}
}

func (cch *Cache) Get(key string) (interface{}, bool) {
	cch.mu.RLock()
	it, ok := cch.items[key]
	cch.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(it.expiry) {
		cch.mu.Lock()
		delete(cch.items, key)
		cch.mu.Unlock()
		return nil, false
	}
	return it.value, true
}

func (cch *Cache) Set(key string, val interface{}, ttl time.Duration) {
	cch.mu.Lock()
	cch.items[key] = cacheItem{value: val, expiry: time.Now().Add(ttl)}
	cch.mu.Unlock()
}

var apiCache = NewCache()

// TTLs per route
var routeTTL = map[string]time.Duration{
	"/api/v1/ssl":          5 * time.Minute,
	"/api/v1/http3":        2 * time.Minute,
	"/api/v1/dns":          2 * time.Minute,
	"/api/v1/ip":           1 * time.Minute,
	"/api/v1/my-ip":        30 * time.Second, // Shorter TTL since it's user-specific
	"/api/v1/web-settings": 1 * time.Minute,
	"/api/v1/email-config": 10 * time.Minute,
	"/api/v1/blocklist":    10 * time.Minute,
	"/api/v1/robots-txt":   10 * time.Minute,
	"/api/v1/sitemap":      10 * time.Minute,
	"/api/v1/og-image":     10 * time.Minute,
}

func cacheKey(route string, q map[string]string) string {
	// build deterministic key
	b := strings.Builder{}
	b.WriteString(route)
	if len(q) > 0 {
		b.WriteString("?")
		first := true
		for k, v := range q {
			if !first {
				b.WriteString("&")
			} else {
				first = false
			}
			b.WriteString(k)
			b.WriteString("=")
			b.WriteString(v)
		}
	}
	return b.String()
}

// NewNetChecker creates a new NetChecker instance
func NewNetChecker() *NetChecker {
	// Standard HTTP client
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	// HTTP/3 client
	http3Transport := &http3.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	http3Client := &http.Client{
		Transport: http3Transport,
		Timeout:   15 * time.Second,
	}

	return &NetChecker{
		httpClient:  httpClient,
		http3Client: http3Client,
	}
}

// CheckSSL checks SSL certificate information
func (nc *NetChecker) CheckSSL(domain string) SSLInfo {
	info := SSLInfo{Domain: domain}

	// Clean domain
	cleanDomain := strings.TrimPrefix(domain, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.Split(cleanDomain, "/")[0]

	// Connect to the domain
	conn, err := tls.Dial("tcp", cleanDomain+":443", &tls.Config{
		ServerName: cleanDomain,
	})
	if err != nil {
		info.Error = fmt.Sprintf("Failed to connect: %v", err)
		return info
	}
	defer conn.Close()

	// Get certificate
	cert := conn.ConnectionState().PeerCertificates[0]
	info.Valid = true
	info.Issuer = cert.Issuer.String()
	info.Subject = cert.Subject.String()
	info.NotBefore = cert.NotBefore
	info.NotAfter = cert.NotAfter
	info.DaysUntilExpiry = int(time.Until(cert.NotAfter).Hours() / 24)
	info.SerialNumber = cert.SerialNumber.String()
	info.SignatureAlg = cert.SignatureAlgorithm.String()
	info.PublicKeyAlg = cert.PublicKeyAlgorithm.String()

	// Handle different public key types
	switch pubKey := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		info.KeySize = pubKey.N.BitLen()
	default:
		info.KeySize = 0 // Unknown key type
	}

	return info
}

// CheckHTTP3 checks HTTP/3 support
func (nc *NetChecker) CheckHTTP3(domain string) HTTP3Info {
	info := HTTP3Info{Domain: domain}

	// Ensure domain has protocol
	if !strings.HasPrefix(domain, "https://") && !strings.HasPrefix(domain, "http://") {
		domain = "https://" + domain
	}

	// Try HTTP/3 request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", domain, nil)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to create request: %v", err)
		return info
	}

	req.Header.Set("User-Agent", "NetCheck-API/1.0")

	resp, err := nc.http3Client.Do(req)
	if err != nil {
		// Try alternative QUIC connection
		cleanDomain := strings.TrimPrefix(domain, "https://")
		cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
		cleanDomain = strings.Split(cleanDomain, "/")[0]

		conn, err := quic.DialAddr(context.Background(), cleanDomain+":443", &tls.Config{
			ServerName: cleanDomain,
		}, &quic.Config{
			HandshakeIdleTimeout: 5 * time.Second,
		})

		if err == nil {
			conn.CloseWithError(0, "test connection")
			info.Supported = true
			info.Protocol = "HTTP/3"
			info.Details = "QUIC connection successful"
		} else {
			info.Supported = false
			info.Details = fmt.Sprintf("HTTP/3 not supported: %v", err)
		}
		return info
	}
	defer resp.Body.Close()

	info.Status = resp.StatusCode
	if resp.ProtoMajor == 3 {
		info.Supported = true
		info.Protocol = resp.Proto
		info.Details = fmt.Sprintf("HTTP/3 supported! Status: %d", resp.StatusCode)
	} else {
		info.Supported = false
		info.Protocol = resp.Proto
		info.Details = fmt.Sprintf("HTTP/3 not supported. Protocol: %s", resp.Proto)
	}

	return info
}

// CheckDNS checks DNS resolution information
func (nc *NetChecker) CheckDNS(domain string) DNSInfo {
	info := DNSInfo{Domain: domain}

	// Clean domain
	cleanDomain := strings.TrimPrefix(domain, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.Split(cleanDomain, "/")[0]

	// A records (IPv4)
	ips, err := net.LookupIP(cleanDomain)
	if err == nil {
		for _, ip := range ips {
			if ip.To4() != nil {
				info.IPv4 = append(info.IPv4, ip.String())
			} else {
				info.IPv6 = append(info.IPv6, ip.String())
			}
		}
	}

	// CNAME records
	cnames, err := net.LookupCNAME(cleanDomain)
	if err == nil && cnames != cleanDomain+"." {
		info.CNAME = append(info.CNAME, cnames)
	}

	// MX records
	mxRecords, err := net.LookupMX(cleanDomain)
	if err == nil {
		for _, mx := range mxRecords {
			info.MX = append(info.MX, fmt.Sprintf("%s (priority: %d)", mx.Host, mx.Pref))
		}
	}

	// TXT records
	txtRecords, err := net.LookupTXT(cleanDomain)
	if err == nil {
		info.TXT = txtRecords
	}

	// NS records
	nsRecords, err := net.LookupNS(cleanDomain)
	if err == nil {
		for _, ns := range nsRecords {
			info.NS = append(info.NS, ns.Host)
		}
	}

	return info
}

// isIPAddress checks if a string is a valid IP address
func isIPAddress(input string) bool {
	return net.ParseIP(input) != nil
}

// parseHSTSHeader parses the Strict-Transport-Security header
func parseHSTSHeader(header string) HSTSInfo {
	info := HSTSInfo{}

	if header == "" {
		info.Details = "HSTS header not present"
		return info
	}

	info.Enabled = true
	info.Directive = header

	// Parse max-age
	maxAgeRegex := regexp.MustCompile(`max-age=(\d+)`)
	maxAgeMatch := maxAgeRegex.FindStringSubmatch(header)
	if len(maxAgeMatch) > 1 {
		if maxAge, err := strconv.Atoi(maxAgeMatch[1]); err == nil {
			info.MaxAge = maxAge
		}
	}

	// Check for includeSubDomains
	if strings.Contains(strings.ToLower(header), "includesubdomains") {
		info.IncludeSubDomains = true
	}

	// Check for preload
	if strings.Contains(strings.ToLower(header), "preload") {
		info.Preload = true
	}

	// Build details
	var details []string
	details = append(details, fmt.Sprintf("Max-Age: %d seconds", info.MaxAge))

	if info.IncludeSubDomains {
		details = append(details, "Includes SubDomains: Yes")
	} else {
		details = append(details, "Includes SubDomains: No")
	}

	if info.Preload {
		details = append(details, "Preload: Yes")
	}

	info.Details = strings.Join(details, ", ")

	return info
}

// CheckIP checks IP address information (accepts both IP addresses and domain names)
func (nc *NetChecker) CheckIP(input string) IPInfo {
	info := IPInfo{Input: input}

	// Clean input
	cleanInput := strings.TrimPrefix(input, "https://")
	cleanInput = strings.TrimPrefix(cleanInput, "http://")
	cleanInput = strings.Split(cleanInput, "/")[0]

	// Check if input is an IP address or a domain
	if isIPAddress(cleanInput) {
		// It's an IP address
		info.IsDomain = false
		info.IP = cleanInput

		// Try to establish connection to check if IP is reachable
		conn, err := net.Dial("tcp", cleanInput+":80")
		if err == nil {
			conn.Close()
		}

		// For demonstration, we'll return basic info
		info.Country = "Unknown"
		info.Region = "Unknown"
		info.City = "Unknown"
		info.ISP = "Unknown"
		info.Organization = "Unknown"
		info.Timezone = "Unknown"

		return info
	} else {
		// It's a domain name - resolve to IPs
		info.IsDomain = true

		ips, err := net.LookupIP(cleanInput)
		if err != nil {
			info.Error = fmt.Sprintf("Failed to resolve domain: %v", err)
			return info
		}

		// Store resolved IPs
		for _, ip := range ips {
			if ip.To4() != nil {
				info.ResolvedIPs = append(info.ResolvedIPs, ip.String())
			}
		}

		if len(info.ResolvedIPs) > 0 {
			// Use the first resolved IP
			info.IP = info.ResolvedIPs[0]

			// Try to establish connection
			conn, err := net.Dial("tcp", info.IP+":80")
			if err == nil {
				conn.Close()
			}

			// For demonstration, we'll return basic info
			info.Country = "Unknown"
			info.Region = "Unknown"
			info.City = "Unknown"
			info.ISP = "Unknown"
			info.Organization = "Unknown"
			info.Timezone = "Unknown"
		} else {
			info.Error = "No IP addresses found for domain"
		}

		return info
	}
}

// CheckWebSettings checks web server settings and headers
func (nc *NetChecker) CheckWebSettings(domain string) WebSettingsInfo {
	info := WebSettingsInfo{Domain: domain}

	// Ensure domain has protocol
	if !strings.HasPrefix(domain, "https://") && !strings.HasPrefix(domain, "http://") {
		domain = "https://" + domain
	}

	start := time.Now()
	resp, err := nc.httpClient.Get(domain)
	info.ResponseTime = time.Since(start).Milliseconds()

	if err != nil {
		info.Error = fmt.Sprintf("Failed to connect: %v", err)
		return info
	}
	defer resp.Body.Close()

	info.StatusCode = resp.StatusCode
	info.Server = resp.Header.Get("Server")
	info.ContentType = resp.Header.Get("Content-Type")
	info.LastModified = resp.Header.Get("Last-Modified")
	info.ETag = resp.Header.Get("ETag")
	info.ContentLength = resp.ContentLength

	// Check for redirects
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		info.RedirectURL = resp.Header.Get("Location")
	}

	// Parse HSTS header
	hstsHeader := resp.Header.Get("Strict-Transport-Security")
	info.HSTS = parseHSTSHeader(hstsHeader)

	// Collect all headers
	info.Headers = make(map[string]string)
	for key, values := range resp.Header {
		info.Headers[key] = strings.Join(values, ", ")
	}

	return info
}

// CheckEmailConfig checks email authentication configuration (SPF, DKIM, DMARC, BIMI)
func (nc *NetChecker) CheckEmailConfig(domain string) EmailConfigInfo {
	info := EmailConfigInfo{Domain: domain}

	// Clean domain
	cleanDomain := strings.TrimPrefix(domain, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.Split(cleanDomain, "/")[0]

	// Check SPF
	info.SPF = nc.CheckSPF(cleanDomain)

	// Check DKIM
	info.DKIM = nc.CheckDKIM(cleanDomain)

	// Check DMARC
	info.DMARC = nc.CheckDMARC(cleanDomain)

	// Check BIMI
	info.BIMI = nc.CheckBIMI(cleanDomain)

	return info
}

// CheckSPF checks SPF (Sender Policy Framework) record
func (nc *NetChecker) CheckSPF(domain string) SPFInfo {
	info := SPFInfo{}

	// SPF records are TXT records for the domain
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to lookup TXT records: %v", err)
		return info
	}

	// Look for SPF record
	for _, txt := range txtRecords {
		if strings.HasPrefix(strings.ToLower(txt), "v=spf1") {
			info.Configured = true
			info.Record = txt
			info.Valid = true
			info.Details = "SPF record found and configured"
			break
		}
	}

	if !info.Configured {
		info.Details = "No SPF record found"
	}

	return info
}

// CheckDKIM checks DKIM (DomainKeys Identified Mail) records
func (nc *NetChecker) CheckDKIM(domain string) DKIMInfo {
	info := DKIMInfo{}

	// Common DKIM selectors to check
	commonSelectors := []string{"default", "dkim", "key1", "selector1", "s1", "s2"}

	// Check for DKIM records in common locations
	found := false
	for _, selector := range commonSelectors {
		dkimDomain := selector + "._domainkey." + domain
		txtRecords, err := net.LookupTXT(dkimDomain)
		if err == nil {
			for _, txt := range txtRecords {
				if strings.HasPrefix(strings.ToLower(txt), "v=dkim1") {
					info.Selectors = append(info.Selectors, selector)
					found = true
				}
			}
		}
	}

	if found {
		info.Configured = true
		info.Valid = true
		info.Details = fmt.Sprintf("DKIM records found for selectors: %v", info.Selectors)
	} else {
		info.Details = "No DKIM records found for common selectors"
	}

	return info
}

// CheckDMARC checks DMARC (Domain-based Message Authentication) record
func (nc *NetChecker) CheckDMARC(domain string) DMARCInfo {
	info := DMARCInfo{}

	// DMARC record is in _dmarc subdomain
	dmarcDomain := "_dmarc." + domain
	txtRecords, err := net.LookupTXT(dmarcDomain)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to lookup DMARC record: %v", err)
		info.Details = "No DMARC record found"
		return info
	}

	// Look for DMARC record
	for _, txt := range txtRecords {
		if strings.HasPrefix(strings.ToLower(txt), "v=dmarc1") {
			info.Configured = true
			info.Record = txt
			info.Valid = true

			// Parse policy
			if strings.Contains(txt, "p=none") {
				info.Policy = "none"
			} else if strings.Contains(txt, "p=quarantine") {
				info.Policy = "quarantine"
			} else if strings.Contains(txt, "p=reject") {
				info.Policy = "reject"
			}

			info.Details = fmt.Sprintf("DMARC record found with policy: %s", info.Policy)
			break
		}
	}

	if !info.Configured {
		info.Details = "No DMARC record found"
	}

	return info
}

// CheckBIMI checks BIMI (Brand Indicators for Message Identification) record
func (nc *NetChecker) CheckBIMI(domain string) BIMIInfo {
	info := BIMIInfo{}

	// BIMI record is in default._bimi subdomain
	bimiDomain := "default._bimi." + domain
	txtRecords, err := net.LookupTXT(bimiDomain)
	if err != nil {
		info.Details = "No BIMI record found"
		return info
	}

	// Look for BIMI record
	for _, txt := range txtRecords {
		if strings.HasPrefix(strings.ToLower(txt), "v=bimi1") {
			info.Configured = true
			info.Record = txt
			info.Valid = true

			// Parse logo URL
			if strings.Contains(txt, "l=") {
				parts := strings.Split(txt, "l=")
				if len(parts) > 1 {
					info.LogoURL = strings.TrimSpace(strings.Split(parts[1], ";")[0])
				}
			}

			info.Details = "BIMI record found"
			if info.LogoURL != "" {
				info.Details = fmt.Sprintf("BIMI record found with logo: %s", info.LogoURL)
			}
			break
		}
	}

	if !info.Configured {
		info.Details = "No BIMI record found"
	}

	return info
}

// CombinedResult holds both SSL and WebSettings info from a single connection
type CombinedResult struct {
	SSL         SSLInfo
	WebSettings WebSettingsInfo
}

// CheckSSLAndWebSettings performs both SSL and WebSettings checks using a single HTTP request
func (nc *NetChecker) CheckSSLAndWebSettings(domain string) CombinedResult {
	var result CombinedResult

	// Clean domain for HTTP request
	fullURL := domain
	if !strings.HasPrefix(domain, "https://") && !strings.HasPrefix(domain, "http://") {
		fullURL = "https://" + domain
	}

	cleanDomain := strings.TrimPrefix(fullURL, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.Split(cleanDomain, "/")[0]

	// Make HTTP request and get certificate from TLS connection
	start := time.Now()
	req, err := http.NewRequest("HEAD", fullURL, nil)
	if err != nil {
		result.WebSettings.Error = fmt.Sprintf("Failed to create request: %v", err)
		result.SSL.Error = result.WebSettings.Error
		return result
	}

	resp, err := nc.httpClient.Do(req)
	result.WebSettings.ResponseTime = time.Since(start).Milliseconds()
	result.WebSettings.Domain = cleanDomain

	if err != nil {
		result.WebSettings.Error = fmt.Sprintf("Failed to connect: %v", err)
		result.SSL.Domain = cleanDomain
		result.SSL.Error = result.WebSettings.Error
		return result
	}
	defer resp.Body.Close()

	// Extract WebSettings info
	result.WebSettings.StatusCode = resp.StatusCode
	result.WebSettings.Server = resp.Header.Get("Server")
	result.WebSettings.ContentType = resp.Header.Get("Content-Type")
	result.WebSettings.LastModified = resp.Header.Get("Last-Modified")
	result.WebSettings.ETag = resp.Header.Get("ETag")
	result.WebSettings.ContentLength = resp.ContentLength

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		result.WebSettings.RedirectURL = resp.Header.Get("Location")
	}

	// Parse HSTS header
	hstsHeader := resp.Header.Get("Strict-Transport-Security")
	result.WebSettings.HSTS = parseHSTSHeader(hstsHeader)

	result.WebSettings.Headers = make(map[string]string)
	for key, values := range resp.Header {
		result.WebSettings.Headers[key] = strings.Join(values, ", ")
	}

	// Extract SSL certificate from TLS connection state
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		result.SSL.Domain = cleanDomain
		result.SSL.Valid = true
		result.SSL.Issuer = cert.Issuer.String()
		result.SSL.Subject = cert.Subject.String()
		result.SSL.NotBefore = cert.NotBefore
		result.SSL.NotAfter = cert.NotAfter
		result.SSL.DaysUntilExpiry = int(time.Until(cert.NotAfter).Hours() / 24)
		result.SSL.SerialNumber = cert.SerialNumber.String()
		result.SSL.SignatureAlg = cert.SignatureAlgorithm.String()
		result.SSL.PublicKeyAlg = cert.PublicKeyAlgorithm.String()

		switch pubKey := cert.PublicKey.(type) {
		case *rsa.PublicKey:
			result.SSL.KeySize = pubKey.N.BitLen()
		default:
			result.SSL.KeySize = 0
		}
	} else {
		result.SSL.Domain = cleanDomain
		result.SSL.Error = "No TLS certificate found"
	}

	return result
}

// API Handlers

func handleSSL(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	// cache
	key := cacheKey("/api/v1/ssl", map[string]string{"domain": domain})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}

	checker := NewNetChecker()
	sslInfo := checker.CheckSSL(domain)
	if ttl, ok := routeTTL["/api/v1/ssl"]; ok {
		apiCache.Set(key, sslInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    sslInfo,
	})
}

func handleHTTP3(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	key := cacheKey("/api/v1/http3", map[string]string{"domain": domain})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}
	checker := NewNetChecker()
	http3Info := checker.CheckHTTP3(domain)
	if ttl, ok := routeTTL["/api/v1/http3"]; ok {
		apiCache.Set(key, http3Info, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    http3Info,
	})
}

func handleDNS(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}
	key := cacheKey("/api/v1/dns", map[string]string{"domain": domain})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}
	checker := NewNetChecker()
	dnsInfo := checker.CheckDNS(domain)
	if ttl, ok := routeTTL["/api/v1/dns"]; ok {
		apiCache.Set(key, dnsInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    dnsInfo,
	})
}

func handleIP(c *gin.Context) {
	// Accept both 'ip' and 'domain' parameters for flexibility
	input := c.Query("ip")
	if input == "" {
		input = c.Query("domain")
	}

	if input == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "IP or domain parameter is required",
		})
		return
	}

	key := cacheKey("/api/v1/ip", map[string]string{"input": input})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}
	checker := NewNetChecker()
	ipInfo := checker.CheckIP(input)
	if ttl, ok := routeTTL["/api/v1/ip"]; ok {
		apiCache.Set(key, ipInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    ipInfo,
	})
}

func handleMyIP(c *gin.Context) {
	// Get the real client IP from the request
	clientIP := getRealClientIP(c)

	if clientIP == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Could not determine client IP address",
		})
		return
	}

	// Use client IP as cache key (per user)
	key := cacheKey("/api/v1/my-ip", map[string]string{"ip": clientIP})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}

	// Get detailed IP information
	checker := NewNetChecker()
	ipInfo := checker.CheckIP(clientIP)

	// Override input field to indicate this is the user's own IP
	ipInfo.Input = "Your IP: " + clientIP

	if ttl, ok := routeTTL["/api/v1/my-ip"]; ok {
		apiCache.Set(key, ipInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    ipInfo,
	})
}

func handleWebSettings(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	key := cacheKey("/api/v1/web-settings", map[string]string{"domain": domain})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}
	checker := NewNetChecker()
	webInfo := checker.CheckWebSettings(domain)
	if ttl, ok := routeTTL["/api/v1/web-settings"]; ok {
		apiCache.Set(key, webInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    webInfo,
	})
}

func handleEmailConfig(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	key := cacheKey("/api/v1/email-config", map[string]string{"domain": domain})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}
	checker := NewNetChecker()
	emailInfo := checker.CheckEmailConfig(domain)
	if ttl, ok := routeTTL["/api/v1/email-config"]; ok {
		apiCache.Set(key, emailInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    emailInfo,
	})
}

// DNS servers with blocklists
var blocklistDNSServers = []DNSServer{
	{Name: "AdGuard", IP: "176.103.130.130"},
	{Name: "AdGuard Family", IP: "176.103.130.132"},
	{Name: "CleanBrowsing Adult", IP: "185.228.168.10"},
	{Name: "CleanBrowsing Family", IP: "185.228.168.168"},
	{Name: "CleanBrowsing Security", IP: "185.228.168.9"},
	{Name: "CloudFlare", IP: "1.1.1.1"},
	{Name: "CloudFlare Family", IP: "1.1.1.3"},
	{Name: "Comodo Secure", IP: "8.26.56.26"},
	{Name: "Google DNS", IP: "8.8.8.8"},
	{Name: "Neustar Family", IP: "156.154.70.3"},
	{Name: "Neustar Protection", IP: "156.154.70.2"},
	{Name: "Norton Family", IP: "199.85.126.20"},
	{Name: "OpenDNS", IP: "208.67.222.222"},
	{Name: "OpenDNS Family", IP: "208.67.222.123"},
	{Name: "Quad9", IP: "9.9.9.9"},
	{Name: "Yandex Family", IP: "77.88.8.7"},
	{Name: "Yandex Safe", IP: "77.88.8.88"},
}

// Known block IPs that DNS servers return when blocking domains
var knownBlockIPs = []string{
	"146.112.61.106", // OpenDNS
	"185.228.168.10", // CleanBrowsing
	"8.26.56.26",     // Comodo
	"208.69.38.170",  // OpenDNS
	"208.69.39.170",  // OpenDNS
	"208.67.222.222", // OpenDNS
	"208.67.222.123", // OpenDNS FamilyShield
	"199.85.126.10",  // Norton
	"199.85.126.20",  // Norton Family
	"156.154.70.22",  // Neustar
	"77.88.8.7",      // Yandex
}

// contains checks if a string slice contains a value
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isDomainBlockedByServer checks if a domain is blocked by a specific DNS server
func (nc *NetChecker) isDomainBlockedByServer(domain, serverIP string) bool {
	// Create a custom DNS resolver for the specific server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to resolve the domain using the server
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", serverIP+":53")
		},
	}

	// Resolve A records
	addrs, err := resolver.LookupIPAddr(ctx, domain)
	if err != nil {
		// If DNS resolution fails, it might indicate the domain is blocked
		return true
	}

	// Check if any returned IP is a known block IP
	for _, addr := range addrs {
		if contains(knownBlockIPs, addr.IP.String()) {
			return true
		}
	}

	return false
}

// CheckBlocklist checks if a domain is blocked by various DNS servers
func (nc *NetChecker) CheckBlocklist(domain string) BlocklistInfo {
	info := BlocklistInfo{Domain: domain}

	// Clean domain
	cleanDomain := strings.TrimPrefix(domain, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.Split(cleanDomain, "/")[0]

	var results []BlocklistResult

	// Check against each DNS server
	for _, server := range blocklistDNSServers {
		isBlocked := nc.isDomainBlockedByServer(cleanDomain, server.IP)
		results = append(results, BlocklistResult{
			Server:    server.Name,
			ServerIP:  server.IP,
			IsBlocked: isBlocked,
		})
	}

	info.Results = results
	return info
}

func handleBlocklist(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	key := cacheKey("/api/v1/blocklist", map[string]string{"domain": domain})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}
	checker := NewNetChecker()
	blocklistInfo := checker.CheckBlocklist(domain)
	if ttl, ok := routeTTL["/api/v1/blocklist"]; ok {
		apiCache.Set(key, blocklistInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    blocklistInfo,
	})
}

func handleHSTS(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	checker := NewNetChecker()

	// Get web settings which includes HSTS
	webInfo := checker.CheckWebSettings(domain)

	// Return just the HSTS info
	hstsInfo := HSTSInfo{
		Enabled:           webInfo.HSTS.Enabled,
		MaxAge:            webInfo.HSTS.MaxAge,
		IncludeSubDomains: webInfo.HSTS.IncludeSubDomains,
		Preload:           webInfo.HSTS.Preload,
		Directive:         webInfo.HSTS.Directive,
		Details:           webInfo.HSTS.Details,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"domain": domain,
			"hsts":   hstsInfo,
		},
	})
}

// CheckRobotsTxt checks robots.txt file for a domain
func (nc *NetChecker) CheckRobotsTxt(domain string) RobotsTxtInfo {
	info := RobotsTxtInfo{Domain: domain}

	// Clean domain
	cleanDomain := strings.TrimPrefix(domain, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.Split(cleanDomain, "/")[0]

	// Build robots.txt URL
	robotsURL := fmt.Sprintf("https://%s/robots.txt", cleanDomain)

	// Try HTTP if HTTPS fails
	resp, err := nc.httpClient.Get(robotsURL)
	if err != nil {
		robotsURL = fmt.Sprintf("http://%s/robots.txt", cleanDomain)
		resp, err = nc.httpClient.Get(robotsURL)
		if err != nil {
			info.Exists = false
			info.Status = "Not Found"
			info.Error = fmt.Sprintf("Failed to fetch robots.txt: %v", err)
			return info
		}
	}
	defer resp.Body.Close()

	info.Exists = true
	info.Status = fmt.Sprintf("HTTP %d", resp.StatusCode)

	// Read content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to read response: %v", err)
		return info
	}

	info.Content = string(body)
	info.Lines = strings.Split(info.Content, "\n")

	// Parse robots.txt content
	currentUserAgent := "*"
	for _, line := range info.Lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		field := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch field {
		case "user-agent":
			currentUserAgent = value
			info.UserAgents = append(info.UserAgents, value)
		case "disallow":
			if value != "" {
				info.Disallowed = append(info.Disallowed, fmt.Sprintf("%s: %s", currentUserAgent, value))
			}
		case "allow":
			if value != "" {
				info.Allowed = append(info.Allowed, fmt.Sprintf("%s: %s", currentUserAgent, value))
			}
		case "sitemap":
			info.Sitemaps = append(info.Sitemaps, value)
		case "crawl-delay":
			if value != "" {
				info.CrawlDelay = value
			}
		}
	}

	return info
}

func handleRobotsTxt(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	key := cacheKey("/api/v1/robots-txt", map[string]string{"domain": domain})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}
	checker := NewNetChecker()
	robotsInfo := checker.CheckRobotsTxt(domain)
	if ttl, ok := routeTTL["/api/v1/robots-txt"]; ok {
		apiCache.Set(key, robotsInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    robotsInfo,
	})
}

// CheckSitemap checks sitemap.xml file for a domain
func (nc *NetChecker) CheckSitemap(domain string) SitemapInfo {
	info := SitemapInfo{Domain: domain}

	// Clean domain
	cleanDomain := strings.TrimPrefix(domain, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.Split(cleanDomain, "/")[0]

	// Common sitemap locations to try
	sitemapPaths := []string{
		"/sitemap.xml",
		"/sitemap_index.xml",
		"/sitemap-index.xml",
		"/sitemaps.xml",
	}

	var sitemapURL string
	var resp *http.Response
	var err error

	// Try each sitemap path
	for _, path := range sitemapPaths {
		sitemapURL = fmt.Sprintf("https://%s%s", cleanDomain, path)
		resp, err = nc.httpClient.Get(sitemapURL)
		if err == nil && resp.StatusCode == 200 {
			info.SitemapURL = sitemapURL
			break
		}
		if resp != nil {
			resp.Body.Close()
		}

		// Try HTTP if HTTPS fails
		sitemapURL = fmt.Sprintf("http://%s%s", cleanDomain, path)
		resp, err = nc.httpClient.Get(sitemapURL)
		if err == nil && resp.StatusCode == 200 {
			info.SitemapURL = sitemapURL
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	if resp == nil || err != nil {
		info.Exists = false
		info.Status = "Not Found"
		info.Error = fmt.Sprintf("Failed to find sitemap: %v", err)
		return info
	}
	defer resp.Body.Close()

	info.Exists = true
	info.Status = fmt.Sprintf("HTTP %d", resp.StatusCode)

	// Read content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to read response: %v", err)
		return info
	}

	// Try to parse as sitemap index first
	var sitemapIndex SitemapIndex
	err = xml.Unmarshal(body, &sitemapIndex)
	if err == nil && len(sitemapIndex.Sitemaps) > 0 {
		info.IsSitemapIndex = true
		info.URLCount = len(sitemapIndex.Sitemaps)
		for _, s := range sitemapIndex.Sitemaps {
			info.SubSitemaps = append(info.SubSitemaps, s.Loc)
		}
		return info
	}

	// Try to parse as regular sitemap
	var sitemap Sitemap
	err = xml.Unmarshal(body, &sitemap)
	if err == nil && len(sitemap.URLs) > 0 {
		info.IsSitemapIndex = false
		info.URLCount = len(sitemap.URLs)

		// Get sample URLs (first 10)
		maxSamples := 10
		if len(sitemap.URLs) < maxSamples {
			maxSamples = len(sitemap.URLs)
		}
		for i := 0; i < maxSamples; i++ {
			info.SampleURLs = append(info.SampleURLs, sitemap.URLs[i].Loc)
			if sitemap.URLs[i].LastMod != "" {
				info.LastModified = append(info.LastModified, sitemap.URLs[i].LastMod)
			}
		}
		return info
	}

	// If parsing failed, still return that we found something
	info.Error = "Sitemap found but could not parse XML structure"
	return info
}

func handleSitemap(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	key := cacheKey("/api/v1/sitemap", map[string]string{"domain": domain})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}
	checker := NewNetChecker()
	sitemapInfo := checker.CheckSitemap(domain)
	if ttl, ok := routeTTL["/api/v1/sitemap"]; ok {
		apiCache.Set(key, sitemapInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    sitemapInfo,
	})
}

// HTMLProxyResponse represents HTML content fetched from a URL
type HTMLProxyResponse struct {
	URL    string `json:"url"`
	HTML   string `json:"html"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

// OGImageInfo represents Open Graph image and metadata information
type OGImageInfo struct {
	URL    string `json:"url"`
	Domain string `json:"domain"`
	Found  bool   `json:"found"`

	// OG Image specific
	ImageURL    string `json:"image_url,omitempty"`
	ImageURLAlt string `json:"image_url_alt,omitempty"`
	ImageSecure string `json:"image_secure,omitempty"`
	ImageWidth  string `json:"image_width,omitempty"`
	ImageHeight string `json:"image_height,omitempty"`
	ImageType   string `json:"image_type,omitempty"`
	Accessible  bool   `json:"accessible"`
	Status      int    `json:"status,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int64  `json:"size,omitempty"`

	// OG Metadata
	OGTitle       string `json:"og_title,omitempty"`
	OGDescription string `json:"og_description,omitempty"`
	OGType        string `json:"og_type,omitempty"`
	OGURL         string `json:"og_url,omitempty"`
	OGSiteName    string `json:"og_site_name,omitempty"`
	OGLocale      string `json:"og_locale,omitempty"`

	// Twitter Card metadata
	TwitterCard        string `json:"twitter_card,omitempty"`
	TwitterSite        string `json:"twitter_site,omitempty"`
	TwitterCreator     string `json:"twitter_creator,omitempty"`
	TwitterTitle       string `json:"twitter_title,omitempty"`
	TwitterDescription string `json:"twitter_description,omitempty"`
	TwitterImage       string `json:"twitter_image,omitempty"`
	TwitterImageAlt    string `json:"twitter_image_alt,omitempty"`

	// Standard meta tags
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`

	// All detected tags
	AllMetaTags    map[string]string `json:"all_meta_tags,omitempty"`
	AllTwitterTags map[string]string `json:"all_twitter_tags,omitempty"`

	Error string `json:"error,omitempty"`
}

func handleHTMLProxy(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "URL parameter is required",
		})
		return
	}

	// Clean and normalize URL
	targetURL := url
	if !strings.HasPrefix(targetURL, "https://") && !strings.HasPrefix(targetURL, "http://") {
		targetURL = "https://" + targetURL
	}

	checker := NewNetChecker()
	resp, err := checker.httpClient.Get(targetURL)
	if err != nil {
		// Try HTTP if HTTPS fails
		if strings.HasPrefix(targetURL, "https://") {
			targetURL = strings.Replace(targetURL, "https://", "http://", 1)
			resp, err = checker.httpClient.Get(targetURL)
		}
		if err != nil {
			c.JSON(http.StatusOK, APIResponse{
				Success: true,
				Data: HTMLProxyResponse{
					URL:   targetURL,
					Error: fmt.Sprintf("Failed to fetch: %v", err),
				},
			})
			return
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Data: HTMLProxyResponse{
				URL:   targetURL,
				Error: fmt.Sprintf("Failed to read response: %v", err),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: HTMLProxyResponse{
			URL:    targetURL,
			HTML:   string(body),
			Status: resp.StatusCode,
		},
	})
}

// CheckOGImage checks Open Graph image tags for a URL
func (nc *NetChecker) CheckOGImage(url string) OGImageInfo {
	info := OGImageInfo{URL: url}

	// Clean and normalize URL
	targetURL := url
	if !strings.HasPrefix(targetURL, "https://") && !strings.HasPrefix(targetURL, "http://") {
		targetURL = "https://" + targetURL
	}

	// Extract domain
	cleanDomain := strings.TrimPrefix(targetURL, "https://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "http://")
	cleanDomain = strings.Split(cleanDomain, "/")[0]
	info.Domain = cleanDomain

	// Fetch HTML content
	resp, err := nc.httpClient.Get(targetURL)
	if err != nil {
		// Try HTTP if HTTPS fails
		if strings.HasPrefix(targetURL, "https://") {
			targetURL = strings.Replace(targetURL, "https://", "http://", 1)
			resp, err = nc.httpClient.Get(targetURL)
		}
		if err != nil {
			info.Error = fmt.Sprintf("Failed to fetch URL: %v", err)
			return info
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		info.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return info
	}

	// Read HTML content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		info.Error = fmt.Sprintf("Failed to read response: %v", err)
		return info
	}

	htmlContent := string(body)
	info.AllMetaTags = make(map[string]string)
	info.AllTwitterTags = make(map[string]string)

	// Parse all OG tags
	allOGTagsRegex := regexp.MustCompile(`(?i)<meta\s+(?:property|name)=["'](og:[^"']+)["']\s+content=["']([^"']+)["']`)
	matches := allOGTagsRegex.FindAllStringSubmatch(htmlContent, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			tagName := strings.ToLower(match[1])
			tagValue := match[2]
			info.AllMetaTags[tagName] = tagValue

			// Extract specific OG fields
			switch tagName {
			case "og:title":
				info.OGTitle = tagValue
			case "og:description":
				info.OGDescription = tagValue
			case "og:type":
				info.OGType = tagValue
			case "og:url":
				info.OGURL = tagValue
			case "og:site_name":
				info.OGSiteName = tagValue
			case "og:locale":
				info.OGLocale = tagValue
			case "og:image":
				info.ImageURL = tagValue
				info.Found = true
			case "og:image:url":
				info.ImageURLAlt = tagValue
			case "og:image:secure_url":
				info.ImageSecure = tagValue
			case "og:image:width":
				info.ImageWidth = tagValue
			case "og:image:height":
				info.ImageHeight = tagValue
			case "og:image:type":
				info.ImageType = tagValue
			}
		}
	}

	// Parse all Twitter Card tags
	allTwitterTagsRegex := regexp.MustCompile(`(?i)<meta\s+(?:property|name)=["'](twitter:[^"']+)["']\s+content=["']([^"']+)["']`)
	twitterMatches := allTwitterTagsRegex.FindAllStringSubmatch(htmlContent, -1)
	for _, match := range twitterMatches {
		if len(match) >= 3 {
			tagName := strings.ToLower(match[1])
			tagValue := match[2]
			info.AllTwitterTags[tagName] = tagValue

			// Extract specific Twitter fields
			switch tagName {
			case "twitter:card":
				info.TwitterCard = tagValue
			case "twitter:site":
				info.TwitterSite = tagValue
			case "twitter:creator":
				info.TwitterCreator = tagValue
			case "twitter:title":
				info.TwitterTitle = tagValue
			case "twitter:description":
				info.TwitterDescription = tagValue
			case "twitter:image":
				if !info.Found {
					info.ImageURL = tagValue
					info.Found = true
				}
				info.TwitterImage = tagValue
			case "twitter:image:alt":
				info.TwitterImageAlt = tagValue
			}
		}
	}

	// Extract standard meta tags (title and description)
	metaTitleRegex := regexp.MustCompile(`(?i)<title>([^<]+)</title>`)
	if match := metaTitleRegex.FindStringSubmatch(htmlContent); len(match) > 1 {
		info.MetaTitle = strings.TrimSpace(match[1])
	}

	metaDescRegex := regexp.MustCompile(`(?i)<meta\s+name=["']description["']\s+content=["']([^"']+)["']`)
	if match := metaDescRegex.FindStringSubmatch(htmlContent); len(match) > 1 {
		info.MetaDescription = match[1]
	}

	// If still no image found, try standard meta tags
	if !info.Found {
		standardImageRegex := regexp.MustCompile(`(?i)<link\s+rel=["']image_src["']\s+href=["']([^"']+)["']`)
		if match := standardImageRegex.FindStringSubmatch(htmlContent); len(match) > 1 {
			info.ImageURL = match[1]
			info.Found = true
		}
	}

	// If we found an image URL, check if it's accessible
	if info.Found && info.ImageURL != "" {
		// Resolve relative URLs
		imageURL := info.ImageURL
		if strings.HasPrefix(imageURL, "//") {
			imageURL = "https:" + imageURL
		} else if strings.HasPrefix(imageURL, "/") {
			imageURL = fmt.Sprintf("https://%s%s", cleanDomain, imageURL)
		} else if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
			// Relative URL without leading slash
			if strings.HasSuffix(targetURL, "/") {
				imageURL = targetURL + imageURL
			} else {
				imageURL = targetURL + "/" + imageURL
			}
		}

		// Try to fetch the image
		imgResp, imgErr := nc.httpClient.Get(imageURL)
		if imgErr != nil {
			// Try HTTP if HTTPS fails
			if strings.HasPrefix(imageURL, "https://") {
				imageURLHTTP := strings.Replace(imageURL, "https://", "http://", 1)
				var imgRespHTTP *http.Response
				imgRespHTTP, imgErr = nc.httpClient.Get(imageURLHTTP)
				if imgErr == nil {
					imgResp = imgRespHTTP
				}
			}
		}

		if imgErr == nil && imgResp != nil {
			defer imgResp.Body.Close()
			info.Accessible = true
			info.Status = imgResp.StatusCode
			info.ContentType = imgResp.Header.Get("Content-Type")
			info.Size = imgResp.ContentLength
		} else {
			info.Accessible = false
			if imgErr != nil {
				info.Error = fmt.Sprintf("Image not accessible: %v", imgErr)
			} else {
				info.Error = "Image not accessible"
			}
		}
	}

	return info
}

func handleOGImage(c *gin.Context) {
	url := c.Query("url")
	domain := c.Query("domain")

	// Support both url and domain parameters
	targetURL := url
	if targetURL == "" && domain != "" {
		targetURL = domain
	}

	if targetURL == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "URL or domain parameter is required",
		})
		return
	}

	key := cacheKey("/api/v1/og-image", map[string]string{"url": targetURL})
	if v, ok := apiCache.Get(key); ok {
		c.JSON(http.StatusOK, APIResponse{Success: true, Data: v})
		return
	}

	checker := NewNetChecker()
	ogInfo := checker.CheckOGImage(targetURL)
	if ttl, ok := routeTTL["/api/v1/og-image"]; ok {
		apiCache.Set(key, ogInfo, ttl)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    ogInfo,
	})
}

func handleComprehensive(c *gin.Context) {
	startTime := time.Now()
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Domain parameter is required",
		})
		return
	}

	checker := NewNetChecker()

	// Run all checks in parallel using goroutines (optimized to reduce network requests)
	type result struct {
		name       string
		data       interface{}
		durationMs int64
	}

	resultChan := make(chan result, 7) // Buffer for 7 results (ssl, web_settings, http3, dns, email_config, robots_txt, sitemap)

	// Run all checks concurrently with timing
	// Combined SSL + WebSettings in a single network request
	go func() {
		checkStart := time.Now()
		combined := checker.CheckSSLAndWebSettings(domain)
		duration := time.Since(checkStart).Milliseconds()

		// Send both results separately
		resultChan <- result{"ssl", combined.SSL, duration}
		resultChan <- result{"web_settings", combined.WebSettings, duration}
	}()

	go func() {
		checkStart := time.Now()
		data := checker.CheckHTTP3(domain)
		duration := time.Since(checkStart).Milliseconds()
		resultChan <- result{"http3", data, duration}
	}()

	go func() {
		checkStart := time.Now()
		data := checker.CheckDNS(domain)
		duration := time.Since(checkStart).Milliseconds()
		resultChan <- result{"dns", data, duration}
	}()

	go func() {
		checkStart := time.Now()
		data := checker.CheckEmailConfig(domain)
		duration := time.Since(checkStart).Milliseconds()
		resultChan <- result{"email_config", data, duration}
	}()

	go func() {
		checkStart := time.Now()
		data := checker.CheckRobotsTxt(domain)
		duration := time.Since(checkStart).Milliseconds()
		resultChan <- result{"robots_txt", data, duration}
	}()

	go func() {
		checkStart := time.Now()
		data := checker.CheckSitemap(domain)
		duration := time.Since(checkStart).Milliseconds()
		resultChan <- result{"sitemap", data, duration}
	}()

	// Collect results (7 items total, but only 6 network operations)
	comprehensiveData := make(map[string]interface{})
	timings := make(map[string]int64)

	for i := 0; i < 7; i++ {
		res := <-resultChan
		comprehensiveData[res.name] = res.data
		timings[res.name] = res.durationMs
	}

	// Calculate total processing time
	totalTime := time.Since(startTime).Milliseconds()
	timings["total"] = totalTime

	// Add timings to response
	comprehensiveData["_meta"] = map[string]interface{}{
		"timings":  timings,
		"domain":   domain,
		"total_ms": totalTime,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    comprehensiveData,
	})
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "NetCheck API is running",
		Data: map[string]interface{}{
			"version": "1.0.0",
			"status":  "healthy",
		},
	})
}

func main() {
	// Set Gin mode based on environment variable (debug for dev, release for production)
	ginMode := getenvDefault("GIN_MODE", "release")
	gin.SetMode(ginMode)

	// Create Gin router
	r := gin.Default()

	// Trust proxy addresses so c.ClientIP() uses X-Forwarded-For from Caddy
	// Includes loopback and common private CIDRs (adjust to your network as needed)
	if err := r.SetTrustedProxies([]string{
		"127.0.0.1",
		"::1",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}); err != nil {
		// If setting proxies fails, continue with default (no trust)
	}

	p := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem: "gin",
	})

	// API secret authentication middleware
	r.Use(func(c *gin.Context) {
		// Allow direct access to /metrics endpoint for monitoring systems
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// In debug mode, skip secret validation entirely
		isDev := getenvDefault("GIN_MODE", "release") == "debug"
		if isDev {
			// Set CORS headers for all requests in dev mode
			requestOrigin := c.Request.Header.Get("Origin")
			if requestOrigin != "" {
				c.Header("Access-Control-Allow-Origin", requestOrigin)
				c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Content-Type, X-Internal-Proxy, X-API-Secret, Authorization")
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			c.Next()
			return
		}

		// In production mode, require secret authentication AND internal proxy header
		// This prevents direct API access - all requests must go through the SvelteKit proxy
		apiSecret := strings.TrimSpace(getenvDefault("API_SECRET_KEY", ""))
		if apiSecret == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "API secret key not configured",
			})
			return
		}

		// Require X-Internal-Proxy header to prevent direct API access
		// Only the SvelteKit server-side proxy sets this header
		internalProxy := c.Request.Header.Get("X-Internal-Proxy")
		if internalProxy != "true" {
			c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
				Success: false,
				Error:   "Direct API access not allowed. Requests must go through the application proxy.",
			})
			return
		}

		// Get secret from header (query parameter removed for security)
		providedSecret := c.Request.Header.Get("X-API-Secret")
		if providedSecret == "" {
			// Try Authorization header with Bearer token as fallback
			authHeader := c.Request.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				providedSecret = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Require secret to match
		if providedSecret != apiSecret {
			c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
				Success: false,
				Error:   "Invalid or missing API secret key",
			})
			return
		}

		// Set CORS headers if Origin is present (for browser requests from frontend)
		requestOrigin := c.Request.Header.Get("Origin")
		if requestOrigin != "" {
			c.Header("Access-Control-Allow-Origin", requestOrigin)
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, X-Internal-Proxy, X-API-Secret, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	p.Use(r)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// API routes
	api := r.Group("/api/v1")
	// Rate limiter: default 60 rpm; heavy routes stricter
	rl := NewRateLimiter(60, map[string]int{
		"/api/v1/comprehensive": 6,
		"/api/v1/robots-txt":    10,
		"/api/v1/sitemap":       10,
		"/api/v1/blocklist":     10,
		"/api/v1/web-settings":  20,
		"/api/v1/og-image":      15,
	})
	api.Use(rl.Middleware())
	{
		api.GET("/health", handleHealth)
		api.GET("/ssl", handleSSL)
		api.GET("/http3", handleHTTP3)
		api.GET("/dns", handleDNS)
		api.GET("/ip", handleIP)
		api.GET("/my-ip", handleMyIP)
		api.GET("/web-settings", handleWebSettings)
		api.GET("/email-config", handleEmailConfig)
		api.GET("/blocklist", handleBlocklist)
		api.GET("/hsts", handleHSTS)
		api.GET("/robots-txt", handleRobotsTxt)
		api.GET("/sitemap", handleSitemap)
		api.GET("/og-image", handleOGImage)
		api.GET("/html-proxy", handleHTMLProxy)
		api.GET("/comprehensive", handleComprehensive)
	}

	// Root endpoint with API documentation
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "NetCheck API",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"health":        "GET /api/v1/health",
				"ssl":           "GET /api/v1/ssl?domain=example.com",
				"http3":         "GET /api/v1/http3?domain=example.com",
				"dns":           "GET /api/v1/dns?domain=example.com",
				"ip":            "GET /api/v1/ip?ip=8.8.8.8 or ?domain=example.com",
				"my-ip":         "GET /api/v1/my-ip (returns your IP address)",
				"web-settings":  "GET /api/v1/web-settings?domain=example.com",
				"email-config":  "GET /api/v1/email-config?domain=example.com",
				"blocklist":     "GET /api/v1/blocklist?domain=example.com",
				"hsts":          "GET /api/v1/hsts?domain=example.com",
				"robots-txt":    "GET /api/v1/robots-txt?domain=example.com",
				"sitemap":       "GET /api/v1/sitemap?domain=example.com",
				"og-image":      "GET /api/v1/og-image?url=https://example.com or ?domain=example.com",
				"comprehensive": "GET /api/v1/comprehensive?domain=example.com",
			},
		})
	})

	// Start server
	mode := getenvDefault("GIN_MODE", "release")
	fmt.Printf("NetCheck API Server starting on :8080 (mode: %s)\n", mode)
	if mode == "debug" {
		fmt.Println("⚠️  DEVELOPMENT MODE: API access allowed without secret authentication")
	}
	fmt.Println("API Documentation available at http://localhost:8080")
	r.Run(":8080")
}
