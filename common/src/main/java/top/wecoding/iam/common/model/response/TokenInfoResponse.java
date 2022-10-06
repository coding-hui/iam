package top.wecoding.iam.common.model.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;

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

  private String id;

  private Long userId;

  private String clientId;

  private String username;

  private String accessToken;

  private Instant issuedAt;

  private Instant expiresAt;
}
