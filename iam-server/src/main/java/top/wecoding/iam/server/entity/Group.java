package top.wecoding.iam.server.entity;

import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.*;
import top.wecoding.mybatis.base.BaseEntity;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_group")
@EqualsAndHashCode(callSuper = true)
public class Group extends BaseEntity {

  @TableId private Long id;

  private String groupId;

  private String tenantId;

  private String groupName;

  private String groupCode;

  private String description;
}
