package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
public interface PasswordRequest {

  @JsonProperty("reset")
  boolean reset();

  @JsonProperty("old_pwd")
  String getOldPwd();

  @JsonProperty("new_pwd")
  String getNewPwd();
}
