package resource

import (
	"net/http"
	"repository"
	"log"
	"github.com/alioygur/gores"
	"github.com/pressly/chi"
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
	r.Get("/:id", c.GetAllCategories)
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
	return
}
