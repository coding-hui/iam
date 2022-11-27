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
@TableName("iam_organization")
@EqualsAndHashCode(callSuper = true)
public class Organization extends BaseEntity {

  private String organizationCode;

  private String organizationName;

  private String description;
}
