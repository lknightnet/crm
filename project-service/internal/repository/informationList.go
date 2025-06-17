package repository

import (
	"project-service/internal/model"
	"project-service/internal/repository/customRepositoryError"
	"project-service/pkg/database"
)

type informationListRepository struct {
	db *database.PostgreSQL
}

func (i *informationListRepository) GetListByProjectID(projectID int) ([]model.InformationList, error) {
	var lists []model.InformationList
	if err := i.db.DB.Where("project_id = ?", projectID).Find(&lists).Error; err != nil {
		return nil, err
	}

	if len(lists) == 0 {
		return nil, customRepositoryError.ErrInformationListNotFound
	}
	return lists, nil
}

func (i *informationListRepository) CreateInformationList(informationList *model.InformationList) error {
	return i.db.DB.Create(informationList).Error
}

func (i *informationListRepository) UpdateInformationList(informationList *model.InformationList) error {
	return i.db.DB.Model(&model.InformationList{}).Where("id = ?", informationList.ID).Updates(informationList).Error

}

func (i *informationListRepository) DeleteInformationList(informationListID int) error {
	return i.db.DB.Delete(&model.InformationList{}, informationListID).Error
}

func newInformationListRepository(db *database.PostgreSQL) *informationListRepository {
	return &informationListRepository{
		db: db,
	}
}
