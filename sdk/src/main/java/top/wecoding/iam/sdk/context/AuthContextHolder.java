package top.wecoding.iam.sdk.context;

import top.wecoding.iam.common.model.AuthInfo;

/**
 * 获取当前线程变量中认证信息
 *
 * @author liuyuhui
 * @qq 1515418211
 */
public class AuthContextHolder {

  private static ThreadLocalAuthContextHolderStrategy strategy;

  static {
    initializeStrategy();
  }

  private static void initializeStrategy() {
    strategy = new ThreadLocalAuthContextHolderStrategy();
  }

  public static void clearContext() {
    strategy.clearContext();
  }

  public static AuthInfo getContext() {
    return strategy.getContext();
  }

  public static void setContext(AuthInfo context) {
    strategy.setContext(context);
  }
}
