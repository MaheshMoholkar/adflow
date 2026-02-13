package service

import (
	"context"
	"encoding/json"

	"callflow/internal/domain/rule"
)

// RuleService provides rule business logic
type RuleService struct {
	ruleRepo rule.Repository
}

// NewRuleService creates a new rule service instance
func NewRuleService(ruleRepo rule.Repository) *RuleService {
	return &RuleService{ruleRepo: ruleRepo}
}

func (s *RuleService) Get(ctx context.Context, userID int64) (*rule.Rule, error) {
	return s.ruleRepo.Get(ctx, userID)
}

func (s *RuleService) Upsert(ctx context.Context, userID int64, data rule.RuleUpdate) (*rule.Rule, error) {
	return s.ruleRepo.Upsert(ctx, userID, data.Config)
}

func (s *RuleService) GetCompiledConfig(ctx context.Context, userID int64) (*rule.RuleConfig, error) {
	r, err := s.ruleRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	var config rule.RuleConfig
	if err := json.Unmarshal(r.Config, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
