package top.wecoding.iam.sdk.context;

import com.alibaba.ttl.TransmittableThreadLocal;
import top.wecoding.iam.common.model.AuthInfo;
import top.wecoding.core.enums.rest.CommonErrorCodeEnum;
import top.wecoding.core.util.AssertUtils;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public class ThreadLocalAuthContextHolderStrategy {

  private static final TransmittableThreadLocal<AuthInfo> contextHolder =
      new TransmittableThreadLocal<>();

  public void clearContext() {
    contextHolder.remove();
  }

  public AuthInfo getContext() {
    return contextHolder.get();
  }

  public void setContext(AuthInfo context) {
    AssertUtils.isNotNull(
        context, CommonErrorCodeEnum.COMMON_ERROR, "AuthInfo Context cannot be null");
    contextHolder.set(context);
  }
}
