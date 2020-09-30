package controllers

import "website_sc/models"

type ReposController interface {
	DatabaseGet(url *models.Pingdom) error
	DatabaseGets(url *[]models.Pingdom) error
	DatabaseSave(url *models.Pingdom) error
	DatabaseCreate(url *models.Pingdom) error
	DatabaseDelete(url *models.Pingdom) error
}

type MonitorRepo struct{}

var repo ReposController

func setRepoController(repoType ReposController) {
	repo = repoType
}
func (rp *MonitorRepo) DatabaseGet(url *models.Pingdom) error {
	return models.DB.Find(&url).Error
}

func (rp *MonitorRepo) DatabaseCreate(url *models.Pingdom) error {
	return models.DB.Create(&url).Error
}
func (rp *MonitorRepo) DatabaseSave(url *models.Pingdom) error {
	return models.DB.Save(&url).Error
}
func (rp *MonitorRepo) DatabaseDelete(url *models.Pingdom) error {
	return models.DB.Delete(&url).Error
}

func (rp *MonitorRepo) DatabaseGets(url *[]models.Pingdom) error {
	return models.DB.First(&url).Error
}
