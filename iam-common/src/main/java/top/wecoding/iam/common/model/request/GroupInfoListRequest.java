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
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class GroupInfoListRequest implements CommonRequest {

  @JsonProperty("group_ids")
  private Set<String> groupIds;
}
