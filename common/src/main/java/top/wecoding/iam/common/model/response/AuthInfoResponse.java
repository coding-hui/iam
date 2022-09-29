package top.wecoding.iam.common.model.response;

import io.swagger.annotations.ApiModelProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import top.wecoding.core.constant.SecurityConstants;

import java.time.LocalDateTime;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class AuthInfoResponse {

  @ApiModelProperty(value = "会话唯一标识")
  private String uuid;

  @ApiModelProperty(value = "令牌")
  private String accessToken;

  @ApiModelProperty(value = "令牌类型")
  private String tokenType;

  @ApiModelProperty(value = "刷新令牌")
  private String refreshToken;

  @ApiModelProperty(value = "租房Id")
  private String tenantId;

  @ApiModelProperty(value = "用户Id")
  private String userId;

  @ApiModelProperty(value = "用户名")
  private String realName;

  @ApiModelProperty(value = "账号名")
  private String account;

  @ApiModelProperty(value = "到期时间")
  private LocalDateTime expiration;

  @ApiModelProperty(value = "有效期")
  private Long expireMillis;

  @ApiModelProperty(value = "客户端Id")
  private String clientId;

  @ApiModelProperty(value = "许可证")
  private String license = SecurityConstants.PROJECT_LICENSE;
}