package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"
	"server/config"
	student_psql "server/models/student_psql"
	"time"
)

// Resolver struct
type Resolver struct{}

// Query resolver for getting all leaderboard records
func (r *queryResolver) GetLeaderboard(ctx context.Context) ([]*student_psql.StudentLeaderboardRecordTable, error) {
	var leaderboard []*student_psql.StudentLeaderboardRecordTable
	if err := config.GetPostgresDBConnection().Find(&leaderboard).Error; err != nil {
		return nil, err
	}
	return leaderboard, nil
}

// Query resolver for getting a leaderboard record by ID
func (r *queryResolver) GetLeaderboardByID(ctx context.Context, id string) (*student_psql.StudentLeaderboardRecordTable, error) {
	var leaderboard student_psql.StudentLeaderboardRecordTable
	if err := config.GetPostgresDBConnection().Where("leaderboard_record_id = ?", id).First(&leaderboard).Error; err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

// Mutation resolver for adding a leaderboard record
func (r *mutationResolver) AddLeaderboard(ctx context.Context, rank int, score float64, domain string, subDomain string, timePeriod string) (*student_psql.StudentLeaderboardRecordTable, error) {
	leaderboard := student_psql.StudentLeaderboardRecordTable{
		Rank:        rank,
		Score:       score,
		Domain:      domain,
		SubDomain:   subDomain,
		TimePeriod:  timePeriod,
		LastUpdated: time.Now(),
	}
	if err := config.GetPostgresDBConnection().Create(&leaderboard).Error; err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

// Mutation resolver for updating a leaderboard record
func (r *mutationResolver) UpdateLeaderboard(ctx context.Context, id string, rank *int, score *float64, domain *string, subDomain *string, timePeriod *string) (*student_psql.StudentLeaderboardRecordTable, error) {
	var leaderboard student_psql.StudentLeaderboardRecordTable
	if err := config.GetPostgresDBConnection().Where("leaderboard_record_id = ?", id).First(&leaderboard).Error; err != nil {
		return nil, err
	}
	if rank != nil {
		leaderboard.Rank = *rank
	}
	if score != nil {
		leaderboard.Score = *score
	}
	if domain != nil {
		leaderboard.Domain = *domain
	}
	if subDomain != nil {
		leaderboard.SubDomain = *subDomain
	}
	if timePeriod != nil {
		leaderboard.TimePeriod = *timePeriod
	}
	leaderboard.LastUpdated = time.Now()
	if err := config.GetPostgresDBConnection().Save(&leaderboard).Error; err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

// Mutation resolver for deleting a leaderboard record
func (r *mutationResolver) DeleteLeaderboard(ctx context.Context, id string) (bool, error) {
	if err := config.GetPostgresDBConnection().Where("leaderboard_record_id = ?", id).Delete(&student_psql.StudentLeaderboardRecordTable{}).Error; err != nil {
		return false, err
	}
	return true, nil
}
