package apiserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SociumR/bmlabs-test/helpers"
	"github.com/SociumR/bmlabs-test/models"
	"github.com/SociumR/bmlabs-test/store"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *apiserver) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := u.Validate(); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := s.store.Insert(&u, store.CollectionUser)

	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	u.ID = id

	s.respond(w, r, http.StatusOK, u)

}

func (s *apiserver) CreateGame(w http.ResponseWriter, r *http.Request) {
	var g models.Game

	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := g.Validate(); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if _, err := s.store.Insert(g, store.CollectionGame); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	s.respond(w, r, http.StatusOK, g)

}

func (s *apiserver) GetUsers(w http.ResponseWriter, r *http.Request) {
	s.PageResp(store.CollectionUser, w, r)
}

func (s *apiserver) GetGames(w http.ResponseWriter, r *http.Request) {
	s.PageResp(store.CollectionGame, w, r)
}

func (s *apiserver) Stat(w http.ResponseWriter, r *http.Request) {

	group := r.URL.Query().Get("query")

	mp := map[string]interface{}{}

	err := json.Unmarshal([]byte(group), &mp)

	m, err := s.store.Group(store.CollectionUser, helpers.MapStringToBSonM(mp))

	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	s.respond(w, r, http.StatusOK, m)
}

func (s *apiserver) PageResp(c string, w http.ResponseWriter, r *http.Request) {
	resp, query := s.PreparePageResponse(r)

	count, err := s.store.Count(0, 0, query, c)

	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	resp.Records = count

	skip := resp.PageLimit * (resp.CurrentPage - 1)

	m, err := s.store.Find(c, helpers.MapStringToBSonD(query), &options.FindOptions{
		Limit: &resp.PageLimit,
		Skip:  &skip,
	})

	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	resp.Data = m

	s.respond(w, r, http.StatusOK, resp)
}

func (s *apiserver) PreparePageResponse(r *http.Request) (models.PageResponse, map[string]interface{}) {

	q := r.URL.Query()

	m := models.PageResponse{
		PageLimit:   100,
		CurrentPage: 1,
		Records:     1,
	}

	if len(q) > 0 {
		pageLimit, err := strconv.Atoi(q.Get("pageLimit"))

		if err == nil {
			q.Del("pageLimit")
			m.PageLimit = int64(pageLimit)
		}

		page, err := strconv.Atoi(q.Get("page"))

		if err == nil {
			q.Del("page")
			m.CurrentPage = int64(page)
		}

	}

	query := map[string]interface{}{}

	for k, v := range q {
		query[k] = v[0]
	}

	return m, query
}

func (s *apiserver) error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	s.respond(w, r, statusCode, map[string]string{"error": err.Error()})
}

func (s *apiserver) respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {

	defer r.Body.Close()

	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
