package top.wecoding.iam.common.enums;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.experimental.Accessors;

/**
 * @author liuyuhui
 * @date 2022/9/11
 * @qq 1515418211
 */
@Getter
@AllArgsConstructor
@Accessors(fluent = true)
public enum AuthType {
  WEB("WEB"),
  API_TOKEN("API_TOKEN");

  private static final Map<String, AuthType> DICT =
      new HashMap<String, AuthType>() {
        {
          Arrays.asList(AuthType.values()).forEach(item -> put(item.code, item));
        }
      };

  private final String code;

  public static AuthType of(String code) {
    return DICT.get(code);
  }

  public boolean is(String code) {
    return this.code.equals(code);
  }
}
