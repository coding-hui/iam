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
 * @date 2022/10/5
 */
@Data
@SuperBuilder
@NoArgsConstructor
@AllArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class CreateOauth2ClientPageRequest extends PageRequest {

  @JsonProperty("client_id")
  private String clientId;

  @JsonProperty("client_fuzzy_name")
  private String clientFuzzyName;
}
