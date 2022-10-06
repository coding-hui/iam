package top.wecoding.iam.common.model.request;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;
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
public class CreateUserRequest {

  @NotBlank private String tenantId;

  @NotBlank private String tenantName;

  @NotBlank private String username;

  @NotBlank private String password;

  private String nickName;

  private String avatar;

  private Date birthday;

  private String gender;

  @Email private String email;

  private String phone;

  private Integer userType;

  private String country;
}
