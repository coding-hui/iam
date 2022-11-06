package top.wecoding.iam.common.model.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonUnwrapped;
import java.util.List;
import java.util.Set;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import top.wecoding.iam.common.entity.UserInfo;
import top.wecoding.iam.common.model.GroupInfo;

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

  @JsonUnwrapped private List<GroupInfo> groups;

  @JsonProperty("permissions")
  private Set<String> permissions;

  @JsonProperty("roles")
  private Set<String> roles;
}
