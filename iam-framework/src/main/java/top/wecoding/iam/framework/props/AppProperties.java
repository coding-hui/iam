package top.wecoding.iam.framework.props;

import org.springframework.boot.context.properties.ConfigurationProperties;

/**
 * @author liuyuhui
 * @date 2022/10/2
 */
@ConfigurationProperties(prefix = AppProperties.PREFIX)
public class AppProperties {

  public static final String PREFIX = "wecoding.iam";

  public static Integer userFailCount = 3;

  public static Long userFailLockTime = 60 * 10L;

  public static Integer getUserFailCount() {
    return userFailCount;
  }

  public static void setUserFailCount(Integer userFailCount) {
    AppProperties.userFailCount = userFailCount;
  }

  public static Long getUserFailLockTime() {
    return userFailLockTime;
  }

  public static void setUserFailLockTime(Long userFailLockTime) {
    AppProperties.userFailLockTime = userFailLockTime;
  }
}
