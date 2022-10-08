package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Date;
import javax.validation.constraints.NotNull;
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
public class UpdateUserRequest {

  @NotNull
  @JsonProperty("user_id")
  private String userId;

  @JsonProperty("nick_name")
  private String nickName;

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

  @JsonProperty("country")
  private String country;
}
