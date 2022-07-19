package task

import (
	"context"

	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/internal/ports"
	"gitlab.com/g6834/team26/task/pkg/api"
)

type Service struct {
	db             ports.TaskDB
	grpcAuth       ports.GrpcAuth
	analyticSender ports.TaskAnalyticSender
}

func New(db ports.TaskDB, grpcAuth ports.GrpcAuth, analyticSender ports.TaskAnalyticSender) *Service {
	return &Service{
		db:             db,
		grpcAuth:       grpcAuth,
		analyticSender: analyticSender,
	}
}

func (s *Service) ListTasks(ctx context.Context, login string) ([]*models.Task, error) {
	t, err := s.db.List(ctx, login)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) RunTask(ctx context.Context, createdTask *models.Task) error { // TODO: отправлять письмо первому согласующему
	err := s.db.Run(ctx, createdTask)
	if err != nil {
		return err
	}

	err = s.analyticSender.ActionTask(ctx, createdTask.UUID, "run", createdTask.InitiatorLogin) // TODO: отправлять сообщение в отдельной горутине
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateTask(ctx context.Context, id, login, name, text string) error {
	err := s.db.Update(ctx, id, login, name, text)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteTask(ctx context.Context, login, id string) error { // TODO: отправлять письма всем участникам об отмене операции
	err := s.db.Delete(ctx, login, id)
	if err != nil {
		return err
	}

	err = s.analyticSender.ActionTask(ctx, id, "delete", "true") // TODO: отправлять сообщение в отдельной горутине
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ApproveTask(ctx context.Context, login, id, approvalLogin string) error { // TODO: отправлять письмо следующему согласующему
	err := s.db.Approve(ctx, login, id, approvalLogin)
	if err != nil {
		return err
	}

	err = s.analyticSender.ActionTask(ctx, id, "approve", "true") // TODO: отправлять сообщение в отдельной горутине
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeclineTask(ctx context.Context, login, id, approvalLogin string) error { // TODO: отправлять письма всем участникам об остановке согласования операции в связи с отклонением одним из участников
	err := s.db.Decline(ctx, login, id, approvalLogin)
	if err != nil {
		return err
	}

	err = s.analyticSender.ActionTask(ctx, id, "approve", "false") // TODO: отправлять сообщение в отдельной горутине
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Validate(ctx context.Context, tokens ports.TokenPair) (*api.AuthResponse, error) {
	grpcResponse, err := s.grpcAuth.Validate(ctx, tokens)
	if err != nil {
		return nil, err
	}
	return grpcResponse, nil
}