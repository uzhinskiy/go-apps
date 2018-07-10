package main

import (
        "flag"
        "fmt"
        "html"
        "log"
        "mime"
        "net/http"
        "net/url"
        "os"
        "path"
        "sort"
        "strconv"
        //"time"

        "github.com/uzhinskiy/lib.go/pconf"
)

type dirListFilesType struct {
        Name     string
        Url      string
        Icon     string
        Size     string
        SizeReal int64
        ModTime  string
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

const (
        iconDir  = " data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABYAAAAWCAYAAADEtGw7AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAAXdEVYdEF1dGhvcgBMYXBvIENhbGFtYW5kcmVp35EaKgAAACl0RVh0RGVzY3JpcHRpb24AQmFzZWQgb2YgSmFrdWIgU3RlaW5lciBkZXNpZ26ghAVzAAACFUlEQVQ4jbXUPUwUQRTA8f/s7R1eQxBjjBTWakJMNIgmklBo54kNDZW9JhYglq4lfsTSgtJCEywUC1qxINHEnIkfhQnyYUQSIh9HgLvdfe9Z3B1wwl0gi5NMdic789u3b98M/KfmqjcPBnIfzOioMy9yypXBJ2/G9w0PDeSi6303fefcjknLSwurb8deFoGLgw9fTe4bPn/h0q5wc2sbhUKh8P7dWHMjzHPelzuPXrcD+NsfiAi7wUsLP2k5eqK5p+9WXdRUGH3x9HR1XANrHRhgcX6Kxfmp+uGa1QwbRhyFJaanJykVw/pgbfOG+q+N3H082vsPrFRdEWX6xyRnOi9ztqsHzDC0fDXFVJCwiEQbxKV14nCNuLTByPPhXBD0ZnyAof7cMOB//Zyvef3J9k46uq6yOPsJqP3UHaECmUwTgHd87rD5QdDts8qN2/ef4XkpTAWcwwEqMX9m8pjpXlMB4H63fTcfuhU+er++jZNOZ/YDNAi+W70gCBRwKoqZJe6AuxcEVv15Kqops8Z53GNzDjZhUYkPCga26ljFhKRwZbluh0XlAGBVqNTlVsRiiWEtrzfA+YBvoCJisdruB8VeYdmEfR84ZKqsFRaiFGGiQhYVq2zQrA94G6V4ZmJi/JTDxUlgzFkYR7OA54CmbDZ77EhL07koltYkbtr3llfXS/mVleKco3x+pCvdSxRxOb8hEP0FgzwzsGRM+0UAAAAASUVORK5CYII="
        iconFile = " data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABYAAAAWCAYAAADEtGw7AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAASdEVYdFRpdGxlAFBhcGVyIFNoZWV0c7mvkfkAAAAXdEVYdEF1dGhvcgBMYXBvIENhbGFtYW5kcmVp35EaKgAAACd0RVh0RGVzY3JpcHRpb24Ad2l0aCBhIEhVR0UgaGVscCBmcm9tIEpha3VihlQHswAAA1NJREFUOI2VlN1rHFUYh593NpvMbjbVtpuoiVEbi/UDrF8gfkRzLUqhxQoF/QsEQS/8G7wpit4JBhXBC01Tq9Qb8UJ7I+JFGwu2aQJNstkk2zaSmuzsnPO+XszMdjabSD3syzkzs/vsw+89c8TMAPjhx9PDQRB8oKrHVa0XjOyZmSEiGyJy2dR+9+p/NePro0dej9llSPbj789OTz5y6NG3DjxwMAiCoOuL3jti52i1ImrLS9GfFy+sq9f3Dfvy+LETuit4+sw3/xx59Vi5GTVR9ahqUqaoejAIgoC+3pBisRczY+bi+ejy7F+L3uvzJ954czUPbqs553tFBOdivPd49agppoqZ4dUTtSLWN24wO3+JjZsbPHn46b5nnnp2DOznL76aLO8I9t4DhqpSX66zeHWJpYUataU69doqq/UGjdXrXF9bp7XlmJmZQUR4cOygPPzQY4eA7yY//7SQ8XryGVraqMGhKmaGpQ3MSlN7M+Xvm+tZUzn8+BOFG+vXXpibn5sAfuoEuyRHNWVtpUEc79pwAKLNFqemT7Fv315eGn+Z0ZH7wytzV17pAjvvk4UZ+wf3tgGJIRgGOftyJWSlvsL8/DxTU9/y4vg4wGvAe9uiSMCG0Vi9hnP+P40B7thzJxMTE2w1NykEBVwcj330ycniO2+/G3eBAaqD+0E6IZZokyaf3kymSmWAarUKBHjnC0AO7BwAgrC21ugyLvQUGBqqdv5hm28g0IqjtuAOUcDgXVUEAZGEk83puhOc+AuCesV51wnOmldIX2eRDCwknxSYXGbMDnXv/S7GAhIUWFleodVqcTujWOxh9L7Rtpx328HOIwiBCMMj9yAi3ZXpZvaWttEyY7d7xiLCcq1OFEW3aVxkbOxAW85tB7vcdhu5dzhnGhDk1tn9RFRRu3Vu75ixpocQwMLVxa5Xum2W2yHZ2UIe7LbtClWNVa0H2DFjELx6ujay3bqKnTOfKgdprsVms3n2wsx56y9XGKjsodI/QH+5QrnUTyksUwpLhH0lwr4wVyXCMKnZ2UvErfi3D09+HAD0SKIzeHr6zNTm5tZQWAqfM9WCZZs/1/n23H6WBYK2YvfHuV/OfQbcLSILkp6pvcAIUNyp8/9jOKBmZs1/AQG5DNFh7ozgAAAAAElFTkSuQmCC"
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
        base := Config["document_root"]
        code := http.StatusOK
        var contentType string

        file := checkFileExist(base, r)

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
                        q := r.URL.Query()
                        order := "N"
                        if val, ok := q["C"]; ok {
                                order = val[0]
                        }
                        fmt.Println(order)
                        DL, _ := getDirList(base+file, file, order)
                        fmt.Fprintf(w, "<html><body><h3>%s</h3>", DL.Title)
                        fmt.Fprintf(w, "<h5><a href='../'>Parent dir</a></h5>")
                        fmt.Fprintf(w, "<table><tr><th>&nbsp;</th><th><a href='?C=N'>Name</a></th><th align='right'><a href='?C=M'>Last modified</a></th><th align='right'><a href='?C=S'>Size</a></th></tr><tr><th colspan='3'><hr></th></tr>")
                        for _, f := range DL.Files {
                                fmt.Fprintf(w, "<tr><td valign='top'><img src='%s' width='22' height='22'/></td><td><a href='%s'>%s</a></td><td align='right'>%s</td><td align='right'>%s</td></tr>", f.Icon, f.Url, f.Name, f.ModTime, f.Size)
                        }
                        fmt.Fprintf(w, "<tr><td colspan='3'><hr></td></tr><tr><td colspan='3'>"+Config["version"]+"</td></tr></table></body></html>")
                case mode&os.ModeSymlink != 0:
                        fmt.Fprintf(w, "symbolic link")
                case mode&os.ModeNamedPipe != 0:
                        fmt.Fprintf(w, "named pipe")
                }
        }

}

