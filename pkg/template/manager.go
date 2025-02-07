package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type Manager struct {
	templates map[string]*template.Template
	funcMap   template.FuncMap
}

func NewManager() *Manager {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"json": func(v interface{}) template.JS {
			b, err := json.Marshal(v)
			if err != nil {
				return template.JS("null")
			}
			return template.JS(b)
		},
	}
	return &Manager{
		templates: make(map[string]*template.Template),
		funcMap:   funcMap,
	}
}

func (m *Manager) InitTemplates() error {
	// 定義基礎布局、頁面和 partial 的路徑
	layoutFile := "../../public/html/layout.html"
	pagesGlob := "../../public/html/pages/*.html"
	partialsGlob := "../../public/html/partials/*.html"

	// 檢查 layout 檔案是否存在
	if _, err := os.Stat(layoutFile); os.IsNotExist(err) {
		return fmt.Errorf("layout file not found: %s", layoutFile)
	}

	// 獲取所有頁面模板
	pages, err := filepath.Glob(pagesGlob)
	if err != nil {
		return err
	}

	// 獲取所有 partial 模板
	partials, err := filepath.Glob(partialsGlob)
	if err != nil {
		return err
	}

	// 為每個頁面創建模板集合
	for _, page := range pages {
		name := filepath.Base(page)
		name = name[:len(name)-5]

		t := template.New(name).Funcs(m.funcMap)

		t, err = t.ParseFiles(layoutFile)
		if err != nil {
			log.Printf("Parse layout File Error: %s", err)
			return err
		}

		t, err = t.ParseFiles(page)
		if err != nil {
			log.Printf("Parse page File Error: %s", err)
			return err
		}

		if len(partials) > 0 {
			t, err = t.ParseFiles(partials...)
			if err != nil {
				log.Printf("Parse Partial File Error: %s", err)
				return err
			}
		}

		m.templates[name] = t
	}

	// 註冊 partial 模板
	for _, partial := range partials {
		name := filepath.Base(partial)
		name = name[:len(name)-5]

		t := template.New(name).Funcs(m.funcMap)
		t, err = t.ParseFiles(layoutFile)
		if err != nil {
			log.Printf("Parse layout File Error: %s", err)
			return err
		}
		t, err = t.ParseFiles(partial)
		if err != nil {
			log.Printf("Parse Partial File Error: %s", err)
			return err
		}

		m.templates[name] = t
	}

	return nil
}

// 取得所有模板
func (m *Manager) GetTemplates() map[string]*template.Template {
	return m.templates
}

// 取得單一模板
func (m *Manager) GetTemplate(name string) (*template.Template, bool) {
	tmpl, ok := m.templates[name]
	return tmpl, ok
}
