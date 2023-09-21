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
		index(m[__fromZh], __fromZh, m[__toZh], __toZh, __date)

		id := train(__no, m[__fromZh], m[__toZh], __date)
		stations := pass(id, m[__fromZh], m[__toZh], __date)
		walk(stations, id, __date)
	},
}

var fw = &cobra.Command{
	Use: "fullWalk",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("start fullWalk __no: %s, __date: %s, from: %s, to: %s", __no, __date, __fromZh, __toZh)
		index(m[__fromZh], __fromZh, m[__toZh], __toZh, __date)
		ids := trainAll(m[__fromZh], m[__toZh], __date)
		for _, id := range ids {
			stations := pass(id, m[__fromZh], m[__toZh], __date)
			fullWalk(stations, id, __date)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(fw)
	rootCmd.AddCommand(wk)
	rootCmd.AddCommand(sv)

	wk.PersistentFlags().StringVarP(&__no, "__no", "n", "", "train __no")
	rootCmd.PersistentFlags().StringVarP(&__date, "__date", "d", "", "__date")
	rootCmd.PersistentFlags().StringVarP(&__fromZh, "from", "f", "", "from")
	rootCmd.PersistentFlags().StringVarP(&__toZh, "to", "t", "", "from")
	_ = rootCmd.MarkFlagRequired("__date")
	_ = rootCmd.MarkFlagRequired("from")
	_ = rootCmd.MarkFlagRequired("to")

	m = mapper()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
