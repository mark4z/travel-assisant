package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	__no     string
	__date   string
	__fromZh string
	__toZh   string
)

var rootCmd = &cobra.Command{
	Use: "travel",
}

//go:embed travel/dist/**
var f embed.FS
var templ *template.Template

func init() {
	templ = template.Must(template.New("").ParseFS(f, "travel/dist/*.html"))
}

var sv = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		m = mapper()

		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.SetHTMLTemplate(templ)

		r.GET("/api/stations", stations)
		r.GET("/api/search", search)
		r.GET("/api/pass", pass)

		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})
		r.GET("/assets/*.", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/fs/travel/dist"+c.Request.URL.Path)
		})
		r.StaticFS("/fs", http.FS(f))

		log.Println("start serve, open http://localhost:8080")
		if err := r.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

var wk = &cobra.Command{
	Use: "walk",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("start walk __no: %s, __date: %s, from: %s, to: %s", __no, __date, __fromZh, __toZh)
		t, _ := findTrainByNo(__no, m[__fromZh], m[__toZh], __date)
		log.Printf("%s\t[%s\t:%s]\ttwo：%s\tone：%s\tspecial：%s\tid:%s\n", t.TrainNo, m[t.FromStation], m[t.ToStation], t.TwoSeat, t.OneSeat, t.SpecialSeat, t.TrainCode)
		stations, _ := findPassStationsByCode(t.TrainCode, m[__fromZh], m[__toZh], __date)
		size := len(stations)
		for i := 1; i < size; i++ {
			t, _ = findTrainByCode(t.TrainCode, m[stations[0].StationName], m[stations[i].StationName], __date)
			log.Printf("%s\t[%s\t:%s]\ttwo：%s\tone：%s\tspecial：%s\tid:%s\n", t.TrainNo, m[t.FromStation], m[t.ToStation], t.TwoSeat, t.OneSeat, t.SpecialSeat, t.TrainCode)
		}
		log.Println("=====================================")
		for i := 0; i < size-1; i++ {
			t, _ = findTrainByCode(t.TrainCode, m[stations[i].StationName], m[stations[size-1].StationName], __date)
			log.Printf("%s\t[%s\t:%s]\ttwo：%s\tone：%s\tspecial：%s\tid:%s\n", t.TrainNo, m[t.FromStation], m[t.ToStation], t.TwoSeat, t.OneSeat, t.SpecialSeat, t.TrainCode)
		}
	},
}

var fw = &cobra.Command{
	Use: "fullWalk",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("start fullWalk no: %s, date: %s, from: %s, to: %s", __no, __date, __fromZh, __toZh)
		trains, _ := findAllTrain(m[__fromZh], m[__toZh], __date)
		for _, t := range trains {
			stations, _ := findPassStationsByCode(t.TrainCode, m[__fromZh], m[__toZh], __date)
			size := len(stations)
			t, _ = findTrainByCode(t.TrainCode, m[stations[0].StationName], m[stations[size-1].StationName], __date)
			log.Printf("%s\t[%s\t:%s\t]\ttwo:%s\tone:%s\tspecial:%s\tid:%s\n", t.TrainNo, m[t.FromStation], m[t.ToStation], t.TwoSeat, t.OneSeat, t.SpecialSeat, t.TrainCode)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(fw)
	rootCmd.AddCommand(wk)
	rootCmd.AddCommand(sv)

	wk.PersistentFlags().StringVarP(&__no, "no", "n", "", "findTrainByNo __no")
	rootCmd.PersistentFlags().StringVarP(&__date, "date", "d", "", "__date")
	rootCmd.PersistentFlags().StringVarP(&__fromZh, "from", "f", "", "from")
	rootCmd.PersistentFlags().StringVarP(&__toZh, "to", "t", "", "from")
	_ = rootCmd.MarkFlagRequired("date")
	_ = rootCmd.MarkFlagRequired("from")
	_ = rootCmd.MarkFlagRequired("to")

	m = mapper()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
