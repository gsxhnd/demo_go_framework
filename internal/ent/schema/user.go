package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"go_sample_code/internal/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Mixins of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Unique().NotEmpty().MaxLen(50).Comment("用户名"),
		field.String("email").Unique().NotEmpty().MaxLen(255).Comment("邮箱"),
		field.String("password").NotEmpty().Sensitive().Comment("密码(哈希存储)"),
		field.String("nickname").MaxLen(100).Optional().Comment("昵称"),
		field.String("avatar").MaxLen(500).Optional().Comment("头像URL"),
		field.String("phone").MaxLen(20).Optional().Comment("手机号"),
		field.Bool("is_active").Default(true).Comment("是否激活"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
		index.Fields("email").Unique(),
	}
}
