package top.wecoding.iam.common.model.request;

import lombok.*;
import top.wecoding.core.model.request.PageRequest;

/**
 * @author liuyuhui
 * @date 2022/10/6
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class TokenInfoPageRequest extends PageRequest {

  private String clientId;

  private String userId;
}
