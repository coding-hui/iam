package top.wecoding.iam.server.entity;

import com.baomidou.mybatisplus.annotation.TableName;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import top.wecoding.mybatis.base.BaseEntity;

/**
 * @author liuyuhui
 * @date 2022/11/6
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_role")
@EqualsAndHashCode(callSuper = true)
public class Role extends BaseEntity {

  private String roleCode;

  private String description;
}
