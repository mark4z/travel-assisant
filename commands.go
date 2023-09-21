package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	no     string
	date   string
	fromZh string
	toZh   string
)

var rootCmd = &cobra.Command{
	Use:              "travel",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var wk = &cobra.Command{
	Use: "walk",
	Run: func(cmd *cobra.Command, args []string) {
		index(m[fromZh], fromZh, m[toZh], toZh, date)

		id := train(no, m[fromZh], m[toZh], date)
		stations := pass(id, m[fromZh], m[toZh], date)
		walk(stations, id, date)
	},
}

var fw = &cobra.Command{
	Use: "fullWalk",
	Run: func(cmd *cobra.Command, args []string) {
		index(m[fromZh], fromZh, m[toZh], toZh, date)
		ids := trainAll(m[fromZh], m[toZh], date)
		for _, id := range ids {
			stations := pass(id, m[fromZh], m[toZh], date)
			fullWalk(stations, id, date)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(fw)
	rootCmd.AddCommand(wk)

	rootCmd.PersistentFlags().StringVarP(&no, "no", "n", "", "train no")
	rootCmd.PersistentFlags().StringVarP(&date, "date", "d", "", "date")
	rootCmd.PersistentFlags().StringVarP(&fromZh, "from", "f", "", "from")
	rootCmd.PersistentFlags().StringVarP(&toZh, "to", "t", "", "from")
	_ = rootCmd.MarkFlagRequired("date")
	_ = rootCmd.MarkFlagRequired("from")
	_ = rootCmd.MarkFlagRequired("to")

	m = mapper()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
