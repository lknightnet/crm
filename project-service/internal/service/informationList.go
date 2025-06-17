package service

import (
	"errors"
	"project-service/internal/model"
	"project-service/internal/repository"
	"project-service/internal/repository/customRepositoryError"
	"project-service/internal/service/customServiceError"
	"project-service/pkg/tg"
)

type informationList struct {
	InformationListRepository repository.InformationListRepository
}

func (i *informationList) GetListByProjectID(projectID int) ([]model.InformationList, error) {
	lists, err := i.InformationListRepository.GetListByProjectID(projectID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrInformationListNotFound) {
			return nil, customServiceError.ErrInformationListNotFound
		}
		tg.SendError(err.Error(), "api/lists/get/project/:id")
		return nil, customServiceError.ErrUnknownError
	}
	return lists, nil
}

func (i *informationList) CreateInformationList(projectID int, key, value string) error {
	list := &model.InformationList{
		ProjectID: projectID,
		Key:       key,
		Value:     value,
	}
	err := i.InformationListRepository.CreateInformationList(list)
	if err != nil {
		tg.SendError(err.Error(), "api/lists/create")
		return customServiceError.ErrUnknownError
	}
	return nil
}

func (i *informationList) UpdateInformationList(listID int, key, value *string) error {
	list := &model.InformationList{
		ID: listID,
	}

	if key != nil {
		list.Key = *key
	}
	if value != nil {
		list.Value = *value
	}

	err := i.InformationListRepository.UpdateInformationList(list)
	if err != nil {
		tg.SendError(err.Error(), "api/lists/update")
		return customServiceError.ErrUnknownError
	}
	return nil
}

func (i *informationList) DeleteInformationList(listID int) error {
	err := i.InformationListRepository.DeleteInformationList(listID)
	if err != nil {
		tg.SendError(err.Error(), "api/list/delete/:id")
		return customServiceError.ErrUnknownError
	}
	return nil
}

func newInformationList(informationListRepository repository.InformationListRepository) *informationList {
	return &informationList{
		InformationListRepository: informationListRepository,
	}
}
