package core

import (
	"github.com/naoina/genmai"
	"github.com/raa0121/go-kokoro-io/models"
	"log"
	"strings"
)

var Db *genmai.DB

type Onunume struct {
	Id        int64  `db:"pk"`
	Name      string `db:"unique" size:"255"`
	Regexp    string `size:"255"`
	Content   string `size:"255"`
	CreatedBy string `size:"255"`
	Enable    bool   `db:"-"`
	Expect    string `size:"255"`
}

func Add(m models.MessageEntity) string {
	var result []Onunume
	s := strings.Split(m.RawContent, " ")
	name, regexp, content := s[2], s[3], s[4]
	if err := Db.Select(&result, Db.Where("name", "=", name)); err != nil {
		log.Fatal(err)
	}
	if len(result) == 0 {
		obj := &Onunume{
			Name: name,
			Regexp: regexp,
			Content: content,
			CreatedBy: m.Profile.ID,
			Enable: true,
		}
		if _, err := Db.Insert(obj); err != nil {
			return "DBにInsertできませんでした"
		}
	} else {
		obj := &Onunume{
			Id: result[0].Id,
			Name: name,
			Regexp: regexp,
			Content: content,
			CreatedBy: m.Profile.ID,
			Enable: true,
		}
		if _, err := Db.Update(obj); err != nil {
			return "更新できませんでした"
		}
		return "更新しました"
	}
	return ""
}

func Remove(m models.MessageEntity) string {
	var result []Onunume
	s := strings.Split(m.RawContent, " ")
	name := s[2]
	if err := Db.Select(&result, Db.Where("name", "=", name)); err != nil {
		log.Fatal(err)
	}
	if len(result) == 0 {
		return "削除対象がありませんでした"
	} else {
		if _, err := Db.Delete(result); err != nil {
			return "削除できませんでした"
		}
		return "削除しました"
	}
	return ""
}

func Help(m models.MessageEntity) string {
	// TODO: Help を書く
	return "Help なんてものはないです"
}
