package top.wecoding.iam.service;

import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.model.request.CreateTenantRequest;
import top.wecoding.iam.model.request.TenantInfoPageRequest;
import top.wecoding.iam.model.request.UpdateTenantRequest;
import top.wecoding.iam.model.response.CreateTenantResponse;
import top.wecoding.iam.model.response.TenantInfoResponse;
import top.wecoding.iam.pojo.Tenant;
import top.wecoding.mybatis.base.BaseService;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface TenantService extends BaseService<Tenant> {

  TenantInfoResponse getInfo(String tenantId);

  CreateTenantResponse create(CreateTenantRequest createTenantRequest);

  void update(String tenantId, UpdateTenantRequest updateTenantRequest);

  void delete(String tenantId);

  PageInfo<TenantInfoResponse> infoPage(TenantInfoPageRequest tenantInfoPageRequest);
}
