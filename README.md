# dev-proxy
penjelasan singkat, ini adalah tool sederhana untuk meredirect local port ke url yang sudah di tentukan dalam file ports.txt dan akan mengirimkan response yang di terima ke client yang merequest.

contoh

8080;https://github.com/              
8081;https://google.com/

contoh konfigurasi di atas, segala request ke url http://localhost:8080/ akan di alihkan ke https://github.com/.

jadi jika http://localhost:8080/books/ maka akan di alihkan ke https://github.com/books/.

Versi 1.1.0
- Fix karakter kosong di file ports.txt
- support profile
contoh 
  > dev-proxy dev   => akan me refer ke file ports-dev.txt
