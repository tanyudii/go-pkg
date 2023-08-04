package actor

import (
	"context"
	gotex "pkg.tanyudii.me/go-pkg/go-tex"
)

type EntityActor interface {
	SetCreator(ctx context.Context) error
	SetEditor(ctx context.Context) error
	SetDestroyer(ctx context.Context) error
}

type EntityCEditor interface {
	SetCreator(ctx context.Context) error
	SetEditor(ctx context.Context) error
}

type EntityCreator interface {
	SetCreator(ctx context.Context) error
}

type EntityHaveActor struct {
	CreatedBy *string
	UpdatedBy *string
	DeletedBy *string
}

func (e *EntityHaveActor) SetCreator(ctx context.Context) error {
	uid, err := gotex.GetUserID(ctx)
	if err != nil {
		return nil
	}
	e.CreatedBy = &uid
	e.UpdatedBy = &uid
	return nil
}

func (e *EntityHaveActor) SetEditor(ctx context.Context) error {
	uid, err := gotex.GetUserID(ctx)
	if err != nil {
		return nil
	}
	e.UpdatedBy = &uid
	return nil
}

func (e *EntityHaveActor) SetDestroyer(ctx context.Context) error {
	uid, err := gotex.GetUserID(ctx)
	if err != nil {
		return nil
	}
	e.DeletedBy = &uid
	return nil
}

type EntityHaveCEditor struct {
	CreatedBy *string
	UpdatedBy *string
}

func (e *EntityHaveCEditor) SetCreator(ctx context.Context) error {
	uid, err := gotex.GetUserID(ctx)
	if err != nil {
		return nil
	}
	e.CreatedBy = &uid
	e.UpdatedBy = &uid
	return nil
}

func (e *EntityHaveCEditor) SetEditor(ctx context.Context) error {
	uid, err := gotex.GetUserID(ctx)
	if err != nil {
		return nil
	}
	e.UpdatedBy = &uid
	return nil
}

type EntityHaveCreator struct {
	CreatedBy *string
}

func (e *EntityHaveCreator) SetCreator(ctx context.Context) error {
	uid, err := gotex.GetUserID(ctx)
	if err != nil {
		return nil
	}
	e.CreatedBy = &uid
	return nil
}
