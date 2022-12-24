package top.wecoding.iam.common.enums;

import java.util.Map;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.experimental.Accessors;
import top.wecoding.iam.common.util.EnumUtil;

/**
 * @author liuyuhui
 * @date 2022/9/11
 */
@Getter
@AllArgsConstructor
@Accessors(fluent = true)
public enum AuthType {
  PASSWORD("PASSWORD"),
  API_TOKEN("API_TOKEN"),
  LDAP("LDAP"),
  AD("AD");

  private static final Map<String, AuthType> DICT =
      EnumUtil.ofDict(AuthType::values, AuthType::code);

  private final String code;

  public static AuthType of(String code) {
    return DICT.get(code);
  }

  public boolean is(String code) {
    return this.code.equals(code);
  }
}
