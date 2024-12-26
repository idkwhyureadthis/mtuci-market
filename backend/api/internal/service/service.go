package service

import (
	"api-service/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Service struct {
}

func New() *Service {
	s := Service{}
	return &s
}

func formToReader(form *multipart.Form) (io.Reader, string, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	for key, values := range form.Value {
		for _, value := range values {
			err := writer.WriteField(key, value)
			if err != nil {
				return nil, "", err
			}
		}
	}
	for fieldName, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				return nil, "", err
			}
			defer file.Close()

			part, err := writer.CreateFormFile(fieldName, fileHeader.Filename)
			if err != nil {
				return nil, "", err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, "", err
			}
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf, writer.FormDataContentType(), nil
}

func (s *Service) OnModeration() (*model.Card, error) {
	resp, err := http.Get("http://db-service:8081/on_moderation")
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("failed to get cards")
	}
	defer resp.Body.Close()
	var card model.Card
	err = json.NewDecoder(resp.Body).Decode(&card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (s *Service) CreateModerator(data struct {
	Name     string `json:"name"`
	Telegram string `json:"telegram"`
	Password string `json:"password"`
	Login    string `json:"login"`
}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post("http://db-service:8081/create_moderator", "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != 200 {
		return errors.New("failed to create moderator")
	}
	return nil
}

func (s *Service) DeleteCard(id string) error {
	resp, err := http.Post(fmt.Sprintf("http://db-service:8081/delete_card/%s", id), "plain/text", nil)
	if resp.StatusCode != 200 || err != nil {
		return errors.New("failed to delete card")
	}
	return nil
}

func (s *Service) ProductsOf(id int) ([]*model.Card, error) {
	cards := []*model.Card{}

	resp, err := http.Get(fmt.Sprintf("http://db-service:8081/cards_of?id=%d", id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&cards)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *Service) GetCards() ([]*model.Card, error) {
	cards := []*model.Card{}

	resp, err := http.Get("http://db-service:8081/moderated")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&cards)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *Service) CreateCard(form *multipart.Form) (*model.Tokens, error) {
	toRet := model.Tokens{}
	access, refresh := form.Value["access"], form.Value["refresh"]
	if len(access) != 1 || len(refresh) != 1 {
		return nil, errors.New("failed to process data")
	}
	userData, id, err := s.Verify(access[0], refresh[0])
	if userData != nil && userData.NewTokens != nil {
		toRet = *userData.NewTokens
	}
	if err != nil {
		return &toRet, errors.New("not authorized")
	}
	formAsReader, contentType, err := formToReader(form)
	if err != nil {
		return &toRet, err
	}
	createCardResp, err := http.Post(fmt.Sprintf("http://db-service:8081/create_card?id=%d", id), contentType, formAsReader)
	if err != nil || createCardResp.StatusCode != 200 {
		return &toRet, errors.New("failed to create card")
	}
	return &toRet, nil
}

func (s *Service) SignUp(login, password, dormNumber, room, name, telegram string) (*model.Tokens, error) {
	tokens := model.Tokens{}
	data := struct {
		Name       string `json:"name"`
		Telegram   string `json:"telegram"`
		Password   string `json:"password"`
		Room       string `json:"room"`
		DormNumber string `json:"dorm_number"`
		Login      string `json:"login"`
	}{
		Login:      login,
		Password:   password,
		DormNumber: dormNumber,
		Room:       room,
		Name:       name,
		Telegram:   telegram,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("http://db-service:8081/create_user", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("failed to create user")
	}
	respData := struct {
		Id int `json:"id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}
	tokenResp, err := http.Get(fmt.Sprintf("http://auth-service:8080/new?id=%d&role=user", respData.Id))
	if err != nil {
		return nil, err
	}
	defer tokenResp.Body.Close()
	err = json.NewDecoder(tokenResp.Body).Decode(&tokens)
	if err != nil {
		return nil, err
	}
	setTokenResp, err := http.Post(fmt.Sprintf("http://db-service:8081/set_refresh?id=%d&refresh=%s", respData.Id, tokens.Refresh), "plain/text", nil)
	if err != nil {
		return nil, err
	}
	if setTokenResp.StatusCode != 200 {
		return nil, errors.New("failed to set refresh token")
	}
	return &tokens, nil
}

func (s *Service) LogIn(login, password string) (*model.Tokens, error) {
	verifyData := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    login,
		Password: password,
	}
	JsonVerifyData, err := json.Marshal(verifyData)
	if err != nil {
		return nil, err
	}
	verifyResp, err := http.Post("http://db-service:8081/check_password", "application/json", bytes.NewBuffer(JsonVerifyData))
	if err != nil {
		return nil, err
	}
	if verifyResp.StatusCode != 200 {
		return nil, errors.New("failed to verify password")

	}
	defer verifyResp.Body.Close()
	respData := struct {
		Status string `json:"status"`
		Id     int    `json:"id"`
	}{}
	err = json.NewDecoder(verifyResp.Body).Decode(&respData)
	if err != nil {
		return nil, errors.New("failed to decode response")
	}
	tokenResp, err := http.Get(fmt.Sprintf("http://auth-service:8080/new?id=%d&role=%s", respData.Id, respData.Status))
	if err != nil {
		return nil, errors.New("failed to get Tokens")
	}
	defer tokenResp.Body.Close()
	tokens := model.Tokens{}
	err = json.NewDecoder(tokenResp.Body).Decode(&tokens)
	if err != nil {
		return nil, errors.New("failed to decode response")
	}
	updateTokenResp, err := http.Post(fmt.Sprintf("http://db-service:8081/set_refresh?id=%d&refresh=%s", respData.Id, tokens.Refresh), "plain/text", nil)
	if err != nil || updateTokenResp.StatusCode != 200 {
		return nil, errors.New("failed to update refresh")
	}
	return &tokens, nil
}

func (s *Service) Verify(access, refresh string) (*model.UserData, int, error) {
	userData := model.UserData{}
	accJson, err := json.Marshal(struct {
		Token string `json:"token"`
	}{Token: access})
	if err != nil {
		return nil, -1, err
	}
	refJson, err := json.Marshal(struct {
		Token string `json:"token"`
	}{Token: refresh})
	if err != nil {
		return nil, -1, err
	}
	accResp, err := http.Post("http://auth-service:8080/verify", "application.json", bytes.NewBuffer(accJson))
	if err != nil {
		return nil, -1, errors.New("failed to verify Tokens")
	}
	defer accResp.Body.Close()
	if accResp.StatusCode == 200 {
		accRespData := struct {
			Id   int    `json:"id"`
			Role string `json:"role"`
		}{}
		err = json.NewDecoder(accResp.Body).Decode(&accRespData)
		if err != nil {
			return nil, -1, errors.New("failed to decode response data")
		}
		userDataResp, err := http.Get(fmt.Sprintf("http://db-service:8081/get_data?id=%d", accRespData.Id))
		if err != nil || userDataResp.StatusCode != 200 {
			return nil, -1, errors.New("failed to get user data")
		}
		err = json.NewDecoder(userDataResp.Body).Decode(&userData)
		if err != nil {
			return nil, -1, err
		}
		userData.NewTokens = nil
		userData.Id = accRespData.Id
		userData.Role = accRespData.Role
		return &userData, accRespData.Id, nil
	}
	refResp, err := http.Post("http://auth-service:8080/verify", "application.json", bytes.NewBuffer(refJson))
	if err != nil {
		return nil, -1, err
	}
	if refResp.StatusCode == 200 {
		defer accResp.Body.Close()
		refRespData := struct {
			Id   int    `json:"id"`
			Role string `json:"role"`
		}{}
		err = json.NewDecoder(refResp.Body).Decode(&refRespData)
		if err != nil {
			return nil, -1, err
		}
		verifyTokenResp, err := http.Get(fmt.Sprintf("http://db-service:8081/verify_refresh?id=%d&refresh=%s", refRespData.Id, refresh))
		if err != nil || verifyTokenResp.StatusCode != 200 {
			return nil, -1, errors.New("previous refresh token provided")
		}
		tokenResp, err := http.Get(fmt.Sprintf("http://auth-service:8080/new?id=%d&role=%s", refRespData.Id, refRespData.Role))
		if err != nil || tokenResp.StatusCode != 200 {
			return nil, -1, errors.New("failed to get Tokens")
		}
		newTokens := model.Tokens{}
		err = json.NewDecoder(tokenResp.Body).Decode(&newTokens)
		if err != nil {
			return nil, -1, err
		}
		userDataResp, err := http.Get(fmt.Sprintf("http://db-service:8081/get_data?id=%d", refRespData.Id))
		if err != nil || userDataResp.StatusCode != 200 {
			return nil, -1, errors.New("failed to get user data")
		}
		err = json.NewDecoder(userDataResp.Body).Decode(&userData)
		if err != nil {
			return nil, -1, err
		}
		setTokenResp, err := http.Post(fmt.Sprintf("http://db-service:8081/set_refresh?id=%d&refresh=%s", refRespData.Id, newTokens.Refresh), "application/json", nil)
		if err != nil || setTokenResp.StatusCode != 200 {
			return nil, -1, errors.New("failed to update refresh")
		}
		userData.NewTokens = &newTokens
		userData.Id = refRespData.Id
		userData.Role = refRespData.Role
		return &userData, refRespData.Id, nil
	}
	return nil, -1, errors.New("not authorized")
}
