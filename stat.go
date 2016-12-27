package main

import (
	"flag"

	"github.com/gorilla/securecookie"
	_ "github.com/lib/pq"

	"database/sql"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/urfave/negroni"
)

func main() {
	dbConn := flag.String("dbconnect", "user:pass@localhost/paddle_stat?sslmode=disable", "the database connect string of Paddle Stat.")
	ver := flag.String("version", "v0.9.0a", "the highest Paddle version number")
	addr := flag.String("addr", ":3000", "http address")
	key := flag.String("key", "", "hash key for cookie")
	flag.Parse()
	fmt.Println(*dbConn)
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s", *dbConn))
	e(err)
	defer db.Close()
	hashKey := []byte(*key)
	secureCookie := securecookie.New(hashKey, nil)
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	mux := http.NewServeMux()
	mux.HandleFunc("/version", func(res http.ResponseWriter, req *http.Request) {
		var uid int64
		const cookieName = "paddle_stat"
		c, err := req.Cookie(cookieName)
		if err != nil { // First time to use Paddle.
			err = db.QueryRow("SELECT NEXTVAL('UID')").Scan(&uid)
			e(err)
			encoded, err := secureCookie.Encode(cookieName, uid)
			e(err)
			http.SetCookie(res, &http.Cookie{
				Name:  cookieName,
				Value: encoded,
				Path:  "/version",
			})
		} else { // Parse cookie to get uid
			e(secureCookie.Decode(cookieName, c.Value, &uid))
		}

		req.ParseForm()
		content := req.Form.Get("content")
		_, err = db.Exec("INSERT INTO Usage(uid, content, type) VALUES($1, $2, 0)", uid, content)
		e(err)
		res.Write([]byte(*ver))
	})
	mux.HandleFunc("/feedback", func(res http.ResponseWriter, req *http.Request) {
		var uid int64
		const cookieName = "paddle_stat"
		c, err := req.Cookie(cookieName)
		if err != nil { // First time to use Paddle.
			err = db.QueryRow("SELECT NEXTVAL('UID')").Scan(&uid)
			e(err)
			encoded, err := secureCookie.Encode(cookieName, uid)
			e(err)
			http.SetCookie(res, &http.Cookie{
				Name:  cookieName,
				Value: encoded,
				Path:  "/feedback",
			})
		} else { // Parse cookie to get uid
			e(secureCookie.Decode(cookieName, c.Value, &uid))
		}
		req.ParseForm()
		keys := []string {
			"system_time", "paddle_version", "github_user", "job_name",
			"duration", "exit_code",
		}
		content_obj := map[string]string {}
		for i := 0; i < len(keys); i++ {
			key := keys[i]
			val := req.Form.Get(key)
			if key != "" {
				content_obj[key] = val
			}
		}

		content, err := json.Marshal(content_obj)
		//content, err := json.Marshal(req.Form)
		e(err)
		fmt.Println("The uid:", uid)
		fmt.Println("content:", string(content))

		_, err = db.Exec("INSERT INTO Usage(uid, content, type) VALUES($1, $2, 1)", uid, content)
		e(err)
		res.Write(content)
	})
	n.UseHandler(mux)
	fmt.Println("listen http on ", *addr)
	http.ListenAndServe(*addr, n)
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
