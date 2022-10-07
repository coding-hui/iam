package top.wecoding.iam.common.model.request;

import top.wecoding.iam.common.util.AuthUtil;

/**
 * @author liuyuhui
 * @date 2022/10/7
 * @qq 1515418211
 */
public interface CommonRequest {

  default String getCurrentTenantId() {
    return AuthUtil.currentTenantId();
  }

  default String getCurrentUserId() {
    return AuthUtil.currentUserId();
  }
}
