package resource

import (
	"net/http"
	"repository"
	"log"
	"github.com/alioygur/gores"
	"github.com/pressly/chi"
	"strconv"
	"encoding/json"
	"model"
)

type CategoryResource struct {
	cr repository.ICategoryRepository
}

func NewCategoryResource(categoryRepository repository.ICategoryRepository) *CategoryResource {
	cr := CategoryResource{
		cr: categoryRepository,
	}

	return &cr
}

func (c *CategoryResource) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", c.GetAllCategories)
	r.Get("/:id", c.GetCategoryById)
	r.Post("/", c.CreateCategory)
	return r
}


func (c *CategoryResource) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.cr.FindAllCategories()
	if err != nil {
		log.Print(err)
		gores.JSON(w, http.StatusInternalServerError, map[string]string{"error":"Internal Error"})
		return
	}
	gores.JSON(w, http.StatusOK, categories)
}

func (c *CategoryResource) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r,"id"))
	if err != nil {
		log.Println(err)
		gores.JSON(w, http.StatusBadRequest, map[string]string{"error":"Invalid ID param"})
		return
	}
	category, err := c.cr.FindCategoryById(id)
	if err != nil {
		log.Println(err)
		gores.JSON(w, http.StatusInternalServerError, map[string]string{"error":"Internal error"})
		return
	}
	gores.JSON(w, http.StatusOK, category)
}

func (c *CategoryResource) CreateCategory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var cat model.Category
	err := decoder.Decode(&cat)
	if err != nil {
		log.Println(err)
		gores.JSON(w, http.StatusBadRequest, map[string]string{"error":"Invalid category payload"})
		return
	}
	id, err := c.cr.CreateCategory(cat)
	if err != nil {
		log.Println(err)
		gores.JSON(w, http.StatusInternalServerError, map[string]string{"error":"Internal error"})
		return
	}
	gores.JSON(w, http.StatusOK, id)
}