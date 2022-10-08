package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Date;
import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

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

  @NotBlank
  @JsonProperty("username")
  private String username;

  @NotBlank
  @JsonProperty("password")
  private String password;

  @JsonProperty("nick_name")
  private String nickName;

  @JsonProperty("avatar")
  private String avatar;

  @JsonProperty("birthday")
  private Date birthday;

  @JsonProperty("gender")
  private String gender;

  @Email
  @JsonProperty("email")
  private String email;

  @JsonProperty("phone")
  private String phone;

  @JsonProperty("user_type")
  private Integer userType;

  @JsonProperty("country")
  private String country;
}
