package top.wecoding.iam.common.model.request;

import lombok.*;
import top.wecoding.core.model.request.PageRequest;

/**
 * @author liuyuhui
 * @date 2022/10/5
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class CreateOauth2ClientPageRequest extends PageRequest {

  private String clientId;

  private String clientFuzzyName;
}
