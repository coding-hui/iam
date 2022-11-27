package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import lombok.experimental.SuperBuilder;
import top.wecoding.commons.core.model.request.PageRequest;

import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
@Data
@SuperBuilder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class GroupInfoPageRequest extends PageRequest implements CommonRequest {

  @JsonProperty("group_ids")
  private Set<String> groupIds;

  @JsonProperty("fuzzy_name")
  private String fuzzyName;
}
