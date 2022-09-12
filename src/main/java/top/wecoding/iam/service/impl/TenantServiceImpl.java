package top.wecoding.iam.service.impl;

import org.springframework.stereotype.Service;
import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.mapper.TenantMapper;
import top.wecoding.iam.model.request.CreateTenantRequest;
import top.wecoding.iam.model.request.TenantInfoPageRequest;
import top.wecoding.iam.model.request.UpdateTenantRequest;
import top.wecoding.iam.model.response.CreateTenantResponse;
import top.wecoding.iam.model.response.TenantInfoResponse;
import top.wecoding.iam.pojo.Tenant;
import top.wecoding.iam.service.TenantService;
import top.wecoding.mybatis.base.BaseServiceImpl;

import javax.annotation.Resource;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Service
public class TenantServiceImpl extends BaseServiceImpl<TenantMapper, Tenant>
    implements TenantService {

  @Resource private TenantMapper tenantMapper;

  @Override
  public TenantInfoResponse getInfo(String tenantId) {
    return null;
  }

  @Override
  public CreateTenantResponse create(CreateTenantRequest createTenantRequest) {
    return null;
  }

  @Override
  public void update(String tenantId, UpdateTenantRequest updateTenantRequest) {}

  @Override
  public void delete(String tenantId) {}

  @Override
  public PageInfo<TenantInfoResponse> infoPage(TenantInfoPageRequest tenantInfoPageRequest) {
    return null;
  }
}
