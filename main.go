package attribute

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	contentful "github.com/contentful-labs/contentful-go"
)

type Env struct {
	Client *contentful.Client
	Space  *contentful.Space
}

type ArticleAttribute struct {
	ArticleID string
	Slug      string
	Title     string
	Env       *Env
}

func New(articleID string, accessToken string, spaceID string) (*ArticleAttribute, error) {
	client := contentful.NewCMA(accessToken)
	space, err := client.Spaces.Get(spaceID)
	if err != nil {
		return nil, err
	}

	env := &Env{
		Client: client,
		Space:  space,
	}

	return &ArticleAttribute{ArticleID: articleID, Env: env}, nil
}

func (a *ArticleAttribute) Get() error {
	// Contentfulから記事情報を取得する
	entry, err := a.entry()
	if err != nil {
		return err
	}

	// Contentfulから記事情報が取得できない場合、処理を終了する
	if entry == nil {
		return errors.New("article not found")
	}

	slug, title, err := entryToArticleAttribute(entry)
	if err != nil {
		return err
	}

	a.Slug = slug
	a.Title = title

	return nil
}

func (a *ArticleAttribute) entry() (*contentful.Entry, error) {
	return a.Env.Client.Entries.Get(a.Env.Space.Sys.ID, a.ArticleID)
}

// toArticleAttribute は、 Contentful から取得した Entry を記事情報に変換する。
func entryToArticleAttribute(entry *contentful.Entry) (string, string, error) {
	var slug, title string
	var err error
	for attr, field := range entry.Fields {
		switch attr {
		case "slug":
			slug, err = fieldToString(field)
			if err != nil {
				fmt.Println(err)
				return "", "", err
			}
		case "title":
			value, err := fieldToString(field)
			if err != nil {
				fmt.Println(err)
				return "", "", err
			}
			title = strings.Replace(value, "/", "／", -1)
		}
	}

	return slug, title, nil
}

// fieldToString は、Contentful の field を文字列に変換する。
func fieldToString(field interface{}) (string, error) {
	type lang struct {
		Ja string `json:"ja"`
	}

	byte, err := json.Marshal(field)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var body lang
	if err := json.Unmarshal(byte, &body); err != nil {
		fmt.Println(err)
		return "", err
	}

	return body.Ja, nil
}
