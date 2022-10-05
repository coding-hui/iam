package top.wecoding.iam.common.pojo;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.experimental.Accessors;
import top.wecoding.core.constant.SecurityConstants;

import java.time.LocalDateTime;

/**
 * 认证成功信息
 *
 * @author liuyuhui
 * @date 2022/5/11
 * @qq 1515418211
 */
@Data
@Builder
@Accessors(chain = true)
@NoArgsConstructor
@AllArgsConstructor
public class AuthInfo {

  private String tenantId;

  private String tenantName;

  private String userId;

  private String username;

  private Integer userType;

  private String authType;

  private String clientId;

  private String accessToken;

  private String tokenType;

  private String refreshToken;

  private Long loginTime;

  private LocalDateTime expiration;

  private Long expireMillis;

  private String license = SecurityConstants.PROJECT_LICENSE;
}
