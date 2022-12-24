package top.wecoding.iam.server.util;

import java.util.regex.Pattern;
import lombok.experimental.UtilityClass;
import top.wecoding.commons.core.enums.SystemErrorCodeEnum;
import top.wecoding.commons.core.exception.IllegalParameterException;

/**
 * @author liuyuhui
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
}