func checkFileExist(base string, r *http.Request) string {
        t := 0
        file := r.URL.Path
        if _, err := os.Stat(base + file); os.IsNotExist(err) {
                t = 1
        }
        if t == 0 {
                return file
        }

        file, _ = url.QueryUnescape(r.RequestURI)
        if _, err := os.Stat(base + file); os.IsNotExist(err) {
                t = 2
        }
        if t == 0 {
                return file
        }

        return file
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

func getDirList(dirname, request string, order string) (dirListType, error) {
        f, err := os.Open(dirname)
        defer f.Close()
        var (
                dirContent dirListType
                files      []dirListFilesType
                size       string
                rsize      int64
                icon       string
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
                        rsize = -1
                        icon = iconDir
                } else {
                        if i.Size()/1024 > 1 {
                                size = strconv.FormatInt(i.Size()/1024, 10) + "K"
                        } else {
                                size = strconv.FormatInt(i.Size(), 10)
                        }
                        rsize = i.Size()
                        icon = iconFile
                }
                files = append(files, dirListFilesType{Name: name, Url: name, ModTime: i.ModTime().Format("2006-01-02 15:04"), Size: size, SizeReal: rsize, Icon: icon})
        }

        switch order {
        case "N":
                sort.Slice(files[:], func(i, j int) bool {
                        return files[i].Name < files[j].Name
                })
        case "M":
                sort.Slice(files[:], func(i, j int) bool {
                        return files[i].ModTime < files[j].ModTime
                })
        case "S":
                sort.Slice(files[:], func(i, j int) bool {
                        return files[i].SizeReal < files[j].SizeReal
                })
        default:
                sort.Slice(files[:], func(i, j int) bool {
                        return files[i].Name < files[j].Name
                })
        }

        dirContent.Files = files

        return dirContent, nil
}
