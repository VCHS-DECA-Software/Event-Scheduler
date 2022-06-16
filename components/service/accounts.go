package handlers

import (
	"context"
	"errors"
	"main/components/db"
	"main/components/proto"
	"main/components/users"
)

type AccountManager struct {
	proto.UnimplementedAccountManagerServer
}

func createGeneric[T users.Type](
	details *proto.AccountDetails,
) (*proto.WithStatus, error) {
	user, err := users.NewAccount(users.Account[T]{
		Username: details.Username,
		Password: details.Password,
	})
	if err != nil {
		return nil, err
	}

	token, err := users.IssueToken[T](user.ID)
	if err != nil {
		return nil, err
	}

	return &proto.WithStatus{
		Credentials: &proto.Credentials{
			Type:  details.Type,
			Id:    user.ID,
			Token: token,
		},
		Status: &proto.Status{
			Code:    proto.StatusCode_OK,
			Message: "account successfully created",
		},
	}, nil
}

/* code duplication, the practice of using generics to differentiate
variations of structs should be handled with a runtime value instead.
may implement this later */
func (a *AccountManager) Create(
	c context.Context, details *proto.AccountDetails,
) (*proto.WithStatus, error) {
	switch details.Type {
	case proto.AccountType_STUDENT:
		return createGeneric[users.Student](details)
	case proto.AccountType_JUDGE:
		return createGeneric[users.Judge](details)
	case proto.AccountType_ADMIN:
		return createGeneric[users.Admin](details)
	}
	return &proto.WithStatus{
		Status: &proto.Status{
			Code:    proto.StatusCode_INVALID_TYPE,
			Message: "given user type is invalid",
		},
	}, nil
}

func authGeneric[T users.Type](
	r *proto.AuthRequest,
) (*proto.WithStatus, error) {
	token, err := users.Authenticate[T](
		r.Details.Username,
		r.Details.Password,
	)
	if err != nil {
		return nil, err
	}
	return &proto.WithStatus{
		Credentials: &proto.Credentials{
			Type:  r.Details.Type,
			Id:    r.Id,
			Token: token,
		},
		Status: &proto.Status{
			Code:    proto.StatusCode_OK,
			Message: "authentication successful",
		},
	}, nil
}

func (a *AccountManager) Authenticate(
	c context.Context, r *proto.AuthRequest,
) (*proto.WithStatus, error) {
	switch r.Details.Type {
	case proto.AccountType_STUDENT:
		return authGeneric[users.Student](r)
	case proto.AccountType_JUDGE:
		return authGeneric[users.Judge](r)
	case proto.AccountType_ADMIN:
		return authGeneric[users.Admin](r)
	}
	return &proto.WithStatus{
		Status: &proto.Status{
			Code:    proto.StatusCode_INVALID_TYPE,
			Message: "given user type is invalid",
		},
	}, nil
}

func updateGeneric[T users.Type](
	r *proto.UpdateRequest,
) (*proto.Status, error) {
	valid, err := users.Validate[T](
		r.Credentials.Id,
		r.Credentials.Token,
	)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid token")
	}

	err = db.Update(&users.Account[T]{
		ID:       r.Credentials.Id,
		Username: r.Details.Username,
		Password: r.Details.Password,
	})
	if err != nil {
		return nil, err
	}

	return &proto.Status{
		Code:    proto.StatusCode_OK,
		Message: "update successful",
	}, nil
}

func (a *AccountManager) Update(
	c context.Context, r *proto.UpdateRequest,
) (*proto.Status, error) {
	switch r.Details.Type {
	case proto.AccountType_STUDENT:
		return updateGeneric[users.Student](r)
	case proto.AccountType_JUDGE:
		return updateGeneric[users.Judge](r)
	case proto.AccountType_ADMIN:
		return updateGeneric[users.Admin](r)
	}
	return &proto.Status{
		Code:    proto.StatusCode_INVALID_TYPE,
		Message: "given user type is invalid",
	}, nil
}

func deleteGeneric[T users.Type](
	r *proto.Credentials,
) (*proto.Status, error) {
	err := db.Delete(&users.Account[T]{
		ID: r.Id,
	})
	if err != nil {
		return nil, err
	}

	return &proto.Status{
		Code:    proto.StatusCode_OK,
		Message: "delete successful",
	}, nil
}

func (a *AccountManager) Delete(
	c context.Context, r *proto.Credentials,
) (*proto.Status, error) {
	switch r.Type {
	case proto.AccountType_STUDENT:
		return deleteGeneric[users.Student](r)
	case proto.AccountType_JUDGE:
		return deleteGeneric[users.Judge](r)
	case proto.AccountType_ADMIN:
		return deleteGeneric[users.Admin](r)
	}
	return &proto.Status{
		Code:    proto.StatusCode_INVALID_TYPE,
		Message: "given user type is invalid",
	}, nil
}
