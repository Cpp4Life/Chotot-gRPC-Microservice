package main

import (
	"Chotot-Microservice/cmd/pb"
	"Chotot-Microservice/constants"
	"Chotot-Microservice/dto"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func entityToProductProto(product dto.ProductDTO, header string) *pb.Product {
	return &pb.Product{
		Id:          product.Id,
		Name:        product.ProductName,
		CatId:       product.CatId,
		TypeId:      product.TypeId,
		Price:       product.Price,
		State:       product.State,
		CreatedTime: time.Now().Unix(),
		ExpiredTime: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Address:     product.Address,
		Content:     product.Content,
		Key:         &pb.UserToken{Header: header},
	}
}

func getUser(client pb.UserServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		header := c.GetHeader(constants.AuthKey)
		res, err := client.GetUser(c, &pb.UserToken{
			Header: header,
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			//return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": res,
		})
	}
}

func updateUser(client pb.UserServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		header := c.GetHeader(constants.AuthKey)
		var user dto.UserUpdateDTO
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		_, err := client.UpdateUser(c, &pb.UpdateRequest{
			Key:      &pb.UserToken{Header: header},
			Username: user.Username,
			Address:  user.Address,
			Email:    user.Email,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": true,
		})
	}
}

func getUserProducts(product pb.ProductServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		header := c.GetHeader(constants.AuthKey)
		stream, err := product.GetProductsByUserId(c, &pb.UserToken{
			Header: header,
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

func createUserProduct(product pb.ProductServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		header := c.GetHeader(constants.AuthKey)
		var productDTO dto.ProductDTO
		if err := c.ShouldBindJSON(&productDTO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := product.CreateProduct(c, entityToProductProto(productDTO, header))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": res,
		})
	}
}

func updateUserProduct(product pb.ProductServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		header := c.GetHeader(constants.AuthKey)
		productId := c.Param(constants.ParamId)
		var productDTO dto.ProductDTO
		if err := c.ShouldBindJSON(&productDTO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		request := entityToProductProto(productDTO, header)
		request.Id = productId
		_, err := product.UpdateProduct(c, request)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": true,
		})
	}
}

func deleteUserProduct(product pb.ProductServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		header := c.GetHeader(constants.AuthKey)
		productId := c.Param(constants.ParamId)
		_, err := product.DeleteProduct(c, &pb.ProductId{
			Id:  productId,
			Key: &pb.UserToken{Header: header},
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": true,
		})
	}
}
