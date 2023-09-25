package cmd

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

func search(c *gin.Context) {
	var req SearchReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from, to, date, no := req.FromStation, req.ToStation, req.Date, req.TrainNo
	if len(no) == 0 {
		trains, err := findAllTrain(from, to, date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, trains)
		return
	}
	log.Printf("from: %s, to: %s, date: %s, no: %s", from, to, date, no)
	if no[0] >= 'A' && no[0] <= 'Z' {
		train, err := findTrainByNo(no, from, to, date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, []*TrainRes{train})
		return
	}
	train, err := findTrainByCode(no, from, to, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, []*TrainRes{train})
}

func pass(c *gin.Context) {
	var req SearchReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from, to, date, no := req.FromStation, req.ToStation, req.Date, req.TrainNo
	log.Printf("from: %s, to: %s, date: %s, no: %s", from, to, date, no)
	if no[0] >= 'A' && no[0] <= 'Z' {
		c.JSON(http.StatusBadRequest, gin.H{"error": "train no error"})
		return
	}
	stationsByCode, err := findPassStationsByCode(no, from, to, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stationsByCode)
}
