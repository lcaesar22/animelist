package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/konaha-gakure/otaku/internal/data"
	"github.com/konaha-gakure/otaku/internal/validator"
)

func (app *application) listAnimesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		Genre []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	// get the title and genre
	input.Title = app.readString(qs, "title", "")
	input.Genre = app.readCSV(qs, "genre", []string{})

	// get the page and page size
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime", "-id", "-year", "-runtime"}

	// executes the validation check on the filter structs
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Error)
		return
	}

	animes, metadata, err := app.models.Animes.GetAll(input.Title, input.Genre, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"animes": animes, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// calls the get method to fetch data
	anime, err := app.models.Animes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"anime": anime}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createAnimeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int          `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genre   []string     `json:"genre"`
	}

	// readJSON decodes the request body into an input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	v := validator.New()

	// checking if checks failed
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Error)
		return
	}

	anime := &data.Anime{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genre:   input.Genre,
	}
	if data.ValidateAnime(v, anime); !v.Valid() {
		app.failedValidationResponse(w, r, v.Error)
	}

	err = app.models.Animes.Insert(anime)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// send a location header to show url
	headers := make(http.Header)
	headers.Add("Location", fmt.Sprintf("/animes/%d", anime.ID))

	// JSON response with a 201 created
	err = app.writeJSON(w, http.StatusCreated, envelope{"anime": anime}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v", input)
}

// update a movie
func (app *application) updateAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// fetches the existing anime record
	anime, err := app.models.Animes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Round-trip locking
	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(int(anime.Version)) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int          `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genre   []string      `json:"genre"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		anime.Title = *input.Title
	}

	if input.Year != nil {
		anime.Year = *input.Year
	}

	if input.Runtime != nil {
		anime.Runtime = *input.Runtime
	}

	if input.Genre != nil {
		anime.Genre = input.Genre
	}

	// validate the updated anime
	v := validator.New()

	if data.ValidateAnime(v, anime); !v.Valid() {
		app.failedValidationResponse(w, r, v.Error)
		return
	}

	// intercept the ErrEditConflict error
	err = app.models.Animes.Update(anime)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	/// writes the JSON response
	err = app.writeJSON(w, http.StatusOK, envelope{"anime": anime}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// delete anime handler
func (app *application) deleteAnimeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// delete the movie from database
	err = app.models.Animes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// sends an error 202 OK status code if successful
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "anime was successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
