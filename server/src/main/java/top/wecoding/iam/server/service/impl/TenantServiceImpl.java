package top.wecoding.iam.server.service.impl;

import com.baomidou.mybatisplus.core.toolkit.IdWorker;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import java.util.Locale;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import top.wecoding.core.enums.iam.UserTypeEnum;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.core.util.PageUtil;
import top.wecoding.iam.common.enums.IamErrorCode;
import top.wecoding.iam.common.model.TenantInfo;
import top.wecoding.iam.common.model.request.CreateTenantRequest;
import top.wecoding.iam.common.model.request.TenantInfoPageRequest;
import top.wecoding.iam.common.model.request.UpdateTenantRequest;
import top.wecoding.iam.common.model.response.CreateTenantResponse;
import top.wecoding.iam.common.model.response.TenantInfoResponse;
import top.wecoding.iam.common.util.AuthUtil;
import top.wecoding.iam.server.convert.TenantConvert;
import top.wecoding.iam.server.entity.Tenant;
import top.wecoding.iam.server.mapper.TenantMapper;
import top.wecoding.iam.server.service.TenantService;
import top.wecoding.iam.server.util.TenantUtil;
import top.wecoding.mybatis.base.BaseServiceImpl;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Service
@RequiredArgsConstructor
public class TenantServiceImpl extends BaseServiceImpl<TenantMapper, Tenant>
    implements TenantService {

  private final TenantMapper tenantMapper;

  @Override
  public TenantInfoResponse getInfo(String tenantId) {
    Tenant tenant = tenantMapper.getByTenantId(tenantId);

    AssertUtils.isNotNull(tenant, IamErrorCode.TENANT_DOES_NOT_EXIST);

    TenantInfo tenantInfo = TenantConvert.INSTANCE.toTenantInfo(tenant);

    return TenantInfoResponse.builder().tenantInfo(tenantInfo).build();
  }

  @Override
  public CreateTenantResponse create(CreateTenantRequest createTenantRequest) {
    String tenantName = createTenantRequest.getTenantName().trim().toLowerCase(Locale.ROOT);

    TenantUtil.checkTenantName(tenantName);

    AssertUtils.isNull(
        tenantMapper.getByTenantName(tenantName), IamErrorCode.TENANT_NAME_IS_ALREADY_OCCUPIED);

    String tenantId = IdWorker.getIdStr();

    int count = tenantMapper.count();

    AssertUtils.isFalse(100 <= count, IamErrorCode.EXCEEDED_MAXIMUM_NUMBER_OF_TENANTS, 100);

    Tenant tenant =
        Tenant.builder()
            .tenantId(tenantId)
            .tenantName(tenantName)
            .loginType(UserTypeEnum.LOCAL.code())
            .annotate(createTenantRequest.getDescription())
            .ownerId(AuthUtil.currentUserId())
            .username(AuthUtil.currentUsername())
            .build();

    AssertUtils.isFalse(1 != tenantMapper.insert(tenant), IamErrorCode.TENANT_ADD_FAILED);

    return CreateTenantResponse.builder().tenantId(tenantId).tenantName(tenantName).build();
  }

  @Override
  public void update(String tenantId, UpdateTenantRequest updateTenantRequest) {
    String annotate = updateTenantRequest.getDescription();

    Tenant tenant = tenantMapper.getByTenantId(tenantId);

    AssertUtils.isNotNull(tenant, IamErrorCode.TENANT_DOES_NOT_EXIST);

    int rows =
        tenantMapper.updateTenantAnnotate(tenant.getId(), annotate, AuthUtil.currentUsername());

    AssertUtils.isFalse(1 != rows, IamErrorCode.TENANT_UPDATE_FAILED);
  }

  @Override
  public void delete(String tenantId) {
    Tenant tenant = tenantMapper.getByTenantId(tenantId);

    AssertUtils.isNotNull(tenant, IamErrorCode.TENANT_DOES_NOT_EXIST);

    AssertUtils.isFalse(
        1 != tenantMapper.deleteById(tenant.getId()), IamErrorCode.TENANT_DELETE_FAILED);
  }

  @Override
  public PageInfo<TenantInfoResponse> infoPage(TenantInfoPageRequest tenantInfoPageRequest) {
    Page<Tenant> pageResult =
        tenantMapper.page(
            PageUtil.getPageFromRequest(tenantInfoPageRequest), tenantInfoPageRequest);

    return PageInfo.of(pageResult.getRecords(), tenantInfoPageRequest, pageResult.getTotal())
        .map(
            (tenant -> {
              TenantInfo tenantInfo = TenantConvert.INSTANCE.toTenantInfo(tenant);
              return TenantInfoResponse.builder().tenantInfo(tenantInfo).build();
            }));
  }
}
