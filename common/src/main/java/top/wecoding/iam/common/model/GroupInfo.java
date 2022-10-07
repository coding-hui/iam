package top.wecoding.iam.common.model;

import com.fasterxml.jackson.annotation.JsonProperty;
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
public class GroupInfo {

  @JsonProperty("group_id")
  private String groupId;

  @JsonProperty("tenant_id")
  private String tenantId;

  @JsonProperty("group_name")
  private String groupName;
}
