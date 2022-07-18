package main

import (
	"Chotot-Microservice/cmd/pb"
	"Chotot-Microservice/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func login(client pb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := client.Login(c, &pb.LoginRequest{
			Phone:    req.Phone,
			Password: req.Passwd,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": res})
	}
}

func register(client pb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := client.Register(c, &pb.RegisterRequest{
			Phone:    req.Phone,
			Username: req.Username,
			Password: req.Passwd,
			Address:  req.Address,
			Email:    req.Email,
			IsAdmin:  req.IsAdmin,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": res})
	}
}
