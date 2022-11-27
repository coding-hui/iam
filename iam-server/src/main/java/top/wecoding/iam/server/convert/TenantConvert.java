package top.wecoding.iam.server.convert;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;
import top.wecoding.iam.common.model.TenantInfo;
import top.wecoding.iam.common.model.request.CreateTenantRequest;
import top.wecoding.iam.common.model.request.UpdateTenantRequest;
import top.wecoding.iam.common.model.response.TenantInfoResponse;
import top.wecoding.iam.server.entity.Tenant;

/**
 * @author liuyuhui
 * @date 2022/10/7
 */
@Mapper
public interface TenantConvert {

  TenantConvert INSTANCE = Mappers.getMapper(TenantConvert.class);

  @Mapping(target = "createTimestamp", expression = "java(tenant.getCreateTime().getTime())")
  TenantInfo toTenantInfo(Tenant tenant);

  TenantInfoResponse toTenantInfoResponse(Tenant tenant);

  TenantInfoResponse toTenantInfoResponse(Tenant tenant, Long createTimestamp);

  Tenant toTenant(CreateTenantRequest createTenantRequest);

  Tenant toTenant(UpdateTenantRequest updateTenantRequest);
}
