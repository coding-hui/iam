package top.wecoding.iam.common.model.request;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import lombok.experimental.SuperBuilder;
import top.wecoding.commons.core.model.request.PageRequest;

/**
 * @author liuyuhui
 * @date 2022/10/6
 */
@Data
@SuperBuilder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class TokenInfoPageRequest extends PageRequest implements CommonRequest {

  @JsonProperty("user_id")
  private Long userId;

  @JsonProperty("client_id")
  private String clientId;
}
