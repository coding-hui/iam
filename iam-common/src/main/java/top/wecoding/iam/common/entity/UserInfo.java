package top.wecoding.iam.common.entity;

import com.fasterxml.jackson.databind.PropertyNamingStrategies;
import com.fasterxml.jackson.databind.annotation.JsonNaming;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.Date;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@JsonNaming(PropertyNamingStrategies.SnakeCaseStrategy.class)
public class UserInfo implements Serializable {

  @Serial private static final long serialVersionUID = -5097425171646833754L;

  private String userId;

  private String tenantId;

  private String name;

  private String tenantName;

  private String username;

  private String nickName;

  private String password;

  private String avatar;

  private Date birthday;

  private String gender;

  private String email;

  private String phone;

  private Integer userType;

  private Integer userState;

  private String country;

  private String company;

  private String address;

  private String province;

  private String city;

  private String streetAddress;

  private String postalCode;

  private String externalId;

  private String createTime;

  private String lastLoginIp;

  private Date lastLoginTime;

  private Integer loginCount;
}
