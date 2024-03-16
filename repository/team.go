package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
)

type TeamRepository interface {
	Get(ctx context.Context, id string) (*model.Team, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object *model.Team) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	List(ctx context.Context, opts *model.TeamListOption) ([]*model.Team, error)
	Set(ctx context.Context, id string, object *model.Team) error
	AddMember(ctx context.Context, teamID string, user *model.User) error
	DeleteMember(ctx context.Context, teamID string, user *model.User) error
}

const teamsCollectionPath = "teams"

func NewTeamRepository(opts *NewRepositoryOption) TeamRepository {
	if opts.dbClient != nil {
		return &gormTeamRepository{
			gormGenericRepository: *newGormGenericRepository[*model.Team](opts.dbClient),
		}
	}
	return &firestoreTeamRepository{
		firestoreGenericRepository: *newFirestoreGenericRepository[*model.Team](opts.firestoreClient, teamsCollectionPath),
	}
}

type firestoreTeamRepository struct {
	firestoreGenericRepository[*model.Team]
}

func (r *firestoreTeamRepository) List(ctx context.Context, opts *model.TeamListOption) ([]*model.Team, error) {
	collectionPath := r.collectionPath
	if collectionPath == "" {
		return nil, errors.New("collection path is empty")
	}
	query := r.client.Collection(collectionPath).Query
	// NOTE: implement where clause
	if opts.UserID != "" {
		query = query.Where("user_id", "==", opts.UserID)
	}
	if opts.OrganizationID != "" {
		query = query.Where("organization_id", "==", opts.OrganizationID)
	}
	return r.firestoreGenericRepository.list(ctx, query)
}

func (r *firestoreTeamRepository) AddMember(ctx context.Context, teamID string, user *model.User) error {
	path := fmt.Sprintf("%s/%s/members", teamsCollectionPath, teamID)
	if _, err := r.client.Collection(path).Doc(user.GetID()).Set(ctx, user); err != nil {
		return errors.Wrap(err, "failed to add member")
	}
	return nil
}
func (r *firestoreTeamRepository) DeleteMember(ctx context.Context, teamID string, user *model.User) error {
	path := fmt.Sprintf("%s/%s/members", teamsCollectionPath, teamID)
	if _, err := r.client.Collection(path).Doc(user.GetID()).Delete(ctx); err != nil {
		return errors.Wrap(err, "failed to delete member")
	}
	return nil
}

func (r *firestoreTeamRepository) Get(ctx context.Context, id string) (*model.Team, error) {
	entity, err := r.firestoreGenericRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	// list members
	path := fmt.Sprintf("%s/%s/members", teamsCollectionPath, id)
	iter := r.client.Collection(path).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var user model.User
		if err := doc.DataTo(&user); err != nil {
			return nil, errors.Wrap(err, "failed to assign")
		}
		user.SetID(doc.Ref.ID)
		entity.Members = append(entity.Members, &user)
	}

	return entity, nil
}

type gormTeamRepository struct {
	gormGenericRepository[*model.Team]
}

func (r *gormTeamRepository) List(ctx context.Context, opts *model.TeamListOption) ([]*model.Team, error) {
	query := r.client
	// NOTE: implement where clause
	if opts.UserID != "" {
		query = query.Where("user_id = ?", opts.UserID)
	}
	if opts.OrganizationID != "" {
		query = query.Where("organization_id = ?", opts.OrganizationID)
	}
	return r.gormGenericRepository.list(query)
}

func (r *gormTeamRepository) AddMember(ctx context.Context, teamID string, user *model.User) error {
	team, err := r.Get(ctx, teamID)
	if err != nil {
		return err
	}
	team.Members = append(team.Members, user)
	if err := r.client.Save(&team).Error; err != nil {
		return err
	}
	return nil
}

func (r *gormTeamRepository) DeleteMember(ctx context.Context, teamID string, user *model.User) error {
	team, err := r.Get(ctx, teamID)
	if err != nil {
		return err
	}
	if err := r.client.Model(&team).Association("Members").Delete(&user); err != nil {
		return err
	}
	return nil
}

func (r *gormTeamRepository) Get(ctx context.Context, id string) (*model.Team, error) {
	var team model.Team

	uid, err := strconv.Atoi(id)
	if err != nil {
		return &team, errors.Wrap(err, "failed to convert id to int")
	}
	r.client.Preload("Members").Find(&team, uid)
	team.SetID(team.GetID())

	for _, member := range team.Members {
		member.SetID(member.GetID())
	}

	return &team, nil
}
