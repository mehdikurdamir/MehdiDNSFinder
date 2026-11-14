<div align="center">
  <img src="assets/icon.png" width="120" alt="MehdiDNSFinder">
  <h1>MehdiDNSFinder</h1>
  <p>ابزار ویندوز برای پیدا کردن DNS مناسب یک سایت یا سرویس — با تست واقعی روی DNSهای ایران و جهان</p>
</div>

---

## فارسی

آدرس یک سرویس را وارد می‌کنید (مثلاً `https://www.youtube.com`). برنامه حدود ۵۵ سرور DNS عمومی (ایرانی و جهانی) را همزمان تست می‌کند و نشان می‌دهد کدام‌یک دسترسی واقعی به آن سرویس می‌دهند.

### روش تست
۱. از هر DNS، آی‌پی دامنه پرسیده می‌شود.
۲. اگر جواب بی‌معنی بود (۰.۰.۰.۰، آی‌پی خصوصی یا صفحه فیلترینگ `10.10.34.x`) به‌عنوان «مسدود» علامت می‌خورد.
۳. در غیر این صورت اتصال واقعی TCP/TLS به آی‌پی دریافتی برقرار و گواهی بررسی می‌شود.

وضعیت‌ها: دسترسی دارد · هشدار گواهی · اتصال ناموفق · مسدود · یافت نشد · DNS بی‌پاسخ

### امکانات
- رابط فارسی با فونت وزیرمتن (داخل برنامه، بدون نیاز به اینترنت)
- نمایش زنده نتایج، مرتب‌سازی بر اساس وضعیت و سرعت
- فیلتر (همه / فقط موفق / ایران / جهان)
- کپی آدرس DNS، افزودن DNS دلخواه، خروجی CSV

### نصب
از بخش [Releases](../../releases) فایل `MehdiDNSFinder-Setup.exe` را دانلود و اجرا کنید.
(چون فایل امضای دیجیتال ندارد، ویندوز ممکن است SmartScreen نشان دهد: More info ← Run anyway)

### نکته‌ها
- بعضی DNSهای ایرانی (رادار، ۴۰۳، هاست ایران) فقط داخل ایران پاسخ می‌دهند.
- نتیجه به اپراتور و اینترنت شما بستگی دارد.
- برنامه DNS ویندوز را تغییر نمی‌دهد؛ آدرس را کپی و در تنظیمات شبکه وارد کنید.

### ساخت از سورس
```bash
sudo apt-get install -y golang-go nsis binutils-mingw-w64-x86-64
./build.sh        # خروجی: dist/MehdiDNSFinder-Setup.exe
```
با هر تگ `v*` (مثلاً `v1.0.0`)، GitHub Actions خودش فایل نصبی را می‌سازد و به Release اضافه می‌کند.

لیست DNSها در `servers.go` است و به‌راحتی قابل ویرایش است.

---

## English

Windows tool that takes a URL and tests ~55 free public DNS servers (Iranian and worldwide) in parallel. For each one it resolves the host, filters out sinkhole answers (0.0.0.0, private IPs, Iran's `10.10.34.x` block page), then performs a real TCP/TLS connection with certificate verification to confirm the service is actually reachable.

Built with Go (standard library only) and a local web UI shown in an Edge/Chrome app window. Cross-compiled to Windows and packaged with NSIS.

```bash
./build.sh   # → dist/MehdiDNSFinder-Setup.exe
```

## License
MIT — see [LICENSE](LICENSE). Bundled font: [Vazirmatn](https://github.com/rastikerdar/vazirmatn), SIL OFL 1.1 ([third_party/Vazirmatn-OFL.txt](third_party/Vazirmatn-OFL.txt)).
