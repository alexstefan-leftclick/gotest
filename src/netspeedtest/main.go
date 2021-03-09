package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/johnsto/speedtest"
	cors "github.com/rs/cors/wrapper/gin"
	"gopkg.in/ddo/go-fast.v0"
	"net/http"
	"time"
)

type SpeedTestResults struct {
	DownloadSpeed string `json:downspeed`
	UploadSpeed   string `json:upspeed`
}

type APIResponse struct {
	Result SpeedTestResults
	Error  error
}

func fastTest() {
	fastCom := fast.New()

	// init
	err := fastCom.Init()

	if err != nil {
		panic(err)
	}

	// get urls
	urls, err := fastCom.GetUrls()

	if err != nil {
		panic(err)
	}

	// measure
	KbpsChan := make(chan float64)
	//c := make(chan []float64)
	var results []float64
	go func() {
		fmt.Println(22222)

		for Kbps := range KbpsChan {
			fmt.Printf("%.2f Kbps %.2f Mbps\n", Kbps, Kbps/1000)
			results = append(results, Kbps)
		}
		fmt.Println(33333)

		//c <- results
	}()

	//ceva := <-c
	fmt.Println(111111)
	//fmt.Println(results)
	err = fastCom.Measure(urls, KbpsChan)

	if err != nil {
		panic(err)
	}

}

func ooaklaTest(client http.Client) (SpeedTestResults, error) {
	// Fetch server list
	settings, err := FetchSettings()
	if err != nil {
		return SpeedTestResults{}, err
	}

	//Get closest server to user
	settings.Servers.SortByDistance()

	//Initiate down and up benchmark
	benchmarkDownload := NewDownloadBenchmark(client, settings.Servers[0])
	benchmarkUpload := NewUploadBenchmark(client, settings.Servers[0])

	//Run benchmark
	rateDownload := RunBenchmark(benchmarkDownload, 4, 16, time.Second*10)
	rateUpload := RunBenchmark(benchmarkUpload, 4, 16, time.Second*10)

	return SpeedTestResults{NiceRate(rateDownload), NiceRate(rateUpload)}, nil
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	v1 := r.Group("/api/v1")
	{
		v1.GET("/speedtest", testPicker)
	}
	r.Run()
}

func testPicker(c *gin.Context) {
	//resp, err := ooaklaTest(http.Client{})
	//c.JSON(http.StatusOK, APIResponse{resp, err})
	fastTest()
}
