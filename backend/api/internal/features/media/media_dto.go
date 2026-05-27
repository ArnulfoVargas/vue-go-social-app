package media

type MediaResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func MediaFromModel(m Media) MediaResponse {
	return MediaResponse{
		ID:  m.ID.Hex(),
		URL: m.URL,
	}
}

func MediasFromModels(ms []Media) []MediaResponse {
	medias := make([]MediaResponse, len(ms))
	for i, m := range ms {
		medias[i] = MediaFromModel(m)
	}
	return medias
}

func ModelsFromPointer(m []*Media) []Media {
	models := make([]Media, len(m))
	for i, model := range m {
		models[i] = *model
	}
	return models
}

func MediasFromModelsPointers(ms []*Media) []MediaResponse {
	medias := make([]MediaResponse, len(ms))
	for i, m := range ms {
		medias[i] = MediaFromModel(*m)
	}
	return medias
}
