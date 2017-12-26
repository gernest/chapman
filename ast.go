package chapman

// NodeType is a javascript ast node type.
type NodeType uint

//common node  types
const (
	Identifier NodeType = iota
	PrivateName
	RegExpLiteral
	NullLiteral
	StringLiteral
	NumericLiteral
	Program
	ExpressionStatement
	BlockStatement
	EmptyStatement
	DebuggerStatement
	WithStatement
	ReturnStatement
	LabeledStatement
	BreakStatement
	ContinueStatement
	IfStatement
	SwitchStatement
	SwitchCase
	ThrowStatement
	TryStatement
	CatchClause
	WhileStatement
	DoWhileStatement
	ForStatement
	ForInStatement
	ForOfStatement
	FunctionDeclaration
	VariableDeclaration
	VariableDeclarator
	Decorator
	Directive
	DirectiveLiteral
	Super
	Import
	ThisExpression
	ArrowFunctionExpression
	YieldExpression
	AwaitExpression
	ArrayExpression
	ObjectExpression
	ObjectProperty
	ObjectMethod
	FunctionExpression
	UnaryExpression
	UpdateExpression
	BinaryExpression
	AssignmentExpression
	SpreadElement
	MemberExpression
	BindExpression
	ConditionalExpression
	CallExpression
	NewExpression
	SequenceExpression
	DoExpression
	TemplateLiteral
	TaggedTemplateExpression
	TemplateElement
	ObjectPattern
	ArrayPattern
	RestElement
	AssignmentPattern
	ClassBody
	ClassMethod
	ClassPrivateMethod
	ClassProperty
	ClassPrivateProperty
	ClassDeclaration
	ClassExpression
	MetaProperty
	ImportDeclaration
	ImportSpecifier
	ImportDefaultSpecifier
	ImportNamespaceSpecifier
	ExportNamedDeclaration
	ExportSpecifier
	ExportDefaultDeclaration
	ExportAllDeclaration
)
