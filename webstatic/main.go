package main

import (
        "log"
        "mime"
        "net/http"
        "os"
        "path"
)

/* функция для обработки подключившихся клиентов */
func requestHandler(w http.ResponseWriter, r *http.Request) {
        file := r.URL.Path
        base := "/srv/es-head"
        /* если отсутствует запрос к конкретному файлу – показать индексный файл */
        if file == "/" {
                file = "/index.html"
        }

        /* если не удалось загрузить нужный файл – показать сообщение о 404-ой ошибке */
        respFile, err := os.OpenFile(base+file, os.O_RDONLY, 0)
        if err != nil {
                log.Println(err)
                file = "/404.html"
                respFile, err = os.OpenFile(base+file, os.O_RDONLY, 0)
        }
        /* считать содержимое файла */
        fi, err := respFile.Stat()
        contentType := mime.TypeByExtension(path.Ext(file))
        var bytes = make([]byte, fi.Size())
        respFile.Read(bytes)
        /* отправить его клиенту */
        w.Header().Set("Content-Type", contentType)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Server", "gohttp/0.1")
        w.Write(bytes)
}

func main() {
        http.HandleFunc("/", requestHandler)
        http.ListenAndServe(":9400", nil)
}
