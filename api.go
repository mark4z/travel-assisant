package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func stations(c *gin.Context) {
	var selects []Select
	for k, v := range m {
		if k[0] > 'Z' || rune(k[0]) < 'A' {
			continue
		}
		selects = append(selects, Select{k, v})
	}
	c.JSON(http.StatusOK, selects)
}

func walk(c *gin.Context) {
	var req WalkReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from, to, date, no := req.FromStation, req.ToStation, req.Date, req.TrainNo
	log.Printf("from: %s, to: %s, date: %s, no: %s", from, to, date, no)
	train := findTrainByNo(no, from, to, date)
	c.JSON(http.StatusOK, train)
}

func fullWalk(c *gin.Context) {
	var req TrainReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from, to, date := req.FromStation, req.ToStation, req.Date
	log.Printf("from: %s, to: %s, date: %s", from, to, date)
	trains := findAllTrain(from, to, date)
	c.JSON(http.StatusOK, trains)
}
