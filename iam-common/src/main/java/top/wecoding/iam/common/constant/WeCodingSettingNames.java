package top.wecoding.iam.common.constant;

/**
 * The names for all the configuration settings.
 *
 * @author liuyuhui
 * @since 0.5
 */
public class WeCodingSettingNames {

  private static final String SETTINGS_NAMESPACE = "wecoding.settings.";

  /** The names for authorization server configuration settings. */
  public static final class AuthorizationServer {

    private static final String AUTHORIZATION_SERVER_SETTINGS_NAMESPACE =
        SETTINGS_NAMESPACE.concat("authorization-server.");

    /** Set the Resource Owner Token Revocation endpoint. */
    public static final String RESOURCE_OWNER_TOKEN_ENDPOINT =
        AUTHORIZATION_SERVER_SETTINGS_NAMESPACE.concat("resource-owner-token-endpoint");

    private AuthorizationServer() {}
  }
}
