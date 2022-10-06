package top.wecoding.iam.server.mybatisplus.handlers;

import lombok.extern.slf4j.Slf4j;
import org.apache.ibatis.type.BaseTypeHandler;
import org.apache.ibatis.type.JdbcType;
import org.apache.ibatis.type.MappedJdbcTypes;
import org.apache.ibatis.type.MappedTypes;
import org.springframework.util.StringUtils;

import java.sql.CallableStatement;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.Arrays;
import java.util.Collection;

/**
 * @author liuyuhui
 * @date 2022/10/5
 * @qq 1515418211
 */
@Slf4j
@MappedTypes({Object.class})
@MappedJdbcTypes(JdbcType.VARCHAR)
public class CollectionToCommaDelimitedStringTypeHandler
    extends BaseTypeHandler<Collection<Object>> {

  @Override
  public void setNonNullParameter(
      PreparedStatement ps, int i, Collection<Object> parameter, JdbcType jdbcType)
      throws SQLException {
    ps.setString(i, toStr(parameter));
  }

  @Override
  public Collection<Object> getNullableResult(ResultSet rs, String columnName) throws SQLException {
    final String str = rs.getString(columnName);
    return StringUtils.hasText(str) ? null : parse(str);
  }

  @Override
  public Collection<Object> getNullableResult(ResultSet rs, int columnIndex) throws SQLException {
    final String str = rs.getString(columnIndex);
    return StringUtils.hasText(str) ? null : parse(str);
  }

  @Override
  public Collection<Object> getNullableResult(CallableStatement cs, int columnIndex)
      throws SQLException {
    final String str = cs.getString(columnIndex);
    return StringUtils.hasText(str) ? null : parse(str);
  }

  protected Collection<Object> parse(String str) {
    String[] strings = StringUtils.delimitedListToStringArray(str, ",");
    return Arrays.asList(strings);
  }

  protected String toStr(Collection<Object> obj) {
    return StringUtils.collectionToCommaDelimitedString(obj);
  }
}
