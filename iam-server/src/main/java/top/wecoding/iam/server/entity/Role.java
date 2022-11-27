package top.wecoding.iam.server.entity;

import com.baomidou.mybatisplus.annotation.TableName;
import lombok.*;
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
