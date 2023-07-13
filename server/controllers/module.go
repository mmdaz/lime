package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmdaz/lime/server/models"
)

func ModulesList(c *gin.Context) {

	allModules, err := models.FindAllModules()
	if err != nil {
		println(err.Error()) // return 500
	}

	// make module names list
	var modulesList []string
	for _, module := range *allModules {
		modulesList = append(modulesList, module.Name)
	}

	c.HTML(http.StatusOK, "modules.html", gin.H{
		"title":   "🔑 Modules",
		"Modules": modulesList,
	})
}

func CreateModule(c *gin.Context) {
	// Bind request body to struct
	request := &requestCreateModule{}
	err := c.BindJSON(request)
	if err != nil {
		respondJSON(c, 400, err.Error())
		return
	}

	module := models.Module{
		Name:   request.Name,
	}
	_, err = module.SaveModule()
	if err != nil {
		respondJSON(c, 500, err.Error())
		return
	}

	respondJSON(c, 200, "Module created")
}
