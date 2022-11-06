package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Set;
import javax.validation.constraints.NotBlank;
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
public class UpdateGroupRequest {

  @NotBlank
  @JsonProperty("group_name")
  private String groupName;

  @JsonProperty("description")
  private String description;

  @NotNull
  @JsonProperty("input_id_set")
  private Set<String> inputIdSet;

  @NotNull
  @JsonProperty("output_id_set")
  private Set<String> outputIdSet;
}
