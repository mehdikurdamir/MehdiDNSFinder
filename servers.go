package main

type DNSServer struct {
	Name  string `json:"name"`
	IP    string `json:"ip"`
	Group string `json:"group"` // iran | world
}

// فهرست DNSهای عمومی. هر کدام در زمان اجرا تست می‌شوند؛ پس اگر سرویسی از کار افتاده باشد
// فقط به‌صورت «بدون پاسخ» نمایش داده می‌شود.
var defaultServers = []DNSServer{
	// ---- ایران ----
	{"شکن (Shecan)", "178.22.122.100", "iran"},
	{"شکن (Shecan) ۲", "185.51.200.2", "iran"},
	{"الکترو (Electro)", "78.157.42.100", "iran"},
	{"الکترو (Electro) ۲", "78.157.42.101", "iran"},
	{"رادار گیم (Radar)", "10.202.10.10", "iran"},
	{"رادار گیم (Radar) ۲", "10.202.10.11", "iran"},
	{"403.online", "10.202.10.202", "iran"},
	{"403.online ۲", "10.202.10.102", "iran"},
	{"بگذر (Begzar)", "185.55.226.26", "iran"},
	{"بگذر (Begzar) ۲", "185.55.225.25", "iran"},
	{"پیشگامان", "5.202.100.100", "iran"},
	{"پیشگامان ۲", "5.202.100.101", "iran"},
	{"شاتل (Shatel)", "85.15.1.14", "iran"},
	{"شاتل (Shatel) ۲", "85.15.1.15", "iran"},
	{"مخابرات (TCI)", "217.218.127.127", "iran"},
	{"مخابرات (TCI) ۲", "217.218.155.155", "iran"},
	{"هاست ایران", "172.29.0.100", "iran"},
	{"هاست ایران ۲", "172.29.2.100", "iran"},
	{"آسیاتک", "194.36.174.161", "iran"},
	{"ایرانسل", "109.69.8.51", "iran"},

	// ---- جهان ----
	{"Google", "8.8.8.8", "world"},
	{"Google ۲", "8.8.4.4", "world"},
	{"Cloudflare", "1.1.1.1", "world"},
	{"Cloudflare ۲", "1.0.0.1", "world"},
	{"Cloudflare (Malware)", "1.1.1.2", "world"},
	{"Quad9", "9.9.9.9", "world"},
	{"Quad9 ۲", "149.112.112.112", "world"},
	{"OpenDNS", "208.67.222.222", "world"},
	{"OpenDNS ۲", "208.67.220.220", "world"},
	{"AdGuard", "94.140.14.14", "world"},
	{"AdGuard ۲", "94.140.15.15", "world"},
	{"CleanBrowsing", "185.228.168.9", "world"},
	{"CleanBrowsing ۲", "185.228.169.9", "world"},
	{"Comodo Secure", "8.26.56.26", "world"},
	{"Comodo Secure ۲", "8.20.247.20", "world"},
	{"Level3", "4.2.2.1", "world"},
	{"Level3 ۲", "4.2.2.2", "world"},
	{"Verisign", "64.6.64.6", "world"},
	{"Verisign ۲", "64.6.65.6", "world"},
	{"DNS.WATCH", "84.200.69.80", "world"},
	{"DNS.WATCH ۲", "84.200.70.40", "world"},
	{"Yandex", "77.88.8.8", "world"},
	{"Yandex ۲", "77.88.8.1", "world"},
	{"Control D", "76.76.2.0", "world"},
	{"Control D ۲", "76.76.10.0", "world"},
	{"Mullvad", "194.242.2.2", "world"},
	{"DNS0.eu", "193.110.81.0", "world"},
	{"CZ.NIC ODVR", "193.17.47.1", "world"},
	{"Neustar UltraDNS", "156.154.70.1", "world"},
	{"SafeDNS", "195.46.39.39", "world"},
	{"Alternate DNS", "76.76.19.19", "world"},
	{"Freenom World", "80.80.80.80", "world"},
	{"Hurricane Electric", "74.82.42.42", "world"},
	{"LibreDNS", "116.202.176.26", "world"},
	{"Gcore", "95.85.95.85", "world"},
}
