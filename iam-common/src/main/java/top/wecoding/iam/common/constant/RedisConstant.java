package top.wecoding.iam.common.constant;

/**
 * @author liuyuhui
 * @date 2022/9/11
 */
public interface RedisConstant {

  String OAUTH_ACCESS_PREFIX = "access_token:*";

  String USER_FAIL_COUNT = "user:fail:count";

  String USER_TOKEN = "token";

  String USER_DETAILS = "user_details";

  String AUTHORIZATION_CONSENT = "token:consent";

  String CLIENT_DETAILS_KEY = "client:details";
}
