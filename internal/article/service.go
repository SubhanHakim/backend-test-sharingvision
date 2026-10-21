package article

import (
	"strings"
	"unicode/utf8"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo} }

func (s *Service) Create(a Article) error {
	if err := validate(a); err != nil {
		return err
	}
	return s.repo.Create(a)
}

func (s *Service) List(limit, offset int) ([]Article, error) { return s.repo.List(limit, offset) }

func (s *Service) Get(id int) (Article, error) { return s.repo.Get(id) }

func (s *Service) Update(id int, a Article) error {
	if err := validate(a); err != nil {
		return err
	}
	return s.repo.Update(id, a)
}

func (s *Service) Delete(id int) error { return s.repo.Delete(id) }

func validate(a Article) error {
	var errs ValidationError
	if minLen(a.Title, 20) {
		errs = append(errs, "title: required, min 20 characters")
	}
	if minLen(a.Content, 200) {
		errs = append(errs, "content: required, min 200 characters")
	}
	if minLen(a.Category, 3) {
		errs = append(errs, "category: required, min 3 characters")
	}
	if s := a.Status; s != "publish" && s != "draft" && s != "thrash" {
		errs = append(errs, "status: required, must be publish, draft or thrash")
	}
	if errs != nil {
		return errs
	}
	return nil
}

// minLen true jika teks (setelah trim) kurang dari n karakter.
func minLen(s string, n int) bool { return utf8.RuneCountInString(strings.TrimSpace(s)) < n }
