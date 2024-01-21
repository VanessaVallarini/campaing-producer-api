package service

import (
	"campaing-producer-service/internal/model"
	"campaing-producer-service/internal/pkg/client"
	"campaing-producer-service/internal/pkg/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type ICampaignService interface {
	GetById(ctx context.Context, campaingId uuid.UUID) (model.Campaing, error)
	List(ctx context.Context, filters model.CampaingListingFilters) ([]model.Campaing, model.Paging, error)
}

type CampaignService struct {
	repository repository.ICampaingRepository
	cache      *cache.Cache
	awsClient  client.IAwsClient
	queueUrl   string
}

func NewCampaignService(repository repository.ICampaingRepository, awsClient client.IAwsClient, queueUrl string) *CampaignService {
	return &CampaignService{
		repository: repository,
		cache:      cache.New(24*time.Hour, 24*time.Hour),
		awsClient:  awsClient,
		queueUrl:   queueUrl,
	}
}

func (c *CampaignService) GetById(ctx context.Context, campaingId uuid.UUID) (model.Campaing, error) {
	if campaing, found := c.getInCache(campaingId); found {
		return campaing, nil
	}

	campaing, err := c.repository.GetById(ctx, campaingId)
	if err != nil {
		return model.Campaing{}, err
	}

	c.setInCache(campaingId, campaing)

	return campaing, nil
}

func (c *CampaignService) List(ctx context.Context, filters model.CampaingListingFilters) ([]model.Campaing, model.Paging, error) {
	var campaings []model.Campaing
	var campaingsIdsOutOfCache []uuid.UUID      //ids de campanhas não encontradas no cache
	var newFilters model.CampaingListingFilters //filstros contendo apenas os ids das campanhas não encontradas no cache
	var campaingsRepository []model.Campaing    //campanhas consultadas na base de dados
	var err error
	var paging model.Paging

	//consulta no cache
	for _, campaingId := range filters.Ids {
		campaing, found := c.getInCache(campaingId)
		if found {
			if filters.Active != nil && campaing.Active == *filters.Active {
				campaings = append(campaings, campaing)
			} else {
				//ids de campanhas não encontradas no cache
				campaingsIdsOutOfCache = append(campaingsIdsOutOfCache, campaingId)
			}
		}
	}

	//faz a busca na base apenas para os ids não encontradso no cache
	if len(campaingsIdsOutOfCache) > 0 {
		newFilters.Ids = campaingsIdsOutOfCache
		newFilters.Page = filters.Page
		newFilters.Size = filters.Size
		newFilters.Active = filters.Active
		campaingsRepository, paging, err = c.repository.List(ctx, newFilters)

		//unindo o que veio do cache com o que veio do repository
		campaings = append(campaings, campaingsRepository...)

		if err != nil {
			return nil, model.Paging{}, err
		}

	}
	return campaings, paging, nil
}

func (c *CampaignService) getInCache(campaingId uuid.UUID) (model.Campaing, bool) {
	cachedCampaing, found := c.cache.Get(campaingId.String())
	if found {
		return cachedCampaing.(model.Campaing), true
	}
	return model.Campaing{}, false
}

func (c *CampaignService) setInCache(campaingId uuid.UUID, campaing model.Campaing) {
	c.cache.SetDefault(campaingId.String(), campaing)
}

func (c *CampaignService) Create(ctx context.Context, req model.CampaingCreatingRequest) error {
	event := model.CampaingCreatingEvent{}
	event.Lat = req.Lat
	event.Long = req.Long
	event.Ip = req.Ip
	event.Action = model.EVENT_ACTION_CREATE

	if err := c.awsClient.SendMessage(ctx, event, &c.queueUrl); err != nil {
		return err
	}

	return nil
}
