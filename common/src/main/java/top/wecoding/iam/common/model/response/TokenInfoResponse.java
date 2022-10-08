package top.wecoding.iam.common.model.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.Instant;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author liuyuhui
 * @date 2022/10/6
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class TokenInfoResponse {

  @JsonProperty("id")
  private String id;

  @JsonProperty("user_id")
  private Long userId;

  @JsonProperty("client_id")
  private String clientId;

  @JsonProperty("username")
  private String username;

  @JsonProperty("access_token")
  private String accessToken;

  @JsonProperty("issued_at")
  private Instant issuedAt;

  @JsonProperty("expires_at")
  private Instant expiresAt;
}
