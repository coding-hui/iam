package top.wecoding.iam.common.model.response;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author liuyuhui
 * @since 0.5
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class LoginResponse {

  @JsonProperty("access_token")
  private String accessToken;

  @JsonProperty("refresh_token")
  @JsonInclude(JsonInclude.Include.NON_NULL)
  private String refreshToken;

  @JsonProperty("token_type")
  @JsonInclude(JsonInclude.Include.NON_NULL)
  private String accessTokenType;

  @JsonProperty("expires_in")
  private Long expiresIn;
}
