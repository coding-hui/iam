package top.wecoding.iam.server.service;

import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.common.model.request.CreateTenantRequest;
import top.wecoding.iam.common.model.request.TenantInfoPageRequest;
import top.wecoding.iam.common.model.request.UpdateTenantRequest;
import top.wecoding.iam.common.model.response.CreateTenantResponse;
import top.wecoding.iam.common.model.response.TenantInfoResponse;
import top.wecoding.iam.server.pojo.Tenant;
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
