package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Set;
import javax.validation.constraints.NotBlank;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.experimental.SuperBuilder;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@SuperBuilder
@NoArgsConstructor
@AllArgsConstructor
public class CreateGroupRequest {

  @NotBlank
  @JsonProperty("group_name")
  private String groupName;

  @JsonProperty("user_ids")
  private Set<String> userIds;
}
