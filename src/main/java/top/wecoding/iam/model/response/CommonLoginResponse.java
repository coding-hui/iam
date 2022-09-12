package top.wecoding.iam.model.response;

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
public class CommonLoginResponse {

  private String tenantId;

  private String tenantName;

  private String userId;

  private String username;

  private Integer userType;
}
