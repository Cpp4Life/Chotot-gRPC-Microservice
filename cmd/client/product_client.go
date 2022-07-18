package main

import (
	"Chotot-Microservice/cmd/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net/http"
)

func getProducts(client pb.ProductServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		stream, err := client.GetProducts(c, &emptypb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var productList []*pb.Product
		for {
			product, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			productList = append(productList, product)
		}
		c.JSON(http.StatusOK, gin.H{
			"response": productList,
		})
	}
}

func searchProducts(client pb.ProductServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		queries := c.Request.URL.Query()
		name := queries.Get("name")
		page := queries.Get("page")

		stream, err := client.SearchProducts(c, &pb.SearchRequest{
			KeyName: name,
			KeyPage: page,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var productList []*pb.Product
		for {
			product, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			productList = append(productList, product)
		}
		c.JSON(http.StatusOK, gin.H{
			"response": productList,
		})
	}
}
