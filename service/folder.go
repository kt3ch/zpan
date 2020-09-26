package service

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/saltbo/gopkg/gormutil"

	"github.com/saltbo/zpan/model"
	"github.com/saltbo/zpan/service/matter"
)

type Folder struct {
	matter.Matter
}

func NewFolder() *Folder {
	return &Folder{}
}

func (f *Folder) Rename(uid int64, alias, name string) error {
	m, err := f.FindUserMatter(uid, alias)
	if err != nil {
		return err
	}

	if _, ok := f.Exist(uid, name, m.Parent); ok {
		return fmt.Errorf("dir already exist a same name file")
	}

	children, err := f.findChildren(m)
	if err != nil {
		return err
	}

	oldParent := fmt.Sprintf("%s%s/", m.Parent, m.Name)
	newParent := fmt.Sprintf("%s%s/", m.Parent, name)
	fc := func(tx *gorm.DB) error {
		for _, v := range children {
			parent := strings.Replace(v.Parent, oldParent, newParent, 1)
			if err := tx.Model(v).Update("parent", parent).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(m).Update("name", name).Error; err != nil {
			return err
		}

		return nil
	}

	return gormutil.DB().Transaction(fc)
}

func (f *Folder) Move(uid int64, alias, parent string) error {
	m, err := f.FindUserMatter(uid, alias)
	if err != nil {
		return err
	}

	if err := f.copyOrMoveValidation(m, uid, parent); err != nil {
		return err
	}

	children, err := f.findChildren(m)
	if err != nil {
		return err
	}

	fc := func(tx *gorm.DB) error {
		for _, v := range children {
			err := tx.Model(v).Update("parent", parent+m.Name+"/").Error
			if err != nil {
				return err
			}
		}

		return tx.Model(m).Update("parent", parent).Error
	}
	return gormutil.DB().Transaction(fc)
}

func (f *Folder) Remove(uid int64, alias string) error {
	m, err := f.FindUserMatter(uid, alias)
	if err != nil {
		return err
	}

	children, err := f.findChildren(m)
	if err != nil {
		return err
	}

	fc := func(tx *gorm.DB) error {
		for _, v := range children {
			err := tx.Delete(v).Error
			if err != nil {
				return err
			}
		}

		return tx.Delete(m).Error
	}
	return gormutil.DB().Transaction(fc)
}

func (f *Folder) findChildren(m *model.Matter) ([]model.Matter, error) {
	var children []model.Matter
	oldParent := fmt.Sprintf("%s%s/", m.Parent, m.Name)
	if err := gormutil.DB().Where("uid=? and parent like ?", m.Uid, oldParent+"%").Find(&children).Error; err != nil {
		return nil, err
	}

	return children, nil
}

func (f *Folder) copyOrMoveValidation(m *model.Matter, uid int64, parent string) error {
	fmt.Println(parent, m.Parent+m.Name)
	if !m.IsDir() {
		return fmt.Errorf("only support direction")
	} else if parent == m.Parent {
		return fmt.Errorf("dir already in the dir")
	} else if parent != "" && strings.HasPrefix(parent, m.Parent+m.Name+"/") {
		return fmt.Errorf("can not move to itself")
	} else if !f.ParentExist(uid, parent) {
		return fmt.Errorf("dir does not exists")
	}

	if _, ok := f.Exist(m.Uid, m.Name, parent); ok {
		return fmt.Errorf("dir already has the same name file")
	}

	return nil
}
