package repository

import (
	"context"

	iamv1alpha1 "github.com/coding-hui/api/iam/v1alpha1"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

// UserRepository defines the user repository interface.
type UserRepository interface {
	Create(ctx context.Context, user *iamv1alpha1.User, opts metav1alpha1.CreateOptions) error
	Update(ctx context.Context, user *iamv1alpha1.User, opts metav1alpha1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1alpha1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1alpha1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1alpha1.GetOptions) (*iamv1alpha1.User, error)
	List(ctx context.Context, opts metav1alpha1.ListOptions) (*iamv1alpha1.UserList, error)
}
