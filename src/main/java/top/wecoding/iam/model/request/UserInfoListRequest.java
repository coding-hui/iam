package top.wecoding.iam.model.request;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

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
public class UserInfoListRequest {

  private String tenantId;

  private Set<String> userIds;

  private String fuzzyName;

  private Integer role;

  private Integer state;

  private Boolean defaultPwd;

  private String usernameSort;
}
