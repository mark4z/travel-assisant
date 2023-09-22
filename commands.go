package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
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

var sv = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

var wk = &cobra.Command{
	Use: "walk",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("start walk __no: %s, __date: %s, from: %s, to: %s", __no, __date, __fromZh, __toZh)
		t := findTrainByNo(__no, m[__fromZh], m[__toZh], __date)
		log.Printf("%s\t[%s\t:%s]\ttwo：%s\tone：%s\tspecial：%s\tid:%s\n", t.TrainNo, m[t.FromStation], m[t.ToStation], t.TwoSeat, t.OneSeat, t.SpecialSeat, t.TrainCode)
		stations := findPassStationsByCode(t.TrainCode, m[__fromZh], m[__toZh], __date)
		size := len(stations)
		for i := 1; i < size; i++ {
			t = findTrainByCode(t.TrainCode, m[stations[0].StationName], m[stations[i].StationName], __date)
			log.Printf("%s\t[%s\t:%s]\ttwo：%s\tone：%s\tspecial：%s\tid:%s\n", t.TrainNo, m[t.FromStation], m[t.ToStation], t.TwoSeat, t.OneSeat, t.SpecialSeat, t.TrainCode)
		}
		log.Println("=====================================")
		for i := 0; i < size-1; i++ {
			t = findTrainByCode(t.TrainCode, m[stations[i].StationName], m[stations[size-1].StationName], __date)
			log.Printf("%s\t[%s\t:%s]\ttwo：%s\tone：%s\tspecial：%s\tid:%s\n", t.TrainNo, m[t.FromStation], m[t.ToStation], t.TwoSeat, t.OneSeat, t.SpecialSeat, t.TrainCode)
		}
	},
}

var fw = &cobra.Command{
	Use: "fullWalk",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("start fullWalk no: %s, date: %s, from: %s, to: %s", __no, __date, __fromZh, __toZh)
		trains := findAllTrain(m[__fromZh], m[__toZh], __date)
		for _, t := range trains {
			stations := findPassStationsByCode(t.TrainCode, m[__fromZh], m[__toZh], __date)
			size := len(stations)
			t = findTrainByCode(t.TrainCode, m[stations[0].StationName], m[stations[size-1].StationName], __date)
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
