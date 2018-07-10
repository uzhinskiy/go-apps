package main

import (
        "flag"
        "fmt"
        "html"
        "log"
        "mime"
        "net/http"
        "os"
        "path"
        "sort"
        "strconv"
        "time"

        "github.com/uzhinskiy/lib.go/pconf"
)

type dirListFilesType struct {
        Name    string                                                                                                                                                      
        Url     string                                                                                                                                                      
        Size    string                                                                                                                                                      
        ModTime time.Time                                                                                                                                                   
}                                                                                                                                                                           
                                                                                                                                                                            
type dirListType struct {                                                                                                                                                   
        Title string
        Files []dirListFilesType
}

var (
        configfile string
        HTTPAddr   string
        err        error
        Config     pconf.ConfigType
)

func init() {
        var addr, port string
        flag.StringVar(&addr, "bind", "", "Address to listen for HTTP requests on")
        flag.StringVar(&port, "port", "8080", "Port to listen for HTTP requests on")
        flag.StringVar(&configfile, "config", "main.cfg", "Read configuration from this file")
        flag.Parse()

        Config = make(pconf.ConfigType)
        err := Config.Parse(configfile)
        if err != nil {
                log.Fatal("Fatal error: ", err)
        }

        log.Println("Read from config ", len(Config), " items:", Config)
        if Config["bind"] != "" {
                addr = Config["bind"]
        }
        if Config["port"] != "" {
                port = Config["port"]
        }
        HTTPAddr = addr + ":" + port
        fmt.Println(HTTPAddr)
}

func main() {
        logTo := os.Stderr
        if logTo, err = os.OpenFile(Config["log_file"], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
                log.Fatal(err)
        }
        defer logTo.Close()
        log.SetOutput(logTo)

        if Config["extended_view"] == "1" {
                http.HandleFunc("/", Index)
        } else {
                http.Handle("/", http.FileServer(http.Dir(Config["document_root"])))
        }
        log.Println("HTTP server listening on: ", HTTPAddr)
        err := http.ListenAndServe(HTTPAddr, nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }

}

/* функция для обработки подключившихся клиентов */
func Index(w http.ResponseWriter, r *http.Request) {
        file := r.URL.Path
        base := Config["document_root"]
        code := http.StatusOK
        var contentType string

        fi, err := os.Lstat(base + file)
        if err != nil {
                w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
                w.Header().Set("Server", Config["version"])
                bytes := []byte(err.Error())
                code = 404
                w.WriteHeader(code)
                w.Write(bytes)

                log.Println(r.RemoteAddr, "\t", r.Method, "\t", r.URL.Path, "\t", code, "\t", r.UserAgent())
        } else {
                switch mode := fi.Mode(); {
                case mode.IsRegular():
                        bytes, err := getFile(base+file, fi.Size())

                        if err != nil {
                                bytes = []byte(err.Error()) //fmt.Fprintf(w, "%d\t%s", 400, err.Error())
                                code = 403
                                contentType = "text/plain; charset=UTF-8"
                        } else {
                                code = http.StatusOK
                                contentType = mime.TypeByExtension(path.Ext(file))
                        }

                        w.Header().Set("Content-Type", contentType)
                        w.Header().Set("Server", Config["version"])

                        log.Println(r.RemoteAddr, "\t", r.Method, "\t", r.URL.Path, "\t", code, "\t", r.UserAgent())
                        w.WriteHeader(code)
                        w.Write(bytes)

                case mode.IsDir():
                        w.Header().Set("Content-Type", "text/html; charset=UTF-8")
                        w.Header().Set("Server", Config["version"])

                        DL, _ := getDirList(base+file, file)
                        fmt.Fprintf(w, "<html><body><h3>%s</h3>", DL.Title)
                        fmt.Fprintf(w, "<h5><a href='../'>Parent dir</a></h5>")
                        fmt.Fprintf(w, "<table><tr><th>Name</th><th>Last modified</th><th>Size</th></tr><tr><th colspan='3'><hr></th></tr>")
                        for _, f := range DL.Files {
                                fmt.Fprintf(w, "<tr><td><a href='%s'>%s</a></td><td>%s</td><td>%s</td></tr>", f.Url, f.Name, f.ModTime.Format("2006-01-02 15:04"), f.Size)
                        }
                        fmt.Fprintf(w, "<tr><td colspan='3'><hr></td></tr><tr><td colspan='3'>"+Config["version"]+"</td></tr></table></body></html>")
                case mode&os.ModeSymlink != 0:
                        fmt.Fprintf(w, "symbolic link")
                case mode&os.ModeNamedPipe != 0:
                        fmt.Fprintf(w, "named pipe")
                }
        }

}

func getFile(fname string, size int64) ([]byte, error) {
        respFile, err := os.OpenFile(fname, os.O_RDONLY, 0)
        defer respFile.Close()
        if err != nil {
                return nil, err
        }
        bytes := make([]byte, size)
        respFile.Read(bytes)
        return bytes, nil
}

func getDirList(dirname, request string) (dirListType, error) {
        f, err := os.Open(dirname)
        defer f.Close()
        var (
                dirContent dirListType
                files      []dirListFilesType
                size       string
        )
        dirContent.Title = "Index of " + dirname

        if err != nil {
                return dirContent, err
        }

        info, err := f.Readdir(0)
        if err != nil {
                return dirContent, err
        }
        for _, i := range info {
                name := html.EscapeString(i.Name())
                if i.IsDir() == true {
                        name = name + "/"
                        size = "-"
                } else {
                        if i.Size()/1024 > 1 {
                                size = strconv.FormatInt(i.Size()/1024, 10) + "K"
                        } else {
                                size = strconv.FormatInt(i.Size(), 10)
                        }
                }
                files = append(files, dirListFilesType{Name: name, Url: name, ModTime: i.ModTime(), Size: size})
        }

        sort.Slice(files[:], func(i, j int) bool {
                return files[i].Name < files[j].Name
        })

        dirContent.Files = files

        return dirContent, nil
}
