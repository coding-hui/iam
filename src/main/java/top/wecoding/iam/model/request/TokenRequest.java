package top.wecoding.iam.model.request;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

/**
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class TokenRequest {

  private String grantType;
  private String account;
  private String password;
  private String refreshToken;
  private Map<String, Object> params = Collections.unmodifiableMap(new HashMap<>());

  public static TokenRequest of(LoginRequest loginRequest) {
    return TokenRequest.builder()
        .account(loginRequest.getAccount())
        .password(loginRequest.getPassword())
        .grantType(loginRequest.getGrantType())
        .refreshToken(loginRequest.getRefreshToken())
        .build();
  }
}
