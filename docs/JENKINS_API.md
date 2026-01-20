# Jenkins Trigger API

Endpoint untuk memicu Jenkins job build tanpa autentikasi.

## Endpoint

```
POST /api/v1/jenkins/trigger
```

## Headers

| Header | Required | Description |
|--------|----------|-------------|
| X-Jenkins-Url | Ya | URL dasar Jenkins (contoh: https://jenkins.example.com) |
| X-Jenkins-Token | Ya | Token autentikasi Jenkins |
| X-Jenkins-Job | Ya | Nama job Jenkins yang akan dipicu |
| X-Jenkins-User | Tidak | Username untuk autentikasi Jenkins (default: sonar) |

## Contoh Request

```bash
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.example.com" \
  -H "X-Jenkins-Token: your-jenkins-token" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

## Contoh Response (Sukses)

```json
{
  "success": true,
  "message": "Jenkins job triggered successfully",
  "job": "plb-sonarqube"
}
```

## Contoh Response (Error)

```json
{
  "error": "Missing or invalid headers"
}
```

## Proses Internal

1. Validasi header yang diterima
2. Mendapatkan crumb dari Jenkins untuk CSRF protection
3. Memanggil endpoint build job Jenkins dengan:
   - Basic authentication menggunakan username (default: "sonar") dan token yang diberikan
   - Header CSRF protection menggunakan crumb yang didapatkan
4. Mengembalikan status sukses atau error

## Catatan

- Endpoint ini tidak memerlukan autentikasi aplikasi
- Pastikan Jenkins URL dapat diakses dari server ini
- Token Jenkins harus memiliki izin yang cukup untuk memicu build