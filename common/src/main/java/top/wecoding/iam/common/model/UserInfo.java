package top.wecoding.iam.common.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Date;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UserInfo {

  @JsonProperty("user_id")
  private String userId;

  @JsonProperty("tenant_id")
  private String tenantId;

  @JsonProperty("tenant_name")
  private String tenantName;

  @JsonProperty("username")
  private String username;

  @JsonProperty("nick_name")
  private String nickName;

  @JsonProperty("password")
  private String password;

  @JsonProperty("avatar")
  private String avatar;

  @JsonProperty("birthday")
  private Date birthday;

  @JsonProperty("gender")
  private String gender;

  @JsonProperty("email")
  private String email;

  @JsonProperty("phone")
  private String phone;

  @JsonProperty("last_login_ip")
  private String lastLoginIp;

  @JsonProperty("last_login_time")
  private Date lastLoginTime;

  @JsonProperty("user_type")
  private Integer userType;

  @JsonProperty("user_state")
  private Integer userState;

  @JsonProperty("default_pwd")
  private Boolean defaultPwd;

  @JsonProperty("infos")
  private String infos;
}
