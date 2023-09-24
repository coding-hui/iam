// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"strings"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/fields"
	metav1 "github.com/coding-hui/common/meta/v1"
)

const (
	DefaultMaxChildrenDepartments = 500
)

// OrganizationService Organization manage api.
type OrganizationService interface {
	CreateOrganization(ctx context.Context, req v1.CreateOrganizationRequest, opts metav1.CreateOptions) error
	UpdateOrganization(ctx context.Context, instanceId string, req v1.UpdateOrganizationRequest, opts metav1.UpdateOptions) error
	DeleteOrganization(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	GetOrganization(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Organization, error)
	DetailOrganization(ctx context.Context, org *model.Organization) (*v1.DetailOrganizationResponse, error)
	ListOrganizations(ctx context.Context, opts metav1.ListOptions) (*v1.OrganizationList, error)
	DisableOrganization(ctx context.Context, instanceId string) error
	EnableOrganization(ctx context.Context, instanceId string) error
	ListDepartments(ctx context.Context, opts metav1.ListOptions) (*v1.DepartmentList, error)
	BatchAddDepartmentMembers(ctx context.Context, department string, batchAddReq v1.BatchAddDepartmentMemberRequest) error
	BatchRemoveDepartmentMembers(ctx context.Context, department string, batchRemoveReq v1.BatchRemoveDepartmentMemberRequest) error
	ListDepartmentMembers(ctx context.Context, department string, opts metav1.ListOptions) (*v1.DepartmentMemberList, error)
	CreateDepartment(ctx context.Context, req v1.CreateDepartmentRequest, opts metav1.CreateOptions) error
	DeleteDepartment(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	UpdateDepartment(ctx context.Context, dept string, req v1.UpdateDepartmentRequest, opts metav1.UpdateOptions) error
	Init(ctx context.Context) error
}

type organizationServiceImpl struct {
	Store repository.Factory `inject:"repository"`
}

// NewOrganizationService new Organization service.
func NewOrganizationService(c config.Config) OrganizationService {
	return &organizationServiceImpl{}
}

// Init initialize default org.
func (o *organizationServiceImpl) Init(ctx context.Context) error {
	old, err := o.Store.OrganizationRepository().GetByName(ctx, model.DefaultOrganization, metav1.GetOptions{})
	if err != nil && !errors.IsCode(err, code.ErrOrgNotFound) {
		return err
	}
	if old != nil {
		return nil
	}
	createReq := v1.CreateOrganizationRequest{
		Name:        model.DefaultOrganization,
		DisplayName: "Built-in Organization",
		WebsiteUrl:  "http://iam.wecoding.top",
		Favicon:     "",
		Disabled:    false,
		Description: "System Build-in Organization",
	}
	err = o.CreateOrganization(ctx, createReq, metav1.CreateOptions{})
	if err != nil {
		return errors.WithMessagef(err, "Failed to initialize system organization")
	}
	log.Info("initialize system organization done")

	return nil
}

func (o *organizationServiceImpl) CreateOrganization(
	ctx context.Context,
	req v1.CreateOrganizationRequest,
	opts metav1.CreateOptions,
) error {
	org := &model.Organization{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		ParentID:    model.RootOrganizationID,
		Ancestors:   model.RootOrganizationID,
		DisplayName: req.DisplayName,
		WebsiteUrl:  req.WebsiteUrl,
		Favicon:     req.Favicon,
		Disabled:    req.Disabled,
		Description: req.Description,
	}
	if org.DisplayName == "" {
		org.DisplayName = req.Name
	}
	err := o.Store.OrganizationRepository().Create(ctx, org, opts)
	if err != nil {
		return err
	}

	return nil
}

func (o *organizationServiceImpl) UpdateOrganization(
	ctx context.Context,
	idOrName string,
	req v1.UpdateOrganizationRequest,
	opts metav1.UpdateOptions,
) error {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, idOrName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if req.DisplayName != "" {
		org.DisplayName = req.DisplayName
	}
	if req.Favicon != "" {
		org.Favicon = req.Favicon
	}
	if req.WebsiteUrl != "" {
		org.WebsiteUrl = req.WebsiteUrl
	}
	if req.Description != "" {
		org.Description = req.Description
	}
	err = o.Store.OrganizationRepository().Update(ctx, org, opts)
	if err != nil {
		return err
	}

	return nil
}

func (o *organizationServiceImpl) DeleteOrganization(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if org.IsSystemBuiltIn() {
		return errors.WithCode(code.ErrCannotDeleteBuiltInOrg, "")
	}
	err = o.Store.OrganizationRepository().DeleteByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return err
	}

	return nil
}

func (o *organizationServiceImpl) GetOrganization(
	ctx context.Context,
	instanceId string,
	opts metav1.GetOptions,
) (*model.Organization, error) {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (o *organizationServiceImpl) DetailOrganization(ctx context.Context, org *model.Organization) (*v1.DetailOrganizationResponse, error) {
	base := *assembler.ConvertModelToOrganizationBase(org)
	membersCount, err := o.Store.OrganizationRepository().CountDepartmentMembers(ctx, base.InstanceID, metav1.ListOptions{})
	if err != nil {
		log.Warnf("Failed to get org [%s] members", base.Name)
	}
	base.MembersCount = membersCount
	return &v1.DetailOrganizationResponse{
		OrganizationBase: base,
	}, nil
}

func (o *organizationServiceImpl) ListOrganizations(ctx context.Context, opts metav1.ListOptions) (*v1.OrganizationList, error) {
	selector, err := fields.ParseSelector(opts.FieldSelector)
	if err != nil {
		return nil, err
	}
	opts.FieldSelector = fields.AndSelectors(selector, fields.OneTermEqualSelector("parentId", model.RootOrganizationID)).String()
	orgRepo := o.Store.OrganizationRepository()

	var orgList []*v1.DetailOrganizationResponse
	organizations, err := orgRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range organizations {
		orgList = append(orgList, &v1.DetailOrganizationResponse{
			OrganizationBase: *assembler.ConvertModelToOrganizationBase(&v),
		})
	}
	count, err := orgRepo.Count(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &v1.OrganizationList{
		Items: orgList,
		ListMeta: metav1.ListMeta{
			TotalCount: count,
		},
	}, nil
}

func (o *organizationServiceImpl) DisableOrganization(ctx context.Context, instanceId string) error {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if org.IsSystemBuiltIn() {
		return errors.WithCode(code.ErrCannotDisableBuiltInOrg, "")
	}
	if org.Disabled {
		return errors.WithCode(code.ErrOrgAlreadyDisabled, "The organization [%s] is already disabled.", org.Name)
	}
	org.Disabled = true

	return o.Store.OrganizationRepository().Update(ctx, org, metav1.UpdateOptions{})
}

func (o *organizationServiceImpl) EnableOrganization(ctx context.Context, instanceId string) error {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if !org.Disabled {
		return errors.WithCode(code.ErrOrgAlreadyEnabled, "The organization [%s] is already enabled.", org.Name)
	}
	org.Disabled = false

	return o.Store.OrganizationRepository().Update(ctx, org, metav1.UpdateOptions{})
}

func (o *organizationServiceImpl) ListDepartments(ctx context.Context, opts metav1.ListOptions) (*v1.DepartmentList, error) {
	orgRepo := o.Store.OrganizationRepository()
	var deptList []*v1.DetailDepartmentResponse
	departments, err := orgRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range departments {
		deptList = append(deptList, convertToDepartmentDetail(&v))
	}
	count, err := orgRepo.Count(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &v1.DepartmentList{
		Items: deptList,
		ListMeta: metav1.ListMeta{
			TotalCount: count,
		},
	}, nil
}

func (o *organizationServiceImpl) BatchAddDepartmentMembers(
	ctx context.Context,
	department string,
	batchAddReq v1.BatchAddDepartmentMemberRequest,
) error {
	dept, err := o.GetOrganization(ctx, department, metav1.GetOptions{})
	if err != nil {
		return err
	}
	var deptMembers []*model.DepartmentMember
	for _, member := range batchAddReq.Members {
		deptMembers = append(deptMembers, &model.DepartmentMember{
			DepartmentID: dept.GetInstanceID(),
			MemberID:     member.MemberID,
		})
	}
	err = o.Store.OrganizationRepository().AddDepartmentMembers(ctx, deptMembers)
	if err != nil {
		return err
	}

	return nil
}

func (o *organizationServiceImpl) BatchRemoveDepartmentMembers(
	ctx context.Context,
	department string,
	batchRemoveReq v1.BatchRemoveDepartmentMemberRequest,
) error {
	dept, err := o.GetOrganization(ctx, department, metav1.GetOptions{})
	if err != nil {
		return err
	}
	var deptMembers []*model.DepartmentMember
	for _, member := range batchRemoveReq.Members {
		deptMembers = append(deptMembers, &model.DepartmentMember{
			DepartmentID: dept.GetInstanceID(),
			MemberID:     member.MemberID,
		})
	}
	err = o.Store.OrganizationRepository().RemoveDepartmentMembers(ctx, deptMembers)
	if err != nil {
		return err
	}

	return nil
}

func (o *organizationServiceImpl) ListDepartmentMembers(
	ctx context.Context,
	department string,
	opts metav1.ListOptions,
) (*v1.DepartmentMemberList, error) {
	items, err := o.Store.OrganizationRepository().ListDepartmentMembers(ctx, department, opts)
	if err != nil {
		return nil, err
	}
	var deptMembers []*v1.DepartmentMember
	for _, item := range items {
		deptMembers = append(deptMembers, &v1.DepartmentMember{
			MemberID:   item.MemberID,
			MemberType: "user",
		})
	}
	totalCount, err := o.Store.OrganizationRepository().CountDepartmentMembers(ctx, department, opts)
	if err != nil {
		return nil, err
	}
	resp := &v1.DepartmentMemberList{
		ListMeta: metav1.ListMeta{
			TotalCount: totalCount,
		},
		Members: deptMembers,
	}

	return resp, nil
}

func (o *organizationServiceImpl) CreateDepartment(
	ctx context.Context,
	req v1.CreateDepartmentRequest,
	opts metav1.CreateOptions,
) error {
	orgRepo := o.Store.OrganizationRepository()

	org, err := orgRepo.GetByInstanceId(ctx, req.OrganizationID, metav1.GetOptions{})
	if err != nil {
		return err
	}

	childCount, err := orgRepo.CountDepartmentByParent(ctx, org.GetInstanceID(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if childCount > DefaultMaxChildrenDepartments {
		return errors.WithCode(code.ErrMaxDepartmentsReached,
			"Organization [%s] creates an upper limit on the number of departments", org.GetName())
	}

	parent, err := orgRepo.GetByInstanceId(ctx, req.ParentID, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Failed to get parent dept [%s]", parent.GetName())
		return err
	}

	err = o.Store.ExecTx(ctx, func(ctx context.Context) error {
		err = orgRepo.UpdateIsLeafState(ctx, parent.GetInstanceID(), false)
		if err != nil {
			log.Errorf("Failed to update parent dept [%s] isLeaf", parent.GetName())
			return err
		}
		dept := assembler.ConvertCreateDeptReqToModel(req, parent)
		if dept.DisplayName == "" {
			dept.DisplayName = req.Name
		}
		err = orgRepo.Create(ctx, dept, opts)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (o *organizationServiceImpl) DeleteDepartment(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	orgRepo := o.Store.OrganizationRepository()

	dept, err := orgRepo.GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}

	count, err := orgRepo.CountDepartmentByParent(ctx, dept.InstanceID, metav1.ListOptions{})
	if err != nil || count > 0 {
		if err != nil {
			return err
		}
		return errors.WithCode(code.ErrSubDepartmentsExist, "[%d] sub-departments exist and cannot be deleted", count)
	}

	err = orgRepo.DeleteByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return err
	}

	if err := o.updateIsLeafStateIfNecessary(dept.ParentID); err != nil {
		log.Errorf("Failed to update parent [%s] isLeaf state", dept.ParentID)
	}

	return nil
}

func (o *organizationServiceImpl) UpdateDepartment(
	ctx context.Context,
	deptId string,
	req v1.UpdateDepartmentRequest,
	opts metav1.UpdateOptions,
) error {
	var (
		orgRepo      = o.Store.OrganizationRepository()
		dept         *model.Organization
		newParent    *model.Organization
		oldParentID  string
		oldAncestors string
		newAncestors string
	)

	_, err := orgRepo.GetByInstanceId(ctx, req.OrganizationID, metav1.GetOptions{})
	if err != nil {
		return err
	}

	newParent, err = orgRepo.GetByInstanceId(ctx, req.ParentID, metav1.GetOptions{})
	if err != nil {
		return err
	}

	dept, err = orgRepo.GetByInstanceId(ctx, deptId, metav1.GetOptions{})
	if err != nil {
		return err
	}

	err = o.Store.ExecTx(ctx, func(ctx context.Context) error {
		oldParentID = dept.ParentID
		oldAncestors = dept.Ancestors
		newAncestors = newParent.Ancestors + "," + newParent.InstanceID
		dept.Ancestors = newAncestors
		if oldAncestors != newAncestors {
			err = o.updateDeptChildren(deptId, newAncestors, oldAncestors)
			if err != nil {
				log.Errorf("Failed to update the [%s] sub department.", dept.GetName())
				return err
			}
		}

		if req.DisplayName != "" {
			dept.DisplayName = req.DisplayName
		}
		if req.WebsiteUrl != "" {
			dept.WebsiteUrl = req.WebsiteUrl
		}
		if req.Description != "" {
			dept.Description = req.Description
		}
		if req.Favicon != "" {
			dept.Favicon = req.Favicon
		}
		dept.ParentID = newParent.GetInstanceID()

		// update dept info
		err = orgRepo.Update(ctx, dept, opts)
		if err != nil {
			return err
		}
		// update new parent isLeaf state
		err = orgRepo.UpdateIsLeafState(ctx, newParent.GetInstanceID(), false)
		if err != nil {
			log.Errorf("Failed to update new parent [%s] isLeaf state", newParent.GetName())
			return err
		}
		return nil
	})

	// update old parent isLeaf state if necessary
	err = o.updateIsLeafStateIfNecessary(oldParentID)
	if err != nil {
		log.Errorf("Failed to update old parent [%s] isLeaf state", oldParentID)
	}

	return err
}

func (o *organizationServiceImpl) updateIsLeafStateIfNecessary(orgOrDept string) error {
	orgRepo := o.Store.OrganizationRepository()
	childCount, err := orgRepo.CountDepartmentByParent(context.Background(), orgOrDept, metav1.ListOptions{})
	if err != nil {
		return err
	}
	// if the parent department has no children, update the parent isLeaf to true
	isLeaf := childCount <= 0
	err = orgRepo.UpdateIsLeafState(context.Background(), orgOrDept, isLeaf)
	if err != nil {
		return err
	}
	return nil
}

func (o *organizationServiceImpl) updateDeptChildren(deptId, newAncestors, oldAncestors string) error {
	children, err := o.Store.OrganizationRepository().ListChildDepartments(context.Background(), deptId, metav1.ListOptions{})
	if err != nil {
		return err
	}
	var needUpdates []*model.Organization
	for _, child := range children {
		child.Ancestors = strings.ReplaceAll(child.Ancestors, oldAncestors, newAncestors)
		needUpdates = append(needUpdates, &child)
	}
	if len(needUpdates) > 0 {
		return o.Store.OrganizationRepository().BatchUpdate(context.Background(), needUpdates, metav1.UpdateOptions{})
	}
	return nil
}

func convertToDepartmentDetail(model *model.Organization) *v1.DetailDepartmentResponse {
	return &v1.DetailDepartmentResponse{
		OrganizationBase: *assembler.ConvertModelToOrganizationBase(model),
		OrganizationID:   model.Owner,
	}
}
