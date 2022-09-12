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
public class UserInfoPageRequest extends PageRequest {

  private String tenantId;

  private Set<String> userIds;

  private String fuzzyName;

  private Integer role;

  private Integer state;

  private Boolean defaultPwd;

  private String usernameSort;
}
