package top.wecoding.iam.common.convert;

import java.util.HashMap;
import java.util.Map;
import org.springframework.core.convert.converter.Converter;
import org.springframework.security.oauth2.core.OAuth2Error;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import org.springframework.util.StringUtils;
import top.wecoding.commons.core.model.R;
import top.wecoding.commons.core.util.JsonUtil;

/**
 * @author liuyuhui
 * @since 0.5
 */
public class RestOAuth2ErrorParametersConverter
    implements Converter<OAuth2Error, Map<String, String>> {

  @Override
  @SuppressWarnings("unchecked")
  public Map<String, String> convert(OAuth2Error oauth2Error) {
    Map<String, String> parameters = new HashMap<>();
    parameters.put(OAuth2ParameterNames.ERROR, oauth2Error.getErrorCode());
    if (StringUtils.hasText(oauth2Error.getDescription())) {
      parameters.put(OAuth2ParameterNames.ERROR_DESCRIPTION, oauth2Error.getDescription());
    }
    if (StringUtils.hasText(oauth2Error.getUri())) {
      parameters.put(OAuth2ParameterNames.ERROR_URI, oauth2Error.getUri());
    }
    return JsonUtil.convertValue(R.ok(parameters), Map.class);
  }
}
