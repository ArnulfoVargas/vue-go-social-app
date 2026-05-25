package dto

import "Server/internal/model"

type Media struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func MediaFromModel(m model.Media) Media {
	return Media{
		ID:  m.ID.Hex(),
		URL: m.URL,
	}
}

func MediasFromModels(ms []model.Media) []Media {
	medias := make([]Media, len(ms))
	for i, m := range ms {
		medias[i] = MediaFromModel(m)
	}
	return medias
}

func ModelsFromPointer(m []*model.Media) []model.Media {
	models := make([]model.Media, len(m))
	for i, model := range m {
		models[i] = *model
	}
	return models
}

func MediasFromModelsPointers(ms []*model.Media) []Media {
	medias := make([]Media, len(ms))
	for i, m := range ms {
		medias[i] = MediaFromModel(*m)
	}
	return medias
}
