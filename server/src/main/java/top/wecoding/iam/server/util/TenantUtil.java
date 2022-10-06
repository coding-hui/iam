package top.wecoding.iam.server.util;

import lombok.experimental.UtilityClass;
import top.wecoding.core.enums.rest.SystemErrorCodeEnum;
import top.wecoding.core.exception.IllegalParameterException;
import top.wecoding.iam.common.model.TenantInfo;
import top.wecoding.iam.server.pojo.Tenant;

import java.util.Date;
import java.util.Optional;
import java.util.regex.Pattern;

/**
 * @author liuyuhui
 * @qq 1515418211
 */
@UtilityClass
public class TenantUtil {

  private final Pattern TENANT_NAME_PATTERN =
      Pattern.compile("([a-z0-9A-Z_][-a-z0-9A-Z_]*)?[a-z0-9A-Z_]");

  public void checkTenantName(String tenantName) {
    int length = tenantName.length();
    if (0 >= length || 30 < length) {
      throw new IllegalParameterException(SystemErrorCodeEnum.PARAM_ERROR);
    }

    if (!TENANT_NAME_PATTERN.matcher(tenantName).matches()) {
      throw new IllegalParameterException(SystemErrorCodeEnum.PARAM_ERROR);
    }
  }

  public TenantInfo ofTenantInfo(Tenant tenant) {
    if (null == tenant) {
      return null;
    }
    return TenantInfo.builder()
        .tenantId(tenant.getTenantId())
        .tenantName(tenant.getTenantName())
        .username(tenant.getUsername())
        .ownerId(tenant.getOwnerId())
        .annotate(tenant.getAnnotate())
        .createTime(tenant.getCreateTime())
        .createTimestamp(
            Optional.ofNullable(tenant.getCreateTime()).map(Date::getTime).orElse(null))
        .build();
  }
}
