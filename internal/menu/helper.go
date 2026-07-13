package menu

import (
	"sort"

	"gin-layout/internal/domain"
)

func toMenuItem(m *domain.Menu) *domain.MenuItem {
	return &domain.MenuItem{
		ID: m.ID, ParentID: m.ParentID, Name: m.Name,
		Code: stringPtrValue(m.Code), Type: string(m.Type),
		Path: m.Path, Redirect: m.Redirect, Component: m.Component,
		Icon: m.Icon, ActiveMenu: m.ActiveMenu, Link: m.Link,
		Query: m.Query, Remark: m.Remark, Sort: m.Sort, Level: m.Level,
		Hidden: m.Hidden, Cache: m.Cache, Affix: m.Affix,
		Breadcrumb: m.Breadcrumb, AlwaysShow: m.AlwaysShow,
		External: m.External, Iframe: m.Iframe, Enabled: m.Enabled,
		Method: m.Method, APIPath: m.APIPath, PermCode: stringPtrValue(m.PermCode),
	}
}

func ToMenuTree(rows []domain.Menu) []domain.MenuItem {
	index := make(map[int64]domain.Menu, len(rows))
	for i := range rows {
		index[rows[i].ID] = rows[i]
	}

	childrenByParent := make(map[int64][]domain.Menu)
	roots := make([]domain.Menu, 0)
	for i := range rows {
		row := rows[i]
		if row.ParentID == nil {
			roots = append(roots, row)
			continue
		}
		if _, ok := index[*row.ParentID]; !ok {
			roots = append(roots, row)
			continue
		}
		childrenByParent[*row.ParentID] = append(childrenByParent[*row.ParentID], row)
	}

	items := make([]domain.MenuItem, 0, len(roots))
	for i := range roots {
		items = append(items, *buildTree(&roots[i], childrenByParent))
	}
	return items
}

func buildTree(row *domain.Menu, childrenByParent map[int64][]domain.Menu) *domain.MenuItem {
	item := toMenuItem(row)
	children := childrenByParent[row.ID]
	item.Children = make([]domain.MenuItem, 0, len(children))
	for i := range children {
		item.Children = append(item.Children, *buildTree(&children[i], childrenByParent))
	}
	return item
}

func sortMenus(rows []domain.Menu) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Sort == rows[j].Sort {
			return rows[i].ID < rows[j].ID
		}
		return rows[i].Sort < rows[j].Sort
	})
}

func applyCreateInput(m *domain.Menu, input CreateMenuReq) {
	m.ParentID = input.ParentID
	m.Code = input.Code
	m.Path = input.Path
	m.Redirect = input.Redirect
	m.Component = input.Component
	m.Icon = input.Icon
	m.ActiveMenu = input.ActiveMenu
	m.Link = input.Link
	m.Query = input.Query
	m.Remark = input.Remark
	m.Sort = input.Sort
	m.Level = input.Level
	if input.Hidden != nil {
		m.Hidden = *input.Hidden
	}
	if input.Cache != nil {
		m.Cache = *input.Cache
	}
	if input.Affix != nil {
		m.Affix = *input.Affix
	}
	if input.Breadcrumb != nil {
		m.Breadcrumb = *input.Breadcrumb
	}
	if input.AlwaysShow != nil {
		m.AlwaysShow = *input.AlwaysShow
	}
	if input.External != nil {
		m.External = *input.External
	}
	if input.Iframe != nil {
		m.Iframe = *input.Iframe
	}
	if input.Enabled != nil {
		m.Enabled = *input.Enabled
	}
	m.Method = input.Method
	m.APIPath = input.APIPath
	m.PermCode = input.PermCode
}

func applyUpdateFields(m *domain.Menu, in UpdateMenuReq) {
	if in.ParentID != nil {
		m.ParentID = in.ParentID
	}
	if in.Name != nil {
		m.Name = *in.Name
	}
	if in.Type != nil {
		m.Type = *in.Type
	}
	if in.Path != nil {
		m.Path = *in.Path
	}
	if in.Redirect != nil {
		m.Redirect = *in.Redirect
	}
	if in.Component != nil {
		m.Component = *in.Component
	}
	if in.Icon != nil {
		m.Icon = *in.Icon
	}
	if in.ActiveMenu != nil {
		m.ActiveMenu = *in.ActiveMenu
	}
	if in.Link != nil {
		m.Link = *in.Link
	}
	if in.Query != nil {
		m.Query = *in.Query
	}
	if in.Remark != nil {
		m.Remark = *in.Remark
	}
	if in.Sort != nil {
		m.Sort = *in.Sort
	}
	if in.Level != nil {
		m.Level = *in.Level
	}
	if in.Hidden != nil {
		m.Hidden = *in.Hidden
	}
	if in.Cache != nil {
		m.Cache = *in.Cache
	}
	if in.Affix != nil {
		m.Affix = *in.Affix
	}
	if in.Breadcrumb != nil {
		m.Breadcrumb = *in.Breadcrumb
	}
	if in.AlwaysShow != nil {
		m.AlwaysShow = *in.AlwaysShow
	}
	if in.External != nil {
		m.External = *in.External
	}
	if in.Iframe != nil {
		m.Iframe = *in.Iframe
	}
	if in.Enabled != nil {
		m.Enabled = *in.Enabled
	}
	if in.Method != nil {
		m.Method = *in.Method
	}
	if in.APIPath != nil {
		m.APIPath = *in.APIPath
	}
	if in.PermCode != nil {
		m.PermCode = in.PermCode
	}
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
