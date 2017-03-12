package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

// var (
// 	redisAddress   = flag.String("redis-address", "127.0.0.1:6379", "Address to the Redis server")
// 	maxConnections = flag.Int("max-connections", 10, "Max connections to Redis")
// )

func BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		accessToken := r.Header.Get("access-token")
		if len(accessToken) > 0 {
			fmt.Println("BasicAuth - Middleware", accessToken)
			h(w, r, ps) // next process
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func GetAccountProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Get Account Profile")
}

func CreateRedis(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := redis.Dial("tcp", os.Getenv("REDIS_DIAL"))
	defer c.Close()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	const bestCarEver = "Lancer Ex"

	c.Do("SET", "best_car_ever", bestCarEver)
	c.Do("SET", "worst_car_ever", "MG3")

	c.Do("EXPIRE", "best_car_ever", 5)
	c.Do("EXPIRE", "worst_car_ever", 10)

	worstCarEver, err := redis.String(c.Do("GET", "worst_car_ever"))
	if err != nil {
		fmt.Println("worstCarEver not found")
	}

	fmt.Fprintln(w, "Best Car Ever is "+bestCarEver)
	fmt.Fprintln(w, "Worst Car Ever is "+worstCarEver)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort := os.Getenv("APP_PORT")

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/account", BasicAuth(GetAccountProfile))
	router.POST("/test/create-redis", CreateRedis)

	fmt.Println("Application is running on " + appPort)
	http.ListenAndServe(":"+appPort, router)
}
