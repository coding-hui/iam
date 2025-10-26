package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
)

// --- fake repositories & factory ---
type fakeFactory struct {
	tplRepo repository.EmailTemplateRepository
	catRepo repository.EmailTemplateCategoryRepository
}

// implement repository.Factory
func (f *fakeFactory) Begin(ctx context.Context) (repository.Transaction, error)         { return f, nil }
func (f *fakeFactory) Commit(ctx context.Context) error                                  { return nil }
func (f *fakeFactory) Rollback(ctx context.Context) error                                { return nil }
func (f *fakeFactory) UserRepository() repository.UserRepository                         { return nil }
func (f *fakeFactory) CasbinRepository() repository.CasbinRepository                     { return nil }
func (f *fakeFactory) ResourceRepository() repository.ResourceRepository                 { return nil }
func (f *fakeFactory) RoleRepository() repository.RoleRepository                         { return nil }
func (f *fakeFactory) PolicyRepository() repository.PolicyRepository                     { return nil }
func (f *fakeFactory) OrganizationRepository() repository.OrganizationRepository         { return nil }
func (f *fakeFactory) IdentityProviderRepository() repository.IdentityProviderRepository { return nil }
func (f *fakeFactory) ApplicationRepository() repository.ApplicationRepository           { return nil }
func (f *fakeFactory) ApiKeyRepository() repository.ApiKeyRepository                     { return nil }
func (f *fakeFactory) DeviceAuthRepository() repository.DeviceAuthRepository             { return nil }
func (f *fakeFactory) EmailTemplateRepository() repository.EmailTemplateRepository       { return f.tplRepo }
func (f *fakeFactory) EmailTemplateCategoryRepository() repository.EmailTemplateCategoryRepository {
	return f.catRepo
}
func (f *fakeFactory) Close() error { return nil }
func (f *fakeFactory) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// email template repo fake
type fakeTplRepo struct {
	store     map[string]*model.EmailTemplate
	idCounter int
}

