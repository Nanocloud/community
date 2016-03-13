/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Nanocloud/community/nanocloud/models/ldap"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

func Delete(c *echo.Context) error {
	userId := c.Param("id")
	if len(userId) == 0 {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "User id needed for deletion",
				},
			},
		})
	}

	user, err := users.GetUser(userId)
	if err != nil {
		return err
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, hash{
			"error": [1]hash{
				hash{
					"detail": "User not found",
				},
			},
		})
	}

	if user.IsAdmin {
		return c.JSON(http.StatusUnauthorized, hash{
			"error": [1]hash{
				hash{
					"detail": "Admins cannot be deleted",
				},
			},
		})
	}

	err = ldap.DeleteAccount(user.Id)
	if err != nil {
		log.Errorf("Unable to delete user in ad: %s", err.Error())
		switch err {
		case ldap.DeleteFailed:
			return c.JSON(http.StatusInternalServerError, hash{
				"error": [1]hash{
					hash{
						"detail": err.Error(),
					},
				},
			})
		case ldap.UnknownUser:
			log.Info("User doesn't exist in AD")
			break
		default:
			return err
		}
	}

	err = users.DeleteUser(user.Id)
	if err != nil {
		log.Errorf("Unable to delete user: ", err.Error())
		return err
	}

	return c.JSON(http.StatusOK, hash{
		"data": hash{
			"success": true,
		},
	})
}

func Disable(userId string) (int, error) {
	if userId == "" {
		return http.StatusNotFound, errors.New("User id needed for desactivation")
	}

	exists, err := users.UserExists(userId)
	if err != nil {
		return http.StatusConflict, err
	}

	if !exists {
		return http.StatusNotFound, errors.New("User not found")
	}

	err = users.DisableUser(userId)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Unable to disable user: " + err.Error())
	}

	return 0, nil
}

func Update(c *echo.Context) error {
	var attr hash

	err := utils.ParseJSONBody(c, &attr)
	if err != nil {
		return nil
	}

	data, ok := attr["data"].(map[string]interface{})
	if ok == false {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "data is missing",
				},
			},
		})
	}

	attributes, ok := data["attributes"].(map[string]interface{})
	if ok == false {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "attributes is missing",
				},
			},
		})
	}

	user, err := users.GetUser(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "Could not find user",
				},
			},
		})
	}

	password, ok := attributes["password"].(string)
	if ok == false || password == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "password is missing",
				},
			},
		})
	}

	error := users.UpdateUserPassword(user.Id, password);
	if error != nil {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": "Cannot set new password",
				},
			},
		})
	}
	return c.JSON(http.StatusOK, hash{
		"data": hash{
			"id":         user.Id,
			"type":       "user",
			"attributes": user,
		},
	})
}

func Get(c *echo.Context) error {
	user := c.Get("user").(*users.User)
	if c.Query("me") == "true" {
		return c.JSON(http.StatusOK, hash{
			"data": hash{
				"id":         user.Id,
				"type":       "user",
				"attributes": user,
			},
		})
	}

	if !user.IsAdmin {
		return c.JSON(http.StatusUnauthorized, hash{
			"error": [1]hash{
				hash{
					"detail": "Unauthorized",
				},
			},
		})
	}

	users, err := users.FindUsers()
	if err != nil {
		return errors.New(
			fmt.Sprintf("unable to get user lists: %s", err.Error()),
		)
	}

	var response = make([]hash, len(*users))
	for i, val := range *users {
		res := hash{
			"id":         val.Id,
			"type":       "user",
			"attributes": val,
		}
		response[i] = res
	}
	return c.JSON(http.StatusOK, hash{"data": response})
}

func Post(c *echo.Context) error {
	var attr hash

	err := utils.ParseJSONBody(c, &attr)
	if err != nil {
		return err
	}

	data, ok := attr["data"].(map[string]interface{})
	if ok == false {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "data is missing",
				},
			},
		})
	}

	attributes, ok := data["attributes"].(map[string]interface{})
	if ok == false {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "attributes is missing",
				},
			},
		})
	}

	email, ok := attributes["email"].(string)
	if ok == false || email == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "email is missing",
				},
			},
		})
	}

	firstName, ok := attributes["first-name"].(string)
	if ok == false || firstName == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "first-name is missing",
				},
			},
		})
	}

	lastName, ok := attributes["last-name"].(string)
	if ok == false || lastName == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "last-name is missing",
				},
			},
		})
	}

	password, ok := attributes["password"].(string)
	if ok == false || password == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "password is missing",
				},
			},
		})
	}

	newUser, err := users.CreateUser(
		true,
		email,
		firstName,
		lastName,
		password,
		false,
	)
	switch err {
	case users.UserDuplicated:
		return c.JSON(http.StatusConflict, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
	case users.UserNotCreated:
		return err
	}

	winpass := utils.RandomString(8) + "s4D+"
	sam, err := ldap.AddUser(newUser.Id, winpass)
	if err != nil {
		return err
	}

	err = users.UpdateUserAd(newUser.Id, sam, winpass)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, hash{
		"data": hash{
			"id":         newUser.Id,
			"type":       "user",
			"attributes": newUser,
		},
	})
}

func UpdatePassword(c *echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "User id needed to modify account",
				},
			},
		})
	}

	var user struct {
		Data struct {
			Password string
		}
	}

	err := utils.ParseJSONBody(c, &user)
	if err != nil {
		return nil
	}

	exists, err := users.UserExists(userId)
	if err != nil {
		log.Errorf("Unable to check user existance: %s", err.Error())
		return err
	}

	if !exists {
		return c.JSON(http.StatusNotFound, hash{
			"error": [1]hash{
				hash{
					"detail": "User not found",
				},
			},
		})
	}

	err = users.UpdateUserPassword(userId, user.Data.Password)
	if err != nil {
		log.Errorf("Unable to update user password: %s", err.Error())
		return err
	}


	return c.JSON(http.StatusOK, hash{
		"data": hash{
			"success": true,
		},
	})
}

func GetUser(c *echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "User id needed to retrieve account informations",
				},
			},
		})
	}

	user, err := users.GetUser(userId)
	if err != nil {
		return err
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, hash{
			"error": [1]hash{
				hash{
					"detail": "User Not Found",
				},
			},
		})
	}

	return c.JSON(http.StatusOK, hash{
		"data": hash{
			"id":         user.Id,
			"type":       "user",
			"attributes": user,
		},
	})
}
