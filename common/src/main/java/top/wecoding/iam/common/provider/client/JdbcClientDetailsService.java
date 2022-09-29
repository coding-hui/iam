package top.wecoding.iam.common.provider.client;

import lombok.AllArgsConstructor;
import org.springframework.jdbc.core.BeanPropertyRowMapper;
import org.springframework.jdbc.core.JdbcTemplate;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.iam.common.provider.ClientDetails;
import top.wecoding.iam.common.provider.ClientDetailsService;

/**
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@AllArgsConstructor
public class JdbcClientDetailsService implements ClientDetailsService {

  private final JdbcTemplate jdbcTemplate;

  @Override
  public ClientDetails loadClientByClientId(String clientId) {
    try {
      return jdbcTemplate.queryForObject(
          SecurityConstants.DEFAULT_SELECT_STATEMENT,
          new BeanPropertyRowMapper<>(BaseClientDetails.class),
          clientId);
    } catch (Exception ex) {
      return null;
    }
  }
}
