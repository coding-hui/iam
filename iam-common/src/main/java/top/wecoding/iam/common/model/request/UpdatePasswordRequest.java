package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author liuyuhui
 * @date 2022/10/6
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UpdatePasswordRequest implements PasswordRequest {

  @NotBlank
  @JsonProperty("old_pwd")
  private String oldPwd;

  @NotBlank
  @JsonProperty("new_pwd")
  private String newPwd;

  @Override
  public boolean reset() {
    return false;
  }

  @Override
  public String getOldPwd() {
    return oldPwd;
  }

  @Override
  public String getNewPwd() {
    return newPwd;
  }
}
