package top.wecoding.iam.common.model;

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

  private String userId;

  private String tenantId;

  private String tenantName;

  private String username;

  private String nickName;

  private String password;

  private String avatar;

  private Date birthday;

  private String gender;

  private String email;

  private String phone;

  private String lastLoginIp;

  private Date lastLoginTime;

  private Integer userType;

  private Integer userState;

  private Boolean defaultPwd;

  private String infos;
}
