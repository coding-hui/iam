package top.wecoding.iam.common.model.response;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class CreateGroupResponse {

  @JsonProperty("group_id")
  private String groupId;

  @JsonProperty("group_name")
  private String groupName;
}
