package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"gitlab.com/slirx/newproj/pkg/logger"
)

func FileServerWithCustom404(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			http.ServeFile(w, r, "web/generated/index.html")
			return
		}
		fsh.ServeHTTP(w, r)
	})
}

func main() {
	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("can not initialize logger: %v", err)
	}

	conf, err := NewConfig("FRONTEND_HTTPD_")
	if err != nil {
		zapLogger.Fatal(err)
	}

	//http.Handle("/", http.FileServer(http.Dir("template/frontend/httpd")))
	//http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web")))) // todo move this to media microservice
	http.Handle("/", FileServerWithCustom404(http.Dir("web/generated")))
	err = http.ListenAndServe(conf.Server.Addr, nil)
	if err != nil {
		zapLogger.Fatal(err)
	}

	//log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/usr/share/doc"))))
}