func newFakeTplRepo() *fakeTplRepo {
	return &fakeTplRepo{store: make(map[string]*model.EmailTemplate), idCounter: 0}
}
func (r *fakeTplRepo) Create(ctx context.Context, tpl *model.EmailTemplate, _ metav1.CreateOptions) error {
	if tpl.InstanceID == "" {
		r.idCounter++
		tpl.InstanceID = fmt.Sprintf("etpl-test-%d", r.idCounter)
	}
	r.store[tpl.InstanceID] = tpl
	return nil
}
func (r *fakeTplRepo) Update(ctx context.Context, tpl *model.EmailTemplate, _ metav1.UpdateOptions) error {
	r.store[tpl.InstanceID] = tpl
	return nil
}
func (r *fakeTplRepo) DeleteByInstanceId(ctx context.Context, instanceId string, _ metav1.DeleteOptions) error {
	delete(r.store, instanceId)
	return nil
}
func (r *fakeTplRepo) GetByInstanceId(ctx context.Context, instanceId string, _ metav1.GetOptions) (*model.EmailTemplate, error) {
	return r.store[instanceId], nil
}
func (r *fakeTplRepo) GetByName(ctx context.Context, name string, _ metav1.GetOptions) (*model.EmailTemplate, error) {
	for _, v := range r.store {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, errors.WithCode(code.ErrRecordNotExist, "template %s not found", name)
}
func (r *fakeTplRepo) List(ctx context.Context, _ metav1.ListOptions) ([]model.EmailTemplate, error) {
	out := make([]model.EmailTemplate, 0, len(r.store))
	for _, v := range r.store {
		out = append(out, *v)
	}
	return out, nil
}
func (r *fakeTplRepo) Count(ctx context.Context, _ metav1.ListOptions) (int64, error) {
	return int64(len(r.store)), nil
}

// versions & variables stubs
func (r *fakeTplRepo) CreateVersion(context.Context, *model.EmailTemplateVersion, metav1.CreateOptions) error {
	return nil
}
func (r *fakeTplRepo) ListVersions(context.Context, string, metav1.ListOptions) ([]model.EmailTemplateVersion, error) {
	return nil, nil
}
func (r *fakeTplRepo) GetVersionByNumber(context.Context, string, int, metav1.GetOptions) (*model.EmailTemplateVersion, error) {
	return nil, nil
}
func (r *fakeTplRepo) DeleteVersion(context.Context, string, int, metav1.DeleteOptions) error {
	return nil
}
func (r *fakeTplRepo) AddVariables(context.Context, []*model.EmailTemplateVariable, metav1.CreateOptions) error {
	return nil
}
func (r *fakeTplRepo) ListVariables(context.Context, string, metav1.ListOptions) ([]model.EmailTemplateVariable, error) {
	return nil, nil
}
func (r *fakeTplRepo) DeleteVariable(context.Context, string, string, metav1.DeleteOptions) error {
	return nil
}
func (r *fakeTplRepo) DeleteVariablesByTemplate(context.Context, string, metav1.DeleteOptions) error {
	return nil
}

// category repo fake
type fakeCatRepo struct {
	store     map[string]*model.EmailTemplateCategory
	idCounter int
}

func newFakeCatRepo() *fakeCatRepo {
	return &fakeCatRepo{store: make(map[string]*model.EmailTemplateCategory), idCounter: 0}
}
func (r *fakeCatRepo) Create(ctx context.Context, cat *model.EmailTemplateCategory, _ metav1.CreateOptions) error {
	if cat.InstanceID == "" {
		r.idCounter++
		cat.InstanceID = fmt.Sprintf("ecat-test-%d", r.idCounter)
	}
	r.store[cat.InstanceID] = cat
	return nil
}
func (r *fakeCatRepo) Update(ctx context.Context, cat *model.EmailTemplateCategory, _ metav1.UpdateOptions) error {
	r.store[cat.InstanceID] = cat
	return nil
}
func (r *fakeCatRepo) DeleteByInstanceId(ctx context.Context, instanceId string, _ metav1.DeleteOptions) error {
	delete(r.store, instanceId)
	return nil
}
func (r *fakeCatRepo) GetByInstanceId(ctx context.Context, instanceId string, _ metav1.GetOptions) (*model.EmailTemplateCategory, error) {
	return r.store[instanceId], nil
}
func (r *fakeCatRepo) GetByName(ctx context.Context, name string, _ metav1.GetOptions) (*model.EmailTemplateCategory, error) {
	for _, v := range r.store {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, errors.WithCode(code.ErrRecordNotExist, "category %s not found", name)
}
func (r *fakeCatRepo) List(ctx context.Context, _ metav1.ListOptions) ([]model.EmailTemplateCategory, error) {
	out := make([]model.EmailTemplateCategory, 0, len(r.store))
	for _, v := range r.store {
		out = append(out, *v)
	}
	return out, nil
}
func (r *fakeCatRepo) Count(ctx context.Context, _ metav1.ListOptions) (int64, error) {
	return int64(len(r.store)), nil
}
func (r *fakeCatRepo) ListChildren(context.Context, string, metav1.ListOptions) ([]model.EmailTemplateCategory, error) {
	return nil, nil
}
func (r *fakeCatRepo) CountChildren(context.Context, string, metav1.ListOptions) (int64, error) {
	return 0, nil
}

// --- tests ---
func TestEmailTemplateService_CreateAndList(t *testing.T) {
	svc := &emailTemplateServiceImpl{Store: &fakeFactory{tplRepo: newFakeTplRepo(), catRepo: newFakeCatRepo()}}

	tpl, err := svc.CreateTemplate(context.Background(), &v1.CreateEmailTemplateRequest{
		Name:    "welcome",
		Subject: "Welcome",
		Content: "<html>...</html>",
		Status:  "draft",
	})
	if err != nil {
		t.Fatalf("create template failed: %v", err)
	}
	if tpl.Subject != "Welcome" || tpl.Status != "draft" {
		t.Fatalf("unexpected template: %+v", tpl)
	}

	list, err := svc.ListTemplates(context.Background(), v1.ListEmailTemplateOptions{})
	if err != nil {
		t.Fatalf("list templates failed: %v", err)
	}
	if list.ListMeta.TotalCount != 1 {
		t.Fatalf("expected total 1, got %d", list.ListMeta.TotalCount)
	}
}

func TestEmailTemplateService_CategoryCRUD(t *testing.T) {
	svc := &emailTemplateServiceImpl{Store: &fakeFactory{tplRepo: newFakeTplRepo(), catRepo: newFakeCatRepo()}}

	cat, err := svc.CreateCategory(context.Background(), &v1.CreateEmailTemplateCategoryRequest{
		Name: "system",
	})
	if err != nil {
		t.Fatalf("create category failed: %v", err)
	}
	if cat.ObjectMeta.Name != "system" {
		t.Fatalf("unexpected category: %+v", cat)
	}

	fetched, err := svc.GetCategory(context.Background(), "ecat-test-1")
	if err != nil {
		t.Fatalf("get category failed: %v", err)
	}
	if fetched.ObjectMeta.InstanceID == "" {
		t.Fatalf("expected instance id set")
	}

	updated, err := svc.UpdateCategory(context.Background(), "ecat-test-1", &v1.UpdateEmailTemplateCategoryRequest{Description: "desc"})
	if err != nil {
		t.Fatalf("update category failed: %v", err)
	}
	if updated.Description != "desc" {
		t.Fatalf("update not applied")
	}
}

func TestEmailTemplateService_Init(t *testing.T) {
	svc := &emailTemplateServiceImpl{Store: &fakeFactory{tplRepo: newFakeTplRepo(), catRepo: newFakeCatRepo()}}

	// Test init method
	err := svc.Init(context.Background())
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Verify both categories were created by listing all categories
	categories, err := svc.ListCategories(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("list categories failed: %v", err)
	}
	if categories.ListMeta.TotalCount != 2 {
		t.Fatalf("expected 2 categories, got %d", categories.ListMeta.TotalCount)
	}

	// Verify System and Custom categories exist with correct types
	var systemFound, customFound bool
	for _, cat := range categories.Items {
		if cat.ObjectMeta.Name == "System" && cat.Type == "system" {
			systemFound = true
		}
		if cat.ObjectMeta.Name == "Custom" && cat.Type == "custom" {
			customFound = true
		}
	}
	if !systemFound {
		t.Fatalf("System category with system type not found")
	}
	if !customFound {
		t.Fatalf("Custom category with custom type not found")
	}

	// Verify default templates were created
	templates, err := svc.ListTemplates(context.Background(), v1.ListEmailTemplateOptions{})
	if err != nil {
		t.Fatalf("list templates failed: %v", err)
	}
	if templates.ListMeta.TotalCount != 2 {
		t.Fatalf("expected 2 default templates, got %d", templates.ListMeta.TotalCount)
	}

	// Test running init again (should not create duplicates)
	err = svc.Init(context.Background())
	if err != nil {
		t.Fatalf("second init failed: %v", err)
	}

	// Verify no duplicates were created
	categories, err = svc.ListCategories(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("list categories after second init failed: %v", err)
	}
	if categories.ListMeta.TotalCount != 2 {
		t.Fatalf("expected 2 categories after second init, got %d", categories.ListMeta.TotalCount)
	}

	templates, err = svc.ListTemplates(context.Background(), v1.ListEmailTemplateOptions{})
	if err != nil {
		t.Fatalf("list templates after second init failed: %v", err)
	}
	if templates.ListMeta.TotalCount != 2 {
		t.Fatalf("expected 2 templates after second init, got %d", templates.ListMeta.TotalCount)
	}
}

func TestEmailTemplateService_DeleteProtection(t *testing.T) {
	svc := &emailTemplateServiceImpl{Store: &fakeFactory{tplRepo: newFakeTplRepo(), catRepo: newFakeCatRepo()}}

	// Initialize to create system templates and categories
	err := svc.Init(context.Background())
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Get the system category
	categories, err := svc.ListCategories(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("list categories failed: %v", err)
	}

	var systemCategoryID, customCategoryID string
	for _, cat := range categories.Items {
		if cat.Type == "system" {
			systemCategoryID = cat.ObjectMeta.InstanceID
		} else if cat.Type == "custom" {
			customCategoryID = cat.ObjectMeta.InstanceID
		}
	}

	// Test: Try to delete system category (should fail)
	err = svc.DeleteCategory(context.Background(), systemCategoryID)
	if err == nil {
		t.Fatalf("expected error when deleting system category, but got none")
	}
	if !errors.IsCode(err, code.ErrCannotDeleteSystemEmailTemplateCategory) {
		t.Fatalf("expected ErrCannotDeleteSystemEmailTemplateCategory, got %v", err)
	}

	// Test: Delete custom category (should succeed)
	err = svc.DeleteCategory(context.Background(), customCategoryID)
	if err != nil {
		t.Fatalf("failed to delete custom category: %v", err)
	}

	// Get system templates
	templates, err := svc.ListTemplates(context.Background(), v1.ListEmailTemplateOptions{})
	if err != nil {
		t.Fatalf("list templates failed: %v", err)
	}

	// Test: Try to delete system templates (should fail)
	for _, tpl := range templates.Items {
		err = svc.DeleteTemplate(context.Background(), tpl.ObjectMeta.InstanceID)
		if err == nil {
			t.Fatalf("expected error when deleting system template %s, but got none", tpl.ObjectMeta.Name)
		}
		// Check for either system template or default template error codes
		if !errors.IsCode(err, code.ErrCannotDeleteSystemEmailTemplate) && !errors.IsCode(err, code.ErrCannotDeleteDefaultEmailTemplate) {
			t.Fatalf("expected ErrCannotDeleteSystemEmailTemplate or ErrCannotDeleteDefaultEmailTemplate for template %s, got %v", tpl.ObjectMeta.Name, err)
		}
	}

	// Test: Create and delete custom template (should succeed)
	customTpl, err := svc.CreateTemplate(context.Background(), &v1.CreateEmailTemplateRequest{
		Name:    "custom-test",
		Subject: "Custom Test",
		Content: "<html>Custom content</html>",
		Status:  "draft",
	})
	if err != nil {
		t.Fatalf("failed to create custom template: %v", err)
	}

	err = svc.DeleteTemplate(context.Background(), customTpl.ObjectMeta.InstanceID)
	if err != nil {
		t.Fatalf("failed to delete custom template: %v", err)
	}
}

func TestEmailTemplateService_ErrorCodes(t *testing.T) {
	svc := &emailTemplateServiceImpl{Store: &fakeFactory{tplRepo: newFakeTplRepo(), catRepo: newFakeCatRepo()}}

	// Initialize to create system templates and categories
	err := svc.Init(context.Background())
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Get the system category and template
	categories, err := svc.ListCategories(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("list categories failed: %v", err)
	}

	var systemCategoryID string
	for _, cat := range categories.Items {
		if cat.Type == "system" {
			systemCategoryID = cat.ObjectMeta.InstanceID
			break
		}
	}

	templates, err := svc.ListTemplates(context.Background(), v1.ListEmailTemplateOptions{})
	if err != nil {
		t.Fatalf("list templates failed: %v", err)
	}

	// Test specific error codes and messages
	// Test system category deletion
	err = svc.DeleteCategory(context.Background(), systemCategoryID)
	if err == nil {
		t.Fatalf("expected error when deleting system category")
	}
	if !errors.IsCode(err, code.ErrCannotDeleteSystemEmailTemplateCategory) {
		t.Fatalf("expected ErrCannotDeleteSystemEmailTemplateCategory, got error code: %d", errors.ParseCoder(err).Code())
	}
	if err.Error() != "Cannot delete system email template category" {
		t.Fatalf("expected English error message 'Cannot delete system email template category', got: %s", err.Error())
	}

	// Test system template deletion
	if len(templates.Items) > 0 {
		err = svc.DeleteTemplate(context.Background(), templates.Items[0].ObjectMeta.InstanceID)
		if err == nil {
			t.Fatalf("expected error when deleting system template")
		}

		// Should be either system template or default template error
		coder := errors.ParseCoder(err)
		if coder.Code() != code.ErrCannotDeleteSystemEmailTemplate && coder.Code() != code.ErrCannotDeleteDefaultEmailTemplate {
			t.Fatalf("expected system or default template error code, got: %d", coder.Code())
		}

		// Check English error message
		expectedMessages := []string{
			"Cannot delete system email template",
			"Cannot delete default email template",
		}
		found := false
		for _, msg := range expectedMessages {
			if err.Error() == msg {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected English error message, got: %s", err.Error())
		}
	}
}
