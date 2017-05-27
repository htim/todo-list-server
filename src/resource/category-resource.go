package resource

import (
	"net/http"
	"repository"
	"log"
	"github.com/pressly/chi"
	"strconv"
	"encoding/json"
	"model"
	"env"
	"github.com/dgrijalva/jwt-go"
)

type CategoryResource struct {
	env *env.Env
	cr  repository.ICategoryRepository
}

func NewCategoryResource(categoryRepository repository.ICategoryRepository, env *env.Env) *CategoryResource {
	cr := CategoryResource{
		env: env,
		cr:  categoryRepository,
	}

	return &cr
}

func (c *CategoryResource) Router(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	for _, middleware := range middlewares {
		r.Use(middleware)
	}
	r.Get("/", c.GetAllCategories)
	r.Get("/:id", c.GetCategoryById)
	r.Post("/", c.CreateCategory)
	r.Put("/:id", c.UpdateCategory)
	r.Delete("/:id", c.DeleteCategory)
	return r
}

func (c *CategoryResource) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.cr.FindAllCategories()
	if err != nil {
		log.Print(err)
		c.env.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Error"})
		return
	}
	c.env.Render.JSON(w, http.StatusOK, categories)
}

func (c *CategoryResource) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	log.Println(user)
	if usr, ok := user.(*jwt.Token); ok {
		if mapClaims, ok := usr.Claims.(jwt.MapClaims); ok {
			login := mapClaims["login"]
			log.Println(login)
		}
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID param"})
		return
	}
	category, err := c.cr.FindCategoryById(id)
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		return
	}
	c.env.Render.JSON(w, http.StatusOK, category)
}

func (c *CategoryResource) CreateCategory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var cat model.Category
	err := decoder.Decode(&cat)
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid category payload"})
		return
	}
	id, err := c.cr.CreateCategory(cat)
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		return
	}
	c.env.Render.JSON(w, http.StatusOK, id)
}

func (c *CategoryResource) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID param"})
		return
	}
	decoder := json.NewDecoder(r.Body)
	var cat model.Category
	err = decoder.Decode(&cat)
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid category payload"})
		return
	}
	cat.ID = id
	err = c.cr.UpdateCategory(cat)
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		return
	}
	c.env.Render.JSON(w, http.StatusOK, "")
}

func (c *CategoryResource) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID param"})
		return
	}
	err = c.cr.DeleteCategory(id)
	if err != nil {
		log.Println(err)
		c.env.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		return
	}
	c.env.Render.JSON(w, http.StatusOK, "")
}
