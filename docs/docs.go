// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/cart": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "Update the quantity of a cart item",
                "parameters": [
                    {
                        "description": "Update cart item request payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.UpdateCartItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Cart item updated successfully",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to update cart item",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "Add a product to the cart",
                "parameters": [
                    {
                        "description": "Add to cart request payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.AddToCartRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Product added to cart successfully",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to add product to cart",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "Remove a cart item",
                "parameters": [
                    {
                        "description": "Remove cart item request payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.RemoveCartItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Cart item removed successfully",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to remove cart item",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/cart/clear": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "Clear the cart",
                "parameters": [
                    {
                        "description": "Clear cart request payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.ClearCartRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Cart cleared successfully",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to clear cart",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/cart/{userID}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "Get the user's cart by user ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User's cart retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/services.GetCartResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid user ID",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch cart",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/products": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Products"
                ],
                "summary": "List all products",
                "responses": {
                    "200": {
                        "description": "List of products",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Product"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve products"
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Products"
                ],
                "summary": "Add a new product to the inventory",
                "parameters": [
                    {
                        "description": "Product details",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.AddProductRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Product added successfully",
                        "schema": {
                            "$ref": "#/definitions/services.AddProductResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to add product",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/products/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Products"
                ],
                "summary": "Get product details by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Product details",
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve product details"
                    }
                }
            },
            "put": {
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "Products"
                ],
                "summary": "Update an existing product in the inventory",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product name",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product description",
                        "name": "description",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product category",
                        "name": "category",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product ingredients",
                        "name": "ingredients",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product nutritional information",
                        "name": "nutritional_info",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product image URLs (comma-separated)",
                        "name": "image_urls",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Product updated successfully"
                    },
                    "400": {
                        "description": "Invalid form data"
                    },
                    "500": {
                        "description": "Failed to update product"
                    }
                }
            },
            "delete": {
                "tags": [
                    "Products"
                ],
                "summary": "Delete a product from the inventory",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Product deleted successfully"
                    },
                    "500": {
                        "description": "Failed to delete product"
                    }
                }
            }
        },
        "/productweight/{productID}/weights": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Product Weights"
                ],
                "summary": "Add a new weight variant for a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "productID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Product weight details",
                        "name": "weight",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.AddProductWeightRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Weight added successfully",
                        "schema": {
                            "$ref": "#/definitions/services.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request body or product ID",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to add weight",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/productweight/{productID}/weights/{weightID}": {
            "put": {
                "tags": [
                    "Product Weights"
                ],
                "summary": "Update an existing weight variant for a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "productID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Weight ID",
                        "name": "weightID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Product weight details",
                        "name": "weight",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.UpdateProductWeightRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Weight updated successfully"
                    },
                    "400": {
                        "description": "Invalid request body or product/weight ID"
                    },
                    "500": {
                        "description": "Failed to update weight"
                    }
                }
            }
        },
        "/purchase": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Purchase"
                ],
                "summary": "Create a purchase",
                "parameters": [
                    {
                        "description": "Purchase request payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PurchaseRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Purchase created successfully",
                        "schema": {
                            "$ref": "#/definitions/models.PurchaseResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to create purchase",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/purchase/{userID}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Purchase"
                ],
                "summary": "Get purchases by user ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Purchases retrieved successfully",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Purchase"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid user ID",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch purchases",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "User object",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User created successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Contact number already exists",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to create user",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Retrieve a user by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve user",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User object",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User updated successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to update user",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "Users"
                ],
                "summary": "Delete a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to delete user",
                        "schema": {
                            "$ref": "#/definitions/services.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Address": {
            "type": "object",
            "properties": {
                "addressID": {
                    "type": "integer"
                },
                "addressLine1": {
                    "type": "string"
                },
                "addressLine2": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                },
                "userID": {
                    "type": "integer"
                },
                "zipCode": {
                    "type": "string"
                }
            }
        },
        "models.CartItem": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "product": {
                    "$ref": "#/definitions/models.ProductWeight"
                },
                "quantity": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.Product": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image_urls": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "ingredients": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "nutritional_info": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "weights": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ProductWeight"
                    }
                }
            }
        },
        "models.ProductWeight": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "price": {
                    "type": "number"
                },
                "productID": {
                    "type": "integer"
                },
                "stockAvailability": {
                    "type": "integer"
                },
                "weight": {
                    "type": "integer"
                }
            }
        },
        "models.Purchase": {
            "type": "object",
            "properties": {
                "address_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PurchaseItem"
                    }
                },
                "payment_id": {
                    "type": "integer"
                },
                "total_price": {
                    "type": "number"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.PurchaseItem": {
            "type": "object",
            "properties": {
                "product_id": {
                    "type": "integer"
                },
                "product_name": {
                    "type": "string"
                },
                "product_price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                },
                "total_price": {
                    "type": "number"
                }
            }
        },
        "models.PurchaseRequest": {
            "type": "object",
            "properties": {
                "address_id": {
                    "type": "integer"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PurchaseItem"
                    }
                },
                "payment_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.PurchaseResponse": {
            "type": "object",
            "properties": {
                "purchase_id": {
                    "type": "integer"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Address"
                    }
                },
                "contactNumber": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "userID": {
                    "type": "integer"
                },
                "verifiedAccount": {
                    "type": "boolean"
                }
            }
        },
        "services.AddProductRequest": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "image_urls": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "ingredients": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "nutritional_info": {
                    "type": "string"
                }
            }
        },
        "services.AddProductResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "product_id": {
                    "type": "integer"
                }
            }
        },
        "services.AddProductWeightRequest": {
            "type": "object",
            "properties": {
                "price": {
                    "type": "number"
                },
                "stock_availability": {
                    "type": "integer"
                },
                "weight": {
                    "type": "integer"
                }
            }
        },
        "services.AddToCartRequest": {
            "type": "object",
            "properties": {
                "product_weight_id": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "services.ClearCartRequest": {
            "type": "object",
            "properties": {
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "services.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "services.GetCartResponse": {
            "type": "object",
            "properties": {
                "cart_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.CartItem"
                    }
                },
                "total_price": {
                    "type": "number"
                }
            }
        },
        "services.RemoveCartItemRequest": {
            "type": "object",
            "properties": {
                "product_weight_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "services.SuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "product_id": {
                    "type": "integer"
                }
            }
        },
        "services.UpdateCartItemRequest": {
            "type": "object",
            "properties": {
                "product_weight_id": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "services.UpdateProductWeightRequest": {
            "type": "object",
            "properties": {
                "price": {
                    "type": "number"
                },
                "stock_availability": {
                    "type": "integer"
                },
                "weight": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}