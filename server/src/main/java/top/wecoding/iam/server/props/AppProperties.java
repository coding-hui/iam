package top.wecoding.iam.server.props;

import org.springframework.boot.context.properties.ConfigurationProperties;

/**
 * @author liuyuhui
 * @date 2022/10/2
 * @qq 1515418211
 */
@ConfigurationProperties(prefix = AppProperties.PREFIX)
public class AppProperties {

  public static final String PREFIX = "wecoding.iam";

  public static Long userFailCount = 3L;

  public static Long userFailLockTime = 60 * 60 * 1000L;

  public static Long getUserFailCount() {
    return userFailCount;
  }

  public static void setUserFailCount(Long userFailCount) {
    AppProperties.userFailCount = userFailCount;
  }

  public static Long getUserFailLockTime() {
    return userFailLockTime;
  }

  public static void setUserFailLockTime(Long userFailLockTime) {
    AppProperties.userFailLockTime = userFailLockTime;
  }
}
