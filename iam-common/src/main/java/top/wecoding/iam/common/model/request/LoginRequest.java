package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
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
public class LoginRequest implements CommonRequest {

  @NotBlank
  @JsonProperty("auth_type")
  private String authType;

  @JsonProperty("client_id")
  private String clientId;

  @JsonProperty("client_secret")
  private String clientSecret;

  @JsonProperty("password_payload")
  private PasswordPayload passwordPayload;

  @JsonProperty("options")
  private Options options;

  @Data
  public static class PasswordPayload {

    @JsonProperty("account")
    private String account;

    @JsonProperty("password")
    private String password;
  }

  @Data
  public static class Options {

    @JsonProperty("captcha_code")
    private String captchaCode;

    @JsonProperty("scope")
    private String scope;
  }
}
