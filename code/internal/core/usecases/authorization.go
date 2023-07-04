package usecases

import (
    "context"
    "github.com/SOAT1StackGoLang/tech-challenge/helpers"
    "github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
    "github.com/google/uuid"
    "go.uber.org/zap"
)

func validateIsAdmin(log *zap.SugaredLogger, uRepo ports.UsersUseCase, ctx context.Context, userID uuid.UUID) error {
    admin, err := uRepo.IsUserAdmin(ctx, userID)
    switch {
    case err != nil:
        log.Errorw(
            "failed checking user is admin",
            zap.String("userID", userID.String()),
            zap.Error(err),
        )
        return err
    case !admin:
        log.Errorw(
            "unauthorized user",
            zap.String("id", userID.String()),
            zap.Error(helpers.ErrUnauthorized),
        )
        err = helpers.ErrUnauthorized
        return err
    }
    return nil
}
