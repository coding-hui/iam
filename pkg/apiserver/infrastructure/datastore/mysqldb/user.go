package mysqldb

import (
	"context"

	"gorm.io/gorm"

	"github.com/coding-hui/common/errors"
	"github.com/wecoding/iam/pkg/apiserver/domain/repository"
	"github.com/wecoding/iam/pkg/apiserver/infrastructure/datastore"
	"github.com/wecoding/iam/pkg/code"
	"github.com/wecoding/iam/pkg/utils/gormutil"

	iamv1alpha1 "github.com/coding-hui/api/iam/v1alpha1"
	"github.com/coding-hui/common/fields"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// newUserRepository new User Repository
func newUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db}
}

// Create creates a new user account.
func (u *userRepositoryImpl) Create(ctx context.Context, user *iamv1alpha1.User, opts metav1alpha1.CreateOptions) error {
	if _, err := u.Get(ctx, user.Name, metav1alpha1.GetOptions{}); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return datastore.ErrRecordExist
	}

	return u.db.WithContext(ctx).Create(&user).Error
}

// Update updates an user account information.
func (u *userRepositoryImpl) Update(ctx context.Context, user *iamv1alpha1.User, opts metav1alpha1.UpdateOptions) error {
	if err := u.db.WithContext(ctx).Save(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}

		return err
	}

	return nil
}

// Delete deletes the user by the user identifier.
func (u *userRepositoryImpl) Delete(ctx context.Context, username string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	err := u.db.WithContext(ctx).Where("name = ?", username).Delete(&iamv1alpha1.User{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}

		return err
	}

	return nil
}

// DeleteCollection batch deletes the users.
func (u *userRepositoryImpl) DeleteCollection(ctx context.Context, usernames []string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.WithContext(ctx).Where("name in (?)", usernames).Delete(&iamv1alpha1.User{}).Error
}

// Get get user
func (u *userRepositoryImpl) Get(ctx context.Context, username string, opts metav1alpha1.GetOptions) (*iamv1alpha1.User, error) {
	user := &iamv1alpha1.User{}
	if username == "" {
		return nil, datastore.ErrPrimaryEmpty
	}
	err := u.db.WithContext(ctx).Where("name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, err
	}

	return user, nil
}

// List list users
func (u *userRepositoryImpl) List(ctx context.Context, opts metav1alpha1.ListOptions) (*iamv1alpha1.UserList, error) {
	list := &iamv1alpha1.UserList{}

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	db := u.db.WithContext(ctx)
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("name")
	if username != "" {
		db.Where("name like ?", "%"+username+"%")
	}
	db.Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&list.Items).
		Offset(-1).
		Limit(-1).
		Count(&list.TotalCount)

	return list, db.Error
}
