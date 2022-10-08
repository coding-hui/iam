package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Set;
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
public class UserInfoListRequest implements CommonRequest {

  @JsonProperty("tenant_id")
  private String tenantId;

  @JsonProperty("user_ids")
  private Set<String> userIds;

  @JsonProperty("fuzzy_name")
  private String fuzzyName;

  @JsonProperty("role")
  private Integer role;

  @JsonProperty("state")
  private Integer state;

  @JsonProperty("default_pwd")
  private Boolean defaultPwd;

  @JsonProperty("username_sort")
  private String usernameSort;
}
