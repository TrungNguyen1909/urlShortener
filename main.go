package main
import (
	"fmt"
	"net/http"
	"log"
	"crypto/rand"
	"io/ioutil"
	"strings"
	"os"
	"net/url"
)
var resolver map[string]string = make(map[string]string)
const alphabet string = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789"
func generateKey() (key string) {
	b := make([]byte, 8)
	rand.Read(b)
	for _,v := range b{
		key += string(alphabet[int(v) % len(alphabet)]);
	}
	return
}
func apiHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	if r.Method == "POST"{
		body, err := ioutil.ReadAll(r.Body)
		if err != nil || body == nil{
			http.Error(w, "Bad Request",http.StatusBadRequest)
			return
		}
		url,err := url.Parse(strings.TrimSpace(string(body)))
		if err != nil || url == nil{
			http.Error(w, "WTF did you just send me?",http.StatusBadRequest)
			return
		}
		if len(url.Scheme) == 0{
			url.Scheme = "http"
		}
		key := generateKey()
		for _, ok := resolver[key]; ok == true;_, ok = resolver[key] {
			key = generateKey()
		}
		log.Println(url)
		resolver[key] = url.String()
		fmt.Fprintf(w,key);
		return;
	}
	http.Error(w, "Bad Request",http.StatusBadRequest)
	return
}
func handler(w http.ResponseWriter, r *http.Request){
	path := r.URL.Path[1:];
	url, ok := resolver[path];
	if ok == false {
		data, err := ioutil.ReadFile(string(path))
		if err != nil{
			http.ServeFile(w,r,"./index.html")
			return;
		}
		w.Write(data);
		return;
	}
	http.Redirect(w, r, url, 302);
	return;
}
func main() {
	os.Chdir("./static")
	http.HandleFunc("/",handler)
	http.HandleFunc("/*",handler)
	http.HandleFunc("/api", apiHandler)
	port, ok := os.LookupEnv("PORT")
	if !ok{
	port = "8890"
	}
	port = ":"+port
	log.Printf("Listening on %s\n",port);
	log.Fatal(http.ListenAndServe(port, nil))
}
