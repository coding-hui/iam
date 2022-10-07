package top.wecoding.iam.common.model.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonUnwrapped;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import top.wecoding.iam.common.model.GroupInfo;
import top.wecoding.iam.common.model.UserInfo;

import java.util.List;
import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UserInfoResponse {

  @JsonUnwrapped private UserInfo userInfo;

  @JsonUnwrapped private List<GroupInfo> groupInfoList;

  @JsonProperty("permissions")
  private Set<String> permissions;

  @JsonProperty("roles")
  private Set<String> roles;
}
