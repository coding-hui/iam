package top.wecoding.iam.model.request;

import lombok.*;
import top.wecoding.core.model.request.PageRequest;

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
@EqualsAndHashCode(callSuper = true)
public class GroupInfoPageRequest extends PageRequest {

  private Set<String> groupIds;

  private String fuzzyName;
}
